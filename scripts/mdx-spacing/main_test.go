package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func repository(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	gitOK(t, root, "init", "--quiet")
	return root
}

func gitOK(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-c", "core.fsmonitor=false"}, args...)...)
	cmd.Dir = root
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
}

func writeFile(t *testing.T, root, path, source string) {
	t.Helper()
	path = filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
}

func check(root string, args ...string) (int, string) {
	var output bytes.Buffer
	code := run(root, args, &output, &output)
	return code, output.String()
}

func TestStagedContentNotWorkingCopy(t *testing.T) {
	root := repository(t)
	path := "src/content/posts/ja/file with spaces.mdx"
	writeFile(t, root, path, "# 見出し\n\n前**内容（補足）**後\n")
	gitOK(t, root, "add", "--", path)
	writeFile(t, root, path, "# 見出し\n\n前**内容（補足）** 後\n")
	if code, output := check(root, "--staged"); code != 1 || !strings.Contains(output, "file with spaces.mdx:3:") {
		t.Fatalf("invalid index content was accepted: %d %s", code, output)
	}
	gitOK(t, root, "add", "--", path)
	writeFile(t, root, path, "# 見出し\n\n前**内容（補足）**後\n")
	if code, output := check(root, "--staged"); code != 0 {
		t.Fatalf("valid index content was rejected: %s", output)
	}
}

func TestEnglishAndDeletedFilesAreIgnored(t *testing.T) {
	root := repository(t)
	path := "src/content/posts/ja/deleted.mdx"
	writeFile(t, root, path, "**正常**です。\n")
	gitOK(t, root, "add", ".")
	gitOK(t, root, "-c", "user.name=MDX tests", "-c", "user.email=mdx@example.invalid", "-c", "commit.gpgsign=false", "commit", "--quiet", "-m", "test fixture")
	if err := os.Remove(filepath.Join(root, path)); err != nil {
		t.Fatal(err)
	}
	writeFile(t, root, "src/content/posts/en/article.mdx", "前**内容（補足）**後\n")
	gitOK(t, root, "add", "-A")
	if code, output := check(root, "--staged"); code != 0 || !strings.Contains(output, "0개 파일") {
		t.Fatalf("English/deletion should be ignored: %d %s", code, output)
	}
}

func TestEmptyIndexAndInvalidOptions(t *testing.T) {
	root := repository(t)
	if code, output := check(root, "--staged"); code != 0 {
		t.Fatal(output)
	}
	for _, args := range [][]string{{"--staged", "--fix"}, {"--staged", "file.mdx"}, {"--unknown"}} {
		if code, _ := check(root, args...); code != 1 {
			t.Fatalf("expected failure for %v", args)
		}
	}
}

func TestFixOnlyChangesWorkingCopy(t *testing.T) {
	root := repository(t)
	path := "src/content/posts/ja/article.mdx"
	original := "前**内容（補足）**後\n"
	writeFile(t, root, path, original)
	gitOK(t, root, "add", "--", path)
	if code, output := check(root, "--fix", path); code != 0 {
		t.Fatal(output)
	}
	if indexed, err := git(root, "show", ":"+path); err != nil || string(indexed) != original {
		t.Fatal("automatic fix changed the index", err)
	}
	if code, _ := check(root, "--staged"); code != 1 {
		t.Fatal("staged content should still fail until explicitly staged again")
	}
}

func TestLefthookIntegration(t *testing.T) {
	lefthook, err := exec.LookPath("lefthook")
	if err != nil {
		t.Skip("lefthook is not installed")
	}
	root := repository(t)
	project, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	if path, err := exec.Command("mise", "which", "lefthook").Output(); err == nil {
		lefthook = strings.TrimSpace(string(path))
	}
	for _, path := range []string{"lefthook.yml", "scripts/go.mod", "scripts/go.sum", "scripts/mdx-spacing/main.go", "scripts/mdx-spacing/rules.go"} {
		data, err := os.ReadFile(filepath.Join(project, path))
		if err != nil {
			t.Fatal(err)
		}
		writeFile(t, root, path, string(data))
	}
	gitOK(t, root, "-c", "user.name=MDX tests", "-c", "user.email=mdx@example.invalid", "-c", "commit.gpgsign=false", "commit", "--quiet", "--allow-empty", "-m", "test fixture")
	invoke := func(args ...string) (string, error) {
		cmd := exec.Command(lefthook, args...)
		cmd.Dir = root
		// Keep the test repository independent of global mise shim defaults.
		toolPath := filepath.Join(runtime.GOROOT(), "bin") + string(os.PathListSeparator) + filepath.Dir(lefthook)
		cmd.Env = append(os.Environ(), "PATH="+toolPath+string(os.PathListSeparator)+os.Getenv("PATH"))
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	if output, err := invoke("install"); err != nil {
		t.Fatal(err, output)
	}
	path := "src/content/posts/ja/article.mdx"
	writeFile(t, root, path, "説明**IDE（開発環境）**です。\n")
	gitOK(t, root, "add", "--", path)
	writeFile(t, root, path, "説明**IDE（開発環境）** です。\n")
	if output, err := invoke("run", "pre-commit"); err == nil || !strings.Contains(output, "MDX 강조 검사 실패") {
		t.Fatalf("bad staged MDX should be blocked: %v\n%s", err, output)
	}
	gitOK(t, root, "add", "--", path)
	if output, err := invoke("run", "pre-commit"); err != nil {
		t.Fatalf("fixed staged MDX should pass: %v\n%s", err, output)
	}
}
