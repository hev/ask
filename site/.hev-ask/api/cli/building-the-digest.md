---
id: "api/cli#building-the-digest"
title: "CLI"
heading: "Building the digest"
group: "API"
order: 1
url: "/docs/api/cli#building-the-digest"
anchor: "building-the-digest"
terms: ["building","digest","producer","commands","site","root","build","emit","keyless","corpus","assemble","tree","verify","report","shard","coverage","migrate","legacy","json","only","calls","model","needing","incremental","hash","gated","clean","zero","times","shot","bounded","section","text","size","before","requires","sharded","flow","authors","context"]
hash: "c8c9b3c2808376bf87bc7128a3a9a84bdde6aa70ac2011e1ccb1debae2976f33"
mode: "source-primary"
facts: [{"kind":"code","literal":"export ANTHROPIC_API_KEY=sk-ant-...\nask digest build                    # claude-opus-4-8 by default","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"export OPENAI_API_KEY=sk-...\nask digest build --provider openai  # gpt-5.1 by default","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"export OPENROUTER_API_KEY=sk-or-...\nask digest build --provider openrouter   # anthropic/claude-opus-4.8 by default","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"ask digest build","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"--provider","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"ask digest corpus","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"ask digest assemble","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":".hev-ask/","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"context","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"summaries","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"suggestions","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":".hev-ask/digest.json","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"ask digest migrate","chunkId":"api/cli#building-the-digest"},{"kind":"value","literal":"digest.json","chunkId":"api/cli#building-the-digest"}]
sources: [{"chunkId":"api/cli#building-the-digest","url":"/docs/api/cli#building-the-digest","anchor":"building-the-digest"}]
---

Producer commands run from the site root to build, emit a keyless corpus, assemble the tree, verify, report shard coverage, and migrate a legacy JSON digest. Only the build calls a model and is the only one needing a key; it is incremental and hash-gated, so a clean tree calls the model zero times, and a one-shot build is bounded by section-text size before it requires the sharded flow. The model authors only context, glossary, summaries, and suggestions; everything else is computed deterministically.
