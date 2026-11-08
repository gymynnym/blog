package main

import (
	"bytes"
	"regexp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
)

type issue struct {
	Offset, Line, Column int
	Message              string
	Fixable              bool
}

type marker struct {
	start, stop int
	char        byte
}

var (
	frontmatter = regexp.MustCompile(`(?ms)\A(?:\x{FEFF})?---[ \t]*\r?\n.*?^---[ \t]*\r?$`)
	mdxComment  = regexp.MustCompile(`(?s)\{/\*.*?\*/\}`)
	mdxImport   = regexp.MustCompile(`(?ms)^(?:import|export)[ \t].*?(?:\r?\n[ \t]*\r?\n|\z)`)
	boldMarker  = regexp.MustCompile(`\*{2,}|_{2,}`)
)

// Mask MDX-only syntax while preserving byte offsets and original line endings.
func markdownSource(source []byte) []byte {
	masked := bytes.Clone(source)
	for _, pattern := range []*regexp.Regexp{frontmatter, mdxComment, mdxImport} {
		for _, span := range pattern.FindAllIndex(masked, -1) {
			for i := span[0]; i < span[1]; i++ {
				if masked[i] != '\n' && masked[i] != '\r' {
					masked[i] = ' '
				}
			}
		}
	}
	return masked
}

func escaped(source []byte, offset int) bool {
	count := 0
	for offset--; offset >= 0 && source[offset] == '\\'; offset-- {
		count++
	}
	return count%2 == 1
}

func space(r rune) bool { return unicode.IsSpace(r) || r == '\uFEFF' }

// CommonMark 0.31.2 §6.2 uses Unicode punctuation AND symbol categories.
// In particular, ） and 。 behave differently from the last letter in 관점.
func punctuation(r rune) bool { return unicode.IsPunct(r) || unicode.IsSymbol(r) }

func before(source []byte, offset int) rune {
	if offset == 0 {
		return ' '
	}
	r, _ := utf8.DecodeLastRune(source[:offset])
	return r
}

func after(source []byte, offset int) rune {
	if offset == len(source) {
		return ' '
	}
	r, _ := utf8.DecodeRune(source[offset:])
	return r
}

func flanking(source []byte, m marker) (open, close bool) {
	prev, next := before(source, m.start), after(source, m.stop)
	left := !space(next) && (!punctuation(next) || space(prev) || punctuation(prev))
	right := !space(prev) && (!punctuation(prev) || space(next) || punctuation(next))
	if m.char == '_' {
		return left && (!right || punctuation(prev)), right && (!left || punctuation(next))
	}
	return left, right
}

func lint(source []byte) []issue {
	masked := markdownSource(source)
	doc := parser.New(parser.WithExtensions(extension.NewGFMParser())).Parse(masked)
	groups := make(map[ast.Node][]marker)

	// Inspect only delimiters that remain literal text after CommonMark parsing.
	// Valid **성능 관점**에서 and adjacent inline code never enter this check.
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n.(type) {
		case *ast.CodeSpan, *ast.CodeBlock,
			*ast.Image, *ast.AutoLink, *ast.RawHTML, *ast.HTMLBlock:
			return ast.WalkSkipChildren, nil
		}
		t, ok := n.(*ast.Text)
		if !ok || t.Value.IsOwned() {
			return ast.WalkContinue, nil
		}
		block := n.Parent()
		for block != nil {
			if _, ok := block.(ast.BlockNode); ok {
				break
			}
			block = block.Parent()
		}
		span := t.Value.Index()
		for _, match := range boldMarker.FindAllIndex(masked[span.Start:span.Stop], -1) {
			start, stop := span.Start+match[0], span.Start+match[1]
			if !escaped(source, start) {
				groups[block] = append(groups[block], marker{start, stop, source[start]})
			}
		}
		return ast.WalkContinue, nil
	})

	issues := make(map[int]issue)
	report := func(offset int, message string, fixable bool) {
		if old, exists := issues[offset]; exists && !old.Fixable {
			return
		}
		prefix := source[:offset]
		lineStart := bytes.LastIndexByte(prefix, '\n') + 1
		issues[offset] = issue{
			Offset: offset, Line: bytes.Count(prefix, []byte{'\n'}) + 1,
			Column:  utf8.RuneCount(source[lineStart:offset]) + 1,
			Message: message, Fixable: fixable,
		}
	}
	unparsed := func(m marker) {
		report(m.start, "굵게 구분자가 일반 텍스트로 남습니다. 구분자를 수정하거나 이스케이프하세요.", false)
	}
	for _, markers := range groups {
		for i := 0; i < len(markers); {
			opener := markers[i]
			if i+1 >= len(markers) || opener.stop-opener.start != 2 {
				unparsed(opener)
				i++
				continue
			}
			closer := markers[i+1]
			if closer.char != opener.char || closer.stop-closer.start != 2 {
				unparsed(opener)
				i++
				continue
			}
			i += 2
			canOpen, _ := flanking(source, opener)
			_, canClose := flanking(source, closer)
			if space(after(source, opener.stop)) || space(before(source, closer.start)) || (canOpen && canClose) {
				unparsed(opener)
				continue
			}
			if !canOpen {
				report(opener.start, "강조 시작 구분자가 해석되지 않습니다. 여기에 공백을 추가하세요.", true)
			}
			if !canClose {
				report(closer.stop, "강조 끝 구분자가 해석되지 않습니다. 여기에 공백을 추가하세요.", true)
			}
		}
	}
	result := make([]issue, 0, len(issues))
	for _, item := range issues {
		result = append(result, item)
	}
	slices.SortFunc(result, func(a, b issue) int { return a.Offset - b.Offset })
	return result
}

func fix(source []byte) []byte {
	issues := lint(source)
	for _, item := range issues {
		if !item.Fixable {
			return source
		}
	}
	if len(issues) == 0 {
		return source
	}
	var out bytes.Buffer
	start := 0
	for _, item := range issues {
		out.Write(source[start:item.Offset])
		out.WriteByte(' ')
		start = item.Offset
	}
	out.Write(source[start:])
	updated := out.Bytes()
	// Never write an automatic edit unless it actually removes the parse failure.
	if len(lint(updated)) != 0 {
		return source
	}
	return updated
}

func targetMDX(path string) bool {
	path = strings.ReplaceAll(path, `\`, "/")
	return (strings.HasPrefix(path, "src/content/posts/ko/") || strings.HasPrefix(path, "src/content/posts/ja/")) && strings.HasSuffix(path, ".mdx")
}
