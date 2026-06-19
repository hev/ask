---
id: "api/mcp#configure"
title: "MCP server"
heading: "Configure"
group: "API"
order: 37
url: "/docs/api/mcp#configure"
anchor: "configure"
terms: ["configure","shows","client","configuration","cases","keyless","setup","pointing","local","checked","tree","deployed","site","endpoint","lets","running","pulled","agent","workspace","anywhere","mcpservers","docs","command","args","digest","hevask","https","repo","reads","including","other","runs"]
hash: "bd48710aaa8e6b45422ed8cee00e44b2004fdbd68d1178fac60bae41e706b728"
mode: "source-primary"
facts: [{"kind":"code","literal":"{\n  \"mcpServers\": {\n    \"docs\": {\n      \"command\": \"ask\",\n      \"args\": [\"--digest-dir\", \".hev-ask\", \"mcp\"]\n    }\n  }\n}","chunkId":"api/mcp#configure"},{"kind":"code","literal":"{\n  \"mcpServers\": {\n    \"hevask\": {\n      \"command\": \"ask\",\n      \"args\": [\"--endpoint\", \"https://hevask.com/api/ask\", \"mcp\"]\n    }\n  }\n}","chunkId":"api/mcp#configure"}]
sources: [{"chunkId":"api/mcp#configure","url":"/docs/api/mcp#configure","anchor":"configure"}]
---

Shows MCP client configuration for two cases: a keyless setup pointing at a local checked-out tree, and one pointing at a deployed site's endpoint, which lets any site running hev ask be pulled into the agent's workspace from anywhere.
