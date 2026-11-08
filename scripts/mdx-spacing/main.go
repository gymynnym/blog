package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

func git(root string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", append([]string{"-c", "core.fsmonitor=false"}, args...)...)
	cmd.Dir = root
	return cmd.Output()
}

func filesToCheck(root string, staged bool, inputs []string) ([]string, error) {
	var files []string
	if staged {
		output, err := git(root, "diff", "--cached", "--name-only", "--diff-filter=ACMR", "-z", "--", "src/content/posts/ko", "src/content/posts/ja")
		if err != nil {
			return nil, fmt.Errorf("staged 파일 조회 실패: %w", err)
		}
		files = strings.Split(string(output), "\x00")
	} else if len(inputs) > 0 {
		for _, input := range inputs {
			if !filepath.IsAbs(input) {
				input = filepath.Join(root, input)
			}
			path, err := filepath.Rel(root, input)
			if err != nil {
				return nil, err
			}
			files = append(files, filepath.ToSlash(path))
		}
	} else {
		for _, locale := range []string{"ko", "ja"} {
			err := filepath.WalkDir(filepath.Join(root, "src/content/posts", locale), func(path string, entry fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if entry.Type().IsRegular() {
					relative, err := filepath.Rel(root, path)
					if err != nil {
						return err
					}
					files = append(files, filepath.ToSlash(relative))
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
	}
	files = slices.DeleteFunc(files, func(path string) bool { return !targetMDX(path) })
	slices.Sort(files)
	return slices.Compact(files), nil
}

func run(root string, args []string, stdout, stderr io.Writer) int {
	var staged, fixFiles bool
	var inputs []string
	for _, arg := range args {
		switch arg {
		case "--staged":
			staged = true
		case "--fix":
			fixFiles = true
		default:
			if strings.HasPrefix(arg, "-") {
				fmt.Fprintf(stderr, "알 수 없는 옵션: %s\n", arg)
				return 1
			}
			inputs = append(inputs, arg)
		}
	}
	if staged && (fixFiles || len(inputs) != 0) {
		fmt.Fprintln(stderr, "--staged는 --fix 또는 파일 경로와 함께 사용할 수 없습니다. 수정 후 직접 stage해 주세요.")
		return 1
	}
	files, err := filesToCheck(root, staged, inputs)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	errors, changed := 0, 0
	for _, file := range files {
		path := filepath.Join(root, filepath.FromSlash(file))
		var source []byte
		if staged {
			source, err = git(root, "show", ":"+file)
		} else {
			source, err = os.ReadFile(path)
		}
		if err != nil {
			fmt.Fprintf(stderr, "%s: 읽기 실패: %v\n", file, err)
			errors++
			continue
		}
		if fixFiles {
			updated := fix(source)
			if !bytes.Equal(source, updated) {
				if err := os.WriteFile(path, updated, 0o644); err != nil {
					fmt.Fprintf(stderr, "%s: 수정 실패: %v\n", file, err)
					errors++
					continue
				}
				changed++
				source = updated
			}
		}
		for _, item := range lint(source) {
			fmt.Fprintf(stderr, "%s:%d:%d: %s\n", file, item.Line, item.Column, item.Message)
			errors++
		}
	}
	if errors != 0 {
		fmt.Fprintf(stderr, "MDX 강조 검사 실패: %d건. pnpm mdx:fix 실행 후 수정된 내용을 stage해 주세요.\n", errors)
		return 1
	}
	suffix := ""
	if staged {
		suffix = " (staged)"
	}
	if fixFiles {
		suffix += fmt.Sprintf(", %d개 수정", changed)
	}
	fmt.Fprintf(stdout, "MDX 강조 검사 통과: %d개 파일%s.\n", len(files), suffix)
	return 0
}

func main() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if top, err := git(root, "rev-parse", "--show-toplevel"); err == nil {
		root = strings.TrimSpace(string(top))
	}
	os.Exit(run(root, os.Args[1:], os.Stdout, os.Stderr))
}
