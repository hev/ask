---
id: "api/search-overlay#theming"
title: "SearchOverlay component"
heading: "Theming"
group: "API"
order: 57
url: "/docs/api/search-overlay#theming"
anchor: "theming"
terms: ["theming","overlay","reads","page","custom","properties","background","text","muted","accent","defining","those","tokens","root","makes","inherit","palette","because","scoped","styles","keyed","variables","matching","site","look","usually","just","override","paper","111111","fafaf5","primary","6b6b66","secondary","signal","e25822","active","state","markup","uses"]
hash: "37bbca09a73293fc3a1fe639cabd35a98647bf221c502cc18b4e1a02585df7cf"
mode: "source-primary"
facts: [{"kind":"code","literal":":root {\n  --paper: #111111;       /* overlay background */\n  --ink: #fafaf5;         /* primary text */\n  --muted: #6b6b66;       /* secondary text */\n  --signal: #e25822;      /* accent / active state */\n}","chunkId":"api/search-overlay#theming"},{"kind":"code","literal":"as-","chunkId":"api/search-overlay#theming"},{"kind":"code","literal":":root","chunkId":"api/search-overlay#theming"}]
sources: [{"chunkId":"api/search-overlay#theming","url":"/docs/api/search-overlay#theming","anchor":"theming"}]
---

The overlay reads your page's CSS custom properties for background, text, muted text, and accent, so defining those tokens on the root makes the overlay inherit your palette. Because its scoped styles are keyed to those variables, matching your site's look is usually just defining the tokens with no overlay CSS to override.
