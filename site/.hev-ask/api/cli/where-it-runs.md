---
id: "api/cli#where-it-runs"
title: "CLI"
heading: "Where it runs"
group: "API"
order: 9
url: "/docs/api/cli#where-it-runs"
anchor: "where-it-runs"
terms: ["runs","producer","commands","locally","filesystem","access","astro","integration","also","build","during","site","present","otherwise","falls","back","committed","tree","deployed","reads","through","virtual","module","without","needing","running","verify","step","every","mechanical","check","generated","anchors","still","match","renderer","produces","digest","anthropic","invokes"]
hash: "3160b534adda5d403ca544d7fa24738f6edc01c050b30b719a2cc04be4a6eede"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask digest build","chunkId":"api/cli#where-it-runs"},{"kind":"code","literal":"astro build","chunkId":"api/cli#where-it-runs"},{"kind":"code","literal":"ANTHROPIC_API_KEY","chunkId":"api/cli#where-it-runs"},{"kind":"code","literal":"virtual:hev-ask/digest","chunkId":"api/cli#where-it-runs"},{"kind":"code","literal":"ask digest verify","chunkId":"api/cli#where-it-runs"}]
sources: [{"chunkId":"api/cli#where-it-runs","url":"/docs/api/cli#where-it-runs","anchor":"where-it-runs"}]
---

Producer commands run locally or in CI with filesystem access, the Astro integration also runs the build during the site build when a key is present and otherwise falls back to the committed tree, and the deployed site reads the committed tree through the virtual module without needing filesystem access. Running the verify step on every build is the mechanical check that generated anchors still match what the renderer produces.
