# hev ask documentation

**hev ask** distills a docs site into an **ask digest**: a compact
**directory** with one small markdown file per section. It's a mirror of your
docs that an agent navigates like any other directory, and that a `⌘K` overlay
answers your readers from. The digest is **host-neutral**: it's built from your
markdown, not your renderer, so it works with [Astro](https://astro.build),
Docusaurus, VitePress, MkDocs, or any folder of markdown. hev ask ships a
turnkey Astro integration, and other frameworks add the same overlay with one
script. You build the digest **offline** with the bundled skill and your Claude
Code subscription, then commit the tree. It's a good fit for technical
documentation, internal wikis, and other medium-sized corpora.

*the ask digest as a directory*

```
$ ask tree

_glossary/                  (+10)
_meta                       Digest metadata
api/
  cli/                      CLI  (+9)
  configuration/            Configuration  (+5)
  digest/                   Digest format  (+7)
  endpoint/                 Search endpoint  (+10)
  mcp/                      MCP server  (+6)
  search-overlay/           SearchOverlay component  (+16)
concepts/                   Concepts
  chunks-and-anchors        Concepts > Chunks and anchors
  the-agentic-search-loop   Concepts > The agentic search loop
  the-ask-digest-directory  Concepts > The ask digest directory
  …                         (+7)
digest-creation/            Digest creation  (+5)
limits/                     Limits  (+8)
quickstart/                 Quick start  (+9)
tradeoffs/                  Tradeoffs  (+7)
```

## One artifact, three readers

The digest is built **inside your coding agent** — the bundled Claude Code
skill writes the `.hev-ask/` tree using your existing **agentic coding
subscription**, so there's no provider API key and no per-build token spend.
Commit the tree, and the same files serve:

- **the `⌘K` overlay**, for humans — instant keyword results as readers type,
  plus a grounded Claude answer on Enter, every result deep-linked to the exact
  heading;
- **the `ask` CLI**, for agents — `tree` to map it, `cat` to read, `facts` for
  grounded literals, `grep` to search, over path keys, keyless, from any shell;
- **the `ask mcp` server**, for agents — one tool that hydrates the whole tree
  to local disk, after which the agent reads it with the file tools it
  already has.

The overlay synthesizes an answer for the human reader. The CLI and MCP hand
an agent the raw files and let its own tools do the rest.

The same tree also makes a ready-made [`llms.txt`](https://llmstxt.org): every
section is a distilled summary with a source deep link, so you can serve the
digest at `/llms.txt` and `/llms-full.txt` for agents that look there first.

> **drop ask into your own CLI**
>
> The digest reads (`tree`, `cat`, `grep`, `facts`) also ship as a
> [Go library](api/cli.md#go-library) in `pkg/ask`, beyond the bundled `ask`
> binary. Mount the dependency-free command group in your own CLI with
> `ask.NewCommandGroup(...)`, or call the lower-level helpers directly:
> `LoadDigest` (from disk or an `embed.FS`), `GetSection`, `SearchDigest`, and
> `ServeMCP`. Every read is keyless and offline, so your tool grounds in the same
> committed tree the overlay and MCP do.

## Who this is for

You're building or maintaining a docs site whose content lives in Markdown or
MDX — on **Astro** (the turnkey integration), **Docusaurus**, **VitePress**,
**MkDocs**, or anything that renders a folder of markdown. You want search and
answers that:

- work out of the box without standing up a service or running a crawler,
- deep-link to the right section instead of dumping the reader at the top of a
  long page,
- can answer a question phrased in the reader's words, not just match keywords, and
- are **queryable by your coding agent**, not only by humans in a browser.

On Astro, one integration covers all of that. Other frameworks drop in the
static overlay and, for the agentic answers, point it at an optional hosted
endpoint. See
[the drop-in overlay](api/search-overlay.md#the-overlay-on-other-frameworks)
for the path on yours.

If you only need keyword search over a static site and never want an API key in
the loop, [Pagefind](https://pagefind.app) is simpler and a great fit — see
[Tradeoffs](tradeoffs.md) for an honest comparison.

## Contents

**Overview**

- [Quick start](quickstart.md) — add search to an Astro site in five minutes.
- [Digest creation](digest-creation.md) — how the tree is built, kept fresh,
  and verified.
- [Concepts](concepts.md) — chunks, anchors, the disclosure ladder, the agentic
  loop, and the ask digest directory.
- [Tradeoffs](tradeoffs.md) — what you're choosing, against Pagefind, Algolia,
  and Orama.
- [Limits](limits.md) — what hev ask deliberately doesn't do.
- [Roadmap & changelog](roadmap.md) — what's planned, and what shipped when.

**API reference**

- [Configuration](api/configuration.md) — every `hevAsk()` option and its default.
- [SearchOverlay component](api/search-overlay.md) — props, keyboard model,
  theming, and the drop-in overlay for non-Astro sites.
- [Search endpoint](api/endpoint.md) — the `/api/ask` route contract.
- [CLI](api/cli.md) — `tree`, `cat`, `facts`, `grep`, and the Go library.
- [MCP server](api/mcp.md) — hand a coding agent the digest as a local directory.
- [Digest format](api/digest.md) — the committed `.hev-ask/` tree, field by field.

**Engineering**

- [RFCs](rfcs/) — design alignment before code.
- [Talks](talks/) — decks and outlines.
