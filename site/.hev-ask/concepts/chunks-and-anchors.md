---
id: "concepts#chunks-and-anchors"
title: "Concepts"
heading: "Chunks and anchors"
group: "Overview"
order: 61
url: "/docs/concepts#chunks-and-anchors"
anchor: "chunks-and-anchors"
terms: ["chunks","anchors","indexes","sections","rather","pages","document","split","headings","configurable","depth","content","before","first","heading","becomes","intro","chunk","carries","built","base","path","slug","anchor","generated","same","slugger","renderer","uses","links","land","actually","exist","framework","adapter","declares","scheme","both","offline","build"]
hash: "7114c39ed148b671179972b86fbdceaa60d0146af253d0bf66722359e2522155"
mode: "agent-primary"
facts: [{"kind":"code","literal":"##","chunkId":"concepts#chunks-and-anchors"},{"kind":"code","literal":"###","chunkId":"concepts#chunks-and-anchors"},{"kind":"code","literal":"basePath + slug + #anchor","chunkId":"concepts#chunks-and-anchors"},{"kind":"code","literal":"github-slugger","chunkId":"concepts#chunks-and-anchors"},{"kind":"code","literal":"{#custom-id}","chunkId":"concepts#chunks-and-anchors"},{"kind":"code","literal":"getCollection","chunkId":"concepts#chunks-and-anchors"},{"kind":"value","literal":"github.com","chunkId":"concepts#chunks-and-anchors"}]
sources: [{"chunkId":"concepts#chunks-and-anchors","url":"/docs/concepts#chunks-and-anchors","anchor":"chunks-and-anchors"}]
---

hev ask indexes sections rather than pages: each document is split on its headings up to a configurable depth, content before the first heading becomes an intro chunk, and each chunk carries a URL built from the base path, slug, and anchor. Anchors are generated with the same slugger the renderer uses so links land on headings that actually exist, each framework adapter declares its slug scheme, and both the offline build and the runtime index chunk through one shared function so anchors agree and the same digest comes out regardless of renderer.
