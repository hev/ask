---
id: "api/search-overlay#docusaurus"
title: "SearchOverlay component"
heading: "Docusaurus"
group: "API"
order: 45
url: "/docs/api/search-overlay#docusaurus"
anchor: "docusaurus"
terms: ["docusaurus","shows","wiring","emit","bundle","static","directory","build","script","overlay","through","site","config","headings","carry","explicit","custom","anchors","slug","mode","honors","verify","checks","against","built","html","export","default","scripts","https","jsdelivr","hevmind","type","module","data","digest","endpoint","docs","example","workers"]
hash: "6eceffab590110ca1344359049f1a1892f1f4d5c2871542003b2a624ada9a945"
mode: "source-primary"
facts: [{"kind":"code","literal":"// docusaurus.config.js\nexport default {\n  scripts: [\n    {\n      src: \"https://cdn.jsdelivr.net/npm/@hevmind/ask/overlay.js\",\n      type: \"module\",\n      \"data-hev-ask-digest\": \"/hev-ask/\",\n      \"data-hev-ask-endpoint\": \"https://docs-ask.example.workers.dev/api/ask\",\n    },\n  ],\n};","chunkId":"api/search-overlay#docusaurus"},{"kind":"code","literal":"// package.json — bundle into the static dir before docusaurus build\n\"scripts\": {\n  \"build\": \"ask digest bundle --out static/hev-ask \u0026\u0026 docusaurus build\"\n}","chunkId":"api/search-overlay#docusaurus"},{"kind":"code","literal":"static/","chunkId":"api/search-overlay#docusaurus"},{"kind":"code","literal":"build","chunkId":"api/search-overlay#docusaurus"},{"kind":"code","literal":"docusaurus.config.js","chunkId":"api/search-overlay#docusaurus"},{"kind":"code","literal":"{#custom-id}","chunkId":"api/search-overlay#docusaurus"},{"kind":"code","literal":"ask digest verify","chunkId":"api/search-overlay#docusaurus"},{"kind":"code","literal":"build/","chunkId":"api/search-overlay#docusaurus"}]
sources: [{"chunkId":"api/search-overlay#docusaurus","url":"/docs/api/search-overlay#docusaurus","anchor":"docusaurus"}]
---

Shows the Docusaurus wiring: emit the bundle into the static directory from the build script and add the overlay script through the site config. Docusaurus headings carry explicit custom-id anchors, which its slug mode honors and verify checks against the built HTML.
