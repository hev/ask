---
id: "limits#recall-has-a-keyword-ceiling"
title: "Limits"
heading: "Recall has a keyword ceiling"
group: "Overview"
order: 84
url: "/docs/limits#recall-has-a-keyword-ceiling"
anchor: "recall-has-a-keyword-ceiling"
terms: ["recall","keyword","ceiling","retrieval","token","overlap","widened","glossary","rather","embeddings","agentic","loop","only","ground","finds","recovers","most","synonym","cases","reader","searching","language","shares","tokens","docs","never","surface","right","section","known","paraphrase","deliberately","built","until","analytics","show","consistent","misses","richer","cheaper"]
hash: "df75d2232fca5a321ba6501fe42757e069ef11d194231d4a0af212994581279a"
mode: "agent-primary"
facts: [{"kind":"code","literal":"k8s","chunkId":"limits#recall-has-a-keyword-ceiling"},{"kind":"code","literal":"kubernetes","chunkId":"limits#recall-has-a-keyword-ceiling"}]
sources: [{"chunkId":"limits#recall-has-a-keyword-ceiling","url":"/docs/limits#recall-has-a-keyword-ceiling","anchor":"recall-has-a-keyword-ceiling"}]
---

Retrieval is keyword token-overlap widened by the glossary rather than embeddings, and the agentic loop can only ground in what retrieval finds. The glossary recovers most synonym cases, but a reader searching in language that shares no tokens with the docs and isn't in the glossary may never surface the right section; embeddings are the known fix for paraphrase recall and are deliberately not built yet, so until analytics show consistent misses, a richer glossary is the cheaper lever.
