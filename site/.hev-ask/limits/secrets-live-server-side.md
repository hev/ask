---
id: "limits#secrets-live-server-side"
title: "Limits"
heading: "Secrets live server-side"
group: "Overview"
order: 70
url: "/docs/limits#secrets-live-server-side"
anchor: "secrets-live-server-side"
terms: ["secrets","live","server","side","agentic","path","needs","anthropicapikey","environment","runs","never","exposed","browser","present","runtime","endpoint","serves","keyword","results","search","degrades","doesn","break","anthropic"]
hash: "2e90f9ec17b623efcd2a207bab1d0cb762dcdac4b14f147a7bea174bdf354a10"
mode: "agent-primary"
facts: [{"kind":"code","literal":"ANTHROPIC_API_KEY","chunkId":"limits#secrets-live-server-side"},{"kind":"code","literal":"/api/ask","chunkId":"limits#secrets-live-server-side"}]
sources: [{"chunkId":"limits#secrets-live-server-side","url":"/docs/limits#secrets-live-server-side","anchor":"secrets-live-server-side"}]
---

Secrets live server-side The agentic path needs ANTHROPICAPIKEY in the server environment that runs /api/ask. The key is never exposed to the browser. If the key isn't present at runtime, the endpoint serves keyword results — search degrades, it doesn't break.
