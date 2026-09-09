# Provenance

**The capture-enforcement gate at the code→push boundary.** Provenance verifies that the reasoning behind a change — what changed, why, and on whose authority — was actually captured in the vault before that change can be pushed. It is how a team stops losing institutional knowledge when an engineer leaves: capture stops being a good intention and becomes a precondition for shipping.

It is the record-gate of the Tolvi stack:

```
Guild (plan) → Bastion (harden plan) → code → Provenance (guard the record) → Canary (guard the tests) → prod
```

## Install

```shell
git clone https://github.com/tolvi-labs/provenance
cd provenance
go build -o provenance ./cmd/provenance
```

Then, in the repo you want to gate:

```shell
provenance init                 # scaffolds provenance.yml
provenance hook install         # optional: fast local pre-flight before CI
```

See [`docs/ci-integration.md`](docs/ci-integration.md) for wiring the CI required check — that's the actual un-skippable gate.

## What it does

Most changes touch nothing Provenance cares about and pass with zero friction. When a diff touches a path an *active* vault decision declares it governs (via `x-governs` on the decision's frontmatter), the change needs a capture: either the diff adds/modifies that decision itself (new rationale), or it adds a `vault/captures/*.md` fragment naming the decision it operates under. Missing either — the push is hard-blocked in CI.

The gate never judges whether the reasoning is *good* — only whether it's *present and consistent with the diff*. Whether the "why" is insightful is the human PR reviewer's call, surfaced via the report Provenance posts to the PR.

## Core loop

```
provenance check --base <ref> --head <ref>
  ├─ 1. Diff        → which files changed in this range
  ├─ 2. Cross-check  → which changed files are governed, by which active decisions
  ├─ 3. Match        → addressed by a new decision or a capture fragment in this diff?
  └─ 4. Report       → pass, or block naming exactly which decision and files are unaddressed
```

The same JSON report Provenance posts to the PR is also the handoff [Canary](https://github.com/tolvi-labs/canary) consumes as its risk/why/bindings input for test selection — no second contract to keep in sync.

## Design principles

- **Conditional, not universal.** A capture is only required on a governed path. A gate that fires on every PR trains people to route around it.
- **Gate on objective completeness, never on quality.** The block is on presence and consistency-with-the-diff. Whether the reasoning is good is the human reviewer's call, always.
- **The un-skippable enforcement point is CI, not the local hook.** `git push --no-verify` bypasses the local hook; a CI required check does not.
- **No reachability engine.** Impact/blast-radius reasoning is [Canary](https://github.com/tolvi-labs/canary)'s job, never Provenance's.

The full rationale is in [`docs/PLAN.md`](docs/PLAN.md) and the decisions behind it live in [`vault/decisions/`](vault/decisions/).

## License

Apache 2.0, in line with the rest of the Tolvi suite.
