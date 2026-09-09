package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tolvi-labs/provenance/internal/gate"
)

// jsonReport mirrors the tolvi-provenance-report-v1 schema.
type jsonReport struct {
	Schema string `json:"schema"`
	Range  struct {
		Base string `json:"base"`
		Head string `json:"head"`
	} `json:"range"`
	Status        string             `json:"status"`
	GovernedFiles []jsonGovernedFile `json:"governed_files"`
	Captures      []jsonCapture      `json:"captures"`
	NewDecisions  []string           `json:"new_decisions"`
	BlockedOn     []jsonBlockedEntry `json:"blocked_on"`
}

type jsonGovernedFile struct {
	Path                string   `json:"path"`
	ImplicatedDecisions []string `json:"implicated_decisions"`
	AddressedBy         []string `json:"addressed_by"`
}

type jsonCapture struct {
	Path        string   `json:"path"`
	Decisions   []string `json:"decisions"`
	Attribution string   `json:"attribution"`
}

type jsonBlockedEntry struct {
	Decision string   `json:"decision"`
	Files    []string `json:"files"`
}

func emptyIfNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

// ToJSON renders a gate.Result as the tolvi-provenance-report-v1 JSON schema.
func ToJSON(result gate.Result, base, head string) ([]byte, error) {
	r := jsonReport{
		Schema:       "tolvi-provenance-report-v1",
		Status:       result.Status,
		NewDecisions: emptyIfNil(result.NewDecisions),
	}
	r.Range.Base = base
	r.Range.Head = head

	for _, gf := range result.GovernedFiles {
		r.GovernedFiles = append(r.GovernedFiles, jsonGovernedFile{
			Path:                gf.Path,
			ImplicatedDecisions: emptyIfNil(gf.ImplicatedDecisions),
			AddressedBy:         emptyIfNil(gf.AddressedBy),
		})
	}
	for _, c := range result.Captures {
		r.Captures = append(r.Captures, jsonCapture{
			Path:        c.Path,
			Decisions:   emptyIfNil(c.Decisions),
			Attribution: c.Attribution,
		})
	}
	for _, b := range result.BlockedOn {
		r.BlockedOn = append(r.BlockedOn, jsonBlockedEntry{
			Decision: b.Decision,
			Files:    emptyIfNil(b.Files),
		})
	}
	if r.GovernedFiles == nil {
		r.GovernedFiles = []jsonGovernedFile{}
	}
	if r.Captures == nil {
		r.Captures = []jsonCapture{}
	}
	if r.BlockedOn == nil {
		r.BlockedOn = []jsonBlockedEntry{}
	}

	return json.MarshalIndent(r, "", "  ")
}

// ToHuman renders a gate.Result as a human-readable report, suitable for
// both terminal output and a PR comment body.
func ToHuman(result gate.Result, base, head string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Provenance check: %s..%s\n", base, head)
	fmt.Fprintf(&b, "Status: %s\n\n", strings.ToUpper(result.Status))

	if len(result.GovernedFiles) == 0 {
		b.WriteString("No governed paths touched.\n")
	} else {
		b.WriteString("Governed files:\n")
		for _, gf := range result.GovernedFiles {
			fmt.Fprintf(&b, "  - %s\n", gf.Path)
			fmt.Fprintf(&b, "      implicated: %s\n", strings.Join(gf.ImplicatedDecisions, ", "))
			if len(gf.AddressedBy) > 0 {
				fmt.Fprintf(&b, "      addressed by: %s\n", strings.Join(gf.AddressedBy, ", "))
			}
		}
	}

	if len(result.Captures) > 0 {
		b.WriteString("\nCaptures in this diff:\n")
		for _, c := range result.Captures {
			fmt.Fprintf(&b, "  - %s (attribution: %s, decisions: %s)\n", c.Path, c.Attribution, strings.Join(c.Decisions, ", "))
		}
	}

	if len(result.BlockedOn) > 0 {
		b.WriteString("\nBLOCKED — the following decisions govern files this diff touches, but nothing in the diff addresses them:\n")
		for _, blocked := range result.BlockedOn {
			fmt.Fprintf(&b, "  - %s\n      files: %s\n", blocked.Decision, strings.Join(blocked.Files, ", "))
		}
		b.WriteString("\nAdd a vault/captures/*.md fragment referencing the decision, or a new decision if this changes the rationale.\n")
	}

	return b.String()
}
