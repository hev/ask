---
id: "api/cli#reading-the-digest-as-a-directory"
title: "CLI"
heading: "Reading the digest as a directory"
group: "API"
order: 7
url: "/docs/api/cli#reading-the-digest-as-a-directory"
anchor: "reading-the-digest-as-a-directory"
terms: ["reading","digest","directory","default","reads","come","local","tree","deployed","site","http","endpoint","passed","path","addressed","mirror","urls","there","verb","operation","real","markdown","files","plain","shell","tools","work","verbs","structured","frontmatter","fuzzy","resolution","glossary","aliases","remote","step","cheapest","first","move","synthesized"]
hash: "f04f8863a9420ba4d49e89d9df17323581843c1d5e16be86991322c0735ab771"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":".hev-ask/","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"--endpoint \u003curl\u003e","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"ls","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"head","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"--endpoint","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"ask tree","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"--depth N","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"--depth all","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"(+N)","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"cat","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"facts","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"ask answer","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"ask grep","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"ask cat","chunkId":"api/cli#reading-the-digest-as-a-directory"}]
sources: [{"chunkId":"api/cli#reading-the-digest-as-a-directory","url":"/docs/api/cli#reading-the-digest-as-a-directory","anchor":"reading-the-digest-as-a-directory"}]
---

By default reads come from the local tree, or from a deployed site's HTTP API when an endpoint is passed, and they are path-addressed to mirror your doc URLs. There is one verb per operation and the tree is real markdown files, so plain shell tools work on it too; the verbs add structured frontmatter, fuzzy path resolution, glossary aliases, and remote reads, with the map step the cheapest first move and the synthesized-answer verb requiring a deployed endpoint.
