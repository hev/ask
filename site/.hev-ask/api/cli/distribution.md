---
id: "api/cli#distribution"
title: "CLI"
heading: "Distribution"
group: "API"
order: 3
url: "/docs/api/cli#distribution"
anchor: "distribution"
terms: ["distribution","package","exposes","single","whose","launcher","resolves","environment","variable","override","first","platform","specific","optional","binary","checked","source","development","fallback","published","installs","packaged","hevaskbinary","installed","monorepo","path"]
hash: "6ca69ef59f845bf9734dc9ec208f3fa63fdd9272703d454891bdbebb56fcda55"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask","chunkId":"api/cli#distribution"},{"kind":"code","literal":"HEV_ASK_BINARY","chunkId":"api/cli#distribution"}]
sources: [{"chunkId":"api/cli#distribution","url":"/docs/api/cli#distribution","anchor":"distribution"}]
---

The npm package exposes a single bin whose launcher resolves an environment-variable override first, then a platform-specific optional binary package, then the checked-out Go source as a development fallback. Published installs use the packaged binary.
