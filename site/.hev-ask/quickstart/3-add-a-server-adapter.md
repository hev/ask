---
id: "quickstart#3-add-a-server-adapter"
title: "Quick start"
heading: "3. Add a server adapter"
group: "Overview"
order: 92
url: "/docs/quickstart#3-add-a-server-adapter"
anchor: "3-add-a-server-adapter"
terms: ["server","adapter","hybrid","matching","host","because","endpoint","route","renders","demand","while","existing","pages","stay","prerendered","only","runs","function","example","uses","cloudflare","astro","config","import","astrojs","export","default","defineconfig","platformproxy","enabled","true","integrations","above","whichever","matches","site"]
hash: "795eff804c3ce56e4f62c96e59090a4c4845680aef5ab40735d8a57b33976dcf"
mode: "agent-primary"
facts: [{"kind":"code","literal":"// astro.config.mjs\nimport cloudflare from \"@astrojs/cloudflare\";\n\nexport default defineConfig({\n  adapter: cloudflare({ platformProxy: { enabled: true } }),\n  // ...integrations as above\n});","chunkId":"quickstart#3-add-a-server-adapter"},{"kind":"code","literal":"/api/ask","chunkId":"quickstart#3-add-a-server-adapter"}]
sources: [{"chunkId":"quickstart#3-add-a-server-adapter","url":"/docs/quickstart#3-add-a-server-adapter","anchor":"3-add-a-server-adapter"}]
---

Add a server or hybrid adapter matching your host because the endpoint route renders on demand while existing pages stay prerendered and only that route runs as a function; the example uses Cloudflare.
