---
id: "api/digest#layout"
title: "Digest format"
heading: "Layout"
group: "API"
order: 22
url: "/docs/api/digest#layout"
anchor: "layout"
terms: ["layout","describes","disk","meta","file","overview","context","suggestions","version","content","hash","glossary","directory","term","section","files","mirroring","paths","heading","level","chunk","underscore","prefixed","entries","sort","first","never","collide","real","slug","whole","artifact","markdown","committed","json","contenthash","digest","aliases","definition","quick"]
hash: "6ebf2bffe75c0ce6e0e4ebce97214673a69cdc69161cbfe9d7e1a584adb512bc"
mode: "source-primary"
facts: [{"kind":"code","literal":".hev-ask/\n  _meta.md                     overview · context · suggestions · version · contentHash\n  _glossary/\n    digest.md                  one file per term: aliases + definition\n  overview/\n    quick-start.md             one file per section, mirroring your doc paths\n    limits.md\n  api/\n    cli.md","chunkId":"api/digest#layout"},{"kind":"code","literal":"_meta","chunkId":"api/digest#layout"},{"kind":"code","literal":"_glossary","chunkId":"api/digest#layout"}]
sources: [{"chunkId":"api/digest#layout","url":"/docs/api/digest#layout","anchor":"layout"}]
---

Describes the on-disk layout: a meta file with overview, context, suggestions, version, and content hash; a glossary directory with one file per term; and per-section files mirroring your doc paths, one per heading-level chunk. Underscore-prefixed non-section entries sort first and never collide with a real doc slug, and the whole artifact is markdown with no committed JSON.
