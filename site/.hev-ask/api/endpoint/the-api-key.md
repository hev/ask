---
id: "api/endpoint#the-api-key"
title: "Search endpoint"
heading: "The API key"
group: "API"
order: 34
url: "/docs/api/endpoint#the-api-key"
anchor: "the-api-key"
terms: ["endpoint","reads","named","configured","provider","adapter","runtime","environment","process","build","time","order","wherever","host","injects","server","secrets","never","sent","browser","anthropic","openai","openrouter","locals","import","meta","anthropicapikey","default","openaiapikey","openrouterapikey","option","cloudflare"]
hash: "6b9f401de0da470a422dab86c8b67ba6461101611461a5f6ad508bcf0c854ece"
mode: "source-primary"
facts: [{"kind":"code","literal":"ANTHROPIC_API_KEY","chunkId":"api/endpoint#the-api-key"},{"kind":"code","literal":"OPENAI_API_KEY","chunkId":"api/endpoint#the-api-key"},{"kind":"code","literal":"OPENROUTER_API_KEY","chunkId":"api/endpoint#the-api-key"},{"kind":"code","literal":"provider","chunkId":"api/endpoint#the-api-key"},{"kind":"code","literal":"locals.runtime.env","chunkId":"api/endpoint#the-api-key"},{"kind":"code","literal":"process.env","chunkId":"api/endpoint#the-api-key"},{"kind":"code","literal":"import.meta.env","chunkId":"api/endpoint#the-api-key"},{"kind":"value","literal":"e.g","chunkId":"api/endpoint#the-api-key"}]
sources: [{"chunkId":"api/endpoint#the-api-key","url":"/docs/api/endpoint#the-api-key","anchor":"the-api-key"}]
---

The endpoint reads the key named by the configured provider from the adapter runtime environment, then the process environment, then the build-time environment, in that order. Set it wherever the host injects server secrets; it is never sent to the browser.
