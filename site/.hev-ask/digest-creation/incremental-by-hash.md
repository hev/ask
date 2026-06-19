---
id: "digest-creation#incremental-by-hash"
title: "Digest creation"
heading: "Incremental by hash"
group: "Overview"
order: 72
url: "/docs/digest-creation#incremental-by-hash"
anchor: "incremental-by-hash"
terms: ["incremental","hash","every","section","file","records","content","distilled","rebuild","distills","only","sections","whose","changed","clean","tree","does","model","work","makes","rebuilding","change","cheap","enough","intended","workflow","rather","chore","distils"]
hash: "55d38a351a62131967527fca13fca8ad624324075d63337614941148ce9174f1"
mode: "agent-primary"
facts: []
sources: [{"chunkId":"digest-creation#incremental-by-hash","url":"/docs/digest-creation#incremental-by-hash","anchor":"incremental-by-hash"}]
---

Every section file records the hash of the content it was distilled from, so a rebuild re-distills only the sections whose hash changed and a clean tree does no model work at all. That makes rebuilding on every content change cheap enough to be the intended workflow rather than a chore.
