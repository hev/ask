---
id: "quickstart#5-build-the-digest"
title: "Quick start"
heading: "5. Build the digest"
group: "Overview"
order: 94
url: "/docs/quickstart#5-build-the-digest"
anchor: "5-build-the-digest"
terms: ["build","digest","offline","built","markdown","tree","commit","gives","loop","context","ranks","keyword","results","supplies","glossary","holds","suggested","questions","recommended","bundled","claude","code","skill","inside","subscription","token","spend","call","choosing","provider","verify","anchors","both","paths","incremental","hash","gated","integration","runs","automatically"]
hash: "c3e971ffc29cc5333a5448cb5e2846615f72af766d664a071dd8b904dd6b09d4"
mode: "agent-primary"
facts: [{"kind":"code","literal":"You: build the hev ask digest\n\nClaude runs:\n  ask digest corpus       # emits the sections to distil\n  …writes context/glossary/summaries/suggestions…\n  ask digest assemble     # writes the .hev-ask/ tree","chunkId":"quickstart#5-build-the-digest"},{"kind":"code","literal":"export ANTHROPIC_API_KEY=sk-ant-...\npnpm exec ask digest build      # writes the .hev-ask/ tree","chunkId":"quickstart#5-build-the-digest"},{"kind":"code","literal":"export OPENAI_API_KEY=sk-...\npnpm exec ask digest build --provider openai","chunkId":"quickstart#5-build-the-digest"},{"kind":"code","literal":"export OPENROUTER_API_KEY=sk-or-...\npnpm exec ask digest build --provider openrouter","chunkId":"quickstart#5-build-the-digest"},{"kind":"code","literal":"pnpm exec ask digest verify     # builds the site, checks every anchor resolves\ngit add .hev-ask","chunkId":"quickstart#5-build-the-digest"},{"kind":"code","literal":"k8s","chunkId":"quickstart#5-build-the-digest"},{"kind":"code","literal":"kubernetes","chunkId":"quickstart#5-build-the-digest"},{"kind":"code","literal":"--provider","chunkId":"quickstart#5-build-the-digest"},{"kind":"code","literal":"astro build","chunkId":"quickstart#5-build-the-digest"},{"kind":"value","literal":"claude.com","chunkId":"quickstart#5-build-the-digest"}]
sources: [{"chunkId":"quickstart#5-build-the-digest","url":"/docs/quickstart#5-build-the-digest","anchor":"5-build-the-digest"}]
---

The digest is an offline-built markdown tree you commit that gives the loop context, ranks keyword results, supplies the glossary, and holds suggested questions. Build it the recommended way with the bundled Claude Code skill inside your subscription (no key, no token spend), or with the one-call CLI build for CI choosing a provider, then verify anchors and commit; both paths are incremental and hash-gated, and the integration runs the build automatically during the site build when a key is present.
