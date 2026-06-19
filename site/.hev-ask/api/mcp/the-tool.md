---
id: "api/mcp#the-tool"
title: "MCP server"
heading: "The tool"
group: "API"
order: 41
url: "/docs/api/mcp#the-tool"
anchor: "the-tool"
terms: ["tool","single","materializes","digest","tree","host","keyed","local","cache","path","returns","title","inline","plus","disk","bodies","facts","call","bootstraps","whole","disclosure","ladder","second","round","trip","agent","uses","native","tools","force","argument","pulls","unconditionally","otherwise","compares","remote","content","hash","downloads","only"]
hash: "54bc7be8d62d5112772761d161b203c851308fffd7751b1c920e8c99e1975b96"
mode: "source-primary"
facts: [{"kind":"code","literal":"tree ~/.cache/hev-ask/hevask.com        # already returned inline by fetch_docs\ncat  ~/.cache/hev-ask/hevask.com/overview/quick-start.md\ngrep -r \"prerender\" ~/.cache/hev-ask/hevask.com","chunkId":"api/mcp#the-tool"},{"kind":"code","literal":"fetch_docs","chunkId":"api/mcp#the-tool"},{"kind":"code","literal":"{ force?: boolean }","chunkId":"api/mcp#the-tool"},{"kind":"code","literal":"{ path, contentHash, sections, tree, upToDate }","chunkId":"api/mcp#the-tool"},{"kind":"code","literal":"~/.cache/hev-ask/hevask.com/","chunkId":"api/mcp#the-tool"},{"kind":"code","literal":"force: true","chunkId":"api/mcp#the-tool"},{"kind":"code","literal":"contentHash","chunkId":"api/mcp#the-tool"},{"kind":"code","literal":"upToDate: true","chunkId":"api/mcp#the-tool"},{"kind":"value","literal":"e.g","chunkId":"api/mcp#the-tool"}]
sources: [{"chunkId":"api/mcp#the-tool","url":"/docs/api/mcp#the-tool","anchor":"the-tool"}]
---

The single tool materializes the digest tree at a host-keyed local cache path and returns the title-tree inline plus on-disk bodies and facts, so one call bootstraps the whole disclosure ladder with no second round-trip and the agent then uses its native tools. A force argument re-pulls unconditionally, otherwise it compares the remote content hash and re-downloads only on mismatch; because the corpus is bounded the whole tree ships compressed in one shot with no per-file delta protocol.
