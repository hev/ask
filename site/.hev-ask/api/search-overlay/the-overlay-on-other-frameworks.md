---
id: "api/search-overlay#the-overlay-on-other-frameworks"
title: "SearchOverlay component"
heading: "The overlay on other frameworks"
group: "API"
order: 55
url: "/docs/api/search-overlay#the-overlay-on-other-frameworks"
anchor: "the-overlay-on-other-frameworks"
terms: ["overlay","other","frameworks","astro","component","distribution","same","palette","ships","prebuilt","site","loads","script","reads","bundled","copy","digest","browser","keyword","search","runs","fully","static","sends","agentic","questions","whatever","endpoint","point","optional","attribute","staying","drop","custom","properties","shares","opener","keyboard","model","theming"]
hash: "294cf6be8257ee487de6d9399df7bb94a4a88fbf3c1865008b451588f489d85e"
mode: "source-primary"
facts: [{"kind":"code","literal":"\u003cscript\n  type=\"module\"\n  src=\"https://cdn.jsdelivr.net/npm/@hevmind/ask/overlay.js\"\n  data-hev-ask-digest=\"/hev-ask/\"\n  data-hev-ask-endpoint=\"https://docs-ask.example.workers.dev/api/ask\"\n\u003e\u003c/script\u003e\n\n\u003cbutton data-hev-ask-open\u003eSearch \u003ckbd\u003e⌘K\u003c/kbd\u003e\u003c/button\u003e","chunkId":"api/search-overlay#the-overlay-on-other-frameworks"},{"kind":"code","literal":"SearchOverlay.astro","chunkId":"api/search-overlay#the-overlay-on-other-frameworks"},{"kind":"code","literal":"@hevmind/ask/overlay","chunkId":"api/search-overlay#the-overlay-on-other-frameworks"},{"kind":"code","literal":"data-hev-ask-digest","chunkId":"api/search-overlay#the-overlay-on-other-frameworks"},{"kind":"code","literal":"ask digest bundle","chunkId":"api/search-overlay#the-overlay-on-other-frameworks"},{"kind":"code","literal":"data-hev-ask-endpoint","chunkId":"api/search-overlay#the-overlay-on-other-frameworks"}]
sources: [{"chunkId":"api/search-overlay#the-overlay-on-other-frameworks","url":"/docs/api/search-overlay#the-overlay-on-other-frameworks","anchor":"the-overlay-on-other-frameworks"}]
---

The Astro component is one distribution of the overlay; the same palette ships as a prebuilt web component any site loads with one script tag from npm or a CDN. It reads a bundled copy of the digest in the browser so keyword search runs fully static, and it sends agentic questions to whatever endpoint you point it at via an optional attribute, with the key staying on the endpoint. The drop-in reads the same CSS custom properties and shares the opener attribute and keyboard model, so theming applies unchanged.
