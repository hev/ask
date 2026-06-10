---
id: "quickstart#3-add-a-server-adapter"
title: "Quick start"
heading: "3. Add a server adapter"
group: "Overview"
order: 76
url: "/docs/quickstart#3-add-a-server-adapter"
anchor: "3-add-a-server-adapter"
terms: ["server","adapter","route","renders","demand","whichever","matches","host","existing","pages","stay","prerendered","only","runs","function","site","uses","cloudflare","astro","config","platformproxy","import","astrojs","export","default","defineconfig","enabled","true","integrations","above"]
hash: "795eff804c3ce56e4f62c96e59090a4c4845680aef5ab40735d8a57b33976dcf"
mode: "agent-primary"
facts: [{"kind":"code","literal":"// astro.config.mjs\nimport cloudflare from \"@astrojs/cloudflare\";\n\nexport default defineConfig({\n  adapter: cloudflare({ platformProxy: { enabled: true } }),\n  // ...integrations as above\n});","chunkId":"quickstart#3-add-a-server-adapter"},{"kind":"code","literal":"/api/ask","chunkId":"quickstart#3-add-a-server-adapter"}]
sources: [{"chunkId":"quickstart#3-add-a-server-adapter","url":"/docs/quickstart#3-add-a-server-adapter","anchor":"3-add-a-server-adapter"}]
---

3. Add a server adapter The /api/ask route renders on demand. Add whichever adapter matches your host (existing pages stay prerendered — only /api/ask runs as a function); this site uses Cloudflare: // astro.config.mjs adapter: cloudflare({ platformProxy: { en...
