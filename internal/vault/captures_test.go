package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadAndValidateCaptures(t *testing.T) {
	vaultDir := t.TempDir()
	writeDecision(t, vaultDir, "2026-07-22-billing-webhook-retries.md", `---
tags: [decision]
date: 2026-07-22
status: active
repo: example
x-governs: ["src/billing/**"]
---

# Some decision
`)

	capDir := filepath.Join(vaultDir, "captures")
	if err := os.MkdirAll(capDir, 0755); err != nil {
		t.Fatal(err)
	}
	good := `---
tags: [capture]
date: 2026-09-08
status: active
repo: example
decisions: ["2026-07-22-billing-webhook-retries"]
attribution: human
---

## What
Widened the webhook retry window.
`
	if err := os.WriteFile(filepath.Join(capDir, "2026-09-08-widen-retry.md"), []byte(good), 0644); err != nil {
		t.Fatal(err)
	}
	bad := `---
tags: [capture]
date: 2026-09-08
status: active
repo: example
decisions: ["does-not-exist"]
attribution: bogus
---
`
	if err := os.WriteFile(filepath.Join(capDir, "2026-09-08-bad-capture.md"), []byte(bad), 0644); err != nil {
		t.Fatal(err)
	}

	decisions, err := LoadDecisions(vaultDir)
	if err != nil {
		t.Fatalf("LoadDecisions failed: %v", err)
	}
	captures, err := LoadCaptures(vaultDir)
	if err != nil {
		t.Fatalf("LoadCaptures failed: %v", err)
	}
	if len(captures) != 2 {
		t.Fatalf("expected 2 captures, got %d", len(captures))
	}

	var goodCapture, badCapture Capture
	for _, c := range captures {
		if strings.Contains(c.Path, "widen-retry") {
			goodCapture = c
		} else {
			badCapture = c
		}
	}

	if problems := ValidateCapture(goodCapture, decisions); len(problems) != 0 {
		t.Fatalf("expected no problems for good capture, got: %v", problems)
	}
	problems := ValidateCapture(badCapture, decisions)
	if len(problems) != 2 {
		t.Fatalf("expected 2 problems for bad capture, got %d: %v", len(problems), problems)
	}
}

func TestLoadCaptures_MissingDirReturnsEmpty(t *testing.T) {
	vaultDir := t.TempDir()
	captures, err := LoadCaptures(vaultDir)
	if err != nil {
		t.Fatalf("expected no error for a missing captures dir, got: %v", err)
	}
	if len(captures) != 0 {
		t.Fatalf("expected 0 captures, got %d", len(captures))
	}
}
