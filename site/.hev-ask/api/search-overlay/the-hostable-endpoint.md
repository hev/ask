---
id: "api/search-overlay#the-hostable-endpoint"
title: "SearchOverlay component"
heading: "The hostable endpoint"
group: "API"
order: 53
url: "/docs/api/search-overlay#the-hostable-endpoint"
anchor: "the-hostable-endpoint"
terms: ["hostable","endpoint","agentic","answers","without","astro","deploy","bounded","answer","loop","standalone","service","cloudflare","worker","node","server","vercel","function","serves","same","contract","holds","side","reads","committed","tree","once","point","number","sites","overlays","keyword","search","stays","browser","only","needs","anything","deployed","scaffold"]
hash: "de2a87abed98b0710faec2697d9c06ee50122f8ab3852fc8e2133a2af0333886"
mode: "source-primary"
facts: [{"kind":"code","literal":"# scaffold and deploy the Worker flavor\nask endpoint init --target cloudflare\nwrangler deploy            # set ANTHROPIC_API_KEY as a secret","chunkId":"api/search-overlay#the-hostable-endpoint"},{"kind":"code","literal":"POST /api/ask","chunkId":"api/search-overlay#the-hostable-endpoint"},{"kind":"code","literal":"ANTHROPIC_API_KEY","chunkId":"api/search-overlay#the-hostable-endpoint"}]
sources: [{"chunkId":"api/search-overlay#the-hostable-endpoint","url":"/docs/api/search-overlay#the-hostable-endpoint","anchor":"the-hostable-endpoint"}]
---

For agentic answers without Astro, deploy the bounded answer loop as a standalone service (a Cloudflare Worker, a Node server, or a Vercel function) that serves the same endpoint contract, holds the key server-side, and reads the committed tree. Deploy it once and point any number of sites' overlays at it; keyword search stays in the browser, so only the answer loop needs anything deployed, and a scaffold command bootstraps the Worker flavor.
