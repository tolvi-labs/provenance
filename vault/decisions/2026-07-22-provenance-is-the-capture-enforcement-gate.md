---
tags: [decision, provenance]
date: 2026-07-22
repo: provenance
status: active
ticket: none
user_impact: high
product_area: Product scope
supersedes: 2026-07-08-provenance-impact-blast-radius-hero
---

# Provenance is the capture-enforcement gate at the code→push boundary, not a standalone impact reporter

**Date:** 2026-07-22
**Repo:** provenance

## Why

[[2026-07-08-provenance-impact-blast-radius-hero]] committed Provenance to a vault-ranked impact/blast-radius report as its hero — a *query* over the durability record. Two things pulled that apart. First, the real product goal is adoption: the tooling has to be easy to install and hard to use wrong, and a report that "records the loop rather than improving it" (the admission carried since [[2026-07-04-provenance-durability-not-compliance]]) does not enforce anything. Second, the impact-report hero *depended on* a cross-language reachability engine, which is the single hardest research problem in the suite, so bundling it into Provenance meant Provenance could never ship until that hard part worked. The reframe resolves both: Provenance becomes the gate that *enforces capture* at the push boundary, and the impact reasoning moves to Canary, which needs it anyway. This is the mechanism by which a company actually stops losing institutional knowledge when an engineer leaves — capture stops being a good intention and becomes a precondition for shipping.

## How

- **Identity: a capture-enforcement gate, not a reporter.** Provenance runs at the code→push boundary and verifies that the required `tolvi sync` fragments for this change exist and are well-formed. The fragments the gate *requires* are the same fragments a correct capture is *made of* (what changed, why, on whose authority, which decisions it touches), so the gate and the record are one mechanism, not two.
- **The objective cross-check keeps declaration honest.** The gate never trusts the declaration alone. It cross-checks the declared fragments against the objective diff — which files/symbols actually changed, and which vault decisions govern those paths. If the diff touches a decision-governed path the sync never addressed, the gate flags it. This is the "change-time enforcement of binding decisions" role named in [[2026-07-08-provenance-impact-blast-radius-hero]], now with a concrete home, and it is what makes the gate un-gameable: you cannot under-declare your way past a path the diff objectively touched.
- **It gates on objective completeness, never on quality.** The gate blocks on *presence* and *consistency-with-the-diff*, never on whether the reasoning is *good*. That distinction is load-bearing: gating on quality is exactly the Goodhart trap that [[2026-07-08-chaos-receipt-wrong-use-case-churn-is-a-lens]] warns corrupts the capture source. Judging whether the "why" is insightful is left to the human PR reviewer, for whom the assembled report is surfaced. Hard-blocking is safe precisely because the bar is objective.
- **Hard block, no per-push skip.** When a required fragment is missing or the diff touches an unaddressed governed path, the push is blocked — there is no acknowledge-and-proceed escape hatch. The only way past is to not run Provenance at all, which is a repo-admin governance act set once at adoption, never a per-engineer, per-push choice.
- **The un-skippable enforcement point is CI, not the local hook.** A local pre-push git hook is bypassable (`git push --no-verify`), so the real gate is a CI *required status check* wired into branch protection, which an engineer cannot skip. The local hook is optional fast pre-flight so the failure is seen before CI. "Turn Provenance off" therefore means a repo admin removes the required check.
- **It still produces the report — assembled, not computed.** From the required fragments the gate assembles a provenance/impact report and (a) posts it to the PR for a human to eyeball risk and (b) hands it to Canary as the risk/why/bindings input to test-selection. The report is *assembled from what was declared during sync*, not *computed from static reachability*, so Provenance stays lean and carries no reachability engine.
- **Per-repo capture contract.** Setup defines, per repo, which fragments a correct capture must contain and which paths are decision-governed. The gate enforces that contract. This is the "define what should be committed for this repo" configuration.
- **Boundary with Canary.** Provenance reasons from the *declared* record and gates the *record*; Canary owns diff→tests via its coverage map and gates the *tests*. Provenance hands Canary declared-provenance + governance, not a reachability-derived impact report. See [[2026-07-22-canary-owns-impact-provenance-is-the-gate]].
- **Reconciliations.** The Article 50 retirement and the durability/attestation identity of [[2026-07-04-provenance-durability-not-compliance]] stand — capture-enforcement is that identity expressed as a gate rather than a report. But that decision's "Guild is the execution surface" ecosystem line and its "later probe, not next build" priority are stale: Guild was relaunched (2026-07-21) and shipped public (2026-07-22), so it is alive as the plan/brainstorm surface, and Provenance is an active reframe, not a deferred probe. The suite line is Guild (plan) → Bastion (harden plan) → code → Provenance (guard the record) → Canary (guard the tests) → prod: three gates, three artifacts, no overlap.

## Outcome

Provenance is committed as the capture-enforcement gate at the code→push boundary: it hard-blocks a push (in CI, via branch protection) when the required `tolvi sync` fragments are missing or inconsistent with the objective diff, gates only on objective completeness and never on reasoning quality, assembles the declared fragments into a report it posts to the PR and hands to Canary, and enforces a per-repo capture contract with no per-push skip. The impact/blast-radius report hero is retired and its reasoning moves to Canary; [[2026-07-08-provenance-impact-blast-radius-hero]] is superseded. `docs/PLAN.md` and `docs/impact-kickoff-prompt.md` still describe the retired impact-report hero and a dead-Guild premise and need a reconciliation pass.
