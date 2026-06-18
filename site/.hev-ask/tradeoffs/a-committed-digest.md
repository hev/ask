---
id: "tradeoffs#a-committed-digest"
title: "Tradeoffs"
heading: "A committed digest"
group: "Overview"
order: 100
url: "/docs/tradeoffs#a-committed-digest"
anchor: "a-committed-digest"
terms: ["committed","digest","generated","offline","markdown","tree","file","section","rather","computed","runtime","hidden","service","upside","reviewable","pull","requests","deterministic","free","read","model","call","request","path","bundled","edge","worker","filesystem","access","directly","navigable","agent","cost","stale","only","regenerates","content","changes","build","runs"]
hash: "110adfc68831d0e1a9dcd3844ea876fb7fbd211e79537b3c2e854618aaecdc8d"
mode: "agent-primary"
facts: [{"kind":"code","literal":"tree","chunkId":"tradeoffs#a-committed-digest"},{"kind":"code","literal":"cat","chunkId":"tradeoffs#a-committed-digest"},{"kind":"code","literal":"grep","chunkId":"tradeoffs#a-committed-digest"}]
sources: [{"chunkId":"tradeoffs#a-committed-digest","url":"/docs/tradeoffs#a-committed-digest","anchor":"a-committed-digest"}]
---

The digest is generated offline and committed to git as a markdown tree, one file per section, rather than computed at runtime or hidden in a service. The upside is that it is reviewable per section in pull requests, deterministic and free to read at runtime with no model call on the request path, bundled into the edge worker with no runtime filesystem access, and directly navigable by an agent; the cost is that it can go stale and only regenerates when content changes and a build runs, though the per-section hash gate makes regeneration cheap and incremental so rebuild-in-CI is the intended workflow.
