// Reading the Cloudflare request context off `Astro.locals` is version-sensitive
// enough to be worth isolating.
//
// `@astrojs/cloudflare` 13 and earlier exposed `locals.runtime.env` and
// `locals.runtime.ctx` as plain values. Adapter 14 (Astro v6+) removed both and
// left throwing getters in their place, so `locals?.runtime?.env ?? {}` does not
// degrade to the fallback — it throws out of the whole request. Every read here
// is therefore guarded, and the adapter-14 replacements are consulted next.

/** Execution context shape the endpoint needs — just `waitUntil`. */
export interface CfContext {
  waitUntil?: (promise: Promise<unknown>) => void;
}

/**
 * Reads `locals.runtime[key]`, absorbing the throwing getters adapter 14 leaves
 * behind. Returns undefined on any adapter that no longer supplies the value.
 */
export function readRuntime<T>(locals: unknown, key: 'env' | 'ctx'): T | undefined {
  try {
    return (locals as { runtime?: Record<string, unknown> })?.runtime?.[key] as T | undefined;
  } catch {
    return undefined;
  }
}

/**
 * Merges every environment the endpoint may run under, most specific last:
 * build-time `import.meta.env`, then `process.env` (Node adapters, and
 * Cloudflare under `nodejs_compat`, where the worker's bindings land there),
 * then Cloudflare's per-request `locals.runtime.env` on adapters that still
 * provide it.
 */
export function resolveEnv(locals: unknown): Record<string, string | undefined> {
  const fromRuntime = readRuntime<Record<string, string>>(locals, 'env') ?? {};
  const fromProcess = (typeof process !== 'undefined' ? process.env : undefined) ?? {};
  const fromImportMeta = (import.meta as { env?: Record<string, string> }).env ?? {};
  return { ...fromImportMeta, ...fromProcess, ...fromRuntime };
}

/**
 * The execution context, from `locals.runtime.ctx` (adapter 13 and earlier) or
 * `locals.cfContext` (adapter 14). Undefined off Cloudflare, where the caller
 * simply awaits instead of handing work to `waitUntil`.
 */
export function resolveCfContext(locals: unknown): CfContext | undefined {
  return readRuntime<CfContext>(locals, 'ctx') ?? (locals as { cfContext?: CfContext })?.cfContext;
}
