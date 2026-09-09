package cmdcheck

import (
	"flag"
	"fmt"
	"os"

	"github.com/tolvi-labs/provenance/internal/config"
	"github.com/tolvi-labs/provenance/internal/gate"
	"github.com/tolvi-labs/provenance/internal/gitutil"
	"github.com/tolvi-labs/provenance/internal/report"
	"github.com/tolvi-labs/provenance/internal/vault"
)

// Run implements `provenance check --base <ref> --head <ref> [--json-out <path>]`.
// It returns the process exit code: 0 for pass, 1 for blocked or a
// malformed capture, 2 for a usage or execution error.
func Run(args []string) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	base := fs.String("base", "", "base ref to diff from")
	head := fs.String("head", "HEAD", "head ref to diff to")
	jsonOut := fs.String("json-out", "", "path to write the machine-readable JSON report")
	repoDir := fs.String("repo", ".", "path to the git repository")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *base == "" {
		fmt.Fprintln(os.Stderr, "provenance check: --base is required")
		return 2
	}

	changedFiles, err := gitutil.ChangedFiles(*repoDir, *base, *head)
	if err != nil {
		fmt.Fprintf(os.Stderr, "provenance check: %v\n", err)
		return 2
	}

	vaultDir := *repoDir + "/vault"
	decisions, err := vault.LoadDecisions(vaultDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "provenance check: %v\n", err)
		return 2
	}
	captures, err := vault.LoadCaptures(vaultDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "provenance check: %v\n", err)
		return 2
	}

	// Only captures added in this diff range count. A capture merged in an
	// earlier diff is history, not an answer to the change under review.
	diffCaptures := capturesInDiff(captures, changedFiles)

	for _, c := range diffCaptures {
		if problems := vault.ValidateCapture(c, decisions); len(problems) > 0 {
			fmt.Fprintf(os.Stderr, "provenance check: malformed capture %s:\n", c.Path)
			for _, p := range problems {
				fmt.Fprintf(os.Stderr, "  - %s\n", p)
			}
			return 1
		}
	}

	cfg, err := config.Load(*repoDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "provenance check: %v\n", err)
		return 2
	}

	for _, d := range decisions {
		if err := config.ValidatePatterns(d.XGoverns); err != nil {
			fmt.Fprintf(os.Stderr, "provenance check: decision %s: x-governs: %v\n", d.Slug, err)
			return 2
		}
	}
	if err := config.ValidatePatterns(cfg.ExemptPaths); err != nil {
		fmt.Fprintf(os.Stderr, "provenance check: exempt_paths: %v\n", err)
		return 2
	}

	result := gate.Check(changedFiles, decisions, diffCaptures, cfg)

	fmt.Println(report.ToHuman(result, *base, *head))

	if *jsonOut != "" {
		raw, err := report.ToJSON(result, *base, *head)
		if err != nil {
			fmt.Fprintf(os.Stderr, "provenance check: %v\n", err)
			return 2
		}
		if err := os.WriteFile(*jsonOut, raw, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "provenance check: writing %s: %v\n", *jsonOut, err)
			return 2
		}
	}

	if result.Status == "blocked" {
		return 1
	}
	return 0
}

// capturesInDiff narrows the vault-wide capture set to only those fragments
// whose file is part of this diff range. Capture.Path is already stored
// repo-root-relative, in the same shape `git diff --name-only` emits.
func capturesInDiff(captures []vault.Capture, changedFiles []string) []vault.Capture {
	changed := make(map[string]bool, len(changedFiles))
	for _, f := range changedFiles {
		changed[f] = true
	}
	var inDiff []vault.Capture
	for _, c := range captures {
		if changed[c.Path] {
			inDiff = append(inDiff, c)
		}
	}
	return inDiff
}
