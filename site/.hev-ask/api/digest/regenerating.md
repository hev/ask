---
id: "api/digest#regenerating"
title: "Digest format"
heading: "Regenerating"
group: "API"
order: 23
url: "/docs/api/digest#regenerating"
anchor: "regenerating"
terms: ["regenerating","rebuild","after","content","changes","commit","result","build","incremental","spending","model","work","only","sections","whose","hash","changed","keyless","skill","call","integration","during","site","present","large","sites","sharded","flow","verify","gates","anchors","coverage","legacy","single","file","digest","exploded","tree","astro","json"]
hash: "18c55bbc27ee8afb9bdda4d22407ee10967f514cfa90fb49fb6506dc60a2a138"
mode: "source-primary"
facts: [{"kind":"code","literal":"hash","chunkId":"api/digest#regenerating"},{"kind":"code","literal":"ask digest build","chunkId":"api/digest#regenerating"},{"kind":"code","literal":"astro build","chunkId":"api/digest#regenerating"},{"kind":"code","literal":"ask digest verify","chunkId":"api/digest#regenerating"},{"kind":"code","literal":"digest.json","chunkId":"api/digest#regenerating"},{"kind":"code","literal":"ask digest migrate","chunkId":"api/digest#regenerating"}]
sources: [{"chunkId":"api/digest#regenerating","url":"/docs/api/digest#regenerating","anchor":"regenerating"}]
---

Rebuild after content changes and commit the result; the build is incremental, spending model work only on sections whose hash changed. You can build with the keyless skill, the one-call CLI build, or the integration during the site build when a key is present; large sites use the sharded flow, verify gates anchors and coverage, and a legacy single-file digest can be exploded into the tree with no model call.
