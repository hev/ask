---
id: "api/search-overlay#bundling-the-static-assets"
title: "SearchOverlay component"
heading: "Bundling the static assets"
group: "API"
order: 44
url: "/docs/api/search-overlay#bundling-the-static-assets"
anchor: "bundling-the-static-assets"
terms: ["bundling","static","assets","command","emits","browser","payload","keyword","index","glossary","suggestions","title","tree","directory","site","serves","during","build","like","step","renders","html","output","served","rather","committed","should","gitignored","while","digest","stays","reviewable","source","truth","regenerating","every","keeps","drifting","deploy","bundle"]
hash: "77f8c50f2486577e256ed8eda9ca5648330bd15f01f56ec275d154b6dfa0019d"
mode: "source-primary"
facts: [{"kind":"code","literal":"ask digest bundle","chunkId":"api/search-overlay#bundling-the-static-assets"},{"kind":"code","literal":".hev-ask/","chunkId":"api/search-overlay#bundling-the-static-assets"}]
sources: [{"chunkId":"api/search-overlay#bundling-the-static-assets","url":"/docs/api/search-overlay#bundling-the-static-assets","anchor":"bundling-the-static-assets"}]
---

The bundling command emits the browser payload (keyword index, glossary, suggestions, and title-tree) into a directory your site serves, run during your build like the step that renders HTML. The output is served rather than committed, so it should be gitignored while the committed digest tree stays the reviewable source of truth; regenerating every build keeps the assets from drifting from what you deploy.
