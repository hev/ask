import assert from 'node:assert/strict';
import test from 'node:test';
import { EMPTY_DIGEST, normalizeDigest } from '../src/digest/schema.ts';

test('normalizeDigest returns an empty v2 digest for junk', () => {
  assert.deepEqual(normalizeDigest(null), EMPTY_DIGEST);
  assert.deepEqual(normalizeDigest('nope'), EMPTY_DIGEST);
  assert.equal(EMPTY_DIGEST.version, 2);
});

test('a v1 artifact degrades to a node-less v2 digest (keeps context + glossary)', () => {
  const v1 = {
    version: 1,
    contentHash: 'abc',
    context: 'Layer concepts.',
    glossary: [{ term: 'pipeline', aliases: ['stream'], definition: 'A flow.' }],
  };
  const digest = normalizeDigest(v1);
  assert.equal(digest.version, 2);
  assert.equal(digest.context, 'Layer concepts.');
  assert.equal(digest.glossary.length, 1);
  assert.deepEqual(digest.nodes, []);
  assert.equal(digest.overview, '');
  assert.deepEqual(digest.suggestions, [], 'a digest without suggestions normalizes to none');
});

test('suggestions normalize to non-empty strings only', () => {
  const digest = normalizeDigest({
    version: 2,
    suggestions: ['How does it work?', '', '  ', 42, null, 'What are the limits?'],
  });
  assert.deepEqual(digest.suggestions, ['How does it work?', 'What are the limits?']);
});

test('v2 nodes round-trip and bad fields are coerced', () => {
  const digest = normalizeDigest({
    version: 2,
    nodes: [
      {
        id: 'concepts#autoscaling',
        kind: 'section',
        title: 'Concepts',
        heading: 'Autoscaling',
        group: 'Overview',
        url: '/docs/concepts#autoscaling',
        summary: 'Scales workers from lag.',
        facts: [{ kind: 'flag', literal: '--max-workers', chunkId: 'concepts#autoscaling' }],
        sources: [{ chunkId: 'concepts#autoscaling', url: '/docs/concepts#autoscaling', anchor: 'autoscaling' }],
        mode: 'bogus-mode',
        terms: ['autoscaling', 'workers'],
      },
      { kind: 'section' }, // no id/url → dropped
    ],
  });
  assert.equal(digest.nodes.length, 1);
  const node = digest.nodes[0];
  assert.equal(node.mode, 'agent-primary', 'invalid mode coerced to default');
  assert.equal(node.facts[0].literal, '--max-workers');
  assert.equal(node.sources[0].anchor, 'autoscaling');
  assert.deepEqual(node.terms, ['autoscaling', 'workers']);
});
