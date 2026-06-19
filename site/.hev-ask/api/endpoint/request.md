---
id: "api/endpoint#request"
title: "Search endpoint"
heading: "Request"
group: "API"
order: 32
url: "/docs/api/endpoint#request"
anchor: "request"
terms: ["request","endpoint","takes","post","json","body","documented","table","query","string","empty","whitespace","returns","results","optional","mode","forces","instant","path","requests","loop","defaulting","agentic","behavior","present","does","autoscaling","work","keyword","field","type","description","search","result","omitted","behaves","like"]
hash: "f2a05f3dd644d4d9abb081702bbbc0fe721849ccf0e70f59742aea2ae18980a0"
mode: "source-primary"
facts: [{"kind":"code","literal":"{\n  \"query\": \"how does autoscaling work\",\n  \"mode\": \"agentic\"\n}","chunkId":"api/endpoint#request"},{"kind":"code","literal":"POST","chunkId":"api/endpoint#request"},{"kind":"code","literal":"query","chunkId":"api/endpoint#request"},{"kind":"code","literal":"string","chunkId":"api/endpoint#request"},{"kind":"code","literal":"mode","chunkId":"api/endpoint#request"},{"kind":"code","literal":"'keyword' \\| 'agentic'","chunkId":"api/endpoint#request"},{"kind":"code","literal":"keyword","chunkId":"api/endpoint#request"},{"kind":"code","literal":"agentic","chunkId":"api/endpoint#request"}]
sources: [{"chunkId":"api/endpoint#request","url":"/docs/api/endpoint#request","anchor":"request"}]
---

The endpoint takes a POST with a JSON body documented in a table: a query string where empty or whitespace returns no results, and an optional mode that forces the instant path or requests the loop, defaulting to agentic behavior when a key is present.
