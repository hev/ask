---
id: "api/endpoint#suggested-questions-get"
title: "Search endpoint"
heading: "Suggested questions (GET)"
group: "API"
order: 33
url: "/docs/api/endpoint#suggested-questions-get"
anchor: "suggested-questions-get"
terms: ["suggested","questions","endpoint","base","query","returns","digest","baked","loop","model","call","overlay","fetches","once","first","open","empty","list","suggestions","simply","means","none","shown","does","stay","fresh","claude","haiku","populate","array","without","just","shows"]
hash: "09b11b93e36f83adbf3806b64893ae0d6e67036b50132a6ffb096ad9e279489b"
mode: "source-primary"
facts: [{"kind":"code","literal":"{\n  \"suggestions\": [\"How does the digest stay fresh?\"],\n  \"model\": \"claude-haiku-4-5\"\n}","chunkId":"api/endpoint#suggested-questions-get"},{"kind":"code","literal":"GET /api/ask","chunkId":"api/endpoint#suggested-questions-get"},{"kind":"code","literal":"suggestions","chunkId":"api/endpoint#suggested-questions-get"}]
sources: [{"chunkId":"api/endpoint#suggested-questions-get","url":"/docs/api/endpoint#suggested-questions-get","anchor":"suggested-questions-get"}]
---

A GET to the endpoint base with no query returns the digest's baked-in suggested questions and the loop model with no model call; the overlay fetches this once on first open when AI is on, and an empty list (no suggestions or no digest) simply means none are shown.
