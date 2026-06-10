---
id: "api/cli#claude-code-skill"
title: "CLI"
heading: "Claude Code skill"
group: "API"
order: 2
url: "/docs/api/cli#claude-code-skill"
anchor: "claude-code-skill"
terms: ["claude","code","skill","bundled","build","digest","builds","tree","without","anthropicapikey","runs","deterministic","producer","seam","sharded","corpus","shards","input","json","manifest","fresh","context","distillat","distillation","shard","distill","synthesis","pass","notes","global","assemble","anthropic","because","model","steps","inside","subscription","costs","tokens","size"]
hash: "443db6e48ac36d708dfd7e6b19f5e65b8f34ee4eb0675c8da98bd86cacfecc96"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask digest corpus --shards-dir .hev-ask/shards    -\u003e input-\u003cid\u003e.json + manifest.json\n...one fresh-context distillation per shard       -\u003e distill-\u003cid\u003e.json...\n...one synthesis pass over the shard notes        -\u003e global.json...\nask digest assemble --input-dir .hev-ask/shards   -\u003e the .hev-ask/ tree","chunkId":"api/cli#claude-code-skill"},{"kind":"code","literal":"build-digest","chunkId":"api/cli#claude-code-skill"},{"kind":"code","literal":"ANTHROPIC_API_KEY","chunkId":"api/cli#claude-code-skill"},{"kind":"code","literal":"ask digest build","chunkId":"api/cli#claude-code-skill"}]
sources: [{"chunkId":"api/cli#claude-code-skill","url":"/docs/api/cli#claude-code-skill","anchor":"claude-code-skill"}]
---

Claude Code skill The bundled build-digest skill builds the tree without using your ANTHROPICAPIKEY. It runs the deterministic producer seam, sharded: ask digest corpus --shards-dir .hev-ask/shards -> input- .json + manifest.json ...one fresh-context distillat...
