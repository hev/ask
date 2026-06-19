---
id: "concepts#the-system-prompt-is-cached"
title: "Concepts"
heading: "The system prompt is cached"
group: "Overview"
order: 68
url: "/docs/concepts#the-system-prompt-is-cached"
anchor: "the-system-prompt-is-cached"
terms: ["system","prompt","cached","title","tree","section","summaries","injected","cache","marker","across","search","rounds","rather","sent","tokens","final","answer","turn","changes","tool","cannot","reuse","last","call","anyway","loop","model","defaults","small","configurable","page","warns","confuse","reader","server","side","synthesis","consumer","coding"]
hash: "740d6a77b01eeb8a97450a731c865c453c8e7aac42a1971856ecacdf05622c39"
mode: "agent-primary"
facts: [{"kind":"code","literal":"cache_control","chunkId":"concepts#the-system-prompt-is-cached"},{"kind":"value","literal":"4.5","chunkId":"concepts#the-system-prompt-is-cached"}]
sources: [{"chunkId":"concepts#the-system-prompt-is-cached","url":"/docs/concepts#the-system-prompt-is-cached","anchor":"the-system-prompt-is-cached"}]
---

The title-tree and section summaries are injected into the system prompt with a cache marker, so across the search rounds it is a prompt-cache hit rather than re-sent tokens; the final answer turn changes the tool set and so cannot reuse that cache, but it is the last call anyway. The loop model defaults to a small model and is configurable, and the page warns not to confuse the reader's server-side synthesis loop with a consumer's own coding agent that navigates the files itself, though both climb the same ladder.
