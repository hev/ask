---
id: "tradeoffs#one-dependency-deliberately"
title: "Tradeoffs"
heading: "One dependency, deliberately"
group: "Overview"
order: 105
url: "/docs/tradeoffs#one-dependency-deliberately"
anchor: "one-dependency-deliberately"
terms: ["dependency","deliberately","aims","near","zero","deliberate","exception","github","slugger","tiny","pure","edge","safe","library","because","generating","heading","anchors","hand","risks","drifting","renderer","shipping","link","fails","page","same","astro","guarantees","byte","identical","framework","adapters","extend","guarantee","their","slug","rules","close","404s"]
hash: "ab20e6d42c28963543737ab48f6363737ae3acca033cbf40114668f79b31f2c6"
mode: "agent-primary"
facts: [{"kind":"code","literal":"github-slugger","chunkId":"tradeoffs#one-dependency-deliberately"},{"kind":"value","literal":"github.com","chunkId":"tradeoffs#one-dependency-deliberately"}]
sources: [{"chunkId":"tradeoffs#one-dependency-deliberately","url":"/docs/tradeoffs#one-dependency-deliberately","anchor":"one-dependency-deliberately"}]
---

hev ask aims to be near zero-dependency with one deliberate exception, github-slugger (a tiny, pure-JS, edge-safe library), because generating heading anchors by hand risks drifting from the renderer and shipping a link that fails to the top of the page; using the same slugger Astro and GitHub use guarantees byte-identical anchors and per-framework adapters extend the same guarantee to their slug rules.
