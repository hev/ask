---
id: "concepts#the-system-prompt-is-cached"
title: "Concepts"
heading: "The system prompt is cached"
group: "Overview"
order: 68
url: "/docs/concepts#the-system-prompt-is-cached"
anchor: "the-system-prompt-is-cached"
terms: ["system","prompt","cached","title","tree","section","summaries","injected","cache","control","marker","across","rounds","rather","sent","tokens","answer","turn","changes","tool","cannot","reuse","search","last","call","anyway","loop","model","defaults","fast","configurable","reader","server","side","should","confused","consumer","coding","agent","navigating"]
hash: "740d6a77b01eeb8a97450a731c865c453c8e7aac42a1971856ecacdf05622c39"
mode: "agent-primary"
facts: [{"kind":"code","literal":"cache_control","chunkId":"concepts#the-system-prompt-is-cached"},{"kind":"value","literal":"4.5","chunkId":"concepts#the-system-prompt-is-cached"}]
sources: [{"chunkId":"concepts#the-system-prompt-is-cached","url":"/docs/concepts#the-system-prompt-is-cached","anchor":"the-system-prompt-is-cached"}]
---

The title-tree and section summaries are injected into the system prompt with a cache-control marker so across rounds it is a prompt-cache hit rather than re-sent tokens; the answer turn changes the tool set so it cannot reuse the search rounds' cache but is the last call anyway. The loop model defaults to a fast model and is configurable, and the reader's server-side loop should not be confused with a consumer's coding agent navigating the files itself, though both climb the same four-rung ladder over the same tree.
