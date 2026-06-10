---
id: "api/cli#where-it-runs"
title: "CLI"
heading: "Where it runs"
group: "API"
order: 9
url: "/docs/api/cli#where-it-runs"
anchor: "where-it-runs"
terms: ["runs","producer","commands","locally","filesystem","access","astro","integration","also","invokes","digest","build","during","anthropicapikey","present","falls","back","committed","tree","command","cannot","anthropic","virtual","verify","deployed","site","reads","through","does","need","anchor","correctness","contract","makes","deep","links","trustworthy","running","every","mechanical"]
hash: "3160b534adda5d403ca544d7fa24738f6edc01c050b30b719a2cc04be4a6eede"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask digest build","chunkId":"api/cli#where-it-runs"},{"kind":"code","literal":"astro build","chunkId":"api/cli#where-it-runs"},{"kind":"code","literal":"ANTHROPIC_API_KEY","chunkId":"api/cli#where-it-runs"},{"kind":"code","literal":"virtual:hev-ask/digest","chunkId":"api/cli#where-it-runs"},{"kind":"code","literal":"ask digest verify","chunkId":"api/cli#where-it-runs"}]
sources: [{"chunkId":"api/cli#where-it-runs","url":"/docs/api/cli#where-it-runs","anchor":"where-it-runs"}]
---

Where it runs The producer commands run locally or in CI with filesystem access. The Astro integration also invokes ask digest build during astro build when ANTHROPICAPIKEY is present, then falls back to the committed tree if the build command cannot run....
