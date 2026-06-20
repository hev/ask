import assert from 'node:assert/strict';
import test from 'node:test';
import { expandQueryTerms } from '../src/digest/expand.ts';
import type { GlossaryEntry } from '../src/digest/schema.ts';

const glossary: GlossaryEntry[] = [
  { term: 'OIDC', aliases: ['openid connect', 'okta', 'azure ad'], definition: '' },
  { term: 'external authentication', aliases: ['git auth', 'oauth provider'], definition: '' },
  { term: 'kubernetes', aliases: ['k8s'], definition: '' },
];

test('expands on a full single-token term and skips noisy short tokens', () => {
  const terms = new Set(expandQueryTerms('oidc setup', glossary));
  assert.ok(terms.has('openid') && terms.has('connect') && terms.has('okta'), 'adds OIDC synonyms');
  assert.ok(!terms.has('ad'), 'drops 2-char tokens like "ad"');
});

test('does NOT expand on a single shared common token', () => {
  // "authentication" alone must not pull in the "external authentication" entry.
  const terms = new Set(expandQueryTerms('oidc authentication', glossary));
  assert.ok(!terms.has('external') && !terms.has('git') && !terms.has('provider'), 'no spurious expansion');
  assert.ok(terms.has('openid'), 'still expands OIDC via its full term');
});

test('expands a multi-word term only when all its tokens are present', () => {
  assert.ok(new Set(expandQueryTerms('external authentication for git', glossary)).has('oauth'), 'full phrase matches');
  assert.ok(!new Set(expandQueryTerms('external workspaces', glossary)).has('oauth'), 'partial phrase does not');
});

test('expands an alias (k8s -> kubernetes)', () => {
  assert.ok(new Set(expandQueryTerms('k8s autoscaling', glossary)).has('kubernetes'));
});
