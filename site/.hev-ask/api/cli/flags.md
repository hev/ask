---
id: "api/cli#flags"
title: "CLI"
heading: "Flags"
group: "API"
order: 4
url: "/docs/api/cli#flags"
anchor: "flags"
terms: ["flags","reference","table","covering","local","digest","directory","remote","endpoint","base","json","output","tree","depth","result","caps","collections","path","chunk","heading","model","provider","override","sharded","mode","directories","shard","size","assemble","input","verify","build","command","skip","strict","options","global","precede","grep","openapi"]
hash: "5d1a8038354ca7c57c1f2e3fbe3b99e9d22881cbbfd064e764ef233dc3e098d3"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask --digest-dir .hev-ask --json grep \"openapi\"\nask --endpoint https://hevask.com/api/ask cat api/endpoint\nask digest build --collection docs --collection guides --chunk-heading-depth 2\nask digest verify --skip-build","chunkId":"api/cli#flags"},{"kind":"code","literal":"--digest-dir \u003cdir\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":".hev-ask","chunkId":"api/cli#flags"},{"kind":"code","literal":"--endpoint \u003curl\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"/api/ask","chunkId":"api/cli#flags"},{"kind":"code","literal":"answer","chunkId":"api/cli#flags"},{"kind":"code","literal":"--json","chunkId":"api/cli#flags"},{"kind":"code","literal":"--depth \u003cn\\|all\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"tree","chunkId":"api/cli#flags"},{"kind":"code","literal":"(+N)","chunkId":"api/cli#flags"},{"kind":"code","literal":"all","chunkId":"api/cli#flags"},{"kind":"code","literal":"--max-results \u003cn\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"grep","chunkId":"api/cli#flags"},{"kind":"code","literal":"--collection \u003cname\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"docs","chunkId":"api/cli#flags"},{"kind":"code","literal":"--base-path \u003cpath\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"/docs/","chunkId":"api/cli#flags"},{"kind":"code","literal":"--content-glob \u003cglob\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"--chunk-heading-depth \u003cn\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"--digest-model \u003cmodel\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"ask digest build","chunkId":"api/cli#flags"},{"kind":"code","literal":"claude-opus-4-8","chunkId":"api/cli#flags"},{"kind":"code","literal":"--provider \u003cname\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"anthropic","chunkId":"api/cli#flags"}]
sources: [{"chunkId":"api/cli#flags","url":"/docs/api/cli#flags","anchor":"flags"}]
---

A reference table of CLI flags covering the local digest directory, the remote endpoint base URL, JSON output, tree depth, result caps, collections and base path, chunk heading depth, the digest model and provider with its base-URL override, sharded-mode directories and shard size, the assemble input directory, and verify's build-command, skip-build, and strict options; global flags precede the command.
