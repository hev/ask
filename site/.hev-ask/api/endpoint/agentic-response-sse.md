---
id: "api/endpoint#agentic-response-sse"
title: "Search endpoint"
heading: "Agentic response (SSE)"
group: "API"
order: 25
url: "/docs/api/endpoint#agentic-response-sse"
anchor: "agentic-response-sse"
terms: ["agentic","response","present","mode","endpoint","streams","answer","named","server","sent","events","search","frames","report","context","model","gathered","single","sources","frame","before","tokens","carries","grounding","token","carry","deltas","done","closes","stream","error","reports","post","failure","since","http","status","already","success","source"]
hash: "366017faa11220e7fdd7545cd885b58cdf2133c5f21be853f4339fe775ebc51a"
mode: "source-primary"
facts: [{"kind":"code","literal":"mode","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"agentic","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"content-type: text/event-stream","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"search","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{ query }","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"sources","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{ sources: Source[], model, mode }","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"token","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{ text }","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"done","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{}","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"error","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{ error }","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"200","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"Source","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{ title, heading?, url, group? }","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"snippet","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"url","chunkId":"api/endpoint#agentic-response-sse"}]
sources: [{"chunkId":"api/endpoint#agentic-response-sse","url":"/docs/api/endpoint#agentic-response-sse","anchor":"agentic-response-sse"}]
---

When a key is present and the mode is agentic, the endpoint streams the answer as named Server-Sent Events: search frames report the context the model gathered, a single sources frame sent before any tokens carries the grounding set and model and mode, token frames carry answer deltas, a done frame closes the stream, and an error frame reports a post-stream failure since the HTTP status is already success. Each source carries title, optional heading, URL, and optional group but no snippet, since the prose carries the substance and links point at the URL.
