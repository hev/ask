---
id: "concepts#the-agentic-search-loop"
title: "Concepts"
heading: "The agentic search loop"
group: "Overview"
order: 66
url: "/docs/concepts#the-agentic-search-loop"
anchor: "the-agentic-search-loop"
terms: ["agentic","search","loop","because","overlay","answers","humans","does","synthesis","agent","would","otherwise","itself","running","same","disclosure","ladder","server","side","memory","tree","phases","gather","phase","model","gets","full","title","plus","tool","open","section","summary","facts","reference","sections","source","text","opening","only"]
hash: "2f6dfa7319e162069240585c70a19b0d2727aa413c0bddd5dd0644f360619892"
mode: "agent-primary"
facts: [{"kind":"code","literal":"open_section({ id })","chunkId":"concepts#the-agentic-search-loop"},{"kind":"code","literal":"facts","chunkId":"concepts#the-agentic-search-loop"},{"kind":"code","literal":"maxIterations","chunkId":"concepts#the-agentic-search-loop"},{"kind":"code","literal":"url","chunkId":"concepts#the-agentic-search-loop"}]
sources: [{"chunkId":"concepts#the-agentic-search-loop","url":"/docs/concepts#the-agentic-search-loop","anchor":"the-agentic-search-loop"}]
---

Because the overlay answers humans, it does the synthesis an agent would otherwise do itself, running the same disclosure ladder server-side over the in-memory tree in two phases. In the gather phase the model gets the full title-tree plus one tool to open a section's summary, facts, and (for reference sections) source text, opening only what it needs up to a round cap and citing only what it opened; in the answer phase the accumulated sources are sent to the overlay and the model is called once more with no tools so it can only write grounded, inline-linked prose. Dropping the tools on the final turn guarantees it answers instead of searching again, and it can ground and link only in what retrieval returned.
