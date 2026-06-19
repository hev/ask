---
id: "api/search-overlay#bundling-the-static-assets"
title: "SearchOverlay component"
heading: "Bundling the static assets"
group: "API"
order: 44
url: "/docs/api/search-overlay#bundling-the-static-assets"
anchor: "bundling-the-static-assets"
terms: ["bundling","static","assets","bundle","command","emits","browser","payload","keyword","index","glossary","suggestions","title","tree","served","directory","part","build","like","html","render","step","output","rather","committed","gitignore","while","keeping","reviewable","source","truth","regenerating","every","keeps","drifting","deploy","digest","site","serves","renders"]
hash: "77f8c50f2486577e256ed8eda9ca5648330bd15f01f56ec275d154b6dfa0019d"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask digest bundle","chunkId":"api/search-overlay#bundling-the-static-assets"},{"kind":"code","literal":".hev-ask/","chunkId":"api/search-overlay#bundling-the-static-assets"}]
sources: [{"chunkId":"api/search-overlay#bundling-the-static-assets","url":"/docs/api/search-overlay#bundling-the-static-assets","anchor":"bundling-the-static-assets"}]
---

The bundle command emits the browser payload (keyword index, glossary, suggestions, and title-tree) into a served directory, run as part of your build like the HTML render step. The output is served rather than committed, so gitignore it while keeping the committed tree as the reviewable source of truth; regenerating every build keeps the assets from drifting from what you deploy.
