---
id: "tradeoffs#one-dependency-deliberately"
title: "Tradeoffs"
heading: "One dependency, deliberately"
group: "Overview"
order: 105
url: "/docs/tradeoffs#one-dependency-deliberately"
anchor: "one-dependency-deliberately"
terms: ["dependency","deliberately","aims","near","zero","deliberate","exception","heading","slugger","tiny","pure","edge","safe","library","generating","anchors","hand","risks","drifting","renderer","shipping","link","404s","page","same","uses","guarantees","byte","identical","framework","adapters","extending","guarantee","their","slug","rules","taken","purpose","github","close"]
hash: "ab20e6d42c28963543737ab48f6363737ae3acca033cbf40114668f79b31f2c6"
mode: "agent-primary"
facts: [{"kind":"code","literal":"github-slugger","chunkId":"tradeoffs#one-dependency-deliberately"},{"kind":"value","literal":"github.com","chunkId":"tradeoffs#one-dependency-deliberately"}]
sources: [{"chunkId":"tradeoffs#one-dependency-deliberately","url":"/docs/tradeoffs#one-dependency-deliberately","anchor":"one-dependency-deliberately"}]
---

hev ask aims to be near zero-dependency with one deliberate exception: the heading slugger, a tiny pure-JS edge-safe library. Generating anchors by hand risks drifting from the renderer and shipping a link that 404s to the top of the page, so using the same slugger the renderer uses guarantees byte-identical anchors, with per-framework adapters extending that guarantee to their slug rules; the dependency was taken on purpose.
