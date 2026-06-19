---
id: "limits#agentic-search-adds-latency"
title: "Limits"
heading: "Agentic search adds latency"
group: "Overview"
order: 81
url: "/docs/limits#agentic-search-adds-latency"
anchor: "agentic-search-adds-latency"
terms: ["agentic","search","adds","latency","path","bounded","configured","number","model","round","trips","worst","case","seconds","instant","nature","keyword","always","available","lane","while","considered","iteration","tuned","down","tighter","ceiling","maxiterations","default","claude","tune","need"]
hash: "39f47a0577b861e2dbe61b168c9018a6b5514c62ce98b8c23fb620d1d0a7e717"
mode: "agent-primary"
facts: [{"kind":"code","literal":"maxIterations","chunkId":"limits#agentic-search-adds-latency"}]
sources: [{"chunkId":"limits#agentic-search-adds-latency","url":"/docs/limits#agentic-search-adds-latency","anchor":"agentic-search-adds-latency"}]
---

The agentic path is bounded by the configured number of model round-trips, so worst case is a few seconds and it is not instant by nature. The keyword path is the always-available instant lane while agentic search is the considered one, and the iteration cap can be tuned down for a tighter ceiling.
