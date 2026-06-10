---
id: "api/configuration#what-the-integration-does"
title: "Configuration"
heading: "What the integration does"
group: "API"
order: 14
url: "/docs/api/configuration#what-the-integration-does"
anchor: "what-the-integration-does"
terms: ["integration","does","astro","starts","config","setup","injects","endpoint","route","pointing","hevmind","prerender","false","renders","demand","registers","virtual","modules","resol","digest","digestdir","collections","build","start","anthropic","resolved","options","committed","tree","globbed","adds","watched","files","reloads","changes","warns","empty","runs","anthropicapikey","present"]
hash: "6040e0af6f965fba16023cf7991f297be8eeaf29ec57c84e789f00bf91f54ba0"
mode: "source-primary"
facts: [{"kind":"code","literal":"astro:config:setup","chunkId":"api/configuration#what-the-integration-does"},{"kind":"code","literal":"endpoint","chunkId":"api/configuration#what-the-integration-does"},{"kind":"code","literal":"@hevmind/ask/endpoint","chunkId":"api/configuration#what-the-integration-does"},{"kind":"code","literal":"prerender: false","chunkId":"api/configuration#what-the-integration-does"},{"kind":"code","literal":"virtual:hev-ask/config","chunkId":"api/configuration#what-the-integration-does"},{"kind":"code","literal":"virtual:hev-ask/digest","chunkId":"api/configuration#what-the-integration-does"},{"kind":"code","literal":"digestDir","chunkId":"api/configuration#what-the-integration-does"},{"kind":"code","literal":"collections","chunkId":"api/configuration#what-the-integration-does"},{"kind":"code","literal":"astro:build:start","chunkId":"api/configuration#what-the-integration-does"},{"kind":"code","literal":"ANTHROPIC_API_KEY","chunkId":"api/configuration#what-the-integration-does"}]
sources: [{"chunkId":"api/configuration#what-the-integration-does","url":"/docs/api/configuration#what-the-integration-does","anchor":"what-the-integration-does"}]
---

What the integration does When Astro starts up (astro:config:setup) the integration: injects the endpoint route, pointing at @hevmind/ask/endpoint, with prerender: false so it renders on demand; registers two virtual modules — virtual:hev-ask/config (the resol...
