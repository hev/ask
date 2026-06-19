---
id: "limits#the-one-shot-digest-build-is-bounded-sharded-builds-are-not"
title: "Limits"
heading: "The one-shot digest build is bounded; sharded builds are not"
group: "Overview"
order: 88
url: "/docs/limits#the-one-shot-digest-build-is-bounded-sharded-builds-are-not"
anchor: "the-one-shot-digest-build-is-bounded-sharded-builds-are-not"
terms: ["shot","digest","build","bounded","sharded","builds","sends","full","cleaned","corpus","model","single","call","fits","typical","docs","site","fails","loudly","past","section","text","size","threshold","beyond","splits","prefix","stable","shards","distilled","context","merged","deterministically","stops","being","window","problem","content","edit","distills"]
hash: "b92a502b0c0bf50078803a9d6b07d38be0501c79206b0246bd30626eb1f8f08f"
mode: "agent-primary"
facts: [{"kind":"code","literal":"ask digest build","chunkId":"limits#the-one-shot-digest-build-is-bounded-sharded-builds-are-not"}]
sources: [{"chunkId":"limits#the-one-shot-digest-build-is-bounded-sharded-builds-are-not","url":"/docs/limits#the-one-shot-digest-build-is-bounded-sharded-builds-are-not","anchor":"the-one-shot-digest-build-is-bounded-sharded-builds-are-not"}]
---

The one-shot build sends the full cleaned corpus to the model in a single call, which fits a typical docs site but fails loudly past a section-text size threshold; beyond that the sharded build splits the corpus into prefix-stable shards each distilled in its own context and merged deterministically, so corpus size stops being a context-window problem and a content edit re-distills only the touched shard. The remaining scale consideration is the runtime prompt, since the agentic path inlines section summaries, so trees with tens of thousands of sections aren't yet a fit for the answer loop, though a coding agent reading the tree over MCP has no such ceiling.
