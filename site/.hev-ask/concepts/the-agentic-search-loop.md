---
id: "concepts#the-agentic-search-loop"
title: "Concepts"
heading: "The agentic search loop"
group: "Overview"
order: 66
url: "/docs/concepts#the-agentic-search-loop"
anchor: "the-agentic-search-loop"
terms: ["agentic","search","loop","overlay","answers","humans","doing","synthesis","agent","would","itself","sending","multi","word","query","bounded","tool","phases","gather","phase","model","given","title","tree","every","section","plus","open","summary","facts","reference","sections","source","text","opening","only","needs","iteration","citing","opened"]
hash: "2f6dfa7319e162069240585c70a19b0d2727aa413c0bddd5dd0644f360619892"
mode: "agent-primary"
facts: [{"kind":"code","literal":"open_section({ id })","chunkId":"concepts#the-agentic-search-loop"},{"kind":"code","literal":"facts","chunkId":"concepts#the-agentic-search-loop"},{"kind":"code","literal":"maxIterations","chunkId":"concepts#the-agentic-search-loop"},{"kind":"code","literal":"url","chunkId":"concepts#the-agentic-search-loop"}]
sources: [{"chunkId":"concepts#the-agentic-search-loop","url":"/docs/concepts#the-agentic-search-loop","anchor":"the-agentic-search-loop"}]
---

The overlay answers humans by doing the synthesis an agent would do itself, sending a multi-word query to a bounded tool-use loop in two phases. In the gather phase the model is given the title-tree of every section plus one tool to open a section's summary, facts, and (for reference sections) source text, opening only what it needs up to the iteration cap and citing only what it opened; in the answer phase the accumulated sources are sent to the overlay for link validation and the model is called once more with no tools so it can only write prose, streamed token-by-token. Dropping the tools on the final turn guarantees it answers rather than searching again, and it can only ground in and link to the sections retrieval returned.
