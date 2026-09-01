# hev ask — Agent Guide

hev ask is a `⌘K` search overlay for docs sites, shipped as the npm package
`@hevmind/ask`. This file is for agents doing engineering, docs, and release work
in this repo. Keep it practical and current.

## What this is

`@hevmind/ask` ships today as an **Astro integration** — the flagship adapter. A
consumer site adds `hevAsk()` to `astro.config`, drops `SearchOverlay.astro` in a
layout, and gets two search paths over its content collection:

- **Keyword (instant, keyless):** debounced token-overlap search over
  heading-level chunks, widened by a glossary. Results deep-link to
  `/docs/page#anchor`.
- **Agentic (on Enter, needs the provider API key):** a bounded Claude tool-use
  loop that issues its own `search` sub-queries, then streams a grounded answer
  (SSE) with inline deep links to the doc sections it drew from.

**Host-neutral core, Astro flagship.** The digest is built from markdown, not a
renderer — the Go core and offline build have zero Astro deps — so the same
digest, overlay, CLI, and MCP work whatever framework ships the docs. Astro is
the only adapter wired end-to-end *today*; extending the same overlay to
Docusaurus/VitePress/MkDocs (a static drop-in + a hostable endpoint) is designed
in **RFC 0004**, not yet shipped. When editing engineering code, treat
multi-framework as in-design; the docs in `docs/` describe the target state per a
working-backward exercise (nothing is deployed). Don't claim a
non-Astro adapter works in code comments or commit messages until it does.

A committed, offline-built **ask digest** (`.hev-ask/`, a markdown tree) gives
the loop domain context, section summaries, grounded facts, source anchors, and a
glossary. It was called the knowledge graph (`kg`) before 0.1; the rename is
complete — the CLI group is `ask digest`, the artifact flag is `--digest-dir`,
and the virtual module is `virtual:hev-ask/digest`. Don't reintroduce `kg` names
or say "knowledge graph" in new code or copy (glossary aliases, old-URL
redirects, and explicit legacy-migration docs are the only places it remains on
purpose).

## Repo layout

```
packages/ui    # the package @hevmind/ask — integration, endpoint, search, digest/, CLI
playground     # minimal Astro site for fast local dev of the package; dogfoods @hevmind/ask
docs           # the product docs (the published surface, read on GitHub)
docs/rfcs      # engineering RFCs (same process as ../layer); design alignment before code
docs/talks     # talk outlines and decks
```

It's a pnpm workspace. `packages/ui` is the only published artifact; `playground`
is its private consumer.

## The audience (informs everything we write)

The reader is an **Astro author evaluating search for a docs site**. Their
questions, in order, are: *What is this and why over Pagefind/Algolia/Orama? How
does it work? What can't it do? What am I trading off? How do I add it in five
minutes? What's the full API?* The contents list in `docs/README.md` is
structured to answer them in that order: Overview (Quick start, Digest creation,
Concepts, Tradeoffs, Limits) then API reference. Per-framework
overlay wiring lives on the SearchOverlay reference page, not an Overview page.

Docs-first is the working principle: when changing the package's public
surface, update the docs in `docs/` in the same change.

## Key facts that are easy to get wrong

- **Corpus = configured content collection(s) only.** No crawler, no external
  index, no non-collection pages.
- **Anchors come from `github-slugger`** (the one non-Astro dependency) to match
  Astro's rendered `id`s byte-for-byte. `ask digest verify` is the CI gate that
  catches drift — keep it green.
- **The digest is a committed markdown tree, hash-gated.** `ask digest build`
  skips the model call when the content hash is unchanged and re-distils only
  changed sections when section hashes are available. Regenerate and commit
  after content changes. It's reviewable on purpose.
- **Everything degrades, nothing hard-fails:** no key at runtime → keyword
  mode; no key at build → keep committed digest tree and warn; no `.hev-ask/`
  tree → empty digest.
- **The endpoint renders on demand** (`prerender: false`), so consumers need a
  server/hybrid adapter. A static-only build can't serve search.
- **Inference is provider-pluggable** (`provider` option: `anthropic` default,
  `openai`, `openrouter`; registry in `packages/ui/src/providers.ts`). Each
  provider reads its own key env var (`ANTHROPIC_API_KEY`, `OPENAI_API_KEY`,
  `OPENROUTER_API_KEY`). OpenAI and OpenRouter share the Chat Completions
  client in `llm-openai.ts`; `providerBaseUrl` points it at any
  OpenAI-compatible endpoint. Default models on Anthropic: loop =
  `claude-haiku-4-5`, digest build = `claude-opus-4-8`.

## Public surface (don't break without a version bump + doc update)

- Default export `hevAsk(options)` — options in `packages/ui/src/types.ts`,
  documented in `docs/api/configuration.md`.
- `@hevmind/ask/components/SearchOverlay.astro` — props `endpoint`, `placeholder`,
  `debounce`; opener attribute `data-hev-ask-open`; localStorage key
  `hev-ask:mode`.
- `@hevmind/ask/endpoint` — `POST /api/ask`: keyword mode returns JSON, agentic
  mode streams SSE (`text/event-stream`). Contract in `api/endpoint.mdx`.
- `ask` bin — read verbs (`tree` [path] [--depth], `cat`, `facts`, `grep`),
  `answer`, `mcp`, and `digest build` / `digest verify`. One verb per operation:
  `tree` maps (absorbs the old `ls`/`sections`), `cat` reads (absorbs
  `head`/`section`/`overview`/glossary `get`), `grep` searches (the old `search`
  was an alias). Flags in `api/cli.mdx`.
- Virtual modules `virtual:hev-ask/config` and `virtual:hev-ask/digest`.

When any of these change, update the matching `docs/api/*.md` page in the same
PR.

## The docs (`docs/`)

- The docs are **plain GitHub-flavored markdown read on GitHub**. There is no
  docs site: hevask.com was retired and now 301s to
  `https://github.com/hev/ask`. Don't reintroduce MDX components, frontmatter,
  or a nav config — a heading and a relative link is the whole toolkit.
- `docs/README.md` is the entry point and holds the contents list; every other
  page hangs off it. Cross-links are repo-relative (`api/cli.md#go-library`), so
  they work on GitHub and in an editor.
- Example endpoints in the docs use `https://docs.example.com/api/ask`. hev ask
  hosts no public endpoint — don't put a real hev hostname in a `--endpoint`
  example.
- ASCII architecture diagrams are inline fenced blocks. The originals were
  generated from `site/src/lib/diagrams.ts`; that file is gone, so edit the
  fences directly.
- `/llms.txt` and `/llms-full.txt` are still generated **by consumers** from
  their own digest — that's a product feature, not something this repo serves.
- **Retired:** the Astro site under `site/`, the Cloudflare Pages project
  `hev-ask`, and the hosted `/api/ask` endpoint. hevask.com and askhev.com stay
  registered and redirect to the GitHub repo via zone-level Redirect Rules in
  the Cloudflare dashboard. If you need the old site, it's in git history before
  the `site/` removal.

## Common commands

```sh
pnpm install                          # workspace install
pnpm dev                              # the playground site, for package dev
pnpm test                             # package unit tests
pnpm typecheck                        # tsc across the workspace
pnpm exec ask digest build           # (from a site dir) rebuild the digest
pnpm exec ask digest verify          # (from a site dir) verify anchors
```

## Before changing the package's public API

1. Update `packages/ui/src/types.ts` and the implementation.
2. Update the matching `docs/api/*.md`.
3. `pnpm test && pnpm typecheck`.
4. If anchors or chunking changed, run `ask digest verify` on `playground/`.
5. Public/breaking changes need a version bump in `packages/ui/package.json`
   (see `README.md` for publishing notes).
