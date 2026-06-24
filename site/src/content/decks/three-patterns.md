---
title: Three Patterns for Better Agentic Search
description: Build the index offline, disclose it progressively, and shape the tool loop around it — with hev ask over the Helm docs as the worked example.
pubDate: 2026-06-25
---

<!-- .slide: class="title" -->

<p class="kicker">agentic search · maven</p>

# Three Patterns for Better Agentic Search

### Build the index offline · disclose it progressively · shape the loop around it

<p class="muted">Adam Hevenor · hev ask — with Doug Turnbull & Trey Grainger · &nbsp; → advance · <code>S</code> notes · <code>ESC</code> overview</p>

Note: One line on who I am, and that every pattern is grounded in a real,
opinionated implementation you can read on disk — over the Helm docs.

---

<p class="kicker">about</p>

## hi. I'm hev.

<div class="about">
<div class="about-photo">
<img src="/hev-portrait.jpg" alt="Adam Hevenor" />
</div>
<div class="about-bio">

**25 years building software systems**

- dot-com trading floors → mobile → **10 years in cloud infra**
- patents in **vector search · distributed logging · observability-as-code**
- led the **AI & search team at Aerospike**

</div>
</div>

Note: ~30s, credibility-first — a hands-on infra + AI builder: AI/search at Aerospike, a vector-search patent. Then into the patterns.

---

<p class="kicker">who this is for</p>

## A 101 — bring a corpus, not a search background

<div class="cols">
<div class="card">

**You'll feel at home if you're**

- new to retrieval or RAG
- a docs writer or maintainer
- curious how agentic search works

</div>
<div class="card signal">

**No need to know** embeddings, NDCG, rerankers, or any search tooling.

**Bring one thing:** a corpus you'd like to search.

</div>
</div>

> Raise your hand anytime — questions are the whole point.

Note: Say the level out loud — this is a 101, no search background assumed. If
you've got a folder of docs you wish were searchable, you're exactly who this is
for. Make it safe to interrupt; invite questions throughout, not just at the end.

---

<p class="kicker">where this started</p>

## Docs were one of the first places we bolted on RAG

```python
# Docs Q&A, the old way: retrieve once, stuff, generate once.
results = index.search(embed(question), k=8)        # opaque top-k, a fixed guess
context = "\n\n".join(r.text for r in results)      # dump the whole result set in

answer = llm.chat(
    system=f"Answer using only this context:\n\n{context}",
    user=question,                                  # one turn, one shot
)
```

<p class="muted">A system prompt, the user's question, and a result set — stuffed into one LLM call. That was state of the art, and for a lot of docs Q&A it still is.</p>

Note: This audience built the retrieval under that index.search. The point
isn't that it's bad — it's the *shape*.

---

<p class="kicker">the old way · the picture</p>

## One straight line: retrieve → stuff → answer

```text
        question
           │
           ▼
    embed → search → top-k chunks
           │
           ▼
   ┌───────────────────────────────┐
   │ ONE prompt                    │
   │ system + question + chunks    │
   └───────────────────────────────┘
           │
           ▼
    one LLM call  →  answer
```

<p class="muted">Retrieve once, stuff everything in, answer in one shot — no looking back.</p>

Note: The same code as a picture. A straight line — one retrieval, one prompt, one call. Hold this shape; every pattern today bends it.

---

<p class="kicker">the old way</p>

## One stuffed turn. Three problems for an agent.

- <code>k=8</code> — a fixed guess; over- or under-fetches every question <!-- .element: class="fragment" -->
- <code>"\n\n".join(...)</code> — the whole result set, **opaque**: the model can't ask for less, more, or other <!-- .element: class="fragment" -->
- one <code>llm.chat(...)</code> — no follow-up, no cost model, **no link home** <!-- .element: class="fragment" -->

<p class="muted fragment">The retrieval is great. The shape hands the agent an oracle it can't reason about.</p>

---

<p class="kicker">the reframe</p>

## Make retrieval legible and cheap to the agent

Not *"maximize recall of an opaque index."* Make the corpus the agent **navigates** — with the cost model it already has.

<div class="cols cols-3">
<div class="card">

**01 · Build offline**

Spend the expensive model before query time.

</div>
<div class="card">

**02 · Disclose progressively**

Reveal in rungs the agent chooses.

</div>
<div class="card signal">

**03 · Shape the loop**

Files and instructions, not a query API.

</div>
</div>

<p class="muted">All three are the <strong>same artifact</strong> — one committed directory.</p>

Note: Agentic search is not RAG-with-a-loop. RAG maximizes recall of an opaque
index and dumps chunks in. These three patterns turn the stuffed turn into
something an agent can work.

---

<!-- .slide: data-background-color="#111111" class="invert" -->

<p class="kicker">pattern 01</p>

# Build the index offline, with a model — and commit it

<p class="muted">Spend the expensive, slow, smart work once, at build time — not on every request.</p>

Note: The inversion of query-time RAG. Move the intelligence to index time;
runtime stays deterministic, free to read, and auditable.

---

<p class="kicker">hev ask · the ask digest</p>

## A committed directory of distilled markdown

<div class="cols">
<div class="card">

**The model writes**

one-line section summaries · the glossary · orientation

</div>
<div class="card signal">

**Code derives**

anchors · verbatim facts · hashes · structure

</div>
</div>

- One file **per section** — a distilled mirror, committed to git
- **Hash-gated** — only changed sections re-distil; a clean tree does no model work

<p class="muted">Code derives the anchors — so the deep links never 404.</p>

Note: Built offline with Opus or the Claude Code skill, and reviewed per
section in PRs like any other change.

---

<!-- .slide: data-background-color="#111111" class="invert" -->

<p class="kicker">▶ demo · helm.hevask.com</p>

## Code-review your search index

```sh
ask tree                       # 536 Helm sections, as a directory
$EDITOR .hev-ask/intro/quickstart/install-a-chart-from-an-oci-registry.md
git diff                       # the index change, right there in the PR
```

<p class="muted">Frontmatter carries the <code>url#anchor</code>, the verbatim <code>helm install … oci://…</code> in <code>facts</code>, the content hash. Anchors are CI-gated against the real renderer — the link can't drift.</p>

Note: "When was the last time you code-reviewed a diff of your search index?"

---

<!-- .slide: data-background-color="#111111" class="invert" -->

<p class="kicker">pattern 02</p>

# Disclose progressively — let the consumer choose the next rung

<p class="muted">Context is the scarce resource. Stuffing top-k spends it blindly.</p>

Note: This is the agentic answer to the chunk-size debate — you don't pick one
granularity, you expose several and let the agent climb.

---

<p class="kicker">hev ask · the digest is a filesystem</p>

## Four rungs, each a larger slice of one file

| rung | verb | returns |
| --- | --- | --- |
| 0 | `tree` / `ls` | titles only — **bounded; can't leak a body** |
| 1 | `head` | the one-line summary |
| 2 | `cat` | the full section |
| 3 | `facts` | verbatim flags & commands |

<p class="muted">A listing reads frontmatter only — so titles carry the whole decision. The agent already knows <code>ls</code> is cheap and <code>cat</code> might be big. The filesystem <em>is</em> the disclosure protocol — we don't invent a new one.</p>

---

<!-- .slide: data-background-color="#111111" class="invert" -->

<p class="kicker">▶ demo · the ladder on 536 sections</p>

## The agent never reads 536 files

```sh
ask tree                                 # titles only — bounded
ask grep "oci"                           # one keyword-in-context line per hit
ask cat   intro/quickstart/install-a-chart-from-an-oci-registry
ask facts intro/quickstart/install-a-chart-from-an-oci-registry
```

<p class="muted">It lists titles, greps, and <code>cat</code>s the one. Cost escalates by <strong>choice</strong>, never by surprise.</p>

Note: "Progressive disclosure isn't a feature I built — it's a property of
using a directory. The cost model is free, because the agent already has it."

---

<!-- .slide: data-background-color="#111111" class="invert" -->

<p class="kicker">pattern 03</p>

# Shape the loop — files and instructions, not a query API

<p class="muted">When the consumer is itself an agent, hand it the corpus — not your reasoning.</p>

---

<p class="kicker">hev ask · two consumers, one ladder</p>

## The overlay synthesizes; the agent self-serves

<div class="cols">
<div class="card">

**Human · ⌘K overlay**

- *gather* — open the sections it needs (≤ 4 rounds, cached)
- *answer* — one more call, **no tools** → a grounded answer

</div>
<div class="card signal">

**Agent · CLI / MCP**

- hydrate the tree to disk
- its **own** `tree`/`cat`/`grep` are the loop

</div>
</div>

<p class="muted">The trick: no tools on the last turn → it answers instead of re-searching.</p>

Note: Agent side has no synthesis loop — MCP hands over files, not answers.
Every claim cites an exact url#anchor; the overlay downgrades any non-provided
link to plain text.

---

<!-- .slide: data-background-color="#111111" class="invert" -->

<p class="kicker">▶ demo · both loops</p>

## Same artifact, two readers

- **Overlay** — type `oci` (instant), then ask *"how do I install a chart from an OCI registry?"* + Enter → watch it open sections, then a grounded answer with a deep link
- **Agent** — hydrate the Helm tree in another repo, `grep` it, watch the coding agent `cat` the section and cite the anchor

<p class="muted">The overlay does the agent's job for the human. The coding agent does it for itself.</p>

---

<!-- .slide: data-background-color="#111111" class="invert" -->

<p class="kicker">the through-line</p>

# It's all one directory

<div class="cols cols-3">
<div class="card">

**Index** *(P1)*

the committed tree

</div>
<div class="card">

**Disclosure** *(P2)*

the tree's own cost model

</div>
<div class="card signal">

**Loop** *(P3)*

native file tools over the same tree

</div>
</div>

<p class="muted">Degrades, never hard-fails: no key → keyword mode · no tree → raw token overlap.</p>

Note: Build it as files and the other two patterns fall out almost for free —
that's why they compose. They're not three systems.

---

<p class="kicker">that's the kit</p>

# Before you reach for a vector store, ask what your agent can navigate

Move intelligence to **index time** · expose granularity, don't pick one · hand agents **corpora, not oracles**.

<p class="muted">The patterns outlive the product.</p>

Note: The mic-drop line — pause here. Appendix Q&A is next if asked; close on "Go further".

---

<p class="kicker">appendix · likely questions</p>

## Press ↓ for the ones this room will ask

<p class="muted">Why no embeddings? · cost & latency · staleness · why not the raw repo? · doesn't llms.txt solve this?</p>

--

### "Why no embeddings?"

Keyword + glossary, **no vector store** — on purpose.

- the glossary recovers a lot of synonym recall; the loop recovers paraphrase via query *reformulation*
- the index stays reviewable, edge-safe, with nothing to keep in sync
- embeddings are **deferred, not designed out** — a clean add as a recall/rerank upgrade

<p class="muted">Honest ceiling: paraphrase recall is the weak spot today.</p>

--

### "Cost & latency?"

- one bounded loop per submitted query, on a small model (Haiku), domain context **prompt-cached** across rounds
- worst case ≈ <code>maxIterations</code> round-trips; the keyword path is the instant fast lane
- the offline build uses Opus — but the hash gate means you pay only when content changes

--

### "Staleness?"

- the committed digest can drift from the docs
- a per-section **hash gate** + a runtime warning when the live hash ≠ the digest's
- intended workflow: rebuild in CI on every content change

--

### "Why not just give the agent the raw repo?"

Distillation is the point.

- bare source is noisy, unanchored, and huge
- the digest is small, **anchored**, fact-extracted, and **titled for the listing decision**
- it's the source *shaped for navigation*

--

### "Doesn't llms.txt already solve this?"

llms.txt asks the right question. The digest answers **what goes in it.**

- a **link-only** map → the agent still fetches and reads full pages
- **raw** `llms-full.txt` → context-stuffing at corpus scale; no anchors, no facts
- **the digest is the better llms.txt** — distilled, section-level, anchored

<p class="muted">Your <code>llms.txt</code> should <em>be</em> your digest — the map is rung 0, the full text is the distilled tree.</p>

---

<!-- .slide: data-background-color="#111111" class="invert" -->

<p class="kicker">thank you</p>

# Go further

<div class="cols cols-3">
<div class="card">

**hev ask**

⌘K search + an agent-readable docs digest.

[hevask.com](https://hevask.com) · open source

</div>
<div class="card">

**hev layer**

A transform runtime for your vector store — for when you *do* reach for one.

[hevlayer.com](https://hevlayer.com)

</div>
<div class="card signal">

**hev mind**

Trainings, workshops & consulting — AI, search & agentic engineering.

[hevmind.com](https://hevmind.com)

</div>
</div>

<p class="muted">Now over to Doug &amp; Trey for questions.</p>

Note: hev ask is what you just saw; hev layer is the runtime for when you *do*
go vector (callback to the close line); hev mind is how we work together —
trainings, workshops, consulting. Thank the room and hand to Doug & Trey.
