---
id: "api/mcp#data-sources"
title: "MCP server"
heading: "Data sources"
group: "API"
order: 38
url: "/docs/api/mcp#data-sources"
anchor: "data-sources"
terms: ["data","sources","server","uses","same","source","resolution","endpoint","downloads","deployed","tree","compressed","archive","otherwise","reads","local","digest","directory","just","rebuilt","becomes","visible","next","fetch","without","restarting","docs","disk","defaulting","fetchdocs"]
hash: "863a2df687abe731858ab1157d344474fb96e257b5e52ec7266642179186a4dc"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask mcp","chunkId":"api/mcp#data-sources"},{"kind":"code","literal":"--endpoint \u003curl\u003e","chunkId":"api/mcp#data-sources"},{"kind":"code","literal":"/api/ask/archive","chunkId":"api/mcp#data-sources"},{"kind":"code","literal":"--digest-dir","chunkId":"api/mcp#data-sources"},{"kind":"code","literal":".hev-ask","chunkId":"api/mcp#data-sources"},{"kind":"code","literal":"fetch_docs","chunkId":"api/mcp#data-sources"}]
sources: [{"chunkId":"api/mcp#data-sources","url":"/docs/api/mcp#data-sources","anchor":"data-sources"}]
---

The MCP server uses the same source resolution as the CLI: an endpoint downloads the deployed tree as a compressed archive, otherwise it reads the local digest directory; a just-rebuilt tree becomes visible on the next fetch without restarting the server.
