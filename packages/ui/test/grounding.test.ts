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
