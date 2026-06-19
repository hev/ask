---
id: "api/endpoint#errors"
title: "Search endpoint"
heading: "Errors"
group: "API"
order: 27
url: "/docs/api/endpoint#errors"
anchor: "errors"
terms: ["errors","table","error","responses","request","invalid","json","found","unknown","read","route","missing","term","section","server","chunk","index","fails","build","failure","during","agentic","stream","arrives","final","event","because","http","status","already","successful","body","cause","wasn","valid","digest","glossary","failed","misconfigured","collection"]
hash: "c67a8ab9b70204c62cb0ae5f63d37f2d21863e6f4fdc6e684b85e3a81e37ff3d"
mode: "source-primary"
facts: [{"kind":"code","literal":"400","chunkId":"api/endpoint#errors"},{"kind":"code","literal":"{ \"error\": \"Invalid JSON body.\" }","chunkId":"api/endpoint#errors"},{"kind":"code","literal":"404","chunkId":"api/endpoint#errors"},{"kind":"code","literal":"{ \"error\": \"…\" }","chunkId":"api/endpoint#errors"},{"kind":"code","literal":"500","chunkId":"api/endpoint#errors"},{"kind":"code","literal":"event: error","chunkId":"api/endpoint#errors"},{"kind":"code","literal":"200","chunkId":"api/endpoint#errors"},{"kind":"code","literal":"error","chunkId":"api/endpoint#errors"},{"kind":"value","literal":"e.g","chunkId":"api/endpoint#errors"}]
sources: [{"chunkId":"api/endpoint#errors","url":"/docs/api/endpoint#errors","anchor":"errors"}]
---

A table of error responses: a bad request for invalid JSON, not-found for an unknown read route or missing term or section, and a server error when the chunk index fails to build. A failure during the agentic stream arrives as a final SSE error event because the HTTP status is already successful.
