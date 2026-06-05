import assert from 'node:assert/strict';
import test from 'node:test';
import { prefilter } from '../src/search/prefilter.ts';
import { tokenize, type Chunk } from '../src/search/chunk.ts';
import type { DigestNode } from '../src/digest/schema.ts';

function chunk(id: string, docSlug: string, text: string): Chunk {
  return {
    id,
    docSlug,
    docTitle: docSlug,
    heading: id,
    anchorId: id.split('#')[1],
    url: `/docs/${docSlug}#${id.split('#')[1]}`,
    text,
    raw: text,
    tokens: new Set(tokenize(`${docSlug} ${text}`)),
  };
}

function node(id: string, terms: string[], summary = ''): DigestNode {
  return {
    id,
    kind: 'section',
    title: id,
    heading: id,
    group: null,
    url: `/docs/${id}`,
    summary,
    facts: [],
    sources: [],
    mode: 'agent-primary',
    terms,
  };
}

const chunks = [
  chunk('a#one', 'a', 'autoscaling is mentioned here once in passing'),
  chunk('b#two', 'b', 'autoscaling is mentioned here once in passing'),
];

test('without nodes, prefilter is plain token overlap (ties break by id)', () => {
  const out = prefilter(chunks, 'autoscaling', [], 5, 2);
  assert.equal(out.length, 2);
  assert.equal(out[0].id, 'a#one', 'equal scores → lexical id order');
});

test('a digest node lifts its section above an incidental mention', () => {
  // b#two is the section the digest considers central to "autoscaling".
  const nodes = [node('a#one', []), node('b#two', ['autoscaling'])];
  const out = prefilter(chunks, 'autoscaling', [], 5, 2, nodes);
  assert.equal(out[0].id, 'b#two', 'node terms boost outranks the incidental mention');
});

test('the boost also fires on summary and fact tokens', () => {
  const nodes = [node('b#two', [], 'all about autoscaling workers')];
  const out = prefilter(chunks, 'autoscaling', [], 5, 2, nodes);
  assert.equal(out[0].id, 'b#two', 'summary tokens count toward the boost');
});

test('a query term that matches nothing returns no candidates', () => {
  assert.deepEqual(prefilter(chunks, 'nonexistentterm', [], 5, 2), []);
});
