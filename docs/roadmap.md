# Roadmap & Changelog

hev ask is at **0.3.4** today. What's planned comes from the engineering
[RFCs](rfcs/), and the changelog below
records what shipped in each tagged release.

## 🗺️ Roadmap

### Up next

- 🌐 **Host-neutral delivery beyond Astro** ([RFC 0004](rfcs/0004-beyond-astro-host-neutral-adapters.md)) — a client-side keyword overlay, an `ask digest bundle` step, and a hostable answer endpoint, so Docusaurus, VitePress, and MkDocs get the same overlay.
- 🧭 **Default digest discovery** ([RFC 0005](rfcs/0005-default-digest-discovery.md)) — `ask` walks up to find the nearest `.hev-ask/` (and honors `$HEV_ASK_DIGEST_DIR`), the way `git` and `cargo` find their roots.
- 🔌 **Pluggable LLM tracing** ([RFC 0007](rfcs/0007-pluggable-llm-tracing.md)) — pick the tracing backend like the inference provider; one OTLP adapter reaches Datadog, LangSmith, and Braintrust, while PostHog stays native.

### Later

- 🧠 **Embeddings for paraphrase recall** — a vector path for questions that share no tokens with your docs; the known fix for the [keyword recall ceiling](limits.md#recall-has-a-keyword-ceiling), [deferred on purpose](tradeoffs.md#keyword-retrieval-not-embeddings) until analytics show it's needed.
- 🔮 **Ask across all local digests** ([RFC 0006](rfcs/0006-ask-across-local-digests.md)) — a registry of every `.hev-ask/` on your machine plus model-routed retrieval, so one question reaches the right docs across projects without choosing a digest first.

## 📦 Changelog

### 0.3.1 – 0.3.4 — June 20, 2026

Keyword ranking overhaul, so results rank by relevance instead of raw token
overlap (see [how ranking works](api/endpoint.md#how-keyword-ranking-works)):

- 🎯 IDF-weighted scoring, BM25 length-normalization, and a 3× boost for heading
  and title matches.
- 🧹 Stopword filtering and full-phrase-only glossary expansion, so common words
  stop flooding the results.
- 🧵 Cleaner answer turn when a large digest is searched rather than inlined.

### 0.3.0 — June 20, 2026

- 🧭 **Big-digest routing.** Large digests are searched instead of inlined,
  keeping the answer loop within budget on sites with thousands of sections.
- 🪶 **One verb per operation** in the [CLI](api/cli.md): `tree` to map, `cat`
  to read, `grep` to search, folding the older `ls`, `sections`, `head`,
  `section`, `overview`, and `search` verbs.
- 🛠️ **Binary-based, cross-agent [`build-digest` skill](digest-creation.md)**,
  so the offline build runs the same way under any coding agent.
- 📈 **Outbound-click analytics** on GitHub links.

### 0.2.0 — June 10, 2026

- 🌳 **The digest became a committed markdown tree** with filesystem read verbs
  ([RFC 0002](rfcs/0002-digest-as-filesystem.md)).
  One small file per section, reviewable per PR. *(breaking)*
- 🔌 **MCP collapsed to one hydrate tool**
  ([RFC 0003](rfcs/0003-mcp-as-hydrate.md)):
  it fetches the whole tree to local disk, then the agent reads it with the file
  tools it already has. *(breaking)*
- 🧪 **Provider-pluggable inference.** OpenAI and OpenRouter joined Anthropic,
  with [tabbed provider examples](api/configuration.md#choosing-a-provider)
  and Shiki syntax highlighting.
- 🏷️ **hevask.com became the canonical domain**, and the homepage pivoted to the
  agent-first pitch.
- 📜 **The engineering RFC process** was established (RFCs 0001–0004).

### 0.1.0 / 0.1.1 — June 5, 2026

The first public release.

- ⌘ **The `⌘K` overlay**: instant keyword search over heading anchors, plus a
  grounded Claude answer streamed on Enter, every result deep-linked to the
  exact section.
- 🌲 **The ask digest and the `ask` CLI**, renamed from the knowledge graph. The
  committed artifact grounds the answer loop and ranks keyword results.
- 💡 **Suggested questions** baked into the digest at build time.
- 📊 **PostHog `$ai_generation` tracing** for the answer loop.
- 📦 **Published to npm** as `@hevmind/ask`, with platform binaries for the
  `ask` CLI.
