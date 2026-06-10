---
id: "api/configuration"
title: "Configuration"
heading: null
group: "API"
order: 10
url: "/docs/api/configuration"
anchor: null
terms: ["every","option","hevask","astro","integration","collections","models","endpoint","basepath","chunking","depth","retrieval","caps","digest","paths","defaults","default","export","hevmind","takes","options","object","config","import","defineconfig","integrations","docs","model","claude","haiku","maxresults","callout","only","effectively","required","everything","else"]
hash: "e7c504dd58a7e10c8e9f034b0d0ce1286da961ec8f56fbe2ad5ad932f90fcbf6"
mode: "source-primary"
facts: [{"kind":"code","literal":"// astro.config.mjs\nimport hevAsk from \"@hevmind/ask\";\n\nexport default defineConfig({\n  integrations: [\n    hevAsk({\n      collections: [\"docs\"],\n      basePath: \"/docs/\",\n      model: \"claude-haiku-4-5\",\n      maxResults: 6,\n    }),\n  ],\n});","chunkId":"api/configuration"},{"kind":"code","literal":"hevAsk()","chunkId":"api/configuration"},{"kind":"code","literal":"@hevmind/ask","chunkId":"api/configuration"},{"kind":"code","literal":"collections","chunkId":"api/configuration"},{"kind":"value","literal":"Callout.astro","chunkId":"api/configuration"}]
sources: [{"chunkId":"api/configuration","url":"/docs/api/configuration","anchor":null}]
---

Every option for the hevAsk() Astro integration: collections, models, endpoint and basePath, chunking depth, retrieval caps, and digest paths, with defaults. The hevAsk() integration is the default export of @hevmind/ask. It takes one options object....
