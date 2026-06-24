# Three Patterns for Better Agentic Search — slide outline

> Working outline. Headlines + on-slide bullets + speaker notes + demo beats.
> Pour the headlines/bullets into Keynote/Slides; keep the notes here.

## Talk meta

- **Event:** Maven — *Three Patterns for Better Agentic Search*
- **When:** Thu June 25, 2026, 5:00 PM UTC · 1 hour · virtual (Zoom) · ~85 registered
- **With:** Doug Turnbull & Trey Grainger (*AI-Powered Search*, Manning) — search-literate audience; they know RAG, BM25, dense/sparse retrieval, rerankers, LTR cold.
- **My role:** patterns + **live demo**. Slides are lean scaffolding; the demo carries the weight.
- **Framing:** vendor-neutral patterns, **hev ask as the running case study** for each.
- **Slides:** built as a reveal.js deck (ported from the `../mind` deck system) at
  `site/src/content/decks/three-patterns.md` → hosted at **`/decks/three-patterns`** on hevask.com.
  This markdown file stays the source of truth for the *narrative + speaker notes*; the deck is the rendered version.
- **Time budget (my portion ≈ 28 min):** cold open 3 · thesis 1.5 · P1 6 · P2 7 · P3 8 · synthesis 1.5 · close 1.

## The through-line (say it early, land it at the end)

> The index, the disclosure mechanism, and the tool surface are **the same artifact** — a committed directory of distilled markdown. You're not building three systems. The filesystem metaphor collapses storage, interface, and transport into one shape an agent already knows how to navigate.

For this audience the contrarian edge is: **agentic search is not RAG-with-a-loop.** RAG maximizes recall of an opaque index and dumps chunks into context. Agentic search makes retrieval *legible and cheap to the agent* so it can decide what to pull and when. The three patterns are how you do that.

---

## Demo corpus: Helm (locked)

- **Site:** `helm.hevask.com` (live, CF Pages). **Repo:** `~/workspace/helm-ask-demo`.
- **Why Helm:** a CNCF project the audience knows, **536 sections** in the tree — big enough that "list titles, don't read everything" is obviously the only sane move.
- **Hero section** for the demos: `intro/quickstart/install-a-chart-from-an-oci-registry` — its `facts` hold the verbatim `helm install … oci://…` output and the `kubectl port-forward` / `curl` checks; deep-links to `/docs/intro/quickstart#install-a-chart-from-an-oci-registry`.
- **Top-level groups (for `tree`):** `intro`, `overview`, `topics`, `chart_template_guide`, `chart_best_practices`, `howto`, `plugins`, `helm` (CLI command reference), `changelog`, `_glossary`.

## Demo prerequisites (verify before the talk — see open items)

- [ ] `ask` CLI on PATH; run from `~/workspace/helm-ask-demo` so it reads that committed `.hev-ask/` tree.
- [ ] Terminal font large; rehearse `ask tree`, `ask grep "oci"`, `ask cat intro/quickstart/install-a-chart-from-an-oci-registry`, `ask facts …` on the hero section.
- [ ] `helm.hevask.com` overlay live and the agentic path keyed (server secret present).
- [ ] A second repo open to show a coding agent hydrating + grep-ing the tree (MCP). **Confirm what actually ships today vs. RFC 0003.**
- [ ] `git diff` on a Helm section file ready to show per-section reviewability.

---

## Slide 1 — Title (shared)

**Three Patterns for Better Agentic Search**
- Adam Hevenor (Hevmind) · with Doug Turnbull & Trey Grainger
- *Say:* one line on who I am and that I'll ground each pattern in a real, opinionated implementation you can read on disk.

## Slide 2 — Cold open: docs were RAG's first home

**Headline:** Documentation was one of the first places we bolted on RAG. Here's the old way.

**On-slide code (the naive RAG turn):**
```python
# Docs Q&A, the old way: retrieve once, stuff, generate once.
results = index.search(embed(question), k=8)        # opaque top-k, a fixed guess
context = "\n\n".join(r.text for r in results)      # dump the whole result set in

answer = llm.chat(
    system=f"Answer the question using only this context:\n\n{context}",
    user=question,                                  # one turn, one shot
)
```
- **The whole turn is: a system prompt + the user prompt + a result set, stuffed into one LLM call.** That's it. That was state of the art, and for a lot of docs Q&A it still is.
- Point at each line as the three things that don't work for an *agent*:
  - `k=8` — a guess that over- or under-fetches every time; one chunk size for every question.
  - `"\n\n".join(...)` — the entire result set dumped in, **opaque**: the model can't ask for less, or more, or something else.
  - one `llm.chat(...)` — no follow-up, no cost model, and **no grounded link back to the source** the answer came from.
- *Say:* this audience built the retrieval under that `index.search`. The bottleneck in *agentic* search was never recall — it's that this shape hands the agent an oracle it can't reason about. We bolted a generator onto a search box and called it done.

## Slide 3 — The reframe

**Headline:** Make retrieval legible and cheap to the agent.
- The fix isn't a better `k` or a better embedding. It's to stop pre-deciding the retrieval and let the **agent** drive it — over a corpus it can navigate, with a cost model it already has.
- Not "maximize recall of an opaque index." Make the corpus something the agent **navigates**, not something we **stuff**.
- *Say:* state the through-line. Three patterns turn that single stuffed turn into something an agent can work: (1) build the index offline with a model, (2) disclose it progressively, (3) shape the loop around it. The punchline is they're one artifact.

---

## Pattern 1 — Build the index offline, with a model, and commit it

### Slide 4 — Pattern 1, the principle

**Headline:** Spend the expensive model *before* query time, not during it.
- Do the distillation **offline** with a strong model; commit the result; review it like code.
- Runtime stays deterministic, free to read, and auditable. The smart, slow, costly work happens once — at build — not on every request.
- *Audience note:* this is the inversion of query-time RAG. Move the intelligence to index time.

### Slide 5 — Pattern 1 on hev ask: the ask digest

**Headline:** A committed directory of distilled markdown.
- `.hev-ask/` — one small markdown+frontmatter file **per section** (heading-level chunk). A distilled *mirror* of the docs.
- **Model authors only the distillation** (summary, glossary, orientation). **Anchors, facts, hashes, structure are derived deterministically in code** — the model never invents a link or a flag.
- **Hash-gated & incremental:** each section carries the hash of the content it was distilled from; a rebuild re-distils only what changed; a clean tree does zero model work.
- Built with a strong model (Opus / a Claude Code skill) **offline**; committed to git.
- *Say:* the split matters — the model does the fuzzy part, code does the part that must be exact. That's why the links never 404.

### Slide 6 — DEMO: the index is reviewable

**Demo beat (~2 min) — in `~/workspace/helm-ask-demo`:**
1. `ask tree` → the whole Helm corpus as a directory (536 sections, ~10 groups).
2. Open `.hev-ask/intro/quickstart/install-a-chart-from-an-oci-registry.md` in the editor → frontmatter: `title`, `url`+`anchor`, the verbatim `helm install … oci://…` output in `facts`, `terms`, `hash`; body = the one-line distilled summary.
3. `git diff` on a section → "this is what a search-index change looks like in code review."
- *Say:* the model wrote that one-line summary; **code** derived the anchor, extracted the verbatim commands into `facts`, and stamped the hash. Your index is a PR artifact — a human sees a section's prose and its grounded commands change together, in the same markdown the docs are written in. Anchors are gated in CI against the real renderer (github-slugger), so the deep link can't 404.
- *Audience hook:* "When was the last time you code-reviewed a diff of your search index?"

---

## Pattern 2 — Disclose progressively; let the consumer choose the next rung

### Slide 7 — Pattern 2, the principle

**Headline:** Don't dump. Reveal in rungs the consumer chooses.
- Context is the scarce resource. Stuffing top-k chunks spends it blindly.
- Better: a **disclosure ladder** — cheap overview first, deeper slices only on demand, each step an explicit choice by the consumer.
- *Audience note:* this is the agentic answer to "chunk size" debates — you don't pick one granularity, you expose several and let the agent climb.

### Slide 8 — Pattern 2 on hev ask: the digest *is* a filesystem

**Headline:** Four rungs, each a strictly larger slice of one file.
- `tree` → titles only · `head` → one-line summary · `cat` → full body · `facts` → grounded literals.
- **A listing reads frontmatter only** — it *structurally cannot leak a body*. Output is bounded by section count, not corpus size. Safe to call speculatively.
- **Titles carry the whole decision** — the only thing seen before opening a file. Authored where the docs provide one, model-synthesized where they don't.
- Why a filesystem: the agent **already has the cost model** — `ls` is cheap, `cat` might be big. We don't teach a new disclosure protocol; the filesystem *is* the protocol.

### Slide 9 — DEMO: climbing the ladder

**Demo beat (~3 min) — the money demo for this crowd, on the Helm tree:**
1. `ask tree` → 536 titles, bounded (frontmatter only — *can't* leak a body).
2. `ask grep "oci"` → one keyword-in-context line per hit (evidence, not a paragraph).
3. `ask cat intro/quickstart/install-a-chart-from-an-oci-registry` → the one section that matters.
4. `ask facts intro/quickstart/install-a-chart-from-an-oci-registry` → the verbatim `helm install … oci://…` commands.
- *Say:* 536 files. The agent never reads 536 files — it lists titles, greps, and `cat`s the one. Watch the cost escalate by *choice*, never by surprise. Same ladder the overlay climbs server-side and a coding agent climbs on its own disk.
- *Audience hook:* "Progressive disclosure isn't a feature I built — it's a property of using a directory. The cost model is free, because the agent already has it."

---

## Pattern 3 — Shape the loop: files + instructions, bounded, grounded

### Slide 10 — Pattern 3, the principle

**Headline:** Give the agent files and instructions, not a query API.
- When your consumer is itself an agent, it doesn't need your reasoning — it needs your **corpus**, in a form its own tools can drive.
- Shape the loop, don't just expose endpoints: bound the iterations, force a grounded answer, make citation non-negotiable.
- *Audience note:* the split — **human consumer → synthesize an answer; agent consumer → hand over the files.**

### Slide 11 — Pattern 3 on hev ask: two loops, one ladder

**Headline:** The overlay synthesizes; the agent self-serves.
- **Human (⌘K overlay):** a bounded tool-use loop. Phase 1 *gather* — model gets the title-tree + `open_section`, opens what it needs (default ≤4 rounds, prompt-cached). Phase 2 *answer* — called once more **with no tools**, streams a grounded answer with inline deep links.
- **The trick:** dropping the tools on the final turn *forces* an answer instead of another search. Citation to exact `url`+`anchor` is enforced; non-provided links are downgraded to plain text.
- **Agent (CLI / MCP):** no synthesis loop at all — hydrate the tree to local disk, then the agent's *own* `tree`/`cat`/`grep` are the loop. MCP = files, not answers. *(hydrate shape is RFC 0003 — confirm shipped state before demoing.)*

### Slide 12 — DEMO: both loops

**Demo beat (~4 min) — on `helm.hevask.com`:**
1. **Overlay:** type `oci` → instant keyword results. Type *"how do I install a chart from an OCI registry?"* + Enter → watch the faint "opening section…" activity lines stream, then the grounded answer assemble with a deep link to `/docs/intro/quickstart#install-a-chart-from-an-oci-registry`. Click it → lands on the exact heading.
2. **Agent:** in a second repo, hydrate the Helm tree and `grep` it — no checkout, no key — then ask the coding agent the same question and watch it `cat` the right section and cite the anchor.
- *Say:* same artifact, same ladder, two consumers. The overlay does the agent's job *for the human*; the coding agent does it for itself.

---

## Slide 13 — Synthesis: one artifact, three patterns

**Headline:** It's all one directory.
- Index (P1) = the committed tree. Disclosure (P2) = the tree's own cost model. Loop (P3) = native file tools over the same tree.
- **Degradation by design:** no key at runtime → keyword mode; no key at build → keep the committed tree; no tree → raw token overlap. Nothing hard-fails.
- *Say:* the reason these compose is they're not three systems. Build it as files and the other two patterns fall out almost for free.

## Slide 14 — Close: what this means for *your* systems

**Headline:** Before you reach for a vector store, ask what your agent can navigate.
- Move intelligence to index time. Expose granularity, don't pick one. Hand agents corpora, not oracles.
- hev ask is the opinionated version; the patterns outlive the product.
- *Say:* hand to Doug/Trey for the broader search lens + Q&A. Plug: hevask.com, open source.

---

## Backup slides / anticipated questions (this audience *will* push)

- **"Why no embeddings? Isn't this just lexical search with extra steps?"**
  Keyword + glossary, no vector store — on purpose. The glossary recovers a lot of synonym recall; the agentic loop recovers paraphrase via query *reformulation*; and the index is reviewable + edge-safe with nothing to keep in sync. Embeddings are **deferred, not designed out** — a reranking/dense recall upgrade is a clean add. Be honest: paraphrase recall has a ceiling today.
- **"Isn't progressive disclosure just multi-hop RAG?"**
  Multi-hop RAG still hides retrieval behind a query API the agent can't cost. Here the *agent* chooses each rung against a cost model it already owns; disclosure is structural (a listing can't leak a body), not prompted.
- **"Cost / latency of the loop?"** One bounded loop per submitted query on a small model (Haiku), domain context prompt-cached across rounds; worst case ≈ maxIterations round-trips. Keyword path is the instant fast lane. Offline build uses Opus but only when content changes.
- **"Staleness?"** Committed digest can drift; hash-gate + a runtime warning when live hash ≠ digest hash; "rebuild in CI on content change" is the intended workflow.
- **"Why not just give the agent the raw repo / let it grep the source?"**
  Distillation is the point: bare source is noisy, unanchored, and huge. The digest is small, anchored, fact-extracted, and titled for the listing decision — it's the source *shaped for navigation*.
- **"Does this scale past medium corpora?"** Designed for medium corpora (docs sites, wikis). Bounded by the configured collection — no crawler. Large-corpus routing is being worked (search-routing instead of inlining).

## Open items for us to resolve

1. **Demo reality check** — verify exactly which surfaces work *today* vs. RFC (esp. MCP-as-hydrate / cross-site fetch in RFC 0003; multi-framework in RFC 0004). Don't demo or claim anything unshipped.
2. ~~**Demo corpus**~~ — **RESOLVED: Helm** (`helm.hevask.com`, `~/workspace/helm-ask-demo`, 536 sections). Familiar CNCF corpus; big enough to make disclosure land.
3. **Slide count vs. demo time** — current draft ~14 slides + 3 demo beats for ~28 min. Trim if my portion is shorter.
4. **Opening hook** — lead with the "code-review your search index" line, or the "retrieve/stuff/pray" failure framing? Pick one cold open.
5. **How much to engage the embeddings debate live** — Doug/Trey may want to dig in; decide whether that's a planned beat or Q&A.
