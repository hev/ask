// The digest is a directory you read like any other: a real `ask tree` of this site.
export const askTreeDiagram = String.raw`$ ask tree

_glossary/                  (+10)
_meta                       Digest metadata
api/
  cli/                      CLI  (+9)
  configuration/            Configuration  (+5)
  digest/                   Digest format  (+7)
  endpoint/                 Search endpoint  (+10)
  mcp/                      MCP server  (+6)
  search-overlay/           SearchOverlay component  (+16)
concepts/                   Concepts
  chunks-and-anchors        Concepts > Chunks and anchors
  the-agentic-search-loop   Concepts > The agentic search loop
  the-ask-digest-directory  Concepts > The ask digest directory
  …                         (+7)
digest-creation/            Digest creation  (+5)
limits/                     Limits  (+8)
quickstart/                 Quick start  (+9)
tradeoffs/                  Tradeoffs  (+7)`;

// The digest is a directory: how it's built, what it is, and the three ways it's read.
export const askMapDiagram = String.raw`  ask digest build
      glob collections → chunk by headings → distil each section (Opus 4.8)
      one markdown file per section · hash-gated, incremental
      │
      ▼
  .hev-ask/                a committed, distilled mirror of your docs
  ├─ _meta.md              overview · context · suggestions · content hash
  ├─ _glossary/
  │  └─ digest.md          a term · its aliases · its definition
  ├─ overview/
  │  ├─ quick-start.md     one markdown + frontmatter file per section:
  │  └─ limits.md          title · summary · body · facts · url#anchor
  └─ api/
     └─ cli.md
      │
      ▼
  read three ways
  ├─ ⌘K overlay · humans   keyword search + a grounded answer · synthesis
  ├─ ask CLI    · agents   tree · ls · head · cat · facts · grep · keyless
  └─ ask mcp    · agents   one tool hydrates the tree; the agent reads it`;

// Progressive disclosure expressed as a directory: four rungs, each a larger slice of one file.
export const askLadderDiagram = String.raw`  progressive disclosure as a directory — each verb a larger slice of one section

  tree [--depth] ▸  titles only, the whole map          ·  cheap, safe to skim first
       │              (it's a directory — real ls / head work too)
       ▼
  cat <path>     ▸  the full distilled section body     ·  opt-in, one section at a time
       │
       ▼
  facts <path>   ▸  verbatim flags / code / identifiers ·  grounded literals + url#anchor
                    + sources + terms                     to cite back to the live page`;

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
