package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultsWhenMissing(t *testing.T) {
	dir := t.TempDir()
	c, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(c.GovernedStatuses) != 1 || c.GovernedStatuses[0] != "active" {
		t.Fatalf("unexpected default governed_statuses: %v", c.GovernedStatuses)
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	content := `governed_statuses: [active, in-progress]
exempt_paths:
  - "**/*.md"
  - "docs/**"
`
	if err := os.WriteFile(filepath.Join(dir, "provenance.yml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	c, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(c.GovernedStatuses) != 2 {
		t.Fatalf("unexpected governed_statuses: %v", c.GovernedStatuses)
	}
	if len(c.ExemptPaths) != 2 {
		t.Fatalf("unexpected exempt_paths: %v", c.ExemptPaths)
	}
}

func TestScaffoldRefusesToOverwrite(t *testing.T) {
	dir := t.TempDir()
	if err := Scaffold(dir); err != nil {
		t.Fatalf("first Scaffold failed: %v", err)
	}
	if err := Scaffold(dir); err == nil {
		t.Fatal("expected second Scaffold to fail, it did not")
	}
}

func TestMatchesAny(t *testing.T) {
	patterns := []string{"**/*.md", "docs/**"}
	cases := map[string]bool{
		"README.md":              true,
		"docs/guide.txt":         true,
		"src/billing/webhook.ts": false,
	}
	for path, want := range cases {
		if got := MatchesAny(path, patterns); got != want {
			t.Errorf("MatchesAny(%q) = %v, want %v", path, got, want)
		}
	}
}
