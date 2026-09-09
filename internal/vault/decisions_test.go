package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func writeDecision(t *testing.T, vaultDir, filename, content string) {
	t.Helper()
	dir := filepath.Join(vaultDir, "decisions")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, filename), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadDecisions(t *testing.T) {
	vaultDir := t.TempDir()
	writeDecision(t, vaultDir, "2026-07-22-billing-webhook-retries.md", `---
tags: [decision]
date: 2026-07-22
status: active
repo: example
x-governs:
  - "src/billing/**"
  - "src/webhooks/stripe.ts"
---

# Some decision
`)
	writeDecision(t, vaultDir, "2025-01-01-old-decision.md", `---
tags: [decision]
date: 2025-01-01
status: superseded
repo: example
---

# An old decision
`)

	decisions, err := LoadDecisions(vaultDir)
	if err != nil {
		t.Fatalf("LoadDecisions failed: %v", err)
	}
	if len(decisions) != 2 {
		t.Fatalf("expected 2 decisions, got %d", len(decisions))
	}

	active := FilterByStatus(decisions, []string{"active"})
	if len(active) != 1 {
		t.Fatalf("expected 1 active decision, got %d", len(active))
	}
	if active[0].Slug != "2026-07-22-billing-webhook-retries" {
		t.Fatalf("unexpected slug: %s", active[0].Slug)
	}
	if len(active[0].XGoverns) != 2 || active[0].XGoverns[0] != "src/billing/**" {
		t.Fatalf("unexpected x-governs: %v", active[0].XGoverns)
	}
}

func TestLoadDecisions_MissingDirReturnsEmpty(t *testing.T) {
	vaultDir := t.TempDir()
	decisions, err := LoadDecisions(vaultDir)
	if err != nil {
		t.Fatalf("expected no error for a missing decisions dir, got: %v", err)
	}
	if len(decisions) != 0 {
		t.Fatalf("expected 0 decisions, got %d", len(decisions))
	}
}
