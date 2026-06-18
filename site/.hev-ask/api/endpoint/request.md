---
id: "api/endpoint#request"
title: "Search endpoint"
heading: "Request"
group: "API"
order: 32
url: "/docs/api/endpoint#request"
anchor: "request"
terms: ["request","requests","post","json","body","carrying","query","optional","mode","keyword","forces","instant","path","agentic","loop","empty","whitespace","returns","result","omitted","behaves","like","present","does","autoscaling","work","string","field","type","description","search"]
hash: "f2a05f3dd644d4d9abb081702bbbc0fe721849ccf0e70f59742aea2ae18980a0"
mode: "source-primary"
facts: [{"kind":"code","literal":"{\n  \"query\": \"how does autoscaling work\",\n  \"mode\": \"agentic\"\n}","chunkId":"api/endpoint#request"},{"kind":"code","literal":"POST","chunkId":"api/endpoint#request"},{"kind":"code","literal":"query","chunkId":"api/endpoint#request"},{"kind":"code","literal":"string","chunkId":"api/endpoint#request"},{"kind":"code","literal":"mode","chunkId":"api/endpoint#request"},{"kind":"code","literal":"'keyword' \\| 'agentic'","chunkId":"api/endpoint#request"},{"kind":"code","literal":"keyword","chunkId":"api/endpoint#request"},{"kind":"code","literal":"agentic","chunkId":"api/endpoint#request"}]
sources: [{"chunkId":"api/endpoint#request","url":"/docs/api/endpoint#request","anchor":"request"}]
---

Requests are a POST with a JSON body carrying the query and an optional mode, where keyword forces the instant path and agentic requests the loop; an empty or whitespace query returns an empty result set and an omitted mode behaves like agentic when a key is present.
