---
id: "digest-creation#two-ways-to-run-the-build"
title: "Digest creation"
heading: "Two ways to run the build"
group: "Overview"
order: 73
url: "/docs/digest-creation#two-ways-to-run-the-build"
anchor: "two-ways-to-run-the-build"
terms: ["ways","build","runs","write","same","tree","under","hash","gate","recommended","bundled","claude","code","skill","builds","inside","subscription","token","spend","sharding","corpus","distilling","shard","fresh","context","size","never","hits","limit","does","distillation","unattended","call","outside","default","provider","flag","full","command","flags"]
hash: "cc8ff5c21b0b7eeda4b14b44b4f1df1a979bf029b005e2bf8c24d334fe803413"
mode: "agent-primary"
facts: [{"kind":"code","literal":"You: build the hev ask digest\n\nClaude runs:\n  ask digest corpus       # emits the sections to distil\n  …writes context/glossary/summaries/suggestions…\n  ask digest assemble     # writes the .hev-ask/ tree","chunkId":"digest-creation#two-ways-to-run-the-build"},{"kind":"code","literal":"export ANTHROPIC_API_KEY=sk-ant-...\npnpm exec ask digest build","chunkId":"digest-creation#two-ways-to-run-the-build"},{"kind":"code","literal":"export OPENAI_API_KEY=sk-...\npnpm exec ask digest build --provider openai","chunkId":"digest-creation#two-ways-to-run-the-build"},{"kind":"code","literal":"export OPENROUTER_API_KEY=sk-or-...\npnpm exec ask digest build --provider openrouter","chunkId":"digest-creation#two-ways-to-run-the-build"},{"kind":"code","literal":"build-digest","chunkId":"digest-creation#two-ways-to-run-the-build"},{"kind":"code","literal":"ANTHROPIC_API_KEY","chunkId":"digest-creation#two-ways-to-run-the-build"},{"kind":"code","literal":"--provider","chunkId":"digest-creation#two-ways-to-run-the-build"},{"kind":"value","literal":"4.8","chunkId":"digest-creation#two-ways-to-run-the-build"}]
sources: [{"chunkId":"digest-creation#two-ways-to-run-the-build","url":"/docs/digest-creation#two-ways-to-run-the-build","anchor":"two-ways-to-run-the-build"}]
---

The build runs two ways that write the same tree under the same hash gate: the recommended bundled Claude Code skill builds inside your subscription with no key and no token spend, sharding the corpus and distilling each shard in a fresh context so size never hits a context limit; and the CLI build does the same distillation unattended in one API call for CI or outside Claude Code, on the default provider or any provider via a flag. The full command set, flags, and sharded flow are in the CLI reference.
