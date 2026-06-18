---
id: "quickstart#enable-agentic-search"
title: "Quick start"
heading: "Enable agentic search"
group: "Overview"
order: 95
url: "/docs/quickstart#enable-agentic-search"
anchor: "enable-agentic-search"
terms: ["enable","agentic","search","provider","server","environment","endpoint","runs","host","secrets","local","file","present","pressing","enter","loop","self","issued","queries","grounded","answer","inline","deep","links","without","returns","keyword","results","default","needs","only","while","openai","openrouter","also","need","their","integration","options","nothing"]
hash: "ef1e17ef641cdc4792581c2c5f7d313aa573f6ccc6cb923f2288488784deac0e"
mode: "agent-primary"
facts: [{"kind":"code","literal":"# the default provider — nothing else to configure\nexport ANTHROPIC_API_KEY=sk-ant-...","chunkId":"quickstart#enable-agentic-search"},{"kind":"code","literal":"# with provider: \"openai\" in the hevAsk() options\nexport OPENAI_API_KEY=sk-...","chunkId":"quickstart#enable-agentic-search"},{"kind":"code","literal":"# with provider: \"openrouter\" in the hevAsk() options\nexport OPENROUTER_API_KEY=sk-or-...","chunkId":"quickstart#enable-agentic-search"},{"kind":"code","literal":"/api/ask","chunkId":"quickstart#enable-agentic-search"},{"kind":"code","literal":".env","chunkId":"quickstart#enable-agentic-search"},{"kind":"code","literal":"provider","chunkId":"quickstart#enable-agentic-search"}]
sources: [{"chunkId":"quickstart#enable-agentic-search","url":"/docs/quickstart#enable-agentic-search","anchor":"enable-agentic-search"}]
---

Set the provider's API key in the server environment where the endpoint runs (host secrets or a local env file); with a key present, pressing Enter runs the agentic loop with self-issued sub-queries, a grounded answer, and inline deep links, and without one Enter returns keyword results. The default provider needs only the key while OpenAI and OpenRouter also need their provider set in the integration options.
