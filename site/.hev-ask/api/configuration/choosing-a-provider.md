---
id: "api/configuration#choosing-a-provider"
title: "Configuration"
heading: "Choosing a provider"
group: "API"
order: 11
url: "/docs/api/configuration#choosing-a-provider"
anchor: "choosing-a-provider"
terms: ["choosing","provider","option","selects","serves","both","runtime","answer","loop","offline","digest","build","identical","search","behavior","format","endpoint","contract","across","providers","only","model","environment","variable","change","openrouter","reaches","every","routes","through","base","override","points","openai","compatible","client","chat","completions","models","must"]
hash: "64020ee9bbb713bf179bf7c6d3f2673b780e483e4c6b7f584c829f690db6acd3"
mode: "source-primary"
facts: [{"kind":"code","literal":"// astro.config.mjs — the default; reads ANTHROPIC_API_KEY\nhevAsk({\n  collections: [\"docs\"],\n  // model defaults to claude-haiku-4-5\n});","chunkId":"api/configuration#choosing-a-provider"},{"kind":"code","literal":"// astro.config.mjs — reads OPENAI_API_KEY\nhevAsk({\n  collections: [\"docs\"],\n  provider: \"openai\",\n  // model defaults to gpt-4.1-mini\n});","chunkId":"api/configuration#choosing-a-provider"},{"kind":"code","literal":"// astro.config.mjs — reads OPENROUTER_API_KEY\nhevAsk({\n  collections: [\"docs\"],\n  provider: \"openrouter\",\n  model: \"anthropic/claude-haiku-4.5\", // or any OpenRouter model slug\n});","chunkId":"api/configuration#choosing-a-provider"},{"kind":"code","literal":"hevAsk({\n  collections: [\"docs\"],\n  provider: \"openai\",\n  providerBaseUrl: \"https://my-gateway.example.com/v1\",\n  model: \"my-model\",\n});","chunkId":"api/configuration#choosing-a-provider"},{"kind":"code","literal":"provider","chunkId":"api/configuration#choosing-a-provider"},{"kind":"code","literal":"OPENROUTER_API_KEY","chunkId":"api/configuration#choosing-a-provider"},{"kind":"code","literal":"providerBaseUrl","chunkId":"api/configuration#choosing-a-provider"},{"kind":"code","literal":"provider: \"openai\"","chunkId":"api/configuration#choosing-a-provider"},{"kind":"code","literal":"endpoint","chunkId":"api/configuration#choosing-a-provider"},{"kind":"code","literal":"/api/ask","chunkId":"api/configuration#choosing-a-provider"},{"kind":"code","literal":"SearchOverlay","chunkId":"api/configuration#choosing-a-provider"}]
sources: [{"chunkId":"api/configuration#choosing-a-provider","url":"/docs/api/configuration#choosing-a-provider","anchor":"choosing-a-provider"}]
---

The provider option selects who serves both the runtime answer loop and the offline digest build, with identical search behavior, digest format, and endpoint contract across providers; only the model and key environment variable change. OpenRouter reaches every model it routes through one key, and a base-URL override points the OpenAI-compatible client at any Chat Completions endpoint. Loop models must support tool calling and the digest builder needs forced tool choice, which current major-provider models satisfy; a changed endpoint route must be mirrored on the overlay component.
