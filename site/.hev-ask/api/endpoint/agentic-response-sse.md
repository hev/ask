---
id: "api/endpoint#agentic-response-sse"
title: "Search endpoint"
heading: "Agentic response (SSE)"
group: "API"
order: 24
url: "/docs/api/endpoint#agentic-response-sse"
anchor: "agentic-response-sse"
terms: ["agentic","response","present","mode","endpoint","responds","content","type","text","event","stream","streams","answer","generated","search","query","sources","source","model","token","done","error","title","heading","group","snippet","named","frame","data","autoscaling","core","concepts","kubernetes","docs","overview","claude","haiku","scales","workers","based"]
hash: "366017faa11220e7fdd7545cd885b58cdf2133c5f21be853f4339fe775ebc51a"
mode: "source-primary"
facts: [{"kind":"code","literal":"mode","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"agentic","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"content-type: text/event-stream","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"search","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{ query }","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"sources","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{ sources: Source[], model, mode }","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"token","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{ text }","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"done","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{}","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"error","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{ error }","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"200","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"Source","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"{ title, heading?, url, group? }","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"snippet","chunkId":"api/endpoint#agentic-response-sse"},{"kind":"code","literal":"url","chunkId":"api/endpoint#agentic-response-sse"}]
sources: [{"chunkId":"api/endpoint#agentic-response-sse","url":"/docs/api/endpoint#agentic-response-sse","anchor":"agentic-response-sse"}]
---

Agentic response (SSE) When a key is present and mode is agentic, the endpoint responds with content-type: text/event-stream and streams the answer as it is generated....
