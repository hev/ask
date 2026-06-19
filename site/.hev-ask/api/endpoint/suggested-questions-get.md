---
id: "api/endpoint#suggested-questions-get"
title: "Search endpoint"
heading: "Suggested questions (GET)"
group: "API"
order: 33
url: "/docs/api/endpoint#suggested-questions-get"
anchor: "suggested-questions-get"
terms: ["suggested","questions","base","route","returns","digest","baked","loop","model","query","call","overlay","fetches","once","first","open","enabled","populate","suggestions","empty","array","including","simply","shows","none","does","stay","fresh","claude","haiku","without","just","means"]
hash: "09b11b93e36f83adbf3806b64893ae0d6e67036b50132a6ffb096ad9e279489b"
mode: "source-primary"
facts: [{"kind":"code","literal":"{\n  \"suggestions\": [\"How does the digest stay fresh?\"],\n  \"model\": \"claude-haiku-4-5\"\n}","chunkId":"api/endpoint#suggested-questions-get"},{"kind":"code","literal":"GET /api/ask","chunkId":"api/endpoint#suggested-questions-get"},{"kind":"code","literal":"suggestions","chunkId":"api/endpoint#suggested-questions-get"}]
sources: [{"chunkId":"api/endpoint#suggested-questions-get","url":"/docs/api/endpoint#suggested-questions-get","anchor":"suggested-questions-get"}]
---

A GET to the base route returns the digest's baked-in suggested questions and the loop model with no query and no model call. The overlay fetches this once on first open when AI is enabled to populate its suggestions; an empty array, including no digest at all, simply shows none.
