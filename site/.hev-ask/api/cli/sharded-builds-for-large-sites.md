---
id: "api/cli#sharded-builds-for-large-sites"
title: "CLI"
heading: "Sharded builds for large sites"
group: "API"
order: 8
url: "/docs/api/cli#sharded-builds-for-large-sites"
anchor: "sharded-builds-for-large-sites"
terms: ["sharded","builds","large","sites","even","incremental","rebuilds","first","build","distils","whole","corpus","flow","splits","along","slug","prefix","boundaries","shard","fresh","context","assembles","merged","tree","removing","size","bound","sharding","stable","because","identity","comes","carries","content","hash","editing","pends","only","stale","distillations"]
hash: "edb0ce1e92994c33baa1814a21cdfb79ae95bdf7c3ffc385218dcd0f79c2d35e"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask digest corpus --shards-dir .hev-ask/shards   # input-\u003cid\u003e.json per shard + manifest.json\nask digest status --shards-dir .hev-ask/shards   # distilled / pending / stale, per shard\n# ...one distillation per shard writes distill-\u003cid\u003e.json; a final pass writes global.json...\nask digest assemble --input-dir .hev-ask/shards  # merge + write the .hev-ask/ tree","chunkId":"api/cli#sharded-builds-for-large-sites"},{"kind":"code","literal":"workers/...","chunkId":"api/cli#sharded-builds-for-large-sites"},{"kind":"code","literal":"pages/...","chunkId":"api/cli#sharded-builds-for-large-sites"},{"kind":"code","literal":"corpus","chunkId":"api/cli#sharded-builds-for-large-sites"},{"kind":"code","literal":"ask digest verify","chunkId":"api/cli#sharded-builds-for-large-sites"},{"kind":"code","literal":"--skip-build","chunkId":"api/cli#sharded-builds-for-large-sites"},{"kind":"code","literal":"_meta.md","chunkId":"api/cli#sharded-builds-for-large-sites"},{"kind":"code","literal":"--strict","chunkId":"api/cli#sharded-builds-for-large-sites"}]
sources: [{"chunkId":"api/cli#sharded-builds-for-large-sites","url":"/docs/api/cli#sharded-builds-for-large-sites","anchor":"sharded-builds-for-large-sites"}]
---

Even with incremental rebuilds a first build distils the whole corpus, so the sharded flow splits it along slug-prefix boundaries, distils each shard in a fresh context, and assembles a merged tree, removing the size bound. Sharding is stable and incremental because shard identity comes from the slug prefix and each carries a content hash, so editing one doc re-pends only its shard; stale distillations are detected and skipped with sections falling back to plain excerpts, and verify gates anchors, coverage, fidelity, and tree integrity.
