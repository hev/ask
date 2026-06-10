---
id: "quickstart#2-register-the-integration"
title: "Quick start"
heading: "2. Register the integration"
group: "Overview"
order: 75
url: "/docs/quickstart#2-register-the-integration"
anchor: "2-register-the-integration"
terms: ["register","integration","astro","config","integrations","hevask","collections","docs","content","collection","name","basepath","slug","prefix","option","must","everything","else","import","defineconfig","hevmind","export","default","configuration","reference"]
hash: "02f816dd61687ca1704bca5caeb6868ea8d531804bd9dc7ca569c05522811b6d"
mode: "agent-primary"
facts: [{"kind":"code","literal":"// astro.config.mjs\nimport { defineConfig } from \"astro/config\";\nimport hevAsk from \"@hevmind/ask\";\n\nexport default defineConfig({\n  integrations: [\n    hevAsk({\n      collections: [\"docs\"],   // your content collection name(s)\n      basePath: \"/docs/\",      // slug → URL prefix: basePath + slug\n    }),\n  ],\n});","chunkId":"quickstart#2-register-the-integration"},{"kind":"code","literal":"collections","chunkId":"quickstart#2-register-the-integration"}]
sources: [{"chunkId":"quickstart#2-register-the-integration","url":"/docs/quickstart#2-register-the-integration","anchor":"2-register-the-integration"}]
---

2. Register the integration // astro.config.mjs integrations: [ hevAsk({ collections: ["docs"], // your content collection name(s) basePath: "/docs/", // slug → URL prefix: basePath + slug }), ], }); collections is the one option you must set; everything else...
