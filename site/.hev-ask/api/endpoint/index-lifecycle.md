---
id: "api/endpoint#index-lifecycle"
title: "Search endpoint"
heading: "Index lifecycle"
group: "API"
order: 28
url: "/docs/api/endpoint#index-lifecycle"
anchor: "index-lifecycle"
terms: ["index","lifecycle","chunk","built","once","server","instance","first","request","cached","process","lifetime","endpoint","also","compares","live","content","hash","against","digest","logs","time","warning","differ","signaling","rebuild","build"]
hash: "8ae12673a130ca300e56e9955e6f0835a89f68ad5a0152b23de3665127f7dee8"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask digest build","chunkId":"api/endpoint#index-lifecycle"}]
sources: [{"chunkId":"api/endpoint#index-lifecycle","url":"/docs/api/endpoint#index-lifecycle","anchor":"index-lifecycle"}]
---

The chunk index is built once per server instance on the first request and cached for the process lifetime; on that first request the endpoint also compares the live content hash against the digest's and logs a one-time warning if they differ, signaling that a rebuild is due.
