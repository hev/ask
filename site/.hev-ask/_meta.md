---
version: 2
generatedAt: "2026-06-09T12:33:33.990Z"
contentHash: "03dd6f09c9146cd785f9a521fd82bb61efe3b6e151968824ad0c380384f253af"
suggestions: ["How do I add hev ask to my Astro site?","How does the agentic answer loop work?","Why keyword search instead of embeddings?","How do I build and refresh the ask digest?","What can't hev ask do?"]
---

## Context

**hev ask** (`@hevmind/ask`) is a `⌘K` search overlay for Astro docs sites, documented at hevask.com — a site that searches itself with the package. Three moving parts: a heading-level chunk index with real rendered anchors (via github-slugger, gated by `ask digest verify`), a committed offline-built **ask digest** (`.hev-ask/`) stored as a markdown tree with glossary, section summaries, hashes, facts, and source anchors, and a bounded Claude tool-use loop. Two retrieval paths: instant keyless **keyword** search (token overlap widened by the digest glossary, deep-linking to `/docs/page#anchor`) and an **agentic** answer on Enter (needs `ANTHROPIC_API_KEY`, streams SSE with inline citations). The digest is hash-gated — rebuilds skip model work when content is unchanged — and buildable three ways: a Claude Code skill (subscription, sharded), the one-shot CLI, or the sharded flow for big sites. Everything degrades instead of hard-failing: no key → keyword mode; no digest tree → plain excerpts. The corpus is only the configured content collection(s) — no crawler. The `POST /api/ask` endpoint renders on demand (keyword JSON vs. agentic SSE; keyless GET read routes for glossary/sections/overview), so a server or hybrid adapter is required. The `ask` CLI exposes `tree`, `ls`, `head`, `cat`, `facts`, `grep`, `answer`, `mcp`, and `digest build`/`verify`. Users compare it against Pagefind, Algolia DocSearch, and Orama; the Tradeoffs and Limits pages answer that directly.
## Overview

## API
- CLI — `api/cli`
- Building the digest — `api/cli#building-the-digest`
- Claude Code skill — `api/cli#claude-code-skill`
- Distribution — `api/cli#distribution`
- Flags — `api/cli#flags`
- Go library — `api/cli#go-library`
- MCP — `api/cli#mcp`
- Reading the digest as a directory — `api/cli#reading-the-digest-as-a-directory`
- Sharded builds for large sites — `api/cli#sharded-builds-for-large-sites`
- Where it runs — `api/cli#where-it-runs`
- Configuration — `api/configuration`
- Options — `api/configuration#options`
- Tuning notes — `api/configuration#tuning-notes`
- TypeScript — `api/configuration#typescript`
- What the integration does — `api/configuration#what-the-integration-does`
- Digest format — `api/digest`
- _meta.md and _glossary/ — `api/digest#_metamd-and-_glossary`
- A section file — `api/digest#a-section-file`
- Degradation — `api/digest#degradation`
- Frontmatter fields — `api/digest#frontmatter-fields`
- How each field is used — `api/digest#how-each-field-is-used`
- Layout — `api/digest#layout`
- Regenerating — `api/digest#regenerating`
- Search endpoint — `api/endpoint`
- Agentic response (SSE) — `api/endpoint#agentic-response-sse`
- Digest reads (GET) — `api/endpoint#digest-reads-get`
- Errors — `api/endpoint#errors`
- Index lifecycle — `api/endpoint#index-lifecycle`
- Keyword response (JSON) — `api/endpoint#keyword-response-json`
- LLM tracing — `api/endpoint#llm-tracing`
- Mode selection — `api/endpoint#mode-selection`
- Request — `api/endpoint#request`
- Suggested questions (GET) — `api/endpoint#suggested-questions-get`
- The API key — `api/endpoint#the-api-key`
- MCP server — `api/mcp`
- Co-location — `api/mcp#co-location`
- Configure — `api/mcp#configure`
- Data sources — `api/mcp#data-sources`
- Protocol surface — `api/mcp#protocol-surface`
- The instructions — `api/mcp#the-instructions`
- The tool — `api/mcp#the-tool`
- SearchOverlay component — `api/search-overlay`
- Keyboard model — `api/search-overlay#keyboard-model`
- Keyword results and deep links — `api/search-overlay#keyword-results-and-deep-links`
- Opening the overlay — `api/search-overlay#opening-the-overlay`
- Props — `api/search-overlay#props`
- Suggested questions — `api/search-overlay#suggested-questions`
- The mode toggle — `api/search-overlay#the-mode-toggle`
- The streamed answer — `api/search-overlay#the-streamed-answer`
- Theming — `api/search-overlay#theming`
## Overview
- Concepts — `concepts`
- Asking is the default — `concepts#asking-is-the-default`
- Chunks and anchors — `concepts#chunks-and-anchors`
- Degradation, by design — `concepts#degradation-by-design`
- Keyword search and the glossary — `concepts#keyword-search-and-the-glossary`
- Progressive disclosure as a directory — `concepts#progressive-disclosure-as-a-directory`
- The agentic search loop — `concepts#the-agentic-search-loop`
- The ask digest directory — `concepts#the-ask-digest-directory`
- The system prompt is cached — `concepts#the-system-prompt-is-cached`
- Two ways to build the tree — `concepts#two-ways-to-build-the-tree`
- Introduction — `index`
- Next steps — `index#next-steps`
- One artifact, three surfaces — `index#one-artifact-three-surfaces`
- Who this is for — `index#who-this-is-for`
- Limits — `limits`
- A server route is required — `limits#a-server-route-is-required`
- Agentic search adds latency — `limits#agentic-search-adds-latency`
- Anchors depend on Astro's slugger — `limits#anchors-depend-on-astros-slugger`
- Frontmatter parsing is a flat-YAML subset — `limits#frontmatter-parsing-is-a-flat-yaml-subset`
- Recall has a keyword ceiling — `limits#recall-has-a-keyword-ceiling`
- Secrets live server-side — `limits#secrets-live-server-side`
- The corpus is your content collection — `limits#the-corpus-is-your-content-collection`
- The one-shot digest build is bounded; sharded builds are not — `limits#the-one-shot-digest-build-is-bounded-sharded-builds-are-not`
- Quick start — `quickstart`
- 1. Install — `quickstart#1-install`
- 2. Register the integration — `quickstart#2-register-the-integration`
- 3. Add a server adapter — `quickstart#3-add-a-server-adapter`
- 4. Render the overlay — `quickstart#4-render-the-overlay`
- 5. Build the digest — `quickstart#5-build-the-digest`
- Enable agentic search — `quickstart#enable-agentic-search`
- Prerequisites — `quickstart#prerequisites`
- Set up keyword search — `quickstart#set-up-keyword-search`
- Verify it works — `quickstart#verify-it-works`
- Tradeoffs — `tradeoffs`
- A committed digest — `tradeoffs#a-committed-digest`
- Cost and latency of agentic search — `tradeoffs#cost-and-latency-of-agentic-search`
- How it compares — `tradeoffs#how-it-compares`
- Keyword retrieval, not embeddings — `tradeoffs#keyword-retrieval-not-embeddings`
- One dependency, deliberately — `tradeoffs#one-dependency-deliberately`
- Two paths instead of one — `tradeoffs#two-paths-instead-of-one`
