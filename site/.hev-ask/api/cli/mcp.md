---
id: "api/cli#mcp"
title: "CLI"
heading: "MCP"
group: "API"
order: 6
url: "/docs/api/cli#mcp"
anchor: "mcp"
terms: ["runs","stdio","server","tool","downloads","whole","digest","tree","local","disk","returns","title","inline","after","agent","reads","files","grep","mcpservers","hevask","command","args","endpoint","https","point","checked","repo","deployed","page","instructions","ships","archive","cache"]
hash: "e9ba5f690d5158797f3f370cbc8f7a816946b64c3da22ea8ece21bc4d9afba91"
mode: "source-primary"
facts: [{"kind":"code","literal":"{\n  \"mcpServers\": {\n    \"hevask\": {\n      \"command\": \"ask\",\n      \"args\": [\"--endpoint\", \"https://hevask.com/api/ask\", \"mcp\"]\n    }\n  }\n}","chunkId":"api/cli#mcp"},{"kind":"code","literal":"ask mcp","chunkId":"api/cli#mcp"},{"kind":"code","literal":"tree","chunkId":"api/cli#mcp"},{"kind":"code","literal":"cat","chunkId":"api/cli#mcp"},{"kind":"code","literal":"grep","chunkId":"api/cli#mcp"}]
sources: [{"chunkId":"api/cli#mcp","url":"/docs/api/cli#mcp","anchor":"mcp"}]
---

MCP ask mcp runs a stdio MCP server with one tool — it downloads the whole digest tree to local disk and returns the title-tree inline, after which the agent reads the files with its own tree/cat/grep....
