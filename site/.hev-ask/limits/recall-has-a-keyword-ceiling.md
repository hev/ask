---
id: "limits#recall-has-a-keyword-ceiling"
title: "Limits"
heading: "Recall has a keyword ceiling"
group: "Overview"
order: 84
url: "/docs/limits#recall-has-a-keyword-ceiling"
anchor: "recall-has-a-keyword-ceiling"
terms: ["recall","keyword","ceiling","retrieval","token","overlap","widened","glossary","rather","embeddings","loop","only","ground","finds","recovers","most","synonym","abbreviation","cases","reader","searching","language","shares","tokens","docs","never","surface","right","section","known","paraphrase","deliberately","built","richer","cheaper","lever","until","analytics","show","questions"]
hash: "df75d2232fca5a321ba6501fe42757e069ef11d194231d4a0af212994581279a"
mode: "agent-primary"
facts: [{"kind":"code","literal":"k8s","chunkId":"limits#recall-has-a-keyword-ceiling"},{"kind":"code","literal":"kubernetes","chunkId":"limits#recall-has-a-keyword-ceiling"}]
sources: [{"chunkId":"limits#recall-has-a-keyword-ceiling","url":"/docs/limits#recall-has-a-keyword-ceiling","anchor":"recall-has-a-keyword-ceiling"}]
---

Retrieval is keyword token-overlap widened by the glossary rather than embeddings, and the loop can only ground in what retrieval finds; the glossary recovers most synonym and abbreviation cases, but a reader searching in language that shares no tokens with your docs and is not in the glossary may never surface the right section. Embeddings are the known fix for paraphrase recall and are deliberately not built yet, so a richer glossary is the cheaper lever until analytics show questions consistently missing.
