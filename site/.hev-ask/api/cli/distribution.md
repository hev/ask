---
id: "api/cli#distribution"
title: "CLI"
heading: "Distribution"
group: "API"
order: 3
url: "/docs/api/cli#distribution"
anchor: "distribution"
terms: ["distribution","package","exposes","single","binary","whose","launcher","resolves","environment","override","first","platform","specific","optional","source","monorepo","development","fallback","published","installs","packaged","hevaskbinary","installed","checked","path"]
hash: "6ca69ef59f845bf9734dc9ec208f3fa63fdd9272703d454891bdbebb56fcda55"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask","chunkId":"api/cli#distribution"},{"kind":"code","literal":"HEV_ASK_BINARY","chunkId":"api/cli#distribution"}]
sources: [{"chunkId":"api/cli#distribution","url":"/docs/api/cli#distribution","anchor":"distribution"}]
---

The npm package exposes a single binary whose launcher resolves an environment override first, then a platform-specific optional binary package, then the Go source in the monorepo as a development fallback; published installs use the packaged binary.
