---
id: "tradeoffs#keyword-retrieval-not-embeddings"
title: "Tradeoffs"
heading: "Keyword retrieval, not embeddings"
group: "Overview"
order: 104
url: "/docs/tradeoffs#keyword-retrieval-not-embeddings"
anchor: "keyword-retrieval-not-embeddings"
terms: ["keyword","retrieval","embeddings","dependency","free","token","overlap","widened","glossary","vector","store","upside","nothing","host","keep","sync","edge","safe","instant","recovering","much","synonym","recall","cost","paraphrase","ceiling","since","agent","only","ground","found","readers","routinely","search","words","sharing","tokens","docs","absent","would"]
hash: "816629bc3f9d42150a53cdbf54324d130714f295a4f5db66f417ac4b8965fa94"
mode: "agent-primary"
facts: []
sources: [{"chunkId":"tradeoffs#keyword-retrieval-not-embeddings","url":"/docs/tradeoffs#keyword-retrieval-not-embeddings","anchor":"keyword-retrieval-not-embeddings"}]
---

Retrieval is dependency-free token overlap widened by the glossary, with no embeddings or vector store. The upside is nothing to host or keep in sync, edge-safe and instant, with the glossary recovering much synonym recall; the cost is a paraphrase-recall ceiling, since the agent can only ground in what keyword retrieval found, so readers who routinely search in words sharing no tokens with the docs and absent from the glossary would be better served by embeddings — an upgrade deferred rather than designed out.
