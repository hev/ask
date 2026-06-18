---
id: "concepts#keyword-search-and-the-glossary"
title: "Concepts"
heading: "Keyword search and the glossary"
group: "Overview"
order: 64
url: "/docs/concepts#keyword-search-and-the-glossary"
anchor: "keyword-search-and-the-glossary"
terms: ["keyword","search","glossary","instant","path","shared","remote","runs","dependency","free","prefilter","chunks","expands","query","term","aliases","matched","tokens","additively","capped","scores","token","overlap","widened","digest","sections","considers","central","rank","above","incidental","body","mentions","caps","results","document","long","page","cannot","monopolize"]
hash: "6496e6a2896ae22bcdf61b3dbba97df885d89427d2f4b266e7553f09e094cb91"
mode: "agent-primary"
facts: [{"kind":"code","literal":"grep","chunkId":"concepts#keyword-search-and-the-glossary"},{"kind":"code","literal":"k8s","chunkId":"concepts#keyword-search-and-the-glossary"},{"kind":"code","literal":"kubernetes","chunkId":"concepts#keyword-search-and-the-glossary"},{"kind":"code","literal":"summary","chunkId":"concepts#keyword-search-and-the-glossary"},{"kind":"code","literal":"terms","chunkId":"concepts#keyword-search-and-the-glossary"},{"kind":"code","literal":"facts","chunkId":"concepts#keyword-search-and-the-glossary"}]
sources: [{"chunkId":"concepts#keyword-search-and-the-glossary","url":"/docs/concepts#keyword-search-and-the-glossary","anchor":"keyword-search-and-the-glossary"}]
---

The instant keyword path, shared with the CLI's remote search, runs a dependency-free prefilter over the chunks: it expands each query term with its glossary aliases and matched-term tokens (additively and capped), scores by token overlap widened by the digest so sections the digest considers central rank above incidental body mentions, caps results per document so one long page cannot monopolize, then excerpts around the first match. It needs no key and no embeddings, and degrades to plain token overlap over raw chunk text when no tree is present.
