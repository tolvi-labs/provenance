---
tags: [decision, provenance]
date: 2026-09-08
repo: provenance
status: active
ticket: none
user_impact: high
product_area: Product scope
---

# Provenance ships as a Go CLI, using a new x-governs decision field and a new capture-fragment vault type

**Date:** 2026-09-08
**Repo:** provenance

## Why

The 2026-07-22 decision committed Provenance's identity (a capture-enforcement gate at the code→push boundary) but not its implementation form or the mechanics of how a capture actually attaches to a change — both were explicitly left as open design questions in `docs/capture-gate-kickoff-prompt.md`. Unlike Guild, Bastion, and Magellan, Provenance has to run with no LLM in the loop at all (a CI required check, an optional local pre-push hook), which rules out the Claude Code skill form factor entirely.

## How

- **Form: a standalone Go binary**, not a Claude Code skill. A single static binary is ideal for a CI gate (no interpreter or venv to install), it's the natural sibling to `tolvi`'s own Go CLI, and Go's stdlib plus two small dependencies (`gopkg.in/yaml.v3` for frontmatter, `github.com/bmatcuk/doublestar/v4` for `**`-glob matching) is enough — no external CLI framework, keeping the dependency surface minimal.
- **Governed paths: a new `x-governs` frontmatter field on decision files**, using the extension point `tolvi-format-v1`'s spec already reserves for exactly this (`x-*` fields always validate, no change to the shared `tolvi` repo or its schema required). Only decisions that opt in gate anything; most never do.
- **Capture linkage: a new vault content type, `vault/captures/YYYY-MM-DD-<slug>.md`**, structured (YAML frontmatter, not prose) because a hard gate needs machine-checkable fields. Two ways a change is captured, matching how engineers actually work: the diff adds/modifies a decision itself (new rationale, self-addressing), or a capture fragment names an existing decision the change operates under.
- **Attribution lives in the capture fragment**, never a commit trailer — Torres Atlantic's hard rule against AI attribution in commit messages made a trailer-based approach a non-starter regardless of technical merit.
- **The requirement is conditional, not universal.** A capture is only required when a diff touches a path some active decision governs. Most PRs touch no governed path and pass with zero friction — a hard gate that fires on every PR trains people to route around it, which is a materially worse failure mode for Provenance than for Magellan's advisory DAG gate, since Provenance actually blocks.
- **The local hook blocks**, unlike `tolvi precommit`'s always-exit-0 nudge (confirmed by reading `tolvi/cli/internal/cli/precommit.go` directly) — deliberately, since installing the hook is itself opt-in and its whole point is honest fast feedback before CI. Still bypassable via `git push --no-verify`; the real un-skippable gate is the CI required check.
- **CI integration is a documented pattern, not a published Action** — the binary stays platform-agnostic rather than assuming GitHub.
- **No separate Canary handoff format.** The same `tolvi-provenance-report-v1` JSON report posted to the PR is what Canary will consume — one contract, not two to keep in sync.
- **Two real gaps surfaced and fixed during the build, both independently verified.** `gate.Check`'s `NewDecisions`/`BlockedOn` fields were originally built from Go maps, so their element order was non-deterministic across runs — a real reproducibility defect for a JSON artifact meant for CI diffing/snapshotting, fixed by sorting both before return. The local pre-push hook's new-branch-push skip (remote SHA all zeros) had no mirror-image handling for a branch-deletion push (local SHA all zeros instead), which would have failed with a confusing internal git error instead of passing cleanly since a deletion has nothing to gate — fixed by widening the skip condition to cover both.

## Outcome

Provenance moved from pre-build to a working v1: `cmd/provenance` (check, init, hook install/uninstall), the `x-governs` convention, the `vault/captures/*.md` fragment type, the cross-check algorithm, JSON + human-readable reporting, and a documented CI integration pattern — all covered by real, executable tests (temp git repos, no mocks). Dogfooding on a real repo, and building Canary's actual consumption of this report, are the natural next steps.
