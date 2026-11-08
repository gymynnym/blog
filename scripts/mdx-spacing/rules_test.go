package main

import (
	"bytes"
	"testing"
)

func TestValidMarkdownNeedsNoSurroundingSpaces(t *testing.T) {
	for _, source := range []string{
		"이번 아티클에서는 **성능 관점**에서 설명하겠습니다.",
		"한글**성능 관점**에서 설명합니다.",
		"日本語**重要**です。",
		"앞`code`뒤",
		"日本語``a`b **literal**``です。",
		"(**강조**)",
		"**強調**。",
		"（**内容（補足）**）",
		"**`code`**",
		"***강조***",
		"[**강조**](https://example.com)",
		"앞[`code`](https://example.com)뒤",
		"**강조**\r\n`code`\r\n",
	} {
		t.Run(source, func(t *testing.T) {
			if got := lint([]byte(source)); len(got) != 0 {
				t.Fatalf("valid Markdown was rejected: %+v", got)
			}
			if got := string(fix([]byte(source))); got != source {
				t.Fatalf("valid source changed: %q", got)
			}
		})
	}
}

func TestBrokenEmphasisOnlyGetsNecessarySpaces(t *testing.T) {
	tests := []struct{ source, want string }{
		{"> ADEとは：**IDE（Integrated Development Environment）**が開発の中心です。", "> ADEとは：**IDE（Integrated Development Environment）** が開発の中心です。"},
		{"앞**(설명)**뒤", "앞 **(설명)** 뒤"},
		{"これは**重要。**です。", "これは**重要。** です。"},
		{"**末尾の`code`**です。", "**末尾の`code`** です。"},
		{"**[link](https://example.com)**です。", "**[link](https://example.com)** です。"},
		{"前__重要__後", "前 __重要__ 後"},
		{"| 項目 |\n| --- |\n| **値（単位）**です |\n", "| 項目 |\n| --- |\n| **値（単位）** です |\n"},
	}
	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			if len(lint([]byte(tt.source))) == 0 {
				t.Fatal("broken emphasis was accepted")
			}
			got := fix([]byte(tt.source))
			if string(got) != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
			if issues := lint(got); len(issues) != 0 {
				t.Fatalf("fixed Markdown was rejected: %+v", issues)
			}
			if !bytes.Equal(fix(got), got) {
				t.Fatal("fix is not idempotent")
			}
		})
	}
}

func TestADEDiagnosticUsesOriginalLineAndColumn(t *testing.T) {
	source := []byte("# 見出し\n\n> ADEとは：**IDE（Integrated Development Environment）**が開発の中心です。")
	got := lint(source)
	if len(got) != 1 || got[0].Line != 3 || got[0].Column != 52 || !got[0].Fixable {
		t.Fatalf("unexpected diagnostic: %+v", got)
	}
}

func TestNonProseIsIgnored(t *testing.T) {
	source := []byte("---\r\ntitle: '前**（題名）**後'\r\n---\r\n\r\n" +
		"```js\nconst x = '前**（値）**後';\n```\n\n" +
		"> ~~~md\n> 前**（値）**後\n> ~~~\n\n" +
		"    前**（コード）**後\n\n" +
		"<Component title=\"前**（値）**後\" />\n\n" +
		"{/* 前**（注釈）**後 */}\n\n" +
		"import X from '前**（値）**後';\n\n" +
		"[リンク](https://example.com/`literal`/**path** \"前**（題名）**後\")\n\n" +
		"https://example.com/**path**/file\n\n" +
		"<https://example.com/**path**>\n\n" +
		"![前**（代替テキスト）**後](image.png)\n\n" +
		"前\\**星\\**後\n\n***\n")
	if got := lint(source); len(got) != 0 {
		t.Fatalf("non-prose syntax was inspected: %+v", got)
	}
	if !bytes.Equal(fix(source), source) {
		t.Fatal("non-prose syntax changed")
	}
}

func TestAmbiguousMarkersAreNotAutomaticallyRewritten(t *testing.T) {
	for _, source := range []string{"日本語 **閉じていない強調", "**하나****둘**", "** whitespace **"} {
		issues := lint([]byte(source))
		if len(issues) == 0 {
			t.Fatalf("unparsed marker was accepted: %q", source)
		}
		if got := string(fix([]byte(source))); got != source {
			t.Fatalf("ambiguous source changed to %q", got)
		}
	}
}

func TestFixPreservesFrontmatterCodeAndLineEndings(t *testing.T) {
	source := []byte("---\r\ntitle: '題名'\r\n---\r\n\r\n前**内容（補足）**後\r\n\r\n```js\r\nconst x = 'a**b**c';\r\n```\r\n")
	want := bytes.Replace(source, []byte("**後"), []byte("** 後"), 1)
	if got := fix(source); !bytes.Equal(got, want) {
		t.Fatalf("unexpected changes: %q", got)
	}
}

func TestTargetMDX(t *testing.T) {
	for _, path := range []string{"src/content/posts/ko/ade.mdx", "src/content/posts/ja/nested/post.mdx"} {
		if !targetMDX(path) {
			t.Fatal(path)
		}
	}
	for _, path := range []string{"src/content/posts/en/ade.mdx", "README.md", "src/content/posts/ja/a.md", ""} {
		if targetMDX(path) {
			t.Fatal(path)
		}
	}
}
