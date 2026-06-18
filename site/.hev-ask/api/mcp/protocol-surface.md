---
id: "api/mcp#protocol-surface"
title: "MCP server"
heading: "Protocol surface"
group: "API"
order: 39
url: "/docs/api/mcp#protocol-surface"
anchor: "protocol-surface"
terms: ["protocol","surface","server","speaks","newline","delimited","json","stdio","handling","initialization","returning","instructions","tool","listing","calls","initialized","notification","unknown","methods","error","failures","result","substantive","behavior","lives","shared","package","keeping","small","initialize","tools","list","call","iserror","true","handles","plus","return","level","keeps"]
hash: "e99d0aabc9e51f725c1d5a66babc38fafa6631b456ce3eeb9647327596ca44ce"
mode: "source-primary"
facts: [{"kind":"code","literal":"initialize","chunkId":"api/mcp#protocol-surface"},{"kind":"code","literal":"instructions","chunkId":"api/mcp#protocol-surface"},{"kind":"code","literal":"tools/list","chunkId":"api/mcp#protocol-surface"},{"kind":"code","literal":"tools/call","chunkId":"api/mcp#protocol-surface"},{"kind":"code","literal":"isError: true","chunkId":"api/mcp#protocol-surface"},{"kind":"code","literal":"pkg/ask","chunkId":"api/mcp#protocol-surface"}]
sources: [{"chunkId":"api/mcp#protocol-surface","url":"/docs/api/mcp#protocol-surface","anchor":"protocol-surface"}]
---

The server speaks newline-delimited JSON-RPC over stdio, handling initialization (returning the instructions), tool listing, tool calls, and the initialized notification, with unknown methods returning a protocol error and tool failures returning an error result. All substantive behavior lives in the shared Go package, keeping the server small.
