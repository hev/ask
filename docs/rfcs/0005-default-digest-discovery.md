# RFC 0005: Default digest discovery — env var, then walk-up

## Summary

Make `ask` **find the digest by default**, the way `git`, `cargo`, and `eslint`
find their roots. Today the read verbs (`tree`, `ls`, `head`, `cat`, `facts`,
`grep`) resolve `--digest-dir` to a literal `.hev-ask` in the *exact* current
directory — so `ask tree` works from `site/` but finds nothing from the repo
root or anywhere else, even though the digest is right there one level up. This
RFC adds a resolution ladder for the read/query path:

```
--digest-dir flag  >  $HEV_ASK_DIGEST_DIR  >  nearest ancestor .hev-ask/ (up to $HOME)  >  ./.hev-ask
```

The explicit flag still always wins, so nothing that works today changes. The
two new rungs are additive: an environment default lets you point `ask` at one
digest from anywhere (dogfooding the hev ask docs from `~`, say), and the
walk-up makes the CLI "just work" inside any docs repo without a flag. The
producer verbs (`digest build`, `verify`, `corpus`, `assemble`) keep their
cwd-relative semantics — you build where your content lives.

## Motivation

**The cwd-literal default is the load-bearing papercut.** The docs already
*claim* the better behavior — `api/cli.mdx` says "By default, `ask` reads the
`.hev-ask/` tree from the current repo" — but the code only checks the current
directory, not the repo. The first thing anyone does after `ask digest build`
is `cd` somewhere and type `ask tree`, get an empty/erroring result, and learn
they have to either stand in the one blessed directory or pass `--digest-dir`
every time. For a tool whose whole pitch is "a directory your agent can
intuitively `cat` and `grep`," needing a path flag to find that directory is the
wrong first impression.

**Repo-scoped tools walk up; this one should too.** `git status`, `cargo build`,
`eslint .`, `rg` all resolve their root by walking ancestors. It's the
unsurprising convention, and it's exactly what "from the current repo" already
promises. Walk-up costs a handful of `stat`s and removes the most common reason
to reach for a flag.

**There's no way to set a default at all.** The CLI reads zero environment
variables today. A maintainer who wants `ask` to default to one digest
everywhere — the canonical dogfooding case for this very repo — has no option
but a shell alias that re-types `--digest-dir /abs/path` on every invocation. An
env var is the standard answer and stays host-neutral: nothing hardcodes "ask
docs," the operator chooses the target.

## Goals

- `ask tree` (and the other read verbs) work from anywhere inside a docs repo
  with no flag, by finding the nearest ancestor `.hev-ask/`.
- A `HEV_ASK_DIGEST_DIR` env var sets the default digest for the read path,
  overridable per-call by `--digest-dir`.
- **Zero behavior change** for any invocation that already passes `--digest-dir`
  or already runs with `.hev-ask/` in cwd.
- Resolution is explainable in one line and discoverable (`ask` surfaces which
  digest answered when asked).

## Non-goals

- **No change to the producer verbs.** `digest build`/`verify`/`corpus`/
  `assemble` keep cwd-relative `.hev-ask` and their own `--out`/`--input`/
  `--shards-dir` defaults. Building is an in-place act tied to your content
  config; it should not wander up the tree.
- **No project config file.** No `.askrc`, no `ask.toml`. The env var plus
  walk-up covers the cases; a config file is a separate, larger decision.
- **No remote/registry resolution.** Discovery is local filesystem only.
- **No multi-digest merging.** Walk-up stops at the first `.hev-ask/` it finds.
- **No cross-digest routing.** Resolving *one* digest is this RFC; asking
  *across* all local digests and letting the model pick the right one is RFC
  0006, which consumes the `which`/registry surface exposed here.

## Design

**One resolution function, applied at the read/query load path.** In
`parseGlobalFlags` (`cmd/ask/main.go`) track whether `--digest-dir` /
`--digest-path` was set explicitly rather than collapsing to the default string
immediately. When a read/query verb is about to `LoadDigest` (`pkg/ask/local.go`)
and no explicit flag was given, resolve in order:

1. **`$HEV_ASK_DIGEST_DIR`**, if set and non-empty. Used verbatim (relative to
   cwd if relative). No existence pre-check — a bad value should surface the same
   "digest not found" error a bad `--digest-dir` does, so the misconfiguration is
   legible.
2. **Walk-up**: from cwd, ascend parent by parent looking for a directory named
   `.hev-ask`; the nearest (closest to cwd) hit wins. Ascend **up to and
   including `$HOME`, then stop** — never resolve a digest above your home
   directory (system paths, other users). `$HOME/.hev-ask`, if present, is thus
   the last rung walk-up checks: a natural **global-default digest**.
3. **`./.hev-ask`** — today's literal fallback, preserved so a present-cwd digest
   resolves identically to now (walk-up finds it on the first iteration anyway;
   this rung just pins the error message to cwd when nothing is found at all).

The flag, when present, short-circuits the whole ladder — explicit always wins,
including the legacy `--digest-path` alias.

**Scope: the load path, not the build path.** Only verbs that *read an existing*
digest get discovery: `tree`, `ls`, `head`, `cat`, `facts`, `grep`, and the
agentic `answer`/`mcp` load. The producer verbs already default their paths
independently (`--out .hev-ask/digest-input.json`, etc.) and are untouched, which
is what keeps "build in place" intact.

**Multi-digest is the expected case, not an error.** Raising the ceiling to
`$HOME` makes it likely that several `.hev-ask/` trees live under it — one per
docs repo, plus an optional `$HOME/.hev-ask` global. Resolution stays
deterministic *by construction*: the walk visits one level at a time and takes
the **nearest** match, so there is always exactly one answer or none — never a
silent pick among equals. The graceful handling is about visibility and choice,
three ways:

- **Always say which digest answered.** When discovery or the env var (not an
  explicit flag) selected the directory, `--json` output carries the resolved
  absolute path, and human output notes it on stderr whenever it isn't
  `./.hev-ask`. With many digests around, "which one answered this?" must never
  be a guess.
- **A `which` verb to see the whole picture — and feed the model.** `ask digest
  which` prints the resolved digest plus the other candidates found along the
  walk (and at `$HOME`). Its `--json` form is the machine-readable seed of the
  digest **registry** RFC 0006 builds on, so a model can enumerate every local
  digest and route across them — `which` is designed to be resolvable by the
  loop, not only read by a human.
- **Never guess when genuinely ambiguous.** The one ambiguous shape is standing
  at a monorepo root with sibling digests *below* you (`docs-a/.hev-ask`,
  `docs-b/.hev-ask`) and none above. Walk-up is ancestors-only, so it correctly
  resolves to *none* — but instead of a bare "not found," the error runs a
  shallow descendant scan and lists the candidates: "found 2 digests below
  (docs-a, docs-b); pass `--digest-dir` or set `HEV_ASK_DIGEST_DIR`." The
  descendant scan feeds diagnostics only — it never silently resolves.

**Walk-up cost.** Bounded by directory depth; a `stat` of `<dir>/.hev-ask` per
level. Same order of work as `git`'s root search. No caching needed.

## Consequences

- **The docs become true.** "Reads the `.hev-ask/` tree from the current repo"
  finally matches the code.
- **New public surface:** one env var (`HEV_ASK_DIGEST_DIR`) and a documented
  resolution order. Additive, non-breaking → **minor version bump**, not major.
- **`$HOME` as feature and footgun.** A `$HOME/.hev-ask` becomes an implicit
  global default that answers from inside any repo without its own digest — handy
  as a personal fallback, surprising if forgotten. The mandatory "resolved to …"
  surfacing and `ask digest which` keep it legible; the env var and flag override
  it. Walk-up never climbs above `$HOME`, so a stray digest in a system path or
  another user's tree can't be picked up.
- **MCP/agentic load** inherits the same resolution for free, since they share
  `LoadDigest`.

## Migration

None required. Every current invocation keeps working:

- `ask --digest-dir X …` — flag wins, unchanged.
- `ask tree` with `.hev-ask/` in cwd — walk-up finds it first iteration,
  identical result.
- `ask tree` with no `.hev-ask/` anywhere — same "not found" error, now phrased
  against the resolved/cwd path.

The dogfooding setup this unlocks, documented in `api/cli.mdx`:

```sh
export HEV_ASK_DIGEST_DIR="$HOME/workspace/hev/ask/site/.hev-ask"
ask tree            # answers from the ask docs, from anywhere
```

## Sequencing

1. Resolution function + `explicitlySet` plumbing in `cmd/ask/main.go`; unit
   tests for each rung and precedence (`cmd/ask/main_test.go`).
2. Wire it into the read/query load path in `pkg/ask/local.go`; stderr/`--json`
   "resolved to" surfacing.
3. Docs: `api/cli.mdx` resolution-order section + `HEV_ASK_DIGEST_DIR` row;
   update the line-22 default description.
4. `make check` (vet, `go test`, pnpm test, typecheck, site check); minor bump
   in `packages/ui/package.json`.

## Open questions

- **Ceiling = `$HOME` (decided).** Walk-up ascends to and includes `$HOME`, then
  stops; `$HOME/.hev-ask` is the global-default rung. Follow-ons: what when
  `$HOME` is unset or cwd is *outside* it (a daemon in `/srv`, CI in `/tmp`)? —
  fall back to a filesystem-root-bounded walk-up, or refuse discovery and require
  the env var? Leaning: bounded fallback, since refusing would regress today's
  cwd-`.hev-ask` case.
- **`which` verb surface.** A dedicated `ask digest which`, or fold the resolved
  path + candidate list into existing `--json`/`facts` output? And how deep is
  the ambiguity-error descendant scan — depth 1–2 to stay cheap, or full subtree?
- **Does `HEV_ASK_DIGEST_DIR` ever seed the producer verbs?** Proposed: no —
  read path only, builds stay explicit. Revisit if dogfooders ask for it.
- **Is the "resolved to …" note always-on or `--verbose`?** Made mandatory above
  for legibility under multi-digest; if it proves noisy for power users, gate the
  human-stderr line behind `--verbose` while keeping the path in `--json` always.

## References

- `cmd/ask/main.go` — `parseGlobalFlags`, where `digestPath` defaults to
  `.hev-ask`; the resolution ladder and `explicitlySet` flag land here.
- `pkg/ask/local.go`, `pkg/ask/tree.go` — `LoadDigest`, the read/query load path
  discovery wraps.
- `site/src/content/docs/api/cli.mdx` — line 22 default description and the
  `--digest-dir` flag-table row (152); the resolution-order + env-var docs.
- `packages/ui/bin/ask-launcher.mjs` — prior art for env-driven resolution
  (`HEV_ASK_BINARY`); `HEV_ASK_DIGEST_DIR` mirrors its precedence shape.
- RFC 0002 (`0002-digest-as-filesystem.md`) — established the `.hev-ask/` tree
  and the read-verb surface this RFC makes discoverable.
- RFC 0006 (follow-on, `0006-ask-across-local-digests.md`) — consumes the
  `which`/registry this RFC exposes to route a question to the right digest
  automatically.
- `CLAUDE.md` — the `--digest-dir` flag and "committed markdown tree" facts this
  RFC extends without breaking.
