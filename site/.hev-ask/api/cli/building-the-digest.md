---
id: "api/cli#building-the-digest"
title: "CLI"
heading: "Building the digest"
group: "API"
order: 1
url: "/docs/api/cli#building-the-digest"
anchor: "building-the-digest"
terms: ["building","digest","producer","commands","site","root","only","build","command","calls","model","needs","provider","incremental","hash","gated","unchanged","sections","skipped","clean","tree","zero","times","shot","bounded","fails","past","size","threshold","instructions","shard","separate","emit","assemble","expose","deterministic","seam","authors","context","glossary"]
hash: "c8c9b3c2808376bf87bc7128a3a9a84bdde6aa70ac2011e1ccb1debae2976f33"
mode: "source-primary"
facts: [{"kind":"code","literal":"export ANTHROPIC_API_KEY=sk-ant-...\nask digest build                    # claude-opus-4-8 by default","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"export OPENAI_API_KEY=sk-...\nask digest build --provider openai  # gpt-5.1 by default","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"export OPENROUTER_API_KEY=sk-or-...\nask digest build --provider openrouter   # anthropic/claude-opus-4.8 by default","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"ask digest build","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"--provider","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"ask digest corpus","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"ask digest assemble","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":".hev-ask/","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"context","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"summaries","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"suggestions","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":".hev-ask/digest.json","chunkId":"api/cli#building-the-digest"},{"kind":"code","literal":"ask digest migrate","chunkId":"api/cli#building-the-digest"},{"kind":"value","literal":"digest.json","chunkId":"api/cli#building-the-digest"}]
sources: [{"chunkId":"api/cli#building-the-digest","url":"/docs/api/cli#building-the-digest","anchor":"building-the-digest"}]
---

Producer commands run from the site root; only the build command calls a model, so only it needs a provider key, and it is incremental and hash-gated so unchanged sections are skipped and a clean tree calls the model zero times. The one-shot build is bounded and fails past a size threshold with instructions to shard; separate emit and assemble commands expose the deterministic seam where the model authors only context, glossary, summaries, and suggestions while everything else is computed in code, and a migrate command explodes a legacy single-file digest into the tree with no model call.
