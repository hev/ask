---
id: "api/cli#flags"
title: "CLI"
heading: "Flags"
group: "API"
order: 4
url: "/docs/api/cli#flags"
anchor: "flags"
terms: ["flags","flag","default","description","digest","local","tree","reads","output","producer","commands","endpoint","remote","base","read","answer","json","grep","openapi","https","hevask","build","collection","docs","guides","chunk","heading","depth","verify","skip","results","name","path","content","glob","model","claude","opus","shards","corpus"]
hash: "946e75684475cf281e7801dcb54f8b8c36639354c5251f37809bccc2f276cfb0"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask --digest-dir .hev-ask --json grep \"openapi\"\nask --endpoint https://hevask.com/api/ask cat api/endpoint\nask digest build --collection docs --collection guides --chunk-heading-depth 2\nask digest verify --skip-build","chunkId":"api/cli#flags"},{"kind":"code","literal":"--digest-dir \u003cdir\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":".hev-ask","chunkId":"api/cli#flags"},{"kind":"code","literal":"--endpoint \u003curl\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"/api/ask","chunkId":"api/cli#flags"},{"kind":"code","literal":"answer","chunkId":"api/cli#flags"},{"kind":"code","literal":"--json","chunkId":"api/cli#flags"},{"kind":"code","literal":"--max-results \u003cn\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"grep","chunkId":"api/cli#flags"},{"kind":"code","literal":"--collection \u003cname\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"docs","chunkId":"api/cli#flags"},{"kind":"code","literal":"--base-path \u003cpath\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"/docs/","chunkId":"api/cli#flags"},{"kind":"code","literal":"--content-glob \u003cglob\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"--chunk-heading-depth \u003cn\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"--digest-model \u003cmodel\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"claude-opus-4-8","chunkId":"api/cli#flags"},{"kind":"code","literal":"ask digest build","chunkId":"api/cli#flags"},{"kind":"code","literal":"--shards-dir \u003cdir\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"corpus","chunkId":"api/cli#flags"},{"kind":"code","literal":"status","chunkId":"api/cli#flags"},{"kind":"code","literal":"manifest.json","chunkId":"api/cli#flags"},{"kind":"code","literal":"--shard-bytes \u003cn\u003e","chunkId":"api/cli#flags"},{"kind":"code","literal":"200000","chunkId":"api/cli#flags"}]
sources: [{"chunkId":"api/cli#flags","url":"/docs/api/cli#flags","anchor":"flags"}]
---

Flags Flag Default Description --digest-dir .hev-ask Local digest tree for reads; output dir for producer commands. --endpoint - Remote /api/ask base URL for read commands and answer....
