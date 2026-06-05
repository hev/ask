import assert from 'node:assert/strict';
import test from 'node:test';
import { assembleDigest, corpusSections, parseEmittedDigest, type CorpusBuild } from '../src/digest/build.ts';
import { tokenize, type Chunk } from '../src/search/chunk.ts';

function chunk(id: string, heading: string, raw: string): Chunk {
  const docSlug = id.split('#')[0];
  return {
    id,
    docSlug,
    docTitle: 'Docs',
    group: 'Overview',
    heading,
    anchorId: id.split('#')[1],
    url: `/docs/${docSlug}#${id.split('#')[1]}`,
    text: raw,
    raw,
    tokens: new Set(tokenize(raw)),
  };
}

const corpus: CorpusBuild = {
  documents: [],
  chunks: [chunk('concepts#flags', 'Flags', 'Pass `--max-workers` to scale workers.')],
  contentHash: 'hash123',
};

test('parseEmittedDigest reads summaries and suggestions, ignoring junk', () => {
  const emitted = parseEmittedDigest({
    context: 'A docs site.',
    glossary: [{ term: 'workers', aliases: ['worker'], definition: 'Processes.' }],
    summaries: [
      { id: 'concepts#flags', summary: '  Scales workers.  ' },
      { id: 'concepts#flags' }, // no summary → dropped
    ],
    suggestions: ['How do I scale?', '', 7],
  });
  assert.equal(emitted.context, 'A docs site.');
  assert.equal(emitted.summaries.length, 1);
  assert.equal(emitted.summaries[0].summary, 'Scales workers.', 'summary is trimmed');
  assert.deepEqual(emitted.suggestions, ['How do I scale?']);
});

test('assembleDigest carries suggestions and derives facts deterministically', () => {
  const digest = assembleDigest(
    {
      context: 'ctx',
      glossary: [],
      summaries: [{ id: 'concepts#flags', summary: 'Scales workers.' }],
      suggestions: ['How do I scale workers?'],
    },
    corpus,
  );
  assert.equal(digest.version, 2);
  assert.equal(digest.contentHash, 'hash123');
  assert.deepEqual(digest.suggestions, ['How do I scale workers?']);
  assert.equal(digest.nodes.length, 1);
  assert.equal(digest.nodes[0].summary, 'Scales workers.');
  // The flag is extracted verbatim by code, never authored by the model.
  assert.ok(
    digest.nodes[0].facts.some((fact) => fact.literal === '--max-workers'),
    'verbatim flag is extracted into facts',
  );
});

test('corpusSections projects chunks into the model-input shape', () => {
  const sections = corpusSections(corpus);
  assert.deepEqual(sections, [
    {
      id: 'concepts#flags',
      url: '/docs/concepts#flags',
      title: 'Docs > Flags',
      text: 'Pass `--max-workers` to scale workers.',
    },
  ]);
});
