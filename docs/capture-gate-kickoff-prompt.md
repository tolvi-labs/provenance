# Capture-gate — build kickoff prompt

Kickoff prompt for Provenance's hero use case (see [[2026-07-22-provenance-is-the-capture-enforcement-gate]]). Deliberately **plan-first, not "go write it"** — hand it to an agent working in this repo, then harden its output through Bastion before any code.

---

```
You are helping design the first build slice of Tolvi Provenance, an open-source tool in the Tolvi Labs suite.

READ FIRST, then summarize what you found before designing:
- vault/decisions/2026-07-22-provenance-is-the-capture-enforcement-gate.md  (the hero you are building toward)
- vault/decisions/2026-07-08-provenance-impact-blast-radius-hero.md          (SUPERSEDED — the retired impact-report hero; do NOT resurrect it)
- vault/decisions/2026-07-04-provenance-durability-not-compliance.md         (Article 50 compliance is retired; do NOT resurrect it — the attestation identity survives)
- vault/decisions/2026-07-08-chaos-receipt-wrong-use-case-churn-is-a-lens.md (why a gate must never attach consequence to reasoning QUALITY — only to capture presence)
- docs/PLAN.md                                                               (go-forward plan)
- The sibling Canary repo (~/tolvi-labs/canary): Provenance hands Canary declared provenance + governance; Canary owns impact reasoning and test selection.
- The Tolvi and Tolvi Solo repos, to match stack, CLI conventions, and vault format.

CONTEXT: Tolvi is a per-repo "decision vault" — plain-Markdown records of a team's decisions, rejected alternatives, and incidents-that-became-rules, retrievable by an agent. A correct capture for a change records what changed, why, on whose authority, and which decisions it touches (attribution: human / AI / AI-and-edited). Provenance is the gate that makes that capture a precondition for pushing.

GOAL — the hero use case:
A capture-enforcement gate at the code→push boundary. Before a change can be pushed, verify that the required capture fragments exist AND are consistent with the objective diff, then hard-block the push if they are not. From the same fragments, assemble a provenance report: post it to the PR for a human to eyeball risk, and hand it to Canary as the risk/governance input to test selection. Primary users: engineering teams that must not lose institutional knowledge when someone leaves. Job-to-be-done: make vault capture something the pipeline requires, not something people intend.

NON-NEGOTIABLES (this is the whole point — don't drift):
1. Gate on objective completeness, NEVER on reasoning quality. The block is on fragment presence and consistency-with-the-diff. Whether the "why" is any good is the human reviewer's call, surfaced via the PR report. A gate that scores quality is the Goodhart trap that corrupts the capture source — see the chaos-receipt decision. Do not build it.
2. The objective cross-check is the teeth. Never trust the declaration alone: cross-check declared fragments against which files/symbols the diff actually changed and which vault decisions govern those paths. Flag a diff that touches a decision-governed path the capture never addressed. This is what makes the gate un-gameable.
3. Hard block, no per-push skip. When a required fragment is missing or a governed path is unaddressed, the push is blocked. There is no acknowledge-and-proceed escape. The only way past is a repo admin not running Provenance at all.
4. The un-skippable enforcement point is CI, not the local hook. A local pre-push hook is bypassable (`git push --no-verify`), so the real gate is a CI required status check wired into branch protection. The local hook is optional fast pre-flight only.
5. Assemble the report, do NOT compute reachability. The report is built from the declared fragments. Provenance carries no reachability engine — impact/reachability and test selection are Canary's, and Provenance hands Canary declared provenance + governance, not a reachability-derived report.
6. Per-repo capture contract. Setup defines which fragments a correct capture must contain and which paths are decision-governed; the gate enforces that contract.

DELIVERABLE — an implementation plan for the thinnest slice, not code yet:
- Architecture + data flow: how it reads a change (git diff over a branch/push range), loads the required capture fragments, runs the objective cross-check (changed paths → governing vault decisions), decides pass/block, and assembles the report.
- The capture-contract config: its on-disk shape, how setup defines required fragments and decision-governed paths, and how the gate resolves them per change.
- The report artifact: what the PR-posted report looks like, and the exact shape of the handoff object Canary consumes.
- CI integration: how a pipeline invokes the gate as a required check, and how the local pre-push hook mirrors it for fast feedback.
- The thinnest MVP that proves the capture-gate thesis first, dogfoodable on one real repo, then what layers on.
- Open design questions and the riskiest assumptions, called out explicitly.

Be honest about where this is hard: mapping changed paths to governing decisions precisely (both false blocks and missed governed paths are costly), keeping attribution accurate (human / AI / AI-and-edited) without trusting the author's word blindly, and defining "consistent with the diff" tightly enough to be un-gameable but loose enough to not block honest work. Do not hand-wave those.
```

---

**How to use:** run this in this repo, then run the returned plan back through Bastion ("harden this plan against the vault") before writing any code. Keep this prompt at design altitude; the follow-up build prompt is just "implement Phase 1 of the hardened plan."
