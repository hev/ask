---
id: "api/search-overlay#suggested-questions"
title: "SearchOverlay component"
heading: "Suggested questions"
group: "API"
order: 51
url: "/docs/api/search-overlay#suggested-questions"
anchor: "suggested-questions"
terms: ["suggested","questions","overlay","fetches","short","list","endpoint","first","time","opens","shows","empty","state","come","digest","baked","suggestions","there","model","call","render","simply","nothing","clicking","suggestion","fills","input","asks","immediately","build","none","extra"]
hash: "7b69c77e3577551acb6b854069a902fa75a2f4e203eaa0343c33db7656710c47"
mode: "source-primary"
facts: [{"kind":"code","literal":"GET /api/ask","chunkId":"api/search-overlay#suggested-questions"},{"kind":"code","literal":"suggestions","chunkId":"api/search-overlay#suggested-questions"}]
sources: [{"chunkId":"api/search-overlay#suggested-questions","url":"/docs/api/search-overlay#suggested-questions","anchor":"suggested-questions"}]
---

When AI is on, the overlay fetches a short list of suggested questions from the endpoint the first time it opens and shows them in the empty state. They come from the digest's baked-in suggestions, so there is no model call to render them, and an empty list simply shows nothing; clicking a suggestion fills the input and asks it immediately.
