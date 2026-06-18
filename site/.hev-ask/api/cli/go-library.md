---
id: "api/cli#go-library"
title: "CLI"
heading: "Go library"
group: "API"
order: 5
url: "/docs/api/cli#go-library"
anchor: "go-library"
terms: ["library","reusable","lets","call","pure","helpers","directly","mount","dependency","free","command","group","inside","lower","level","loading","tree","disk","embedded","filesystem","listing","summaries","reading","sections","searching","talking","remote","endpoint","serving","newcommandgroup","commandoptions","digestdir","string","overview","quick","start","stdin","stdout","stderr","loaddigest"]
hash: "fdda297bbdabdbd250da161cec7574e09209c99f50b254923354681e1bef245c"
mode: "source-primary"
facts: [{"kind":"code","literal":"group := ask.NewCommandGroup(ask.CommandOptions{\n\tDigestDir: \".hev-ask\",\n})\nerr := group.Run(ctx, []string{\"cat\", \"overview/quick-start\"}, os.Stdin, os.Stdout, os.Stderr)","chunkId":"api/cli#go-library"},{"kind":"code","literal":"pkg/ask","chunkId":"api/cli#go-library"},{"kind":"code","literal":"LoadDigest","chunkId":"api/cli#go-library"},{"kind":"code","literal":"embed.FS","chunkId":"api/cli#go-library"},{"kind":"code","literal":"ListSectionSummaries","chunkId":"api/cli#go-library"},{"kind":"code","literal":"GetSection","chunkId":"api/cli#go-library"},{"kind":"code","literal":"SearchDigest","chunkId":"api/cli#go-library"},{"kind":"code","literal":"NewEndpointClient","chunkId":"api/cli#go-library"},{"kind":"code","literal":"ServeMCP","chunkId":"api/cli#go-library"}]
sources: [{"chunkId":"api/cli#go-library","url":"/docs/api/cli#go-library","anchor":"go-library"}]
---

A reusable Go API lets you call pure helpers directly or mount the dependency-free command group inside your own CLI, with lower-level helpers for loading the tree (from disk or an embedded filesystem), listing summaries, reading sections, searching, talking to a remote endpoint, and serving MCP.
