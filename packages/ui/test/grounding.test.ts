import assert from 'node:assert/strict';
import test from 'node:test';
import { runAgenticAnswerLoop, toolUse, INLINE_DIGEST_BUDGET, type CallClaude, type StreamClaude } from '../src/search/loop.ts';
import { chunkDocument } from '../src/search/chunk.ts';
import { EMPTY_DIGEST } from '../src/digest/schema.ts';
import { buildNodes } from '../src/digest/build.ts';

// Fictional product facts: the shared prompt must not encode consumer names.
const chunks = chunkDocument({slug:'backend', title:'Backend', body:'For `kind: archive` only, `Ready=False` means `FormatMismatch`. Other backends are not covered.'}, '/docs/');
const nodes = buildNodes(chunks, new Map());
const config = {model:'test',maxIterations:2,candidatePerSearch:5,perDocCap:2,maxResults:6,answerMaxTokens:512};
for (const path of ['legacy','digest','routed']) test(`${path} carries applicability policy and qualified source into the final answer`, async () => {
  let calls=0, streamed=false;
  const check = (system: unknown) => {
    const text=JSON.stringify(system);
    assert.match(text,/Preserve all documented applicability conditions/);
    assert.match(text,/Do not broaden conditional behavior/);
    assert.match(text,/Do not infer installation steps, their absence, or generated output objects/);
    assert.match(text,/Omit unsupported details/);
  };
  const call: CallClaude = async opts => {
    check(opts.system);
    calls++;
    return calls===1
      ? {stop_reason:'tool_use',content:[path==='legacy' ? toolUse('s','search',{query:'archive FormatMismatch'}) : toolUse('o','open_section',{id:'backend'})]}
      : {stop_reason:'end_turn',content:[{type:'text',text:'Ready.'}]};
  };
  const stream: StreamClaude = async function* (opts) {
    streamed=true; check(opts.system);
    const context=JSON.stringify(opts.messages);
    assert.match(context,/kind: archive/);
    assert.match(context,/only/);
    assert.match(context,/FormatMismatch/);
    yield {type:'text',text:'Fixture answer; no model accuracy claim.'};
  };
  const digest={...EMPTY_DIGEST,nodes:path==='legacy'?[]:nodes,overview:path==='routed'?'x'.repeat(INLINE_DIGEST_BUDGET+1):''};
  const events=[];
  for await (const event of runAgenticAnswerLoop({apiKey:'test',query:'archive FormatMismatch',chunks,digest,config,call,stream})) events.push(event);
  assert.ok(streamed);assert.equal(events.at(-1)?.type,'done');
});


// These tests inspect requests, not whether a real model follows the policy.
const installCases = [
  { name: 'inventory-only', body: 'The Lumen installer includes the catalog and worker resources. Run `lumen install`.', literals: ['catalog', 'worker', 'lumen install'] },
  { name: 'scoped-procedure', body: 'For fresh installations only, `lumen install` applies the catalog without a separate action. Upgrades require `lumen migrate` before installation.', literals: ['fresh installations only', 'without a separate action', 'Upgrades require', 'lumen migrate'] },
];
for (const path of ['legacy', 'digest', 'routed']) for (const fixture of installCases) {
  test(`${path} preserves ${fixture.name} evidence in the final installation request`, async (t) => {
    const installChunks = chunkDocument({ slug: 'install', title: 'Install', body: fixture.body }, '/docs/');
    const installNodes = buildNodes(installChunks, new Map());
    const transcript: unknown[] = [];
    let calls = 0;
    let streams = 0;
    const provisional = 'Provisional retrieval prose: check procedural support before answering.';
    const call: CallClaude = async (opts) => {
      // Capture only explicit request fields, never apiKey or transport options.
      transcript.push(structuredClone({ phase: 'retrieval', system: opts.system, messages: opts.messages }));
      calls++;
      return calls === 1
        ? { stop_reason: 'tool_use', content: [
            { type: 'text', text: provisional },
            path === 'legacy' ? toolUse('install-search', 'search', { query: 'Lumen install catalog' })
              : toolUse('install-open', 'open_section', { id: installNodes[0].id }),
          ] }
        : { stop_reason: 'end_turn', content: [{ type: 'text', text: 'Ready.' }] };
    };
    const stream: StreamClaude = async function* (opts) {
      streams++;
      const request = structuredClone({ phase: 'answer', system: opts.system, messages: opts.messages });
      transcript.push(request);
      const policy = JSON.stringify(opts.system);
      assert.match(policy, /A list of installed resources does not by itself establish/);
      assert.match(policy, /require explicit procedural support/);
      assert.match(policy, /do not establish the requirement or its absence/);
      assert.match(policy, /Preserve explicitly documented guarantees with their applicability conditions/);
      const evidence = JSON.stringify(opts.messages);
      for (const literal of fixture.literals) assert.ok(evidence.includes(literal), literal);
      // Existing transcript behavior is observed, not changed by this candidate.
      assert.equal(evidence.includes(provisional), path !== 'routed');
      yield { type: 'text', text: 'Synthetic transport output, not an evaluated answer.' };
    };
    const digest = { ...EMPTY_DIGEST, nodes: path === 'legacy' ? [] : installNodes,
      overview: path === 'routed' ? 'x'.repeat(INLINE_DIGEST_BUDGET + 1) : '' };
    const events = [];
    for await (const event of runAgenticAnswerLoop({ apiKey: 'synthetic-unused',
      query: 'What does Lumen install? Is a separate catalog action required?',
      chunks: installChunks, digest, config, call, stream })) events.push(event);
    assert.equal(streams, 1);
    assert.equal(events.at(-1)?.type, 'done');
    assert.ok(JSON.stringify(transcript).includes('tool_result'));
    t.diagnostic(JSON.stringify({ path, fixture: fixture.name, requests: transcript }));
  });
}
