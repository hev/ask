---
id: "api/endpoint#the-api-key"
title: "Search endpoint"
heading: "The API key"
group: "API"
order: 34
url: "/docs/api/endpoint#the-api-key"
anchor: "the-api-key"
terms: ["endpoint","reads","named","configured","provider","looking","order","adapter","runtime","environment","process","import","time","wherever","host","injects","server","secrets","never","sent","browser","anthropic","openai","openrouter","locals","meta","anthropicapikey","default","openaiapikey","openrouterapikey","option","cloudflare"]
hash: "6b9f401de0da470a422dab86c8b67ba6461101611461a5f6ad508bcf0c854ece"
mode: "source-primary"
facts: [{"kind":"code","literal":"ANTHROPIC_API_KEY","chunkId":"api/endpoint#the-api-key"},{"kind":"code","literal":"OPENAI_API_KEY","chunkId":"api/endpoint#the-api-key"},{"kind":"code","literal":"OPENROUTER_API_KEY","chunkId":"api/endpoint#the-api-key"},{"kind":"code","literal":"provider","chunkId":"api/endpoint#the-api-key"},{"kind":"code","literal":"locals.runtime.env","chunkId":"api/endpoint#the-api-key"},{"kind":"code","literal":"process.env","chunkId":"api/endpoint#the-api-key"},{"kind":"code","literal":"import.meta.env","chunkId":"api/endpoint#the-api-key"},{"kind":"value","literal":"e.g","chunkId":"api/endpoint#the-api-key"}]
sources: [{"chunkId":"api/endpoint#the-api-key","url":"/docs/api/endpoint#the-api-key","anchor":"the-api-key"}]
---

The endpoint reads the key named by the configured provider, looking in order at the adapter runtime environment, the process environment, and the import-time environment; you set it wherever your host injects server secrets and it is never sent to the browser.
