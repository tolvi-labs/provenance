---
tags: [decision, provenance]
date: 2026-07-08
repo: provenance
status: superseded
superseded_by: 2026-07-22-provenance-is-the-capture-enforcement-gate
ticket: none
user_impact: none
product_area: Product scope
---

# Provenance's hero use case is a vault-ranked impact/blast-radius report for eng+QA; attestation is the durable substrate under it

**Date:** 2026-07-08
**Repo:** provenance

> **Superseded 2026-07-22** by [[2026-07-22-provenance-is-the-capture-enforcement-gate]]. The impact/blast-radius report is retired as Provenance's hero and its reasoning moves to Canary; Provenance becomes the capture-enforcement gate at the code→push boundary. The attestation substrate and the "change-time enforcement of binding decisions" role below survive and carry into the superseding decision. The "Guild is killed" premise in this doc is itself stale — Guild was relaunched 2026-07-21 and shipped public 2026-07-22.

## Why

[[2026-07-04-provenance-durability-not-compliance]] retired the Article 50 framing and reframed Provenance as a durability/attestation layer ("who wrote this, why, and on whose authority") — but left it a later, peripheral probe that "records the loop rather than improving it." A parallel Tolvi Labs positioning pass (personal brand vault, 2026-07-08) gave Provenance a hero use case that *does* improve the daily loop, which is exactly what the durability framing admitted it lacked. This banks that sharpening and the ecosystem change that rode in with it (Guild killed).

## How

- **The hero use case: a vault-ranked impact / blast-radius report.** Point Provenance at a change set (a feature branch, or a release-candidate tag) and it emits a report of exactly what the change impacts, ranked by risk, so a team tests only what needs testing instead of re-running everything. Users: engineers, QA, release managers. JTBD: cut test time and surface risk before ship. The durability record is the substrate; the impact report is the first valuable query over it, and it's the loop-improving job the durability framing was missing.
- **The moat is the vault, not static analysis.** Reachability/dependency diffing alone is a commodity. The differentiator is cross-referencing changed code against the vault so the report reads "this touches the billing path DEC-184 marked audit-scoped and the endpoint from the stale-cache incident — test these first," not "X imports Y." Static analysis finds *what* changed; the vault explains *why* it's risky and ranks it.
- **Attribution is the bridge, kept honest.** The who/why/authority attestation from [[2026-07-04-provenance-durability-not-compliance]] rides along as a first-class part of the report. It keeps the name "Provenance" honest and preserves the durable/record value (and the honest in-scope audit sliver) as a *secondary* benefit, never the headline. Compliance stays retired.
- **Two run modes:** (a) against an RC tag → a scoped regression test plan; (b) against a working branch → "what did I touch, what should I test."
- **Boundary with Bastion.** Bastion judges a *plan* (pre-code intent); Provenance judges a *diff* (post-code artifact). Provenance never evaluates plans and never writes code — it reports on changes that already exist.
- **Ecosystem update — Guild is killed (brand vault, 2026-07-08).** The 07-04 ecosystem thesis named Guild the execution surface. Guild was killed as a standing directive/agent layer that adds maintenance work against Tolvi's capture-as-side-effect thesis. Its enforcement function redistributes: forward-intent and invariants are captured as ordinary vault decisions, enforced at plan-time by Bastion and at **change-time by Provenance** (flagging a diff that touches a path governed by a binding decision it doesn't satisfy). This strengthens Provenance's role rather than weakening it. NOTE: `docs/PLAN.md` and the 07-04 decisions still reference Guild as live and should be reconciled in this repo when convenient — left as-is here to avoid rewriting history mid-stream.
- **Priority.** Elevated from "later probe, not next build" to an active, evidence-first probe, now that it has a loop-improving hero. First build is a thin capture+query slice that proves the vault-ranked-risk thesis, dogfooded on a real product. Kickoff prompt lives at `docs/impact-kickoff-prompt.md`.
- **Consistency with [[2026-07-08-chaos-receipt-wrong-use-case-churn-is-a-lens]].** The churn/engineer-defense use case stays shelved; impact/test-scoping is a sounder hero than churn-as-receipt, and any churn signal remains a private, non-consequential lens only.

## Outcome

Provenance's committed hero is a vault-ranked impact/blast-radius report that scopes testing and surfaces risk, sitting on the who/why/authority attestation substrate; it is elevated to an active evidence-first probe; Guild's change-time enforcement folds into it; and the build kickoff is captured at `docs/impact-kickoff-prompt.md`.
