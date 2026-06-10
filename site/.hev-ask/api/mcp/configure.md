---
id: "api/mcp#configure"
title: "MCP server"
heading: "Configure"
group: "API"
order: 36
url: "/docs/api/mcp#configure"
anchor: "configure"
terms: ["configure","checked","repo","keyless","reads","local","tree","mcpservers","docs","command","args","digest","deployed","site","including","other","runs","pulled","agent","workspac","hevask","endpoint","https","workspace","anywhere"]
hash: "bd48710aaa8e6b45422ed8cee00e44b2004fdbd68d1178fac60bae41e706b728"
mode: "source-primary"
facts: [{"kind":"code","literal":"{\n  \"mcpServers\": {\n    \"docs\": {\n      \"command\": \"ask\",\n      \"args\": [\"--digest-dir\", \".hev-ask\", \"mcp\"]\n    }\n  }\n}","chunkId":"api/mcp#configure"},{"kind":"code","literal":"{\n  \"mcpServers\": {\n    \"hevask\": {\n      \"command\": \"ask\",\n      \"args\": [\"--endpoint\", \"https://hevask.com/api/ask\", \"mcp\"]\n    }\n  }\n}","chunkId":"api/mcp#configure"}]
sources: [{"chunkId":"api/mcp#configure","url":"/docs/api/mcp#configure","anchor":"configure"}]
---

Configure For a checked-out repo (keyless, reads the local tree): { "mcpServers": { "docs": { "command": "ask", "args": ["--digest-dir", ".hev-ask", "mcp"] } } } For a deployed site — including any other site that runs hev ask, pulled into the agent's workspac...
