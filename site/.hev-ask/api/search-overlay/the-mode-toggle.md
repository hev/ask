---
id: "api/search-overlay#the-mode-toggle"
title: "SearchOverlay component"
heading: "The mode toggle"
group: "API"
order: 54
url: "/docs/api/search-overlay#the-mode-toggle"
anchor: "the-mode-toggle"
terms: ["mode","toggle","overlay","persists","enter","preference","localstorage","readers","flip","keyword","only","never","trigger","model","call","space","just","searches","phrase","suggestions","show","choice","survives","reloads","agentic","under","suggested","questions","shown"]
hash: "f5b67968083cdfbfa00b247a46ae07a527001d5feb7c83f6ef6b3cd40deac99e"
mode: "source-primary"
facts: [{"kind":"code","literal":"localStorage","chunkId":"api/search-overlay#the-mode-toggle"},{"kind":"code","literal":"hev-ask:mode","chunkId":"api/search-overlay#the-mode-toggle"},{"kind":"code","literal":"agentic","chunkId":"api/search-overlay#the-mode-toggle"},{"kind":"code","literal":"keyword","chunkId":"api/search-overlay#the-mode-toggle"}]
sources: [{"chunkId":"api/search-overlay#the-mode-toggle","url":"/docs/api/search-overlay#the-mode-toggle","anchor":"the-mode-toggle"}]
---

The overlay persists an AI-on-Enter preference in localStorage. Readers who flip it to keyword-only never trigger a model call (a space just searches a phrase and no suggestions show), and the choice survives reloads.
