---
id: "api/mcp#data-sources"
title: "MCP server"
heading: "Data sources"
group: "API"
order: 37
url: "/docs/api/mcp#data-sources"
anchor: "data-sources"
terms: ["data","sources","uses","same","resolution","endpoint","downloads","deployed","tree","compressed","archive","otherwise","reads","digest","disk","defaulting","fetch","docs","just","rebuilt","visible","next","fetchdocs","without","restarting","server"]
hash: "863a2df687abe731858ab1157d344474fb96e257b5e52ec7266642179186a4dc"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask mcp","chunkId":"api/mcp#data-sources"},{"kind":"code","literal":"--endpoint \u003curl\u003e","chunkId":"api/mcp#data-sources"},{"kind":"code","literal":"/api/ask/archive","chunkId":"api/mcp#data-sources"},{"kind":"code","literal":"--digest-dir","chunkId":"api/mcp#data-sources"},{"kind":"code","literal":".hev-ask","chunkId":"api/mcp#data-sources"},{"kind":"code","literal":"fetch_docs","chunkId":"api/mcp#data-sources"}]
sources: [{"chunkId":"api/mcp#data-sources","url":"/docs/api/mcp#data-sources","anchor":"data-sources"}]
---

Data sources ask mcp uses the same resolution as the CLI: --endpoint downloads the deployed tree as a compressed archive from /api/ask/archive. Otherwise it reads --digest-dir from disk, defaulting to .hev-ask....
