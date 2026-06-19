---
id: "api/cli#sharded-builds-for-large-sites"
title: "CLI"
heading: "Sharded builds for large sites"
group: "API"
order: 8
url: "/docs/api/cli#sharded-builds-for-large-sites"
anchor: "sharded-builds-for-large-sites"
terms: ["sharded","builds","large","sites","flow","removes","single","context","bound","splitting","corpus","along","slug","prefix","boundaries","distilling","shard","independently","merging","assembly","sharding","stable","incremental","editing","pends","only","owns","stale","distillations","detected","skipped","warning","while","affected","sections","fall","back","plain","excerpts","tree"]
hash: "edb0ce1e92994c33baa1814a21cdfb79ae95bdf7c3ffc385218dcd0f79c2d35e"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask digest corpus --shards-dir .hev-ask/shards   # input-\u003cid\u003e.json per shard + manifest.json\nask digest status --shards-dir .hev-ask/shards   # distilled / pending / stale, per shard\n# ...one distillation per shard writes distill-\u003cid\u003e.json; a final pass writes global.json...\nask digest assemble --input-dir .hev-ask/shards  # merge + write the .hev-ask/ tree","chunkId":"api/cli#sharded-builds-for-large-sites"},{"kind":"code","literal":"workers/...","chunkId":"api/cli#sharded-builds-for-large-sites"},{"kind":"code","literal":"pages/...","chunkId":"api/cli#sharded-builds-for-large-sites"},{"kind":"code","literal":"corpus","chunkId":"api/cli#sharded-builds-for-large-sites"},{"kind":"code","literal":"ask digest verify","chunkId":"api/cli#sharded-builds-for-large-sites"},{"kind":"code","literal":"--skip-build","chunkId":"api/cli#sharded-builds-for-large-sites"},{"kind":"code","literal":"_meta.md","chunkId":"api/cli#sharded-builds-for-large-sites"},{"kind":"code","literal":"--strict","chunkId":"api/cli#sharded-builds-for-large-sites"}]
sources: [{"chunkId":"api/cli#sharded-builds-for-large-sites","url":"/docs/api/cli#sharded-builds-for-large-sites","anchor":"sharded-builds-for-large-sites"}]
---

The sharded flow removes the single-context bound by splitting the corpus along slug-prefix boundaries, distilling each shard independently, and merging on assembly. Sharding is stable and incremental: editing one doc re-pends only the shard that owns it, stale distillations are detected and skipped with a warning while affected sections fall back to plain excerpts, and the tree stays usable throughout. Verify builds the site and checks rendered anchors, coverage, literal fidelity, and tree integrity, with anchor drift always fatal and the rest warnings unless strict mode is set.
