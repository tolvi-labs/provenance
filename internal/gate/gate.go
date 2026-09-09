package gate

import (
	"github.com/tolvi-labs/provenance/internal/config"
	"github.com/tolvi-labs/provenance/internal/vault"
)

// GovernedFile is one changed file that matches at least one governing
// decision's x-governs pattern.
type GovernedFile struct {
	Path                string
	ImplicatedDecisions []string
	AddressedBy         []string
}

// BlockedEntry names one decision that implicated a file but was never
// addressed by anything in this diff.
type BlockedEntry struct {
	Decision string
	Files    []string
}

// Result is the outcome of one cross-check run.
type Result struct {
	Status        string // "pass" or "blocked"
	GovernedFiles []GovernedFile
	Captures      []vault.Capture
	NewDecisions  []string
	BlockedOn     []BlockedEntry
}

// Check runs the cross-check: which changed files are governed, whether
// they're addressed by a capture or a new decision in this diff, and what
// remains unaddressed.
func Check(changedFiles []string, decisions []vault.Decision, captures []vault.Capture, cfg config.Config) Result {
	governing := vault.FilterByStatus(decisions, cfg.GovernedStatuses)

	// A decision added/modified in this diff is self-addressing.
	newDecisionSlugs := make(map[string]bool)
	for _, d := range governing {
		for _, f := range changedFiles {
			if f == "vault/decisions/"+d.Slug+".md" {
				newDecisionSlugs[d.Slug] = true
			}
		}
	}

	// A decision named in any capture fragment in this diff is addressed.
	addressedByCapture := make(map[string][]string) // decision slug -> capture paths
	for _, c := range captures {
		for _, slug := range c.Decisions {
			addressedByCapture[slug] = append(addressedByCapture[slug], c.Path)
		}
	}

	var governedFiles []GovernedFile
	blockedByDecision := make(map[string][]string) // decision slug -> files that triggered it

	for _, f := range changedFiles {
		if config.MatchesAny(f, cfg.ExemptPaths) {
			continue
		}
		var implicated []string
		var addressed []string
		for _, d := range governing {
			if !config.MatchesAny(f, d.XGoverns) {
				continue
			}
			implicated = append(implicated, d.Slug)
			if newDecisionSlugs[d.Slug] || len(addressedByCapture[d.Slug]) > 0 {
				addressed = append(addressed, d.Slug)
			} else {
				blockedByDecision[d.Slug] = append(blockedByDecision[d.Slug], f)
			}
		}
		if len(implicated) > 0 {
			governedFiles = append(governedFiles, GovernedFile{
				Path:                f,
				ImplicatedDecisions: implicated,
				AddressedBy:         addressed,
			})
		}
	}

	var blockedOn []BlockedEntry
	for slug, files := range blockedByDecision {
		blockedOn = append(blockedOn, BlockedEntry{Decision: slug, Files: files})
	}

	var newDecisions []string
	for slug := range newDecisionSlugs {
		newDecisions = append(newDecisions, slug)
	}

	status := "pass"
	if len(blockedOn) > 0 {
		status = "blocked"
	}

	return Result{
		Status:        status,
		GovernedFiles: governedFiles,
		Captures:      captures,
		NewDecisions:  newDecisions,
		BlockedOn:     blockedOn,
	}
}
