# RFC 0007: Pluggable LLM tracing — one seam, OTLP as the bridge, PostHog kept native

## Summary

Make LLM tracing **pluggable** the way inference already is. Today the agentic
answer path emits traces to PostHog and only PostHog: `observability.ts` speaks
PostHog's capture API directly (`POST /i/v0/e/` with `$ai_generation` /
`$ai_span` / `$ai_trace` events) and `telemetryFromEnv` reads only `POSTHOG_*`.
Consumers who run their observability on Datadog, LangSmith, Braintrust,
Langfuse, Honeycomb, or anything else get nothing from the answer loop.

The abstraction seam is already in the right place — the loop and endpoint only
ever touch a small `Telemetry` interface (`generation` / `span` / `trace`), and
nothing PostHog-specific leaks past it. So this is a **refactor behind a stable
interface**, not a rewrite of the call sites. This RFC:

1. **Generalizes the seam into a tracer registry** modeled byte-for-byte on the
   inference `PROVIDERS` registry in `providers.ts`: a record of backends, each
   with its env keys and a `make(opts): Telemetry` factory, plus a `tracing`
   config option in `types.ts`. PostHog moves into `tracing/posthog.ts`
   unchanged. **Existing PostHog users see zero change.**
2. **Adds one OpenTelemetry (OTLP/HTTP + GenAI semantic conventions) adapter**
   as the universal bridge. The leverage is the whole point: **Datadog,
   LangSmith, and Braintrust — the three backends that motivated this — all
   already ingest OTLP GenAI spans**, each differing only by endpoint and auth
   header. One fetch-based emitter, configured by the standard `OTEL_EXPORTER_*`
   env vars, reaches all of them and the rest of the OTel ecosystem.
3. **Supports fan-out** — a composite sink so a site can ship to PostHog *and*
   an OTLP backend at once — which the interface makes nearly free.

Native per-vendor adapters (Datadog LLM Obs, LangSmith `/runs`, Braintrust logs)
are explicitly **deferred** to a later phase, written only if the OTLP mapping
proves too lossy for a given vendor's UI. We do not build N integrations when
one bridge covers the named targets.

This grows the public option surface, so per the public-surface rule it ships
with docs (`api/endpoint.mdx` "LLM tracing", `api/configuration.mdx`) in the
same change and a digest rebuild. Opened from the brainstorm in issue #1.

## Motivation

**Tracing is the one runtime integration that isn't pluggable.** Inference is
provider-pluggable (`provider: anthropic | openai | openrouter`, registry in
`providers.ts`, per-provider key env var). Tracing is hardcoded. A consumer who
already standardized their observability on Datadog or LangSmith can configure
their model provider freely but cannot see a single agentic trace unless they
also run PostHog. That asymmetry is arbitrary — and it's the most common reason
a serious docs team can't adopt the answer path with confidence.

**The seam is already cut; only the implementation is coupled.** The audit
behind this RFC is small and reassuring. The loop (`search/loop.ts`) and the
endpoint (`endpoint.ts`) touch only this:

```ts
interface Telemetry {
  readonly traceId: string;
  generation(event: GenerationEvent): void;
  span(event: SpanEvent): void;
  trace(event: TraceEvent): void;
}
```

`GenerationEvent` / `SpanEvent` / `TraceEvent` are vendor-neutral shapes (model,
tokens, latency, input/output, ok). Every PostHog-ism — the `$ai_*` property
names, the `/i/v0/e/` path, `$process_person_profile` — lives *inside*
`makeTelemetry`, and every env-ism lives inside `telemetryFromEnv`. Generalizing
the backend means swapping the implementation behind that interface and teaching
env-resolution to pick one. The ~dozen `telemetry.generation()` call sites in
the loop don't move.

**OTel turns "many backends" into "one adapter."** This is the load-bearing
discovery. The three backends named in issue #1 are not three integrations:

| Backend | OTLP traces endpoint | Auth / routing header |
| --- | --- | --- |
| Datadog | Agent OTLP intake (`/v1/traces`) or DD OTLP intake | DD agent / `DD-API-KEY` |
| LangSmith | `https://api.smith.langchain.com/otel/v1/traces` | `x-api-key`, `Langsmith-Project` |
| Braintrust | `https://api.braintrust.dev/otel/v1/traces` | `Authorization: Bearer`, `x-bt-parent` |

All three accept **OTLP/HTTP carrying GenAI-semantic-convention spans**, and so
do Langfuse, Honeycomb, Grafana, New Relic, Arize Phoenix, and Jaeger. A single
OTLP/HTTP-JSON emitter — still zero-dep, still `fetch`-only — reaches the entire
set; the only per-vendor difference is the endpoint URL and a header, which the
standard `OTEL_EXPORTER_OTLP_ENDPOINT` / `OTEL_EXPORTER_OTLP_HEADERS` env vars
already express. We get breadth for the cost of one adapter, not four.

**The constraints rule out the obvious approach.** `observability.ts`
deliberately speaks the capture API over `fetch` rather than importing
`posthog-node`, because the endpoint runs on Cloudflare and must stay
edge-friendly and dependency-free. That single fact disqualifies "just import
each vendor's SDK" — `@datadog/...`, `langsmith`, and `braintrust` are heavy,
Node-only, and would break the Worker. Every adapter has to be fetch-only. OTLP
over HTTP/JSON (a plain POST of a `ResourceSpans` payload) fits perfectly;
vendor SDKs don't. The constraint is what makes OTLP the right answer, not just
a convenient one.

## Goals

- **Tracing is selected like the inference provider** — an explicit `tracing`
  option plus env, defaulting to today's PostHog behavior, validated through a
  registry that mirrors `providers.ts`.
- **One OTLP/GenAI adapter unlocks Datadog, LangSmith, Braintrust, and the
  broader OTel ecosystem**, configured by standard `OTEL_EXPORTER_*` env vars,
  with zero new runtime dependencies.
- **PostHog stays native and byte-identical.** Its `$ai_*` format isn't OTel; we
  keep the working adapter and existing `POSTHOG_*` env vars verbatim.
- **Fan-out to multiple backends** from one request (PostHog *and* OTLP), since
  the interface makes a composite sink trivial and retrofitting it later is
  awkward.
- **Everything still degrades, nothing hard-fails.** No tracing configured →
  no-op sink; a misconfigured backend never breaks the answer stream.
- **Edge-friendly and dependency-free**, exactly as today: `fetch` only, no
  vendor SDKs, `waitUntil`-safe.
- **Docs-first.** The user-visible shape (the `tracing` option, the env vars,
  the per-vendor OTLP recipes) lands in the docs in the same change.

## Non-goals

- **No native Datadog / LangSmith / Braintrust adapters in this RFC.** OTLP
  covers all three. A native adapter is written later *only* if a vendor's OTLP
  mapping is too lossy for its UI (Phase 3), and only as a new registry entry —
  no core change.
- **No new trace data.** Same `generation` / `span` / `trace` events, same
  model/token/latency/content the loop already produces. This is a delivery
  refactor, not a richer instrumentation pass.
- **No vendor SDK dependencies, ever.** `@opentelemetry/*` included — OTLP/HTTP
  is JSON over POST; we hand-build the `ResourceSpans` payload, as we hand-build
  PostHog's today.
- **No change to the `captureContent` model.** `off | redacted | full` and the
  existing redaction of `tool_result` bodies carry over to every backend; OTLP
  maps it onto whether prompt/completion attributes are attached.
- **No client-side or build-time tracing.** This is the runtime answer loop
  only. The offline digest builder's model calls are out of scope.
- **No always-on collector requirement.** We POST OTLP/HTTP directly to a
  vendor's intake or a user-run collector; we do not bundle or require an
  OpenTelemetry Collector.

## Design

### Docs first: settle the option shape

Per the docs-first contract, the user-visible surface lands before code:

- `api/configuration.mdx` gains the `tracing` option (shape below).
- `api/endpoint.mdx`'s `## LLM tracing` section grows from "set `POSTHOG_KEY`"
  to "select a backend," with a recipe per backend (PostHog, OTLP-generic, and
  the Datadog / LangSmith / Braintrust endpoint+header trio).
- CLAUDE.md's "Key facts" note that tracing is backend-pluggable like inference.

Disagreement about the surface is resolved by editing those pages, not the TS.

### The tracer registry — mirror `providers.ts`

A new `packages/ui/src/tracing/` module with a registry shaped exactly like the
inference one:

```ts
// tracing/registry.ts
export interface TracerInfo {
  name: TracerName;                 // 'posthog' | 'otlp'
  label: string;                    // human label for logs/errors
  /** Env vars that, if present, auto-select this backend. */
  envKeys: string[];               // posthog: POSTHOG_KEY/POSTHOG_API_KEY
                                    // otlp: OTEL_EXPORTER_OTLP_ENDPOINT
  make(opts: TelemetryOptions): Telemetry;
}

export const TRACERS: Record<TracerName, TracerInfo> = { posthog, otlp };

export function resolveTracerName(value?: string): TracerName { /* default posthog */ }
```

`tracing/posthog.ts` holds today's `makeTelemetry` body verbatim — same `$ai_*`
events, same `/i/v0/e/`, same redaction. `observability.ts` keeps exporting the
`Telemetry` / `GenerationEvent` / `SpanEvent` / `TraceEvent` types (the seam)
and the no-op sink, and becomes the place that assembles a backend from config +
env. The loop and endpoint imports are unchanged.

### The `tracing` config option

A new field in `types.ts`, mirroring how `provider` reads:

```ts
interface HevAskOptions {
  // ...existing...
  /**
   * LLM tracing backend(s) for the agentic loop. A string selects one backend;
   * an array fans out to several. Omitted → auto-detect from env (POSTHOG_* →
   * posthog, OTEL_EXPORTER_OTLP_ENDPOINT → otlp); none present → no-op.
   * Secrets always come from env, never config.
   */
  tracing?: TracerName | TracerName[] | { backend: TracerName | TracerName[]; captureContent?: CaptureMode };
}
```

Selection precedence (resolved in `observability.ts`):

1. Explicit `tracing` option, if set.
2. Else **auto-detect**: each registry entry whose `envKeys` are present is
   selected. `POSTHOG_KEY` alone keeps doing exactly what it does today;
   `OTEL_EXPORTER_OTLP_ENDPOINT` alone selects OTLP; both present → fan out to
   both.
3. Else → no-op sink.

Auto-detect (vs the explicit-with-default that `providers.ts` uses) is the right
call for tracing specifically: it's optional, secondary, and almost always
keyed off "is the secret present," so a consumer who sets the standard OTel env
var gets tracing with no config edit. Explicit `tracing` is the escape hatch
when env is ambiguous or fan-out needs narrowing.

### The OTLP / GenAI adapter

`tracing/otlp.ts` emits OTLP/HTTP-JSON spans following the OpenTelemetry **GenAI
semantic conventions**. Mapping from our events:

- `generation` → a span, kind `CLIENT`, with
  `gen_ai.operation.name=chat`, `gen_ai.system` (the inference provider),
  `gen_ai.request.model`, `gen_ai.usage.input_tokens`,
  `gen_ai.usage.output_tokens`, and `gen_ai.response.*`. Prompt/completion go on
  span events (`gen_ai.system.message` / `gen_ai.choice`) or
  `gen_ai.prompt`/`gen_ai.completion` attributes — gated by `captureContent`
  exactly as PostHog gates `$ai_input`/`$ai_output_choices`.
- `span` → a generic internal span, parent/child via span IDs.
- `trace` → the root span; its end is when we flush.

Three mechanical concerns this adapter owns, all small and self-contained:

- **ID format.** OTel needs a 16-byte hex trace ID and 8-byte hex span IDs; we
  mint UUIDs (PostHog-friendly). The adapter derives OTel IDs (hash/truncate the
  UUIDs, or mint native ones and keep a map) so the seam's `traceId` contract is
  unchanged for callers.
- **Timestamps.** OTLP wants start/end in nanoseconds; our events carry
  `latencyMs`. We synthesize `endTimeUnixNano = now`,
  `startTimeUnixNano = now − latency`. (`Date.now()` is available at runtime.)
- **Batching.** PostHog emits one POST per event; OTLP buffers spans in closure
  state and flushes **one** `ResourceSpans` POST on `trace()` (the root, always
  last), handed to `waitUntil`. One request → one OTLP export.

Endpoint and auth come from the standard env vars — `OTEL_EXPORTER_OTLP_ENDPOINT`
and `OTEL_EXPORTER_OTLP_HEADERS` (and `OTEL_EXPORTER_OTLP_PROTOCOL`, of which we
support `http/json`). That is exactly enough to point at Datadog's agent intake,
LangSmith's `/otel`, or Braintrust's `/otel` — the per-vendor difference is data,
not code.

### Fan-out (composite sink)

When selection yields more than one backend, `observability.ts` wraps them in a
composite `Telemetry` whose `generation`/`span`/`trace` call each child in turn
(each child still fire-and-forget + `waitUntil`-guarded). `traceId` is shared
across children so a single request correlates everywhere. One sink per backend
failing never affects the others or the answer stream.

### Env vars, end to end

| Backend | Selected by | Config |
| --- | --- | --- |
| PostHog (native) | `POSTHOG_KEY` / `POSTHOG_API_KEY` | `POSTHOG_HOST`, `POSTHOG_CAPTURE_CONTENT` (unchanged) |
| OTLP (generic) | `OTEL_EXPORTER_OTLP_ENDPOINT` | `OTEL_EXPORTER_OTLP_HEADERS`, `OTEL_EXPORTER_OTLP_PROTOCOL` |
| → Datadog | OTLP endpoint = agent intake | `DD-API-KEY` (+ site) via OTLP headers |
| → LangSmith | OTLP endpoint = `…/otel/v1/traces` | `x-api-key`, `Langsmith-Project` via OTLP headers |
| → Braintrust | OTLP endpoint = `…/otel/v1/traces` | `Authorization: Bearer …`, `x-bt-parent` via OTLP headers |

`HEV_ASK_TRACING=posthog|otlp|posthog,otlp` is the explicit override when
auto-detect is ambiguous, mirroring the `tracing` option for the env-only case.

## Consequences

- **The answer loop traces to wherever the team already looks.** Datadog,
  LangSmith, and Braintrust users adopt the agentic path without standing up
  PostHog — the asymmetry with inference providers is gone.
- **One adapter, broad reach.** The OTLP path covers the named targets *and* the
  rest of the OTel ecosystem, so future "can it trace to X?" is usually "yes, if
  X speaks OTLP" — no new code.
- **Public surface grows** by the `tracing` option, the `OTEL_EXPORTER_*` /
  `HEV_ASK_TRACING` env vars, and the `tracing/` module's backend names. Additive
  → **minor version bump** with a doc update under the public-surface rule.
- **Two trace shapes to keep correct.** PostHog `$ai_*` and OTLP GenAI spans
  must both stay faithful as the loop's events evolve. Contract tests for each
  (asserting the emitted payload shape against a captured `fetch`) are the guard.
- **More observability docs to own** — a recipe per backend, including the
  endpoint+header trio for the OTLP vendors. The table above is the honest scope.
- **No new dependencies, no edge regressions.** The OTLP adapter is `fetch` +
  JSON, like PostHog; the Worker bundle and cold-start are unaffected.

## Migration

Fully additive. The default behavior is unchanged: a site with `POSTHOG_KEY`
set and no `tracing` option traces to PostHog exactly as before, byte-for-byte;
a site with nothing set is a no-op, as before. New behavior appears only when a
consumer sets `OTEL_EXPORTER_OTLP_ENDPOINT`, adds the `tracing` option, or sets
`HEV_ASK_TRACING`. Existing `observability.test.ts` cases stay green; new tests
cover registry selection, auto-detect, fan-out, and the OTLP payload shape.
Ships as **v0.4.0** (minor: additive option + env + adapter), following 0.3.0
(RFC 0004). No digest re-distillation needed beyond a normal rebuild for the
doc edits.

## Sequencing

1. **Docs PR (working-backward).** `tracing` option in `api/configuration.mdx`;
   rewrite `## LLM tracing` in `api/endpoint.mdx` with the backend table and
   per-vendor recipes; CLAUDE.md note. Settles the surface.
2. **Registry + PostHog extraction (no behavior change).** Create `tracing/`,
   move PostHog into `tracing/posthog.ts`, add `TRACERS` + `resolveTracerName`,
   generalize `telemetryFromEnv` to select via the registry, add the `tracing`
   option to `types.ts`. Auto-detect defaults preserve current behavior; tests
   stay green. This is the load-bearing refactor — everything else is an entry.
3. **OTLP/GenAI adapter.** `tracing/otlp.ts`: ID/timestamp normalization,
   GenAI-convention span mapping, buffer-and-flush on `trace()`, standard
   `OTEL_EXPORTER_*` env. Contract test asserting the `ResourceSpans` payload.
4. **Fan-out.** Composite sink; selection returns a list; multi-backend test.
5. **Verify against real backends.** Point the OTLP adapter at LangSmith and
   Braintrust `/otel` endpoints (and a Datadog agent) with a throwaway key; confirm
   spans land and read sensibly. This spike answers whether Phase 3 is needed.
6. **Ship.** `pnpm test && pnpm typecheck && pnpm --filter hev-ask-site check`,
   `ask digest verify` on `site/`, minor bump, deploy.

Steps 2–4 are one release (**v0.4.0**); step 5 decides whether a Phase 3
(native Datadog/LangSmith/Braintrust adapters) is ever scheduled — written only
where OTLP fidelity proves insufficient, each as a new registry entry.

## Open questions

- **Auto-detect vs explicit-default.** This RFC proposes auto-detect-from-env
  (set the OTel var, get tracing) because tracing is keyed off secret presence;
  `providers.ts` uses explicit-with-default instead. Is the inconsistency with
  the inference registry worth the lower-friction onboarding? Leaning yes —
  tracing is optional and secondary in a way provider selection isn't.
- **Is OTLP genuinely enough for all three named vendors today,** or does any of
  Datadog / LangSmith / Braintrust need its native format to light up the good
  parts of its UI (e.g. Braintrust eval scoring, LangSmith run trees, Datadog LLM
  Obs clustering)? Step 5 is the spike that answers this and decides whether
  Phase 3 exists.
- **ID strategy.** Derive OTel IDs by hashing our UUIDs (stable, one-way) vs mint
  native OTel IDs and expose them as the seam's `traceId` (cleaner OTel, changes
  what PostHog sees as the trace id). Leaning hash-derive so PostHog and OTLP
  share a legible correlation id.
- **`gen_ai.system` value.** GenAI conventions expect a known system string
  (`anthropic`, `openai`); we pass our `provider` name, which mostly aligns —
  but `openrouter` isn't a GenAI `system`. Map it to the underlying model's
  vendor, or emit `openrouter` and accept a non-standard value?
- **Prompt/completion encoding under OTel.** Span events
  (`gen_ai.system.message` / `gen_ai.choice`, the newer convention) vs
  `gen_ai.prompt`/`gen_ai.completion` attributes (OpenLLMetry, wider current
  support). Pick by what the named vendors actually render today; revisit as the
  convention stabilizes.
- **Fan-out scope.** Support an arbitrary list of backends, or cap at "one native
  + one OTLP"? Arbitrary is simplest to reason about and the interface is
  indifferent; capping only buys a marginal guardrail.
- **OTLP transport.** `http/json` only (simplest, dependency-free) vs also
  `http/protobuf` (smaller, more universal collector support but needs protobuf
  encoding). Lean `http/json` first; most intakes accept it.

## References

- Issue #1 — the brainstorm this RFC formalizes (current PostHog coupling, the
  native-vs-OTel decision, the open questions).
- `packages/ui/src/observability.ts` — the `Telemetry` seam + PostHog
  implementation + `telemetryFromEnv`; the file this RFC factors into `tracing/`.
- `packages/ui/src/search/loop.ts` — the only `telemetry.generation()` call
  sites; unchanged by this RFC.
- `packages/ui/src/endpoint.ts` (`resolveTelemetry`) — the wiring and the
  Cloudflare `waitUntil` plumbing every adapter relies on.
- `packages/ui/src/providers.ts` — the inference registry pattern the tracer
  registry mirrors (`PROVIDERS`, `envKey`, `resolveProviderName`).
- `packages/ui/src/types.ts` — where the `tracing` option lands beside
  `provider`.
- `packages/ui/test/observability.test.ts` — the existing PostHog contract tests
  that must stay green and the model for the OTLP/fan-out tests.
- `site/src/content/docs/api/endpoint.mdx` (`## LLM tracing`),
  `site/src/content/docs/api/configuration.mdx` — the docs updated in the same
  change.
- RFC 0001 (`0001-embeddable-ask-command.md`) — the docs-first contract and the
  public-surface rule this RFC follows.
