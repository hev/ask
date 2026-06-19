---
id: "api/cli#flags"
title: "CLI"
heading: "Flags"
group: "API"
order: 4
url: "/docs/api/cli#flags"
anchor: "flags"
terms: ["flags","reference","table","every","flag","default","purpose","covering","digest","paths","remote","endpoint","reads","json","output","tree","depth","result","caps","collection","chunking","controls","model","provider","selection","sharding","inputs","verify","options","global","come","before","subcommand","grep","openapi","https","hevask","build","docs","guides"]
hash: "5d1a8038354ca7c57c1f2e3fbe3b99e9d22881cbbfd064e764ef233dc3e098d3"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask --digest-dir .hev-ask --json grep \"openapi\"\nask --endpoint https://hevask.com/api/ask cat api/endpoint\nask digest build --collection docs --collection guides --chunk-heading-depth 2\nask digest verify --skip-build","chunkId":"api/cli#flags"},{"kind":"code","literal":"--digest-dir \u003cdir\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":".hev-ask","chunkId":"api/cli#flags"},{"kind":"code","literal":"--endpoint \u003curl\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"/api/ask","chunkId":"api/cli#flags"},{"kind":"code","literal":"answer","chunkId":"api/cli#flags"},{"kind":"code","literal":"--json","chunkId":"api/cli#flags"},{"kind":"code","literal":"--depth \u003cn\\|all\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"tree","chunkId":"api/cli#flags"},{"kind":"code","literal":"(+N)","chunkId":"api/cli#flags"},{"kind":"code","literal":"all","chunkId":"api/cli#flags"},{"kind":"code","literal":"--max-results \u003cn\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"grep","chunkId":"api/cli#flags"},{"kind":"code","literal":"--collection \u003cname\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"docs","chunkId":"api/cli#flags"},{"kind":"code","literal":"--base-path \u003cpath\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"/docs/","chunkId":"api/cli#flags"},{"kind":"code","literal":"--content-glob \u003cglob\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"--chunk-heading-depth \u003cn\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"--digest-model \u003cmodel\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"ask digest build","chunkId":"api/cli#flags"},{"kind":"code","literal":"claude-opus-4-8","chunkId":"api/cli#flags"},{"kind":"code","literal":"--provider \u003cname\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"anthropic","chunkId":"api/cli#flags"}]
sources: [{"chunkId":"api/cli#flags","url":"/docs/api/cli#flags","anchor":"flags"}]
---

A reference table of every CLI flag with its default and purpose, covering digest paths, remote endpoint reads, JSON output, tree depth, result caps, collection and chunking controls, the digest model and provider selection, sharding inputs, and verify options. Global flags come before the subcommand.
