# Tolvi Provenance — Go-Forward Plan

**Status:** Pre-build. Reframed from a vault-ranked impact report to the capture-enforcement gate at the code→push boundary. Active. The impact-reasoning hero has moved to Canary; the Article 50 compliance framing remains retired.
**Last updated:** 2026-07-22

## One-line

The gate that makes vault capture a precondition for shipping: before a change can be pushed, Provenance verifies that the reasoning behind it — what changed, why, and on whose authority — was actually captured, and hard-blocks the push if it was not.

## The reframe (why the framing changed)

Provenance's prior hero (2026-07-08) was a vault-ranked impact / blast-radius report — a query over the durability record. Two problems retired it. It depended on a cross-language reachability engine, the hardest research problem in the suite, so Provenance could never ship until that worked. And a report "records the loop rather than improving it" — it enforces nothing, which is the exact gap the durability framing had already admitted. The reframe makes Provenance the gate that *enforces capture*, and moves the impact reasoning to Canary, which needs reachability anyway. The durability/attestation identity (who wrote this, why, on whose authority) survives — capture-enforcement is that identity expressed as a gate rather than a report. See [[2026-07-22-provenance-is-the-capture-enforcement-gate]].

Earlier framings, for the record: the original marketing-site v1 sold Provenance as EU AI Act Article 50 compliance for AI-generated code. That was retired on 2026-07-04 as factually unsupported (the Commission's draft Article 50 guidelines exempt source code at para. 64); see [[2026-07-04-provenance-durability-not-compliance]]. Compliance stays retired here.

## What it is

- **A capture-enforcement gate, not a reporter.** Provenance runs at the code→push boundary and verifies that the required capture fragments for a change exist and are well-formed.
- **The gate and the record are one mechanism.** The fragments the gate requires — what changed, why, on whose authority, which decisions it touches — are the same fragments a correct capture is made of. The gate does not add a second artifact; it enforces the one that should already exist.
- **The objective cross-check keeps declaration honest.** The gate never trusts the declaration alone. It cross-checks the declared fragments against the objective diff — which files and symbols actually changed, and which vault decisions govern those paths. A diff that touches a decision-governed path the capture never addressed is flagged. This is what makes the gate un-gameable: you cannot under-declare your way past a path the diff objectively touched.
- **It gates on objective completeness, never on quality.** The block is on presence and consistency-with-the-diff, never on whether the reasoning is good. Judging whether the "why" is insightful is left to the human reviewer, for whom the assembled report is surfaced on the PR. Hard-blocking is safe precisely because the bar is objective; the moment a gate scores reasoning quality it becomes the thing that corrupts the capture source.
- **It assembles the report from declared fragments.** From the required fragments the gate assembles a provenance report and (a) posts it to the PR for a human to eyeball risk and (b) hands it to Canary as the risk / why / bindings input to test selection. The report is assembled from what was declared, not computed from static reachability, so Provenance carries no reachability engine.
- **Attribution is a first-class fragment.** Who authored each change (human / AI / AI-and-edited) is part of what a correct capture records, keeping the name honest.

## Enforcement point

- **CI required check is the un-skippable gate.** A local pre-push git hook is bypassable (`git push --no-verify`), so the real gate is a CI required status check wired into branch protection, which an engineer cannot skip.
- **The local hook is optional fast pre-flight** — it lets the engineer fail before CI does, but it is not the enforcement point.
- **No per-push skip.** The only way past is a repo admin choosing not to run Provenance, which is a governance act set once at adoption, never a per-engineer, per-push choice.

## Per-repo capture contract

Setup defines, per repo, which fragments a correct capture must contain and which paths are decision-governed. The gate enforces that contract. This is the "define what should be committed for this repo" configuration.

## Boundary with the suite

Guild (plan) → Bastion (harden plan) → code → **Provenance (guard the record)** → Canary (guard the tests) → prod. Three gates, three artifacts, no overlap: Bastion judges the plan (pre-code), Provenance judges the record (code→push), Canary judges the tests (PR / merge / release). Provenance hands Canary declared provenance + governance, not a reachability-derived impact report. See [[2026-07-22-canary-owns-impact-provenance-is-the-gate]].

## Scope

- **v1:** the gate itself — required-fragment presence + the objective diff cross-check + hard-block in CI — plus the assembled PR report and the Canary handoff. Attribution is in scope as a first-class fragment.
- **Not v1:** any reachability-derived impact analysis (that is Canary's), and any attempt to judge reasoning quality (that is the human reviewer's, by design).

## Ecosystem role

Provenance is the record-gate beneath the daily loop: the vault captures decisions, Guild plans against them, Bastion hardens the plan, and Provenance enforces at push time that the change's reasoning actually landed in the vault before the code moves on. It turns capture from a discipline people intend into one the pipeline requires.

## Immediate cleanup

This plan and the build kickoff prompt were reconciled onto the capture-gate thesis on 2026-07-22; the retired impact-report hero and the stale "Guild is dead" premise are removed here (Guild was relaunched 2026-07-21 and shipped public 2026-07-22). The marketing-site copy for Provenance still sells the old Article 50 framing and needs the same reconciliation before the page goes public.

## Next step

Run `docs/capture-gate-kickoff-prompt.md` in this repo to produce a plan for the thinnest slice, then harden it through Bastion before any code.
