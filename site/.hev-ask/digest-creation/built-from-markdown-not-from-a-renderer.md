---
id: "digest-creation#built-from-markdown-not-from-a-renderer"
title: "Digest creation"
heading: "Built from markdown, not from a renderer"
group: "Overview"
order: 71
url: "/docs/digest-creation#built-from-markdown-not-from-a-renderer"
anchor: "built-from-markdown-not-from-a-renderer"
terms: ["built","markdown","renderer","build","never","imports","framework","reads","files","chunks","headings","derives","anchors","code","writes","tree","same","artifact","comes","whether","astro","docusaurus","vitepress","mkdocs","nothing","renders","pages","differs","host","only","runs","during","present","step","elsewhere","wiring","overlay","site","separate","smaller"]
hash: "473b58f15073f724b5961c7a62a2f9f09cabd49ec8b8e42e7d692e7ae6ec4461"
mode: "agent-primary"
facts: [{"kind":"code","literal":"astro build","chunkId":"digest-creation#built-from-markdown-not-from-a-renderer"}]
sources: [{"chunkId":"digest-creation#built-from-markdown-not-from-a-renderer","url":"/docs/digest-creation#built-from-markdown-not-from-a-renderer","anchor":"built-from-markdown-not-from-a-renderer"}]
---

The build never imports your framework: it reads files, chunks on headings, derives anchors in code, and writes the tree, so the same artifact comes out whether Astro, Docusaurus, VitePress, MkDocs, or nothing renders the pages. What differs per host is only when the build runs (during the Astro build when a key is present, or as a build/CI step elsewhere), and wiring the overlay into a non-Astro site is a separate smaller job.
