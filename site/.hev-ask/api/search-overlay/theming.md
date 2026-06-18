---
id: "api/search-overlay#theming"
title: "SearchOverlay component"
heading: "Theming"
group: "API"
order: 57
url: "/docs/api/search-overlay#theming"
anchor: "theming"
terms: ["theming","overlay","uses","scoped","class","prefix","reads","page","custom","properties","background","primary","secondary","text","accent","defining","those","tokens","root","usually","matches","site","look","override","paper","111111","fafaf5","muted","6b6b66","signal","e25822","active","state","markup","define","these","values","shown","inherits","palette"]
hash: "37bbca09a73293fc3a1fe639cabd35a98647bf221c502cc18b4e1a02585df7cf"
mode: "source-primary"
facts: [{"kind":"code","literal":":root {\n  --paper: #111111;       /* overlay background */\n  --ink: #fafaf5;         /* primary text */\n  --muted: #6b6b66;       /* secondary text */\n  --signal: #e25822;      /* accent / active state */\n}","chunkId":"api/search-overlay#theming"},{"kind":"code","literal":"as-","chunkId":"api/search-overlay#theming"},{"kind":"code","literal":":root","chunkId":"api/search-overlay#theming"}]
sources: [{"chunkId":"api/search-overlay#theming","url":"/docs/api/search-overlay#theming","anchor":"theming"}]
---

The overlay uses a scoped class prefix and reads your page's CSS custom properties for background, primary and secondary text, and accent, so defining those tokens on the root usually matches your site's look with no overlay CSS to override.
