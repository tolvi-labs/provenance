# Provenance

**The capture-enforcement gate at the code→push boundary.** Provenance verifies that the reasoning behind a change — what changed, why, and on whose authority — was actually captured in the vault before that change can be pushed. It is how a team stops losing institutional knowledge when an engineer leaves: capture stops being a good intention and becomes a precondition for shipping.

> **Status: pre-build.** The direction is committed in [`vault/decisions/`](vault/decisions/) and [`docs/PLAN.md`](docs/PLAN.md). No code yet.

It is the record-gate of the Tolvi stack:

```
Guild (plan) → Bastion (harden plan) → code → Provenance (guard the record) → Canary (guard the tests) → prod
```

## What it will do

At the code→push boundary, Provenance checks that the required capture fragments for a change exist and are consistent with the actual diff — cross-checking the declared record against which files and symbols changed, and which vault decisions govern those paths. If the record is missing, or the diff touches a decision-governed path the capture never addressed, it **hard-blocks the push**.

It gates on objective completeness, never on the quality of the reasoning. Whether the "why" is any good is a question for a human reviewer, to whom Provenance surfaces the assembled record on the PR. From the same fragments it assembles a provenance report, which it posts to the PR for a human to eyeball risk and hands to [Canary](https://github.com/tolvi-labs/canary) as the risk-and-governance input to test selection.

## Design principles

- **Objective, not subjective.** The gate enforces that capture happened and that it matches the diff. It never scores whether the reasoning is insightful — that is the trap that corrupts the record it exists to protect.
- **No per-push skip.** The only way past is a repo admin choosing not to run Provenance at all. The un-skippable enforcement point is a CI required check wired into branch protection; a local pre-push hook is optional fast pre-flight, and because it is bypassable it is not the gate.
- **Per-repo capture contract.** Setup defines, per repo, which fragments a correct capture must contain and which paths are decision-governed.

The full rationale is in [`docs/PLAN.md`](docs/PLAN.md) and the decisions behind it live in [`vault/decisions/`](vault/decisions/).

## License

Apache 2.0, in line with the rest of the Tolvi suite.
