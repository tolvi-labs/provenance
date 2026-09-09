package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/tolvi-labs/provenance/internal/gate"
)

func TestToJSON_Shape(t *testing.T) {
	result := gate.Result{
		Status: "blocked",
		GovernedFiles: []gate.GovernedFile{
			{Path: "src/billing/webhook.go", ImplicatedDecisions: []string{"dec-1"}},
		},
		BlockedOn: []gate.BlockedEntry{
			{Decision: "dec-1", Files: []string{"src/billing/webhook.go"}},
		},
	}
	raw, err := ToJSON(result, "abc", "def")
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if parsed["schema"] != "tolvi-provenance-report-v1" {
		t.Fatalf("unexpected schema field: %v", parsed["schema"])
	}
	if parsed["status"] != "blocked" {
		t.Fatalf("unexpected status field: %v", parsed["status"])
	}
	rng, ok := parsed["range"].(map[string]interface{})
	if !ok || rng["base"] != "abc" || rng["head"] != "def" {
		t.Fatalf("unexpected range field: %v", parsed["range"])
	}
}

func TestToHuman_MentionsBlockedDecision(t *testing.T) {
	result := gate.Result{
		Status: "blocked",
		BlockedOn: []gate.BlockedEntry{
			{Decision: "dec-1", Files: []string{"src/billing/webhook.go"}},
		},
	}
	out := ToHuman(result, "abc", "def")
	if !strings.Contains(out, "BLOCKED") {
		t.Fatalf("expected BLOCKED in output, got: %s", out)
	}
	if !strings.Contains(out, "dec-1") {
		t.Fatalf("expected decision slug in output, got: %s", out)
	}
	if !strings.Contains(out, "src/billing/webhook.go") {
		t.Fatalf("expected file path in output, got: %s", out)
	}
}

func TestToHuman_PassingHasNoBlockedSection(t *testing.T) {
	result := gate.Result{Status: "pass"}
	out := ToHuman(result, "abc", "def")
	if strings.Contains(out, "BLOCKED") {
		t.Fatalf("did not expect BLOCKED in a passing report, got: %s", out)
	}
}
