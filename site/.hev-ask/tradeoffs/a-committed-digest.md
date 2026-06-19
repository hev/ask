---
id: "tradeoffs#a-committed-digest"
title: "Tradeoffs"
heading: "A committed digest"
group: "Overview"
order: 100
url: "/docs/tradeoffs#a-committed-digest"
anchor: "a-committed-digest"
terms: ["committed","digest","generated","offline","markdown","tree","rather","computed","runtime","hidden","service","upside","reviewable","section","pull","requests","deterministic","free","read","bundled","edge","worker","without","filesystem","access","directly","navigable","agent","cost","stale","regenerating","only","content","change","build","though","warns","hash","mismatch","gate"]
hash: "110adfc68831d0e1a9dcd3844ea876fb7fbd211e79537b3c2e854618aaecdc8d"
mode: "agent-primary"
facts: [{"kind":"code","literal":"tree","chunkId":"tradeoffs#a-committed-digest"},{"kind":"code","literal":"cat","chunkId":"tradeoffs#a-committed-digest"},{"kind":"code","literal":"grep","chunkId":"tradeoffs#a-committed-digest"}]
sources: [{"chunkId":"tradeoffs#a-committed-digest","url":"/docs/tradeoffs#a-committed-digest","anchor":"a-committed-digest"}]
---

The digest is generated offline and committed to git as a markdown tree rather than computed at runtime or hidden in a service. The upside is that it is reviewable per section in pull requests, deterministic and free to read at runtime, bundled into the edge worker without filesystem access, and directly navigable by an agent; the cost is that it can go stale, regenerating only on a content change and build, though the runtime warns on hash mismatch and the per-section hash gate makes rebuilding on every content change in CI the intended workflow.
