---
id: "concepts#host-neutral-one-digest-any-framework"
title: "Concepts"
heading: "Host-neutral: one digest, any framework"
group: "Overview"
order: 63
url: "/docs/concepts#host-neutral-one-digest-any-framework"
anchor: "host-neutral-one-digest-any-framework"
terms: ["host","neutral","digest","framework","everything","builds","happens","before","renderer","touches","docs","reading","markdown","filesystem","chunking","headings","deriving","anchors","code","writing","tree","without","importing","astro","same","artifact","comes","whatever","renders","pages","differs","only","thin","adapter","glue","batteries","included","while","every","other"]
hash: "b670f196f7770314128233d4a8ce13b0b168f901639a3697bcaf9c119bc1a17a"
mode: "agent-primary"
facts: [{"kind":"code","literal":".hev-ask/","chunkId":"concepts#host-neutral-one-digest-any-framework"},{"kind":"code","literal":"hevAsk()","chunkId":"concepts#host-neutral-one-digest-any-framework"},{"kind":"code","literal":"astro build","chunkId":"concepts#host-neutral-one-digest-any-framework"},{"kind":"code","literal":"/api/ask","chunkId":"concepts#host-neutral-one-digest-any-framework"},{"kind":"code","literal":"SearchOverlay.astro","chunkId":"concepts#host-neutral-one-digest-any-framework"},{"kind":"code","literal":"\u003cscript\u003e","chunkId":"concepts#host-neutral-one-digest-any-framework"},{"kind":"code","literal":"tree","chunkId":"concepts#host-neutral-one-digest-any-framework"},{"kind":"code","literal":"cat","chunkId":"concepts#host-neutral-one-digest-any-framework"},{"kind":"code","literal":"grep","chunkId":"concepts#host-neutral-one-digest-any-framework"}]
sources: [{"chunkId":"concepts#host-neutral-one-digest-any-framework","url":"/docs/concepts#host-neutral-one-digest-any-framework","anchor":"host-neutral-one-digest-any-framework"}]
---

Everything that builds the digest happens before a renderer touches your docs, reading markdown off the filesystem, chunking on headings, deriving anchors in code, and writing the tree without importing Astro, so the same artifact comes out whatever framework renders the pages. What differs per framework is only the thin adapter glue; Astro's is batteries-included while every other framework uses two host-neutral primitives (the fully static keyword overlay and the standalone hostable agentic endpoint), and the CLI and MCP surfaces are already host-neutral because they read the committed tree directly.
