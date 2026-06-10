---
id: "api/endpoint#digest-reads-get"
title: "Search endpoint"
heading: "Digest reads (GET)"
group: "API"
order: 25
url: "/docs/api/endpoint#digest-reads-get"
anchor: "digest-reads-get"
terms: ["digest","reads","these","routes","read","virtual","never","call","model","require","route","response","glossary","terms","glossaryentry","term","matched","alias","error","found","sections","sectionsummary","group","digestnode","overview","string","context","archive","title","heading","2fcli","23flags","endpoint","head","content","hash","section","summaries","filtered","full"]
hash: "e462a9453db8dbf2c1ce2211b0989c34a851fefa9e2dc075bb461f12817eff77"
mode: "source-primary"
facts: [{"kind":"code","literal":"{ \"error\": \"Not found.\" }","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"virtual:hev-ask/digest","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"GET /api/ask/glossary","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"{ \"terms\": GlossaryEntry[] }","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"GET /api/ask/glossary/{term}","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"GlossaryEntry","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"GET /api/ask/sections","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"{ \"sections\": SectionSummary[] }","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"GET /api/ask/sections?group=API","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"GET /api/ask/sections/{id}","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"DigestNode","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"GET /api/ask/overview","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"{ \"overview\": string, \"context\": string }","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"GET /api/ask/archive","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":".hev-ask/","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"SectionSummary","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"{ id, title, heading, group, url }","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"/api/ask/sections/api%2Fcli%23flags","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"ask mcp --endpoint","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"HEAD /api/ask/archive","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"x-hev-ask-content-hash","chunkId":"api/endpoint#digest-reads-get"},{"kind":"code","literal":"404","chunkId":"api/endpoint#digest-reads-get"}]
sources: [{"chunkId":"api/endpoint#digest-reads-get","url":"/docs/api/endpoint#digest-reads-get","anchor":"digest-reads-get"}]
---

Digest reads (GET) These routes read virtual:hev-ask/digest, never call a model, and never require an API key: Route Response GET /api/ask/glossary { "terms": GlossaryEntry[] } GET /api/ask/glossary/{term} one GlossaryEntry, matched by term or alias GET /api/a...
