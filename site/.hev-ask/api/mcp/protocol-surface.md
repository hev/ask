---
id: "api/mcp#protocol-surface"
title: "MCP server"
heading: "Protocol surface"
group: "API"
order: 39
url: "/docs/api/mcp#protocol-surface"
anchor: "protocol-surface"
terms: ["protocol","surface","server","speaks","newline","delimited","json","stdio","handles","initialization","returning","instructions","tool","listing","calls","plus","initialized","notification","unknown","methods","return","error","failures","result","keeping","thin","because","substantive","behavior","lives","shared","core","command","group","reuse","initialize","tools","list","call","iserror"]
hash: "e99d0aabc9e51f725c1d5a66babc38fafa6631b456ce3eeb9647327596ca44ce"
mode: "source-primary"
facts: [{"kind":"code","literal":"initialize","chunkId":"api/mcp#protocol-surface"},{"kind":"code","literal":"instructions","chunkId":"api/mcp#protocol-surface"},{"kind":"code","literal":"tools/list","chunkId":"api/mcp#protocol-surface"},{"kind":"code","literal":"tools/call","chunkId":"api/mcp#protocol-surface"},{"kind":"code","literal":"isError: true","chunkId":"api/mcp#protocol-surface"},{"kind":"code","literal":"pkg/ask","chunkId":"api/mcp#protocol-surface"}]
sources: [{"chunkId":"api/mcp#protocol-surface","url":"/docs/api/mcp#protocol-surface","anchor":"protocol-surface"}]
---

The server speaks newline-delimited JSON-RPC over stdio and handles initialization (returning the instructions), tool listing, and tool calls, plus the initialized notification. Unknown methods return a protocol error and tool failures return an error result, keeping the server thin because all substantive behavior lives in the shared core that the CLI, command group, and MCP server reuse.
