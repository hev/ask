# RFC 0006: Ask across all local digests — a model-routed, machine-local docs oracle

## Summary

Turn the set of `.hev-ask/` digests on one machine into a **single queryable
surface**. RFC 0005 makes one digest *discoverable*; this RFC makes *all of them*
collectively *askable*. A registry enumerates the local digests with their
routing metadata; given a question and no chosen digest, the agentic loop (and
the MCP server) **routes** to the relevant digest(s) automatically — reads each
digest's orientation summary, shortlists by relevance, searches within, and
answers with deep links that resolve to each digest's own site. The model picks
the digest; the user just asks.

The signal this needs already exists: every `_meta.md` carries a domain Context
paragraph, suggested questions, a glossary, and an overview map — that *is* the
routing index. What's missing is (1) a registry that lists the digests, (2) a
routing step, and (3) a site origin per digest so cross-digest links resolve.
This RFC adds exactly those three, reusing the existing keyword scorer and
agentic loop rather than introducing new retrieval machinery.

## Motivation

**Knowledge is local and structured, but siloed per digest.** A developer with
five docs repos has five digests — each one a faithful, agent-readable
distillation — but every `ask`, overlay, and `ask mcp` server is scoped to one.
To answer "how does auth work in project Y" you must already know it's project Y,
`cd` there, and ask. The routing a person does in their head ("that's a Y
question") is precisely what the digests' domain summaries could do for the
model.

**This is the natural top of the discovery ladder.** RFC 0005 resolves a *single*
digest and, in doing so, already has to *enumerate* the candidates under `$HOME`.
Promoting that enumeration to a machine-readable registry and adding a routing
step is a small step with a large payoff: the digests stop being five tools and
become one.

**It pays off twice — at the terminal and in the agent.** At my terminal, `ask`
answers across *my* projects without my having to remember which repo holds the
answer — that's the local CLI doing the routing over my `$HOME` digests. In a
coding agent, one `ask mcp` server lets Claude Code (or any MCP client) draw on a
*lot* of projects — including external repos I've cloned and digested — routing
automatically. Same registry, same routing; the CLI is the front door for what's
mine, MCP for federating many. "Have the model find the right doc in the right
digest" is, concretely, whichever model is asking reading the registry's
orientation metadata and choosing where to search. That is what "resolvable by
the model automatically" means here.

## Goals

- Enumerate all known local digests as a machine-readable registry carrying the
  routing metadata (domain summary, glossary terms, suggested questions, site
  origin, content hash).
- A question with no digest specified is routed by the model to the relevant
  digest(s), searched, and answered — with deep links that resolve to each
  digest's own site.
- Routing is cheap (metadata-first; keyword shortlist before full search) and
  legible (the answer states which digest(s) it drew from).
- **Full parity between the CLI and MCP.** Every capability — list, route,
  ask-across — is a first-class CLI command over my `$HOME` digests *and* a
  mirrored MCP tool; neither surface is a lesser cousin. The CLI is the way I
  query *my own* projects; MCP is the way an agent client federates *many
  external* ones. Same core, two front doors.
- Degrade cleanly: one digest behaves exactly like today; zero digests gives a
  clear "nothing registered" message.

## Non-goals

- **No remote or central index.** Local filesystem only; no upload, no hosted
  registry.
- **No new retrieval stack.** Routing reuses the existing token-overlap scorer
  and glossary widening; embeddings are explicitly out of scope (open question).
- **No taxonomy merging.** We route, search, and cite *per* digest; we do not
  reconcile two projects' glossaries into one answer space.
- **Not single-digest resolution.** That is RFC 0005; this RFC consumes it.
- **Not the browser overlay.** The overlay is inherently one-site; cross-digest
  asking is a CLI/MCP capability.

## Design

### The registry

A machine-readable list of known digests, at
`${XDG_CONFIG_HOME:-~/.config}/hev-ask/registry.json`. Each entry:

- `path` — absolute path to the `.hev-ask/` tree.
- `name` / `domain` — short label distilled from `_meta.md`'s Context.
- `origin` — the doc site's base URL **(new)**; see "Cross-digest deep links."
- `contentHash` — from `_meta.md`, for staleness.
- `summary`, top `glossary` terms, `suggestions` — the routing signal, copied
  (small) from `_meta.md` so routing needs no full-tree read.

Populated three ways, in increasing explicitness: **auto-register on `ask digest
build`** (a repo enrolls itself when it builds), **`ask digest register [dir]` /
`unregister`** for manual control, and an optional bounded **`ask digest scan`**
that walks `$HOME` once (skipping `node_modules`, `.git`, vendored trees) to
bootstrap an existing machine. The registry is a cache, not a source of truth:
entries with a missing `path` are pruned lazily, and a stale `contentHash` is
flagged, never trusted silently. This is the machine-readable form of RFC 0005's
`ask digest which`.

### Two scopes, one capability set: mine (`$HOME`) and external (MCP)

The same registry and the same operations serve two usage patterns, and the
split is the point of this RFC's framing:

- **Mine — the `$HOME` set, via the CLI.** The digests under my home directory
  are the projects *I develop*. `ask digest list` shows them, `ask answer` routes
  across them. This is the everyday local loop: I'm in a terminal, I ask, the CLI
  answers from my own projects. No server, no client, no MCP.
- **Many external — via MCP.** When I want an *agent* client (Claude Code,
  Cursor) to draw on a *lot* of projects I don't own — OSS repos I've cloned and
  digested for reference — MCP is the federation surface: point the client at one
  `ask mcp` server and it lists and routes across the whole set.

Both read the same registry; the difference is breadth and who's asking. To keep
the two legible, each entry carries a `scope` tag (`home` for an auto-discovered
`$HOME` digest, `external` for one registered from elsewhere), and `list` /
routing take an optional `--scope home|external|all` (CLI default: `home`; MCP
default: `all`). Same core, filtered by intent.

### Routing — model-first, keyword-shortlisted

Given a question and no explicit `--digest-dir`:

1. **Shortlist.** Score the query against each registry entry's `summary` +
   `glossary` with the existing token-overlap scorer → top-K digests. Cheap,
   deterministic, and it bounds the work before any model call.
2. **Route.** Hand the loop a `list_digests` tool that returns the shortlist's
   orientation summaries; the model chooses which to search — it may pick one or
   several. (Simpler v1 fallback: auto-search the top-K and let the loop's own
   `search` sub-queries fan across them.)
3. **Search + answer.** The loop's `search` gains a `digest` argument naming
   which digest to query; every result is tagged with its source digest, and the
   final answer names the digest(s) it drew from.

Routing leans on metadata the digests already carry, so the added cost is one
shortlist pass plus the orientation summaries in context — not a second index.

### Cross-digest deep links

Today deep links are **relative** (`/docs/page#anchor`) — fine for one site,
broken across many. Each registry entry records the digest's `origin` (site base
URL), captured by `ask digest build` (which already knows `basePath`; it gains
the configured site origin). A cross-digest answer renders each citation against
its digest's `origin`, so a link into project Y resolves to
`https://projY.example/docs/page#anchor`, not a path that only works on the
current site.

### Surfaces — CLI and MCP as peers

Every capability exists on both surfaces over the same registry. The CLI is not
a thin demo of the MCP tools; it is the primary way I work with my own projects.

- **CLI (mine).** Three commands, mirroring the three MCP tools:
  - `ask digest list [--scope home|external|all]` — the local equivalent of
    `list_digests`; defaults to my `$HOME` digests.
  - `ask answer "q" [--scope …]` — routes across the registry when no
    `--digest-dir` is set and more than one digest is in scope. The explicit-flag
    and single-digest paths are unchanged (RFC 0005).
  - `ask cat/grep/tree --digest <name>` — the read verbs, addressed by registry
    name instead of a path, so I can reach into any of my projects without `cd`.
- **MCP (many external).** The server, scoped to one digest today (`fetch_docs`
  to hydrate, `answer` to answer), gains a **`list_digests`** tool and a `digest`
  selector on `fetch_docs`/`answer`. An MCP client thus enumerates the registry,
  reads the orientation metadata, and routes — extending RFC 0003's hydrate model
  from one digest to many. `ask mcp --all` (or a registry-backed default) serves
  the registry instead of a single tree; `--scope` bounds what it exposes.

The two stay in lockstep by sharing one registry and one routing function — a new
capability lands as a CLI command and an MCP tool in the same change, never one
without the other.

## Consequences

- **hev ask becomes a machine-local, multi-project docs oracle** — answerable
  from the terminal across my own projects and from an agent client across many
  external ones, the clearest expression yet of "your docs as a directory your
  agent can read."
- **New public surface:** the registry file; `ask digest register` /
  `unregister` / `scan` / `which --json`; routing in `ask answer`; MCP
  `list_digests` and digest-scoped `fetch_docs`/`answer`. Additive → **minor
  bump**, but sizable; needs its own docs page (a "Local digests" or "Multi-repo"
  concept + API rows).
- **A new captured field, `origin`,** flows from build config into the digest
  meta and registry. Existing digests gain it on next build; absent it,
  cross-digest links degrade to relative (same-site only) with a warning.
- **Privacy note:** the registry enumerates local repo paths and domain summaries.
  Local-only, but worth saying plainly in the docs.
- **Performance:** routing adds a shortlist pass bounded by registry size; the
  per-digest summaries keep model context small.

## Migration

Additive. Single-digest flows — flag, env var, walk-up, one `ask mcp` — are
unchanged. The registry is opt-in and self-populating (auto-register on build),
so the capability appears as repos are rebuilt; nothing breaks if it stays empty.
`origin` is backfilled on the next `ask digest build` per repo.

## Sequencing

Depends on RFC 0005 (the discovery primitive and `which`).

1. Registry format + `register` / `unregister` / `which --json`; auto-register on
   `ask digest build`.
2. Capture `origin` in build config → `_meta.md` → registry entry.
3. Routing: keyword shortlist over registry, then `list_digests` tool + `search`
   digest-scoping; cross-digest citation rendering.
4. CLI `ask answer` registry routing.
5. MCP `list_digests` + digest-scoped `fetch_docs`/`answer`; `ask mcp --all`.
6. Docs: a multi-repo concept page + API rows; `make check`; minor version bump.

## Open questions

- **Registry population default.** Auto-on-build + explicit register (proposed)
  vs an opt-out `$HOME` scan. How aggressive should `scan` be, and does it
  respect `.gitignore`/ignore vendored trees by default?
- **What makes a digest "external"?** Path alone is weak — a cloned OSS repo I
  only *reference* still lives under `$HOME`. Is `scope` inferred (e.g. inside a
  configured "projects" root = `home`, elsewhere = `external`), set at register
  time, or just a tag I can override? The CLI-vs-MCP default (`home` vs `all`)
  should stay sensible whatever drives the tag.
- **Routing fidelity.** Is keyword-shortlist-then-model enough, or does good
  cross-domain routing eventually want embeddings over the per-digest summaries?
  (Out of scope for v1; flagged so we don't paint into a corner.)
- **Answer shape.** When two digests both contribute, is the answer one merged
  narrative with per-digest citations, or sectioned by digest? Citation model
  needs to make the source digest unambiguous either way.
- **`origin` source.** Read from the consumer's site config at build time, a new
  `hevAsk({ origin })` option, or inferred? Required for cross-digest links;
  optional (same-site degrade) otherwise.
- **Registry concurrency.** Multiple `ask` processes writing the registry —
  last-writer-wins on a small file, or a lock? Lean atomic-replace, no lock.
- **Staleness policy.** Flag-and-use stale entries, or refuse and prompt a
  rebuild? Lean flag-and-use, mirroring 0005's "legible, not blocking" stance.

## References

- RFC 0005 (`0005-default-digest-discovery.md`) — the discovery ladder and
  `which`/registry surface this RFC consumes and promotes to machine-readable.
- RFC 0003 (`0003-mcp-as-hydrate.md`) — the single-digest hydrate model
  (`fetch_docs`, `answer`) this RFC extends to multi-digest.
- `pkg/ask/mcp.go` — the MCP server (`tools/list`, `tools/call`; `fetch_docs`,
  `answer`) gaining `list_digests` and a `digest` selector.
- `site/.hev-ask/_meta.md` — the per-digest Context / suggestions / overview that
  serve as the routing signal; the missing `origin` field is added here.
- `packages/ui/src/llm.ts` — the loop's tool schema (`input_schema`); `search`
  gains a `digest` argument.
- `CLAUDE.md` — the corpus/anchor/host-neutral facts the registry and routing
  must preserve (corpus = configured collections only; links via github-slugger).
