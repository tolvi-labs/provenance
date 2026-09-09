package hook

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallAndUninstall(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".git", "hooks"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := Install(dir, false); err != nil {
		t.Fatalf("Install failed: %v", err)
	}
	info, err := os.Stat(hookPath(dir))
	if err != nil {
		t.Fatalf("expected hook file: %v", err)
	}
	if info.Mode()&0111 == 0 {
		t.Fatalf("expected hook to be executable, mode = %v", info.Mode())
	}

	if err := Install(dir, false); err == nil {
		t.Fatal("expected second Install without force to fail")
	}
	if err := Install(dir, true); err != nil {
		t.Fatalf("Install with force=true should succeed: %v", err)
	}

	if err := Uninstall(dir); err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}
	if _, err := os.Stat(hookPath(dir)); !os.IsNotExist(err) {
		t.Fatalf("expected hook to be removed, got err: %v", err)
	}
}

func TestUninstall_RefusesForeignHook(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".git", "hooks"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hookPath(dir), []byte("#!/bin/sh\necho not provenance\n"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := Uninstall(dir); err == nil {
		t.Fatal("expected Uninstall to refuse a foreign hook")
	}
	if _, err := os.Stat(hookPath(dir)); err != nil {
		t.Fatal("foreign hook should still exist")
	}
}
