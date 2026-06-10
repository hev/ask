---
id: "api/endpoint#the-api-key"
title: "Search endpoint"
heading: "The API key"
group: "API"
order: 33
url: "/docs/api/endpoint#the-api-key"
anchor: "the-api-key"
terms: ["endpoint","resolves","anthropicapikey","order","adapter","runtime","locals","cloudflare","process","import","meta","wherever","host","injects","server","secrets","never","sent","browser","anthropic"]
hash: "58add94e9c0749c8e4375c9547959baeef1548db0d57c8e5c99ecb5392a54247"
mode: "source-primary"
facts: [{"kind":"code","literal":"ANTHROPIC_API_KEY","chunkId":"api/endpoint#the-api-key"},{"kind":"code","literal":"locals.runtime.env","chunkId":"api/endpoint#the-api-key"},{"kind":"code","literal":"process.env","chunkId":"api/endpoint#the-api-key"},{"kind":"code","literal":"import.meta.env","chunkId":"api/endpoint#the-api-key"},{"kind":"value","literal":"e.g","chunkId":"api/endpoint#the-api-key"}]
sources: [{"chunkId":"api/endpoint#the-api-key","url":"/docs/api/endpoint#the-api-key","anchor":"the-api-key"}]
---

The API key The endpoint resolves ANTHROPICAPIKEY from, in order: the adapter runtime env (locals.runtime.env, e.g. Cloudflare), process.env, then import.meta.env. Set it wherever your host injects server secrets; it is never sent to the browser.
