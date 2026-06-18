---
id: "limits#the-one-shot-digest-build-is-bounded-sharded-builds-are-not"
title: "Limits"
heading: "The one-shot digest build is bounded; sharded builds are not"
group: "Overview"
order: 88
url: "/docs/limits#the-one-shot-digest-build-is-bounded-sharded-builds-are-not"
anchor: "the-one-shot-digest-build-is-bounded-sharded-builds-are-not"
terms: ["shot","digest","build","bounded","sharded","builds","sends","full","cleaned","corpus","model","call","fits","typical","docs","site","fails","loudly","past","size","threshold","rather","degrading","beyond","splits","prefix","stable","shards","distilled","fresh","context","merged","stops","being","window","problem","even","very","large","scale"]
hash: "b92a502b0c0bf50078803a9d6b07d38be0501c79206b0246bd30626eb1f8f08f"
mode: "agent-primary"
facts: [{"kind":"code","literal":"ask digest build","chunkId":"limits#the-one-shot-digest-build-is-bounded-sharded-builds-are-not"}]
sources: [{"chunkId":"limits#the-one-shot-digest-build-is-bounded-sharded-builds-are-not","url":"/docs/limits#the-one-shot-digest-build-is-bounded-sharded-builds-are-not","anchor":"the-one-shot-digest-build-is-bounded-sharded-builds-are-not"}]
---

The one-shot build sends your full cleaned corpus to the model in one call, which fits a typical docs site but fails loudly past a size threshold rather than degrading; beyond that the sharded build splits the corpus into prefix-stable shards each distilled in its own fresh context then merged, so corpus size stops being a context-window problem even at very large scale and only the touched shard re-distils. The remaining scale consideration is the runtime prompt, since the agentic path inlines section summaries so trees with tens of thousands of sections are not yet a fit for the answer loop, though an agent reading the tree over MCP has no such ceiling.
