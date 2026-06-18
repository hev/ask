---
id: "api/configuration"
title: "Configuration"
heading: null
group: "API"
order: 10
url: "/docs/api/configuration"
anchor: null
terms: ["integration","package","default","export","takes","options","object","only","content","collection","list","effectively","required","everything","else","including","base","path","model","result","caps","astro","config","import","hevask","hevmind","defineconfig","integrations","collections","docs","basepath","claude","haiku","maxresults","callout","codetabs","every","option","models","endpoint"]
hash: "e7c504dd58a7e10c8e9f034b0d0ce1286da961ec8f56fbe2ad5ad932f90fcbf6"
mode: "source-primary"
facts: [{"kind":"code","literal":"// astro.config.mjs\nimport hevAsk from \"@hevmind/ask\";\n\nexport default defineConfig({\n  integrations: [\n    hevAsk({\n      collections: [\"docs\"],\n      basePath: \"/docs/\",\n      model: \"claude-haiku-4-5\",\n      maxResults: 6,\n    }),\n  ],\n});","chunkId":"api/configuration"},{"kind":"code","literal":"hevAsk()","chunkId":"api/configuration"},{"kind":"code","literal":"@hevmind/ask","chunkId":"api/configuration"},{"kind":"code","literal":"collections","chunkId":"api/configuration"},{"kind":"value","literal":"Callout.astro","chunkId":"api/configuration"},{"kind":"value","literal":"CodeTabs.astro","chunkId":"api/configuration"}]
sources: [{"chunkId":"api/configuration","url":"/docs/api/configuration","anchor":null}]
---

The integration is the package's default export and takes one options object in which only the content collection list is effectively required; everything else has a default, including base path, model, and result caps.
