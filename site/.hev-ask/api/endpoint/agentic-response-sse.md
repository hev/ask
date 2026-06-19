---
id: "api/endpoint#agentic-response-sse"
title: "Search endpoint"
heading: "Agentic response (SSE)"
group: "API"
order: 25
url: "/docs/api/endpoint#agentic-response-sse"
anchor: "agentic-response-sse"
terms: ["agentic","response","present","mode","requested","endpoint","streams","answer","named","frames","whose","example","payloads","meaning","documented","table","search","context","model","gathered","time","grounding","source","streamed","text","deltas","completion","post","stream","error","carries","title","optional","heading","group","snippet","since","prose","substance","links"]
hash: "366017faa11220e7fdd7545cd885b58cdf2133c5f21be853f4339fe775ebc51a"
mode: "source-primary"
facts: [{"kind":"code","literal":"mode","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"agentic","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"content-type: text/event-stream","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"search","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{ query }","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"sources","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{ sources: Source[], model, mode }","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"token","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{ text }","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"done","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{}","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"error","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{ error }","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"200","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"Source","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{ title, heading?, url, group? }","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"snippet","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"url","chunkId":"api/endpoint#agentic-response-sse"}]
sources: [{"chunkId":"api/endpoint#agentic-response-sse","url":"/docs/api/endpoint#agentic-response-sse","anchor":"agentic-response-sse"}]
---

When a key is present and agentic mode is requested, the endpoint streams the answer as named SSE frames whose example payloads and meaning are documented in a table: search context the model gathered, the one-time grounding source set, streamed answer-text deltas, completion, and a post-stream error. A source carries title, optional heading, URL, and group but no snippet, since the prose carries the substance and links point at the URL.
