---
id: "api/cli#go-library"
title: "CLI"
heading: "Go library"
group: "API"
order: 5
url: "/docs/api/cli#go-library"
anchor: "go-library"
terms: ["library","reusable","lives","pure","helpers","want","direct","control","mount","dependency","free","command","group","newcommandgroup","commandoptions","digestdir","strin","string","overview","quick","start","stdin","stdout","stderr","loaddigest","embed","listsectionsummaries","getsection","searchdigest","newendpointclient","servemcp","lower","level","include","reads","tree"]
hash: "fdda297bbdabdbd250da161cec7574e09209c99f50b254923354681e1bef245c"
mode: "source-primary"
facts: [{"kind":"code","literal":"group := ask.NewCommandGroup(ask.CommandOptions{\n\tDigestDir: \".hev-ask\",\n})\nerr := group.Run(ctx, []string{\"cat\", \"overview/quick-start\"}, os.Stdin, os.Stdout, os.Stderr)","chunkId":"api/cli#go-library"},{"kind":"code","literal":"pkg/ask","chunkId":"api/cli#go-library"},{"kind":"code","literal":"LoadDigest","chunkId":"api/cli#go-library"},{"kind":"code","literal":"embed.FS","chunkId":"api/cli#go-library"},{"kind":"code","literal":"ListSectionSummaries","chunkId":"api/cli#go-library"},{"kind":"code","literal":"GetSection","chunkId":"api/cli#go-library"},{"kind":"code","literal":"SearchDigest","chunkId":"api/cli#go-library"},{"kind":"code","literal":"NewEndpointClient","chunkId":"api/cli#go-library"},{"kind":"code","literal":"ServeMCP","chunkId":"api/cli#go-library"}]
sources: [{"chunkId":"api/cli#go-library","url":"/docs/api/cli#go-library","anchor":"go-library"}]
---

Go library The reusable Go API lives in pkg/ask. Use pure helpers when you want direct control, or mount the dependency-free command group in your own CLI. group := ask.NewCommandGroup(ask.CommandOptions{ DigestDir: ".hev-ask", }) err := group.Run(ctx, []strin...
