package cmdcheck

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

func writeBillingDecision(t *testing.T, dir string) {
	t.Helper()
	decisionsDir := filepath.Join(dir, "vault", "decisions")
	if err := os.MkdirAll(decisionsDir, 0755); err != nil {
		t.Fatal(err)
	}
	decision := `---
tags: [decision]
date: 2026-07-22
status: active
repo: example
x-governs: ["src/billing/**"]
---

# Billing decision
`
	if err := os.WriteFile(filepath.Join(decisionsDir, "2026-07-22-billing.md"), []byte(decision), 0644); err != nil {
		t.Fatal(err)
	}
}

func touchBillingFile(t *testing.T, dir string) {
	t.Helper()
	billingDir := filepath.Join(dir, "src", "billing")
	if err := os.MkdirAll(billingDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(billingDir, "webhook.go"), []byte("package billing"), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestRun_BlocksOnUnaddressedGovernedPath(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init", "-q")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "test")

	writeBillingDecision(t, dir)
	runGit(t, dir, "add", "-A")
	runGit(t, dir, "commit", "-q", "-m", "base")
	base := gitRevParse(t, dir)

	touchBillingFile(t, dir)
	runGit(t, dir, "add", "-A")
	runGit(t, dir, "commit", "-q", "-m", "touch billing, no capture")
	head := gitRevParse(t, dir)

	jsonOut := filepath.Join(t.TempDir(), "report.json")
	code := Run([]string{"--repo", dir, "--base", base, "--head", head, "--json-out", jsonOut})
	if code != 1 {
		t.Fatalf("expected exit code 1 (blocked), got %d", code)
	}
	if _, err := os.Stat(jsonOut); err != nil {
		t.Fatalf("expected JSON report to be written: %v", err)
	}
}

func TestRun_PassesWhenAddressedByCapture(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init", "-q")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "test")

	writeBillingDecision(t, dir)
	runGit(t, dir, "add", "-A")
	runGit(t, dir, "commit", "-q", "-m", "base")
	base := gitRevParse(t, dir)

	touchBillingFile(t, dir)
	capturesDir := filepath.Join(dir, "vault", "captures")
	if err := os.MkdirAll(capturesDir, 0755); err != nil {
		t.Fatal(err)
	}
	capture := `---
tags: [capture]
date: 2026-09-08
status: active
repo: example
decisions: ["2026-07-22-billing"]
attribution: human
---

## What
Widened the retry window.
`
	if err := os.WriteFile(filepath.Join(capturesDir, "2026-09-08-widen.md"), []byte(capture), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "-A")
	runGit(t, dir, "commit", "-q", "-m", "touch billing, with capture")
	head := gitRevParse(t, dir)

	code := Run([]string{"--repo", dir, "--base", base, "--head", head})
	if code != 0 {
		t.Fatalf("expected exit code 0 (pass), got %d", code)
	}
}

func TestRun_RejectsMissingBase(t *testing.T) {
	code := Run([]string{"--repo", "."})
	if code != 2 {
		t.Fatalf("expected exit code 2 (usage error) when --base is missing, got %d", code)
	}
}
