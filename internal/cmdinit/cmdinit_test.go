package cmdinit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun_ScaffoldsConfig(t *testing.T) {
	dir := t.TempDir()
	code := Run([]string{"--repo", dir})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if _, err := os.Stat(filepath.Join(dir, "provenance.yml")); err != nil {
		t.Fatalf("expected provenance.yml to exist: %v", err)
	}
}

func TestRun_RefusesToOverwrite(t *testing.T) {
	dir := t.TempDir()
	if code := Run([]string{"--repo", dir}); code != 0 {
		t.Fatalf("first run: expected 0, got %d", code)
	}
	if code := Run([]string{"--repo", dir}); code == 0 {
		t.Fatal("second run: expected non-zero (refuse overwrite), got 0")
	}
}
