import { strict as assert } from 'node:assert';
import { test } from 'node:test';
import { readRuntime, resolveCfContext, resolveEnv } from '../src/runtime-env.ts';

/** Adapter 13 and earlier: plain values on `locals.runtime`. */
function legacyLocals(env: Record<string, string>, ctx?: unknown) {
  return { runtime: { env, ctx } };
}

/**
 * Adapter 14 (Astro v6+): `runtime` survives, but every property is a getter
 * that throws. Mirrors `@astrojs/cloudflare/dist/utils/cf-helpers.js`.
 */
function adapter14Locals(cfContext?: unknown) {
  const locals: Record<string, unknown> = { cfContext };
  Object.defineProperty(locals, 'runtime', {
    enumerable: false,
    value: {
      get env(): never {
        throw new Error('Astro.locals.runtime.env has been removed in Astro v6.');
      },
      get ctx(): never {
        throw new Error('Astro.locals.runtime.ctx has been removed in Astro v6.');
      },
    },
  });
  return locals;
}

test('runtime env is read from locals when the adapter still provides it', () => {
  assert.equal(resolveEnv(legacyLocals({ OPENROUTER_API_KEY: 'sk-or-legacy' })).OPENROUTER_API_KEY, 'sk-or-legacy');
});

test('a throwing locals.runtime getter does not take down the request', () => {
  // The whole point: `locals?.runtime?.env ?? {}` threw here, 500ing every
  // request — including keyword search, which needs no key at all.
  assert.doesNotThrow(() => resolveEnv(adapter14Locals()));
  assert.equal(readRuntime(adapter14Locals(), 'env'), undefined);
});

test('under adapter 14 the key comes from process.env instead', () => {
  process.env.HEV_ASK_TEST_KEY = 'sk-or-from-process';
  try {
    assert.equal(resolveEnv(adapter14Locals()).HEV_ASK_TEST_KEY, 'sk-or-from-process');
  } finally {
    delete process.env.HEV_ASK_TEST_KEY;
  }
});

test('locals.runtime.env outranks process.env where both exist', () => {
  process.env.HEV_ASK_TEST_KEY = 'from-process';
  try {
    const env = resolveEnv(legacyLocals({ HEV_ASK_TEST_KEY: 'from-runtime' }));
    assert.equal(env.HEV_ASK_TEST_KEY, 'from-runtime');
  } finally {
    delete process.env.HEV_ASK_TEST_KEY;
  }
});

test('execution context is found on either adapter generation', () => {
  const legacy = { waitUntil: () => {} };
  assert.equal(resolveCfContext(legacyLocals({}, legacy)), legacy);

  const modern = { waitUntil: () => {} };
  assert.equal(resolveCfContext(adapter14Locals(modern)), modern);
});

test('missing or non-Cloudflare locals resolve to no context', () => {
  assert.equal(resolveCfContext(undefined), undefined);
  assert.equal(resolveCfContext({}), undefined);
  assert.doesNotThrow(() => resolveEnv(undefined));
});
