---
name: build-digest
description: >-
  Build the @hevmind/ask ask digest (the committed .hev-ask/ markdown tree) for a
  docs site using your own agent subscription instead of a provider API key. Use
  when asked to build, rebuild, or refresh the hev ask digest, search index, or
  knowledge graph (KG), or after docs content changes. Shards the corpus, distils
  each shard in a fresh context, then assembles and verifies.
---

# Build the hev ask digest

`@hevmind/ask`'s agentic loop, keyword ranking, and suggested questions run off a
committed digest tree at `.hev-ask/`. Only the **distillation** needs a model —
the CLI computes the node structure, verbatim facts, overview map, and content
hashes deterministically. This skill does that distillation here in your
subscription, so it costs **no provider API tokens**.

`ask` is the `@hevmind/ask` binary: install it on PATH, or it resolves via the
package bin / `HEV_ASK_BINARY` (see `api/cli.mdx`). Run every command from the
**site root** (the dir whose config registers `hevAsk()`). If the integration
overrides any content flags (`--collection`, `--base-path`,
`--chunk-heading-depth`, `--content-glob`, `--digest-dir`), pass the same ones to
`corpus` and `assemble` — they must match.

The corpus splits into ~200KB **shards** (`--shard-bytes` tunes it), each
distilled in its own context, so corpus size is never a context limit. State
lives in `.hev-ask/shards/`, so the build resumes and refreshes incrementally:
after edits, only changed shards re-distil.

**Never read a shard input file yourself** — they hold the full corpus text.
Work from command output and `status`; the per-shard agents read the shards.

## Steps

1. **Shard.** `ask digest corpus --shards-dir .hev-ask/shards`
   Deterministic and keyless. Reports `(N sections, M shards, P pending,
   up-to-date|needs-rebuild)`. Safe to re-run; this is the refresh mechanism.

2. **Check.** If corpus said `up-to-date` AND `0 pending`, the digest already
   matches the content — **stop and tell the user nothing needs rebuilding.**
   Otherwise `ask digest status --shards-dir .hev-ask/shards` lists the
   `pending`/`stale` shards (both need distilling).

3. **Distil each pending/stale shard in a fresh context.** Spawn one agent per
   shard (a few in parallel; don't read shards yourself). Give each this prompt
   with `<id>` filled in:

   > Read ONLY `.hev-ask/shards/input-<id>.json` (from the site root): it has
   > `shardId`, `shardHash`, and a `sections` array of `{ id, url, title, text }`.
   > Write `.hev-ask/shards/distill-<id>.json` with exactly this shape:
   >
   > ```json
   > {
   >   "shardHash": "<copy shardHash verbatim>",
   >   "notes": "5-10 lines: what this shard covers, its key concepts, and how users phrase them.",
   >   "glossary": [{ "term": "...", "aliases": ["..."], "definition": "One line." }],
   >   "summaries": [{ "id": "<exact section id>", "summary": "1-3 sentences." }]
   > }
   > ```
   >
   > - One `summaries` entry for every `id` in `sections` — exact ids, no more, no fewer.
   > - Summaries are what the search agent reasons from: faithful, self-contained.
   >   Paraphrase prose; **never restate code, flags, commands, or identifiers** —
   >   the CLI extracts those verbatim, and they'd only drift if retyped.
   > - `glossary`: ≤10 terms a real user would actually type (aliases like
   >   `k8s`→`kubernetes`); one-line definitions. The CLI dedupes and caps them.
   > - `notes` is not user-facing; it feeds the global synthesis pass.
   > - Reply with just the shard id and how many summaries you wrote.

   Interrupted? Re-run from step 1 — disk is the source of truth and `status`
   shows what's left.

4. **Synthesize the global context.** Once every shard is distilled, read the
   `notes` field from each `.hev-ask/shards/distill-*.json` (small — never the
   full files) and write `.hev-ask/shards/global.json`:

   ```json
   {
     "context": "Compact markdown orientation: what the product/site is, its core concepts and feature areas, and how users talk about them.",
     "suggestions": ["A natural question a reader might type that these docs answer."],
     "glossary": []
   }
   ```

   `suggestions`: 3-5 questions phrased the way a user would ask them, each
   genuinely answerable from these docs (they show in the overlay on open).

5. **Assemble.** `ask digest assemble --input-dir .hev-ask/shards`
   Merges every shard distillation with the global synthesis, derives the
   deterministic parts, and writes the `.hev-ask/` tree. Undistilled shards fall
   back to excerpts and are reported — usable mid-wave, but aim for 0 pending
   before committing.

6. **Verify.** `ask digest verify` — anchor drift is fatal; coverage/fidelity
   warnings are informational.

7. **Commit.** The shards dir is the local resume/refresh cache — keep it on disk
   but out of git, and drop the bulky input files (`corpus` regenerates them):

   ```sh
   rm -f .hev-ask/shards/input-*.json
   git check-ignore -q .hev-ask/shards || echo ".hev-ask/shards/" >> .gitignore
   git add .gitignore .hev-ask
   ```

   Only the `.hev-ask/` tree is committed; `.hev-ask/shards/` stays local.

## Notes

- A small site may produce a single shard — distil it yourself instead of
  spawning an agent (its input is small enough to read directly).
- If `corpus` finds no content, you're likely not in the site root, or the
  collection isn't named `docs` — pass `--collection <name>`.
