---
id: "api/mcp#co-location"
title: "MCP server"
heading: "Co-location"
group: "API"
order: 36
url: "/docs/api/mcp#co-location"
anchor: "co-location"
terms: ["location","hydrating","tree","disk","assumes","server","agent","file","tools","share","host","true","default","stdio","transport","intentionally","exposes","only","path","since","remote","cannot","read","cache","would","need","separate","resource","tool","fallback","rather","returning","useless","local","hydrate","stdin","stdout","because","useful"]
hash: "0dd5d436dbfe12945ae2a67bf0100156e9645f9c1e9a8afe293b1b83a7348757"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask mcp","chunkId":"api/mcp#co-location"}]
sources: [{"chunkId":"api/mcp#co-location","url":"/docs/api/mcp#co-location","anchor":"co-location"}]
---

Hydrating the tree to disk assumes the MCP server and the agent's file tools share a host, which is true for the default stdio transport; this server intentionally exposes only that path, since a remote transport where the agent cannot read the server's cache would need a separate resource or tool fallback rather than returning a useless local path.
