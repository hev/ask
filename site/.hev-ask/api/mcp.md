---
id: "api/mcp"
title: "MCP server"
heading: null
group: "API"
order: 35
url: "/docs/api/mcp"
anchor: null
terms: ["server","stdio","model","context","protocol","tool","plus","instructions","downloads","whole","digest","tree","local","disk","tell","agent","navigate","read","search","file","tools","cite","every","claim","section","anchor","hands","directory","rather","reimplementing","because","consumer","already","needs","corpus","synthesis","there","deliberately","answer","grep"]
hash: "ff242827eef613609690f7b116ad40577a561e412599b78d9afd465f7436aaac"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask mcp","chunkId":"api/mcp"},{"kind":"code","literal":"tree","chunkId":"api/mcp"},{"kind":"code","literal":"cat","chunkId":"api/mcp"},{"kind":"code","literal":"grep","chunkId":"api/mcp"},{"kind":"code","literal":"url","chunkId":"api/mcp"},{"kind":"code","literal":"anchor","chunkId":"api/mcp"},{"kind":"code","literal":"ls","chunkId":"api/mcp"},{"kind":"code","literal":"answer","chunkId":"api/mcp"},{"kind":"value","literal":"Callout.astro","chunkId":"api/mcp"}]
sources: [{"chunkId":"api/mcp","url":"/docs/api/mcp","anchor":null}]
---

The MCP server is a stdio Model Context Protocol server that is one tool plus instructions: the tool downloads the whole digest tree to local disk and the instructions tell the agent to navigate, read, and search it with its own file tools and to cite every claim with the section's URL and anchor. It hands the agent a directory rather than reimplementing a search API, because an MCP consumer is already an agent that needs the corpus, not synthesis, so there is deliberately no answer tool.
