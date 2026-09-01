# hev ask

A `⌘K` search overlay and an agent-readable docs digest, shipped as the npm
package [`@hevmind/ask`](https://www.npmjs.com/package/@hevmind/ask). hev ask
distills a docs site into a committed **ask digest** — a markdown tree your
coding agent can `tree`, `cat`, and `grep`, and that a `⌘K` overlay answers your
readers from.

**Using hev ask?** Everything you need is in **[`docs/`](docs)** — what it is,
how it compares, the five-minute [quick start](docs/quickstart.md), and the full
[API reference](docs/README.md#api-reference).

hevask.com is retired; it redirects here. The docs live in the repo now, and
`hevask.com/api/ask` no longer serves the hosted endpoint — point `--endpoint`
at your own site.

The rest of this file is for **contributing** to hev ask.

## What's in here

hev ask is a host-neutral core with a flagship Astro adapter:

- **Go core** (`github.com/hev/ask`, `cmd/ask`, `pkg/ask`) — the `ask` CLI and
  the offline digest builder. No Astro dependency; it works from markdown, not a
  renderer.
- **`@hevmind/ask`** (`packages/ui`) — the Astro integration, the `/api/ask`
  endpoint, the `SearchOverlay` component, and the `ask` bin (prebuilt Go
  binaries). The only published artifact.

Astro is wired end-to-end today. Other frameworks (Docusaurus, VitePress,
MkDocs) are designed in
[RFC 0004](docs/rfcs/0004-beyond-astro-host-neutral-adapters.md), not yet
shipped — don't describe a non-Astro adapter as working until it is.

## Repo layout

```
cmd/ask        # the `ask` CLI entrypoint (Go)
pkg/ask        # the Go core: digest build, read verbs, MCP
packages/ui    # @hevmind/ask — Astro integration, endpoint, SearchOverlay, bin
playground     # minimal Astro site for fast local dev of the package
docs           # the product docs, the RFCs, and talk material
```

It's a pnpm workspace with a Go module alongside. `packages/ui` is the only
published package; `playground` is its private consumer and the dogfooding
surface.

## Develop

```sh
pnpm install                          # workspace install
make build                            # build the ask binary into ./bin
make install                          # install ask into ~/.local/bin for iteration
pnpm dev                              # the playground site, for package dev
```

Add a provider key (`ANTHROPIC_API_KEY`, or `OPENAI_API_KEY` /
`OPENROUTER_API_KEY`) to a `.env` to exercise the agentic answer loop and
offline digest builds. Without one, everything degrades: search falls back to
keyword mode and the committed digest tree is used as-is.

## Test and check

`make check` runs the full gauntlet (Go vet and tests, the package tests, and
typecheck, which includes the playground's `astro check`). Or run the pieces:

```sh
make test                             # go test ./... + pnpm test
pnpm typecheck                        # tsc across the workspace + astro check
pnpm digest:verify                    # CI anchor gate
```

Run `ask digest verify` whenever anchors or chunking change — it's the CI gate
that catches deep-link drift.

## RFCs

Engineering changes that touch the design start as an RFC in
[`docs/rfcs`](docs/rfcs) — same process as hev layer. Align on the design before
writing code. Changes to the public surface update the matching page under
[`docs/`](docs) in the same PR.

## Releasing

`packages/ui` is published to npm as `@hevmind/ask`.

1. Set the semver in `packages/ui/package.json`.
2. `make npm-binaries` to cross-compile the per-platform binary packages.
3. `make check` and `pnpm digest:verify`.
4. Dry-run the tarball: `pnpm --filter @hevmind/ask pack --dry-run`.
5. Publish from `packages/ui`: `pnpm publish --access public`. **`pnpm`, never
   `npm`.** The optional deps use `workspace:*`, and only pnpm rewrites that to
   a pinned version on publish. `npm publish` uploads the literal string, and
   every install of that version dies with `EUNSUPPORTEDPROTOCOL` — the package
   is not partially broken, it is uninstallable. This ate 0.3.5. Verify after
   publishing: `npm view @hevmind/ask@X.Y.Z optionalDependencies` must show
   version numbers, not `workspace:*`.
6. Tag and cut the GitHub release: `git tag vX.Y.Z && git push --tags`, then
   `gh release create vX.Y.Z`.
7. Bump consumers (`playground/`, `../layer/site`) to the new version.
