package gate

import (
	"testing"

	"github.com/tolvi-labs/provenance/internal/config"
	"github.com/tolvi-labs/provenance/internal/vault"
)

func testConfig() config.Config {
	return config.Config{
		GovernedStatuses: []string{"active"},
		ExemptPaths:      []string{"**/*.md"},
	}
}

func TestCheck_NoGovernedFilesPasses(t *testing.T) {
	decisions := []vault.Decision{
		{Slug: "dec-1", Status: "active", XGoverns: []string{"src/billing/**"}},
	}
	result := Check([]string{"src/other/file.go"}, decisions, nil, testConfig())
	if result.Status != "pass" {
		t.Fatalf("expected pass, got %s (blocked on: %+v)", result.Status, result.BlockedOn)
	}
}

func TestCheck_AddressedByCapturePasses(t *testing.T) {
	decisions := []vault.Decision{
		{Slug: "dec-1", Status: "active", XGoverns: []string{"src/billing/**"}},
	}
	captures := []vault.Capture{
		{Path: "vault/captures/2026-09-08-x.md", Decisions: []string{"dec-1"}, Attribution: "human"},
	}
	result := Check([]string{"src/billing/webhook.go"}, decisions, captures, testConfig())
	if result.Status != "pass" {
		t.Fatalf("expected pass, got %s (blocked on: %+v)", result.Status, result.BlockedOn)
	}
}

func TestCheck_AddressedByNewDecisionPasses(t *testing.T) {
	decisions := []vault.Decision{
		{Slug: "2026-09-08-new-dec", Status: "active", XGoverns: []string{"src/billing/**"}},
	}
	changed := []string{"src/billing/webhook.go", "vault/decisions/2026-09-08-new-dec.md"}
	result := Check(changed, decisions, nil, testConfig())
	if result.Status != "pass" {
		t.Fatalf("expected pass, got %s (blocked on: %+v)", result.Status, result.BlockedOn)
	}
	if len(result.NewDecisions) != 1 || result.NewDecisions[0] != "2026-09-08-new-dec" {
		t.Fatalf("expected NewDecisions to include the slug, got %v", result.NewDecisions)
	}
}

func TestCheck_UnaddressedBlocks(t *testing.T) {
	decisions := []vault.Decision{
		{Slug: "dec-1", Status: "active", XGoverns: []string{"src/billing/**"}},
	}
	result := Check([]string{"src/billing/webhook.go"}, decisions, nil, testConfig())
	if result.Status != "blocked" {
		t.Fatalf("expected blocked, got %s", result.Status)
	}
	if len(result.BlockedOn) != 1 || result.BlockedOn[0].Decision != "dec-1" {
		t.Fatalf("unexpected BlockedOn: %+v", result.BlockedOn)
	}
	if len(result.BlockedOn[0].Files) != 1 || result.BlockedOn[0].Files[0] != "src/billing/webhook.go" {
		t.Fatalf("unexpected blocked files: %v", result.BlockedOn[0].Files)
	}
}

func TestCheck_ExemptPathSkipped(t *testing.T) {
	decisions := []vault.Decision{
		{Slug: "dec-1", Status: "active", XGoverns: []string{"**/*.md"}},
	}
	result := Check([]string{"README.md"}, decisions, nil, testConfig())
	if result.Status != "pass" {
		t.Fatalf("expected pass (exempt path), got %s (blocked on: %+v)", result.Status, result.BlockedOn)
	}
}

func TestCheck_NonGovernedStatusIgnored(t *testing.T) {
	decisions := []vault.Decision{
		{Slug: "dec-1", Status: "superseded", XGoverns: []string{"src/billing/**"}},
	}
	result := Check([]string{"src/billing/webhook.go"}, decisions, nil, testConfig())
	if result.Status != "pass" {
		t.Fatalf("expected pass (superseded decision doesn't gate), got %s (blocked on: %+v)", result.Status, result.BlockedOn)
	}
}
