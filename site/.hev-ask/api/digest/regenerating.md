---
id: "api/digest#regenerating"
title: "Digest format"
heading: "Regenerating"
group: "API"
order: 23
url: "/docs/api/digest#regenerating"
anchor: "regenerating"
terms: ["regenerating","rebuild","after","content","changes","commit","result","build","incremental","only","spends","model","work","changed","sections","bundled","skill","command","call","integration","also","runs","during","site","present","large","sites","sharded","flow","verify","step","gates","anchors","integrity","migrate","explodes","legacy","single","file","digest"]
hash: "18c55bbc27ee8afb9bdda4d22407ee10967f514cfa90fb49fb6506dc60a2a138"
mode: "source-primary"
facts: [{"kind":"code","literal":"hash","chunkId":"api/digest#regenerating"},{"kind":"code","literal":"ask digest build","chunkId":"api/digest#regenerating"},{"kind":"code","literal":"astro build","chunkId":"api/digest#regenerating"},{"kind":"code","literal":"ask digest verify","chunkId":"api/digest#regenerating"},{"kind":"code","literal":"digest.json","chunkId":"api/digest#regenerating"},{"kind":"code","literal":"ask digest migrate","chunkId":"api/digest#regenerating"}]
sources: [{"chunkId":"api/digest#regenerating","url":"/docs/api/digest#regenerating","anchor":"regenerating"}]
---

Rebuild after content changes and commit the result; the build is incremental so it only spends model work on changed sections. You can build it with the bundled skill (no key) or the build command (one model call), the integration also runs it during the site build when a key is present, large sites use the sharded flow, the verify step gates anchors and integrity, and a migrate command explodes a legacy single-file digest into the tree with no model call.
