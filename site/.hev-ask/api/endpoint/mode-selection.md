---
id: "api/endpoint#mode-selection"
title: "Search endpoint"
heading: "Mode selection"
group: "API"
order: 31
url: "/docs/api/endpoint#mode-selection"
anchor: "mode-selection"
terms: ["mode","selection","endpoint","decides","empty","query","returns","keyword","json","explicit","request","missing","agentic","without","plus","downgrade","warning","otherwise","streams","answer","there","unavailable","error","path","overlay","branches","response","content","type","handle","both","shapes","results","model","stream","downgrades","says","reader","still","gets"]
hash: "268456691571b15ea83b4fb47d86fa28225dcc41a170a2483ba5d5e97a419665"
mode: "source-primary"
facts: [{"kind":"code","literal":"{ results: [], query: \"\", model, mode: \"keyword\" }","chunkId":"api/endpoint#mode-selection"},{"kind":"code","literal":"mode: \"keyword\"","chunkId":"api/endpoint#mode-selection"},{"kind":"code","literal":"mode: \"agentic\"","chunkId":"api/endpoint#mode-selection"},{"kind":"code","literal":"warning","chunkId":"api/endpoint#mode-selection"},{"kind":"code","literal":"content-type","chunkId":"api/endpoint#mode-selection"}]
sources: [{"chunkId":"api/endpoint#mode-selection","url":"/docs/api/endpoint#mode-selection","anchor":"mode-selection"}]
---

The endpoint decides what to run: an empty query returns empty keyword JSON, an explicit keyword request or a missing key returns keyword JSON, an agentic request without a key returns keyword JSON plus a downgrade warning, and otherwise it streams the agentic answer. There is no unavailable-AI error path, and the overlay branches on the response content-type to handle both shapes from the one endpoint.
