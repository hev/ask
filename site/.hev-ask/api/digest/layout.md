---
id: "api/digest#layout"
title: "Digest format"
heading: "Layout"
group: "API"
order: 22
url: "/docs/api/digest#layout"
anchor: "layout"
terms: ["layout","describes","disk","meta","file","glossary","directory","term","files","section","markdown","under","directories","mirror","paths","underscore","prefixed","entries","sort","first","never","collide","real","slug","there","committed","json","whole","artifact","overview","context","suggestions","version","contenthash","digest","aliases","definition","quick","start","mirroring"]
hash: "6ebf2bffe75c0ce6e0e4ebce97214673a69cdc69161cbfe9d7e1a584adb512bc"
mode: "source-primary"
facts: [{"kind":"code","literal":".hev-ask/\n  _meta.md                     overview · context · suggestions · version · contentHash\n  _glossary/\n    digest.md                  one file per term: aliases + definition\n  overview/\n    quick-start.md             one file per section, mirroring your doc paths\n    limits.md\n  api/\n    cli.md","chunkId":"api/digest#layout"},{"kind":"code","literal":"_meta","chunkId":"api/digest#layout"},{"kind":"code","literal":"_glossary","chunkId":"api/digest#layout"}]
sources: [{"chunkId":"api/digest#layout","url":"/docs/api/digest#layout","anchor":"layout"}]
---

Describes the on-disk layout: a meta file, a glossary directory of per-term files, and per-section markdown files under directories that mirror the doc paths. Underscore-prefixed non-section entries sort first and never collide with a real slug, and there is no committed JSON; the whole artifact is markdown.
