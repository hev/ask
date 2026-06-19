---
id: "api/cli#mcp"
title: "CLI"
heading: "MCP"
group: "API"
order: 6
url: "/docs/api/cli#mcp"
anchor: "mcp"
terms: ["command","runs","stdio","server","tool","downloads","whole","digest","tree","local","disk","returns","title","inline","after","agent","reads","files","tools","point","checked","repo","deployed","endpoint","configuration","mcpservers","hevask","args","https","grep","page","instructions","ships","archive","cache"]
hash: "e9ba5f690d5158797f3f370cbc8f7a816946b64c3da22ea8ece21bc4d9afba91"
mode: "source-primary"
facts: [{"kind":"code","literal":"{\n  \"mcpServers\": {\n    \"hevask\": {\n      \"command\": \"ask\",\n      \"args\": [\"--endpoint\", \"https://hevask.com/api/ask\", \"mcp\"]\n    }\n  }\n}","chunkId":"api/cli#mcp"},{"kind":"code","literal":"ask mcp","chunkId":"api/cli#mcp"},{"kind":"code","literal":"tree","chunkId":"api/cli#mcp"},{"kind":"code","literal":"cat","chunkId":"api/cli#mcp"},{"kind":"code","literal":"grep","chunkId":"api/cli#mcp"}]
sources: [{"chunkId":"api/cli#mcp","url":"/docs/api/cli#mcp","anchor":"mcp"}]
---

The CLI's MCP command runs a stdio server with one tool that downloads the whole digest tree to local disk and returns the title-tree inline, after which the agent reads the files with its own tools. It can point at a checked-out repo or a deployed endpoint via configuration.
