package gitutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

func gitRevParse(t *testing.T, dir string) string {
	t.Helper()
	out, err := exec.Command("git", "-C", dir, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatalf("git rev-parse failed: %v", err)
	}
	s := string(out)
	return s[:len(s)-1]
}

func TestChangedFiles(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init", "-q")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "test")

	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("one"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "a.txt")
	runGit(t, dir, "commit", "-q", "-m", "base")
	base := gitRevParse(t, dir)

	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("two"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("one-changed"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "-A")
	runGit(t, dir, "commit", "-q", "-m", "head")
	head := gitRevParse(t, dir)

	files, err := ChangedFiles(dir, base, head)
	if err != nil {
		t.Fatalf("ChangedFiles failed: %v", err)
	}

	want := map[string]bool{"a.txt": true, "b.txt": true}
	if len(files) != len(want) {
		t.Fatalf("expected %d files, got %d: %v", len(want), len(files), files)
	}
	for _, f := range files {
		if !want[f] {
			t.Fatalf("unexpected file in diff: %s", f)
		}
	}
}

func TestParsePrePushStdin(t *testing.T) {
	input := "refs/heads/main abc123 refs/heads/main def456\n"
	refs, err := ParsePrePushStdin(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].LocalSHA != "abc123" || refs[0].RemoteSHA != "def456" {
		t.Fatalf("unexpected parse: %+v", refs[0])
	}
}

func TestParsePrePushStdin_RejectsMalformed(t *testing.T) {
	if _, err := ParsePrePushStdin("not-well-formed\n"); err == nil {
		t.Fatal("expected an error for a malformed line")
	}
}
