---
id: "api/endpoint#suggested-questions-get"
title: "Search endpoint"
heading: "Suggested questions (GET)"
group: "API"
order: 32
url: "/docs/api/endpoint#suggested-questions-get"
anchor: "suggested-questions-get"
terms: ["suggested","questions","returns","digest","baked","suggestions","loop","model","query","call","does","stay","fresh","claude","haiku","overlay","fetches","once","first","open","populate","empty","array","graph","without","just","means","shows","none"]
hash: "365b4f536408592ba0f2e389bafa8fbc2827905bc3787ca00efcb19eafd919ba"
mode: "source-primary"
facts: [{"kind":"code","literal":"{\n  \"suggestions\": [\"How does the digest stay fresh?\"],\n  \"model\": \"claude-haiku-4-5\"\n}","chunkId":"api/endpoint#suggested-questions-get"},{"kind":"code","literal":"GET /api/ask","chunkId":"api/endpoint#suggested-questions-get"},{"kind":"code","literal":"suggestions","chunkId":"api/endpoint#suggested-questions-get"}]
sources: [{"chunkId":"api/endpoint#suggested-questions-get","url":"/docs/api/endpoint#suggested-questions-get","anchor":"suggested-questions-get"}]
---

Suggested questions (GET) GET /api/ask returns the digest's baked-in suggestions and the loop model — no query, no model call: { "suggestions": ["How does the digest stay fresh?"], "model": "claude-haiku-4-5" } The overlay fetches this once on first open (when...
