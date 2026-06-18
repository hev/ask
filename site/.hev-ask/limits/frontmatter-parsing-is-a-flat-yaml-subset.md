---
id: "limits#frontmatter-parsing-is-a-flat-yaml-subset"
title: "Limits"
heading: "Frontmatter parsing is a flat-YAML subset"
group: "Overview"
order: 83
url: "/docs/limits#frontmatter-parsing-is-a-flat-yaml-subset"
anchor: "frontmatter-parsing-is-a-flat-yaml-subset"
terms: ["frontmatter","parsing","flat","yaml","subset","offline","build","parses","small","splitter","rather","full","parser","handling","common","docs","schema","string","number","fields","nested","structures","only","affects","reading","files","disk","since","astro","runtime","index","uses","collection","honors","real","getcollection","handles","aren","supported","time"]
hash: "da3ebf7a72c862c9bb94804726792203ebd2dad95887c4c094bf492223e719b8"
mode: "agent-primary"
facts: [{"kind":"code","literal":"getCollection","chunkId":"limits#frontmatter-parsing-is-a-flat-yaml-subset"}]
sources: [{"chunkId":"limits#frontmatter-parsing-is-a-flat-yaml-subset","url":"/docs/limits#frontmatter-parsing-is-a-flat-yaml-subset","anchor":"frontmatter-parsing-is-a-flat-yaml-subset"}]
---

The offline build parses frontmatter with a small flat-YAML splitter rather than a full YAML parser, handling the common docs schema of string and number fields but not nested structures; this only affects the offline build reading files from disk, since on Astro the runtime index uses the collection API that honors your real schema.
