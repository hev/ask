---
id: "api/search-overlay#keyboard-model"
title: "SearchOverlay component"
heading: "Keyboard model"
group: "API"
order: 46
url: "/docs/api/search-overlay#keyboard-model"
anchor: "keyboard-model"
terms: ["keyboard","model","overlay","first","number","words","typed","decides","path","opening","shows","suggested","questions","word","runs","debounced","keyless","keyword","search","result","auto","active","typing","space","switches","mode","enter","sends","question","agentic","loop","arrow","keys","move","selection","escape","closes","footer","hint","reflects"]
hash: "ee1d21380babc93166dcbdb53a1180c1fd71608113d3f633a28f81b7c40dd1b6"
mode: "source-primary"
facts: [{"kind":"code","literal":"Tab","chunkId":"api/search-overlay#keyboard-model"},{"kind":"code","literal":"ANTHROPIC_API_KEY","chunkId":"api/search-overlay#keyboard-model"}]
sources: [{"chunkId":"api/search-overlay#keyboard-model","url":"/docs/api/search-overlay#keyboard-model","anchor":"keyboard-model"}]
---

The overlay is ask-first and the number of words typed decides the path: opening with AI on shows suggested questions, one word runs debounced keyless keyword search with the first result auto-active, and typing a space switches to ask mode where Enter sends the question to the agentic loop. Arrow keys move the keyword selection, Escape closes, and the footer hint reflects the mode; with no server key, asking returns keyword results plus a shown warning rather than escalating.
