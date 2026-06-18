---
id: "limits#secrets-live-server-side"
title: "Limits"
heading: "Secrets live server-side"
group: "Overview"
order: 85
url: "/docs/limits#secrets-live-server-side"
anchor: "secrets-live-server-side"
terms: ["secrets","live","server","side","agentic","path","needs","provider","environment","runs","endpoint","never","exposed","browser","absent","runtime","serves","keyword","results","search","degrades","rather","breaks","anthropic","openai","openrouter","anthropicapikey","default","openaiapikey","openrouterapikey","option","present","doesn","break"]
hash: "2adf59c2589370b5a968f7e95da1bd1c86a15fdc4f3dce2c0f994addb33152bc"
mode: "agent-primary"
facts: [{"kind":"code","literal":"ANTHROPIC_API_KEY","chunkId":"limits#secrets-live-server-side"},{"kind":"code","literal":"OPENAI_API_KEY","chunkId":"limits#secrets-live-server-side"},{"kind":"code","literal":"OPENROUTER_API_KEY","chunkId":"limits#secrets-live-server-side"},{"kind":"code","literal":"provider","chunkId":"limits#secrets-live-server-side"},{"kind":"code","literal":"/api/ask","chunkId":"limits#secrets-live-server-side"}]
sources: [{"chunkId":"limits#secrets-live-server-side","url":"/docs/limits#secrets-live-server-side","anchor":"secrets-live-server-side"}]
---

The agentic path needs the provider's API key in the server environment that runs the endpoint, the key is never exposed to the browser, and if it is absent at runtime the endpoint serves keyword results, so search degrades rather than breaks.
