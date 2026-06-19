---
id: "api/endpoint#keyword-response-json"
title: "Search endpoint"
heading: "Keyword response (JSON)"
group: "API"
order: 29
url: "/docs/api/endpoint#keyword-response-json"
anchor: "keyword-response-json"
terms: ["keyword","response","json","mode","returns","successful","envelope","whose","fields","documented","table","ranked","results","title","optional","heading","group","snippet","echoed","query","configured","loop","model","warning","agentic","requested","result","carries","deep","link","anchor","appended","except","document","intro","chunk","concepts","search","docs","overview"]
hash: "bbc5f7240f935311df97464c2bf6780dc597f246cfa94d8bfd9176bfe72cdd40"
mode: "source-primary"
facts: [{"kind":"code","literal":"{\n  \"results\": [\n    {\n      \"title\": \"Concepts\",\n      \"heading\": \"The agentic search loop\",\n      \"url\": \"/docs/concepts#the-agentic-search-loop\",\n      \"group\": \"Overview\",\n      \"snippet\": \"When the reader presses Enter, the query goes to a bounded loop…\"\n    }\n  ],\n  \"query\": \"how does agentic search work\",\n  \"model\": \"claude-haiku-4-5\",\n  \"mode\": \"keyword\"\n}","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"200","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"results","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"Result[]","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"title","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"heading?","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"url","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"group?","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"snippet","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"query","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"string","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"model","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"mode","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"'keyword'","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"warning","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"string?","chunkId":"api/endpoint#keyword-response-json"},{"kind":"code","literal":"#anchor","chunkId":"api/endpoint#keyword-response-json"}]
sources: [{"chunkId":"api/endpoint#keyword-response-json","url":"/docs/api/endpoint#keyword-response-json","anchor":"keyword-response-json"}]
---

Keyword mode returns a successful JSON envelope whose fields are documented in a table: ranked results with title, optional heading, URL, group, and snippet; the echoed query; the configured loop model; the mode that ran; and an optional warning when agentic was requested with no key. The result URL carries the deep link with anchor appended except for a document's intro chunk.
