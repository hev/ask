---
id: "api/cli#reading-the-digest-as-a-directory"
title: "CLI"
heading: "Reading the digest as a directory"
group: "API"
order: 7
url: "/docs/api/cli#reading-the-digest-as-a-directory"
anchor: "reading-the-digest-as-a-directory"
terms: ["reading","digest","directory","default","reads","local","tree","remote","endpoint","flag","deployed","site","http","instead","path","addressed","mirror","urls","mapping","cheap","first","step","returning","titles","only","defaulting","levels","deep","scoped","deepened","depth","truncated","directories","reporting","hidden","count","glossary","table","collapsed","until"]
hash: "90c1c568baea1522a00a439ee7e9874bb7988865a71bdd8af7759e1b4660f504"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":".hev-ask/","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"--endpoint \u003curl\u003e","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"ls","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"head","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"--endpoint","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"ask tree","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"--depth N","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"--depth all","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"(+N)","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"_glossary/","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"ask tree _glossary","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"cat","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"facts","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"ask answer","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"ask grep","chunkId":"api/cli#reading-the-digest-as-a-directory"},{"kind":"code","literal":"ask cat","chunkId":"api/cli#reading-the-digest-as-a-directory"}]
sources: [{"chunkId":"api/cli#reading-the-digest-as-a-directory","url":"/docs/api/cli#reading-the-digest-as-a-directory","anchor":"reading-the-digest-as-a-directory"}]
---

By default the CLI reads the local tree; a remote endpoint flag reads a deployed site's HTTP API instead, with path-addressed reads that mirror doc URLs. Mapping the tree is the cheap first step, returning titles only and defaulting to two levels deep, scoped and deepened by a path and depth flag, with truncated directories reporting a hidden count; the glossary table is collapsed until you scope into it. Because the tree is real markdown, plain shell tools work on it, while the verbs add frontmatter, fuzzy path resolution, glossary aliases, and remote reads; the synthesized reply requires a remote endpoint, so keyless local retrieval uses search and read.
