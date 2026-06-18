---
id: "tradeoffs#keyword-retrieval-not-embeddings"
title: "Tradeoffs"
heading: "Keyword retrieval, not embeddings"
group: "Overview"
order: 104
url: "/docs/tradeoffs#keyword-retrieval-not-embeddings"
anchor: "keyword-retrieval-not-embeddings"
terms: ["keyword","retrieval","embeddings","dependency","free","token","overlap","widened","glossary","vector","store","upside","nothing","host","keep","sync","edge","safe","instant","recovering","much","synonym","recall","would","give","cost","ceiling","paraphrase","since","agent","only","ground","found","better","readers","routinely","search","words","sharing","tokens"]
hash: "816629bc3f9d42150a53cdbf54324d130714f295a4f5db66f417ac4b8965fa94"
mode: "agent-primary"
facts: []
sources: [{"chunkId":"tradeoffs#keyword-retrieval-not-embeddings","url":"/docs/tradeoffs#keyword-retrieval-not-embeddings","anchor":"keyword-retrieval-not-embeddings"}]
---

Retrieval is dependency-free token overlap widened by the glossary with no embeddings or vector store, so the upside is nothing to host or keep in sync, edge-safe and instant, with the glossary recovering much of the synonym recall embeddings would give; the cost is a ceiling on paraphrase recall, since the agent can only ground in what keyword retrieval found, and embeddings would do better for readers who routinely search in words sharing no tokens with your docs and absent from the glossary. That upgrade is deferred, not designed out.
