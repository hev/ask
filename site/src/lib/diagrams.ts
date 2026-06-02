// Build-time vs runtime: where each part of hev ask does its work.
export const askMapDiagram = String.raw`       BUILD TIME (CLI / Skill)                RUNTIME (edge)
  ╔═════════════════════════════════╗    ╔════════════════════════════════════╗░
  ║  ask kg build                   ║    ║  /api/ask   (prerender: false)     ║░
  ║                                 ║    ║                                    ║░
  ║  glob src/content/docs/**       ║    ║  ┌─── keyword mode · no key ────┐  ║░
  ║   → chunk by heading            ║    ║  │ prefilter chunks + glossary  │  ║░
  ║   → sha256 content hash         ║    ║  └──────────────────────────────┘  ║░
  ║   → Opus 4.8 builds the graph   ║    ║  ┌──── agentic loop · Haiku ────┐  ║░
  ║   → write .hev-ask/kg.json      ║    ║  │ system: kg.context (cached)  │  ║░
  ║                                 ║    ║  │ tool: search(q)  ≤ 4 times   │  ║░
  ╚═════════════════════════════════╝    ║  │ then stream answer, no tools │  ║░
   ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░    ║  │ grounded in /page#anchor     │  ║░
                                         ║  └──────────────────────────────┘  ║░
     .hev-ask/kg.json  (committed)       ╚════════════════════════════════════╝░
           │                              ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
           ▼
     virtual:hev-ask/kg  ──── bundled into /api/ask ────▶`;

// What the reader experiences at the overlay.
export const askFlowDiagram = String.raw`           type a query                    press Enter
                 │                               │
                 ▼                               ▼
      ┌───────────────────┐            ┌─────────────────┐
      │ keyword           │            │ agentic loop    │
      │ instant · keyless │            │ search → answer │
      └───────────────────┘            │ needs API key   │
                 │                     └─────────────────┘
                 ▼                               │
      ┌───────────────────┐                      ▼
      │ section results   │        ┌──────────────────────────┐
      │ /docs/page#anchor │        │ streamed answer with     │
      └───────────────────┘        │ inline /docs/page#anchor │
                                   └──────────────────────────┘`;
