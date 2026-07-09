---
tags: [decision, provenance]
date: 2026-07-08
repo: provenance
status: active
ticket: none
user_impact: none
product_area: Product scope
---

# The engineer-defense "chaos receipt" is the wrong use case for Provenance; decision churn is a lens, not a product

**Date:** 2026-07-08
**Repo:** provenance

## Why

A common engineering story motivates this: an engineer is blamed for being slow while the person above them churns requirements (reversed decisions, mid-sprint scope, priority thrash, meeting load). The instinct was to build tooling that records that churn and protects the engineer, possibly as Provenance's first, sympathetic use case. We brainstormed it end to end and concluded the framing fails and is the wrong thing to anchor Provenance on. This doc banks the analysis so it is not re-litigated, and so future Provenance/Tolvi work does not reach for the same tempting-but-broken shape. The driver is product-strategy, not code.

## How

**The concept we explored, and how far we took it.** An engineer-owned, private-by-default "churn ledger" that emits an on-demand "reality receipt" the engineer could send upward when accused. Capture model was passive + opt-in annotation (read existing events — Tolvi decision supersessions, git history, PM board transitions — plus lightweight timestamped tags for verbal chaos that leaves no trace). Credibility mechanic was a two-tier receipt: Tier 1 hard facts re-derivable from the accuser's own systems, Tier 2 engineer annotations that are cryptographically signed and time-anchored so they cannot be forged or backdated. Home was to be Provenance's staged seed (append-only, time-anchored), with Bastion feeding reversal events in. The access model held up technically: signing binds identity, append-only hash-chaining prevents silent edits, time-anchoring prevents backdating, and local-first encrypted-at-rest storage means there is literally no central server a manager can query, so "only the owner can extract" is the absence of any other reader rather than a flippable permission.

**Why it fails — the framing, not the implementation.** The access model was sound; the *purpose* was not. Five attacks, one root:
- **Tools amplify intent, they do not install it.** The tool's value is gated entirely on the other party being reasonable. A visibility tool helps a manager who wants to see and lacks data; it does nothing for one who does not want to see, and pulling a cryptographic receipt on that manager reads as insubordination and accelerates the breakdown. The managers who most need the mirror are the least likely to accept it. The motivating story only resolved because the founder was already reasonable — survivorship.
- **Burden of proof and an adversarial, reactive frame.** A defense artifact puts the burden on the accused and is deployed mid-conflict, the worst moment to introduce evidence.
- **Selection-neutrality is unfixable in a defense artifact.** The engineer curates what to annotate (tags the chaos, not their own mistakes), and everyone knows it. Provenance can guarantee authenticity, integrity, and non-backdating; it cannot guarantee neutral selection.
- **Subject/beneficiary mismatch — the mirror lands on the wrong person.** In the actual story the churn is verbal (a call, a reversal spoken aloud); the *engineer* then implements it. So any tool-triggered signal fires on the engineer executing the churn, not the manager causing it. The culprit never touches the tooling, so never sees the mirror; the victim, who already knows, is the one shown it.
- **Goodhart poisons the brain.** The moment "reversals within N days" becomes a visible or consequential metric, people route around it: stop recording decisions so there is nothing to reverse, wait out the window, or relabel a reversal as a fresh decision. The gaming suppresses honest decision capture, which is Tolvi's entire value. (This is the original "do not add noise to the brain" worry arriving from a new direction: the danger is not storing churn data, it is attaching any social consequence to it.)

**The two problems hiding under one story.** We kept conflating them:
1. **In-workflow decision churn** — re-litigating settled decisions inside the engineering process (an agent proposing a reversal, a lead overturning their team's call in a plan/PR, a reopened ticket). Tool-visible, agent-amplified, dead-on Tolvi's "stop re-litigating settled decisions" thesis.
2. **Top-down verbal chaos** — a non-technical authority creating churn through meetings and calls. Tool-invisible; the culprit never touches the tooling. This is the emotionally resonant problem, and it is an org-power problem, not a tooling problem. No honest tool fixes a manager who will not look.

The tractable problem and the resonant problem are not the same problem. Letting the story drive the design means building for #2, which does not yield to software.

**What survives (deferred, not built).** Only a small, private, non-consequential feature for #1: a calibrated, rare, in-the-moment reversal-cost nudge shown to whoever makes the change *through the tooling* — for example, at plan-time via Bastion, "this reverses DEC-184, made 5 days ago, ~N of work already built on it; confirm?" It is most valuable against agent-amplified re-litigation. It is closer to a UI affordance on Tolvi's existing retrieval path than a new product. Constraints if it is ever built: it must never become a metric anyone is measured by; it must be calibrated hard against alert-fatigue (fires only on recent + high-sunk-cost + non-author reversals) or it trains people to ignore Tolvi; and it stays decision-level, never person-level.

**Implication for Provenance.** Do not anchor Provenance's identity on the chaos-receipt / engineer-defense use case. Provenance's honest job is unchanged from [[2026-07-04-provenance-durability-not-compliance]]: a vault-native decision/authorship attestation and durability layer ("who wrote this, why, and on whose authority"). Decision churn is at most one derived lens over that record, never the reason it exists, and never a surfaced KPI.

## Outcome

The engineer-defense "chaos receipt" is shelved with nothing built; Provenance stays a decision-attestation/durability layer; and a hard guardrail is recorded for future Tolvi/Provenance work — any churn signal must remain private and non-consequential (never a metric, receipt, or review input) or it corrupts the decision-capture source it reads from.
