---
tags: [decision, provenance]
date: 2026-07-04
repo: provenance
status: active
ticket: none
user_impact: none
product_area: Positioning
---

# Provenance is a durability play, not an EU AI Act Article 50 compliance play

**Date:** 2026-07-04
**Repo:** provenance

## Why

Provenance was framed on the marketing site as the way to satisfy EU AI Act Article 50 for AI-generated code, with the whole v1 urgency resting on that regulatory hook. Researching the actual regulation shows the hook does not attach: Article 50 governs user-facing synthetic content, and the Commission's own draft guidelines explicitly exempt source code. Building — and worse, selling — a compliance product on a premise the regulator contradicts is both a wasted bet and a real liability with the regulated buyers it targets. The durable value of Provenance is elsewhere.

## How

- **The legal finding.** Article 50 is a transparency obligation on providers/deployers of generative AI covering synthetic audio/image/video/text (deepfakes, AI-media, chatbot disclosure), applicable 2 Aug 2026. The Commission's **draft Article 50 guidelines (8 May 2026) exempt source code at para. 64** — source code is treated as a "technical output rather than human-facing content," and the exemption covers source code proper including inline comments and docstrings. Standalone natural-language outputs (README files, product descriptions, explanatory text) re-enter Article 50(2), but that is an edge, not the product. AI code-assistance tools for professional developers may also clear the "obviousness" threshold, and there is a B2B technical-output exception. Caveat: draft guidance, could shift; not legal advice — but specific and repeated enough that a deadline-driven compliance build is not defensible.
  - Sources: artificialintelligenceact.eu/article/50/ ; Bird & Bird "Reading the Commission's Draft Article 50 Guidelines" (para. 64 source-code exemption) ; Covington/Global Policy Watch "10 Takeaways" (obviousness + B2B exemptions) ; digital-strategy.ec.europa.eu draft guidelines.
- **The reframe.** Provenance's real, durable value is the story the site itself defers to 2027: drift detection, mental-model preservation, refactoring confidence — "who wrote this, why, and on whose authority" for the team's future self. That is a correctness/context play, on-thesis, and it compounds with the vault as the audit-grade record beneath the stack. Positioning leads with durability; any regulatory angle is confined to the honest in-scope sliver and never sold as Article 50 compliance for code.
- **Priority.** Real but peripheral to the daily build loop — it records the loop rather than improving it. With arbitrary dates and no demand pull, it sits behind Forge (ship) and Guild (correctness probe), and when taken up follows the same evidence-first path (thin capture+query slice dogfooded on a real product).
- **Cleanup, independent of build order.** The live copy selling Article 50 compliance and showing those requirements "shipped" should be pulled or reframed onto durability before the page goes public.
- **Ecosystem thesis.** Tolvi captures decisions, Guild is the execution surface, Provenance records authorship + reasoning per change and cross-references the vault decisions that shaped it. See [[2026-07-04-guild-correctness-recall-reframe]] and [[2026-07-04-forge-defer-ship-until-usage-data]].

## Outcome

Provenance's committed positioning is a vault-native durability layer; the Article 50 compliance framing is retired as factually unsupported, and the product is a later evidence-first probe rather than a deadline build.
