---
id: "api/mcp#configure"
title: "MCP server"
heading: "Configure"
group: "API"
order: 37
url: "/docs/api/mcp#configure"
anchor: "configure"
terms: ["configure","example","server","configurations","cases","checked","repo","reads","local","tree","keylessly","deployed","site","including","other","running","pulled","agent","workspace","endpoint","mcpservers","docs","command","args","digest","hevask","https","keyless","runs","anywhere"]
hash: "bd48710aaa8e6b45422ed8cee00e44b2004fdbd68d1178fac60bae41e706b728"
mode: "source-primary"
facts: [{"kind":"code","literal":"{\n  \"mcpServers\": {\n    \"docs\": {\n      \"command\": \"ask\",\n      \"args\": [\"--digest-dir\", \".hev-ask\", \"mcp\"]\n    }\n  }\n}","chunkId":"api/mcp#configure"},{"kind":"code","literal":"{\n  \"mcpServers\": {\n    \"hevask\": {\n      \"command\": \"ask\",\n      \"args\": [\"--endpoint\", \"https://hevask.com/api/ask\", \"mcp\"]\n    }\n  }\n}","chunkId":"api/mcp#configure"}]
sources: [{"chunkId":"api/mcp#configure","url":"/docs/api/mcp#configure","anchor":"configure"}]
---

Example MCP server configurations for two cases: a checked-out repo that reads the local tree keylessly, and a deployed site (including any other site running hev ask) pulled into the agent's workspace via its endpoint.
