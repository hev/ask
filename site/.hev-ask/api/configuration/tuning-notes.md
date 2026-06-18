---
id: "api/configuration#tuning-notes"
title: "Configuration"
heading: "Tuning notes"
group: "API"
order: 13
url: "/docs/api/configuration#tuning-notes"
anchor: "tuning-notes"
terms: ["tuning","notes","guidance","search","knobs","raise","chunk","heading","depth","finer","anchors","long","pages","lower","subsections","small","iteration","tighter","latency","multi","part","questions","adjust","document","control","result","spread","candidate","count","trade","recall","against","noise","tokens","chunkheadingdepth","maxiterations","perdoccap","candidatepersearch","granularity","default"]
hash: "357904627c549ccc59120393a2316912522af1027d5a3d35be21390520aadac2"
mode: "source-primary"
facts: [{"kind":"code","literal":"chunkHeadingDepth","chunkId":"api/configuration#tuning-notes"},{"kind":"code","literal":"###","chunkId":"api/configuration#tuning-notes"},{"kind":"code","literal":"maxIterations","chunkId":"api/configuration#tuning-notes"},{"kind":"code","literal":"perDocCap","chunkId":"api/configuration#tuning-notes"},{"kind":"code","literal":"candidatePerSearch","chunkId":"api/configuration#tuning-notes"}]
sources: [{"chunkId":"api/configuration#tuning-notes","url":"/docs/api/configuration#tuning-notes","anchor":"tuning-notes"}]
---

Guidance on tuning the search knobs: raise chunk heading depth for finer anchors on long pages or lower it if subsections are too small, lower the iteration cap for tighter latency or raise it for multi-part questions, adjust the per-document cap to control result spread, and adjust the per-search candidate count to trade recall against noise and tokens.
