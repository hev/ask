import {
  callClaude,
  streamClaude,
  type AnthropicMessage,
  type AnthropicResponse,
  type AnthropicTextBlock,
  type AnthropicTool,
  type AnthropicToolResultBlock,
  type AnthropicUsage,
  type CallClaudeOptions,
  type StreamEvent,
} from '../llm.ts';
import type { Digest, DigestNode } from '../digest/schema';
import type { Source } from '../components/markdown.ts';
import type { Chunk } from './chunk';
import { prefilter, type Candidate } from './prefilter.ts';
import { makeTelemetry, type Telemetry } from '../observability.ts';

export interface SearchLoopConfig {
  model: string;
  maxIterations: number;
  candidatePerSearch: number;
  perDocCap: number;
  maxResults: number;
  answerMaxTokens: number;
}

export type { Source };

/** High-level events the endpoint forwards to the SSE stream. */
export type AgenticEvent =
  | { type: 'search'; query: string }
  | { type: 'sources'; sources: Source[] }
  | { type: 'token'; text: string }
  | { type: 'done' };

export type CallClaude = typeof callClaude;
export type StreamClaude = typeof streamClaude;

export interface AnswerLoopArgs {
  apiKey: string;
  query: string;
  chunks: Chunk[];
  digest: Digest;
  config: SearchLoopConfig;
  signal?: AbortSignal;
  call?: CallClaude;
  stream?: StreamClaude;
  /** PostHog LLM observability sink. Defaults to a no-op. */
  telemetry?: Telemetry;
}

function randomSpanId(): string {
  const c = (globalThis as { crypto?: { randomUUID?: () => string } }).crypto;
  return c?.randomUUID ? c.randomUUID() : `span-${Date.now()}`;
}

/**
 * Run one non-streaming tool turn and emit an `$ai_generation` for it. The input
 * snapshot is taken before the call so the trace shows exactly what the model saw.
 */
async function tracedCall(
  call: CallClaude,
  opts: CallClaudeOptions,
  telemetry: Telemetry,
  step: number,
): Promise<AnthropicResponse> {
  const startedAt = Date.now();
  const input = opts.messages.slice();
  const response = await call(opts);
  telemetry.generation({
    spanId: randomSpanId(),
    spanName: `turn ${step + 1}`,
    model: opts.model,
    input,
    output: response.content,
    usage: response.usage,
    latencyMs: Date.now() - startedAt,
    httpStatus: 200,
  });
  return response;
}

/**
 * Wrap the streamed answer turn: forward every event untouched, accumulate the
 * answer text and token usage, then emit a single `$ai_generation` at the end.
 */
async function* tracedStream(
  stream: StreamClaude,
  opts: CallClaudeOptions,
  telemetry: Telemetry,
): AsyncGenerator<StreamEvent> {
  const startedAt = Date.now();
  const input = opts.messages.slice();
  let text = '';
  let usage: AnthropicUsage | undefined;
  for await (const event of stream(opts)) {
    if (event.type === 'text') text += event.text;
    else if (event.type === 'stop') usage = event.usage;
    yield event;
  }
  telemetry.generation({
    spanId: randomSpanId(),
    spanName: 'answer',
    model: opts.model,
    input,
    output: [{ type: 'text', text }],
    usage,
    latencyMs: Date.now() - startedAt,
    httpStatus: 200,
  });
}

/**
 * Cap on the characters the digest path inlines into the system prompt (the
 * `<map>` + `<summaries>` blocks). Below it, every section summary is inlined so
 * the agent navigates from a complete map — best for small/medium sites. Above
 * it (large docs, e.g. a CLI/API reference with thousands of sections), inlining
 * everything would blow the context window, so the loop switches to search-routed
 * navigation: a compact page map plus a search tool that surfaces ids on demand.
 * ~200 KB ≈ ~50k tokens; a ~500-section site stays fully inlined as before.
 */
export const INLINE_DIGEST_BUDGET = 200_000;

/** Cheap estimate of what `buildDigestSystemPrompt` would inline, without building it. */
export function digestInlineSize(digest: Digest): number {
  let size = digest.overview.length;
  for (const node of digest.nodes) size += node.id.length + node.summary.length + 24;
  return size;
}

/**
 * Entry point. When the committed digest carries distilled `nodes`, the agent
 * navigates that shadow digest: small digests are inlined whole (digest path);
 * digests larger than {@link INLINE_DIGEST_BUDGET} are navigated by search so the
 * prompt stays bounded (routed path). A node-less (v1 / degraded) digest falls
 * back to the original keyword-search loop, unchanged.
 */
export async function* runAgenticAnswerLoop(args: AnswerLoopArgs): AsyncGenerator<AgenticEvent> {
  if (args.digest.nodes && args.digest.nodes.length > 0) {
    if (digestInlineSize(args.digest) <= INLINE_DIGEST_BUDGET) {
      yield* digestAnswerLoop(args);
    } else {
      yield* routedDigestAnswerLoop(args);
    }
  } else {
    yield* legacyAnswerLoop(args);
  }
}

// ---------------------------------------------------------------------------
// Graph path: navigate the distilled shadow digest and answer from it.
// ---------------------------------------------------------------------------

const OPEN_SECTION_TOOL: AnthropicTool = {
  name: 'open_section',
  description:
    'Open a documentation section by its id (taken from the map) to read its distilled summary, its exact facts (flags, code, identifiers), and — for reference sections — its source text. Open every section you draw your answer from; you may only cite sections you have opened.',
  input_schema: {
    type: 'object',
    properties: {
      id: {
        type: 'string',
        description: 'The exact section id from the map, e.g. "concepts#kubernetes-autoscaling".',
      },
    },
    required: ['id'],
  },
};

async function* digestAnswerLoop({
  apiKey,
  query,
  chunks,
  digest,
  config,
  signal,
  call = callClaude,
  stream = streamClaude,
  telemetry = makeTelemetry(),
}: AnswerLoopArgs): AsyncGenerator<AgenticEvent> {
  const byId = new Map(chunks.map((chunk) => [chunk.id, chunk]));
  const nodesById = new Map(digest.nodes.map((node) => [node.id, node]));
  const opened = new Map<string, DigestNode>();
  const messages: AnthropicMessage[] = [{ role: 'user', content: `Query: ${query}` }];
  const system = buildDigestSystemPrompt(digest);

  const open = (id: string): DigestNode | null => {
    const node = nodesById.get(id);
    if (node) opened.set(id, node);
    return node ?? null;
  };

  // Phase 1: bounded loop of section opens (non-streaming tool turns).
  for (let i = 0; i < config.maxIterations; i += 1) {
    const response = await tracedCall(
      call,
      {
        apiKey,
        model: config.model,
        system,
        messages,
        tools: [OPEN_SECTION_TOOL],
        toolChoice: { type: 'auto' },
        maxTokens: 1024,
        signal,
      },
      telemetry,
      i,
    );

    messages.push({ role: 'assistant', content: response.content });
    const toolResults: AnthropicToolResultBlock[] = [];

    for (const block of response.content) {
      if (block.type !== 'tool_use' || block.name !== 'open_section') continue;
      const id = normalizeId(block.input);
      const node = open(id);
      if (node) yield { type: 'search', query: node.heading ?? node.title };
      toolResults.push({
        type: 'tool_result',
        tool_use_id: block.id,
        content: node
          ? JSON.stringify(openSectionResult(node, byId))
          : JSON.stringify({ error: `No section "${id}". Use an exact id from the map.` }),
      });
    }

    if (!toolResults.length) break; // model is ready to answer
    messages.push({ role: 'user', content: toolResults });
  }

  // Fallback: ground the answer even if the model opened nothing, by opening the
  // best keyword matches over the map.
  if (!opened.size) {
    for (const candidate of prefilter(chunks, query, digest.glossary, config.maxResults, config.perDocCap, digest.nodes)) {
      const node = open(candidate.id);
      if (node) yield { type: 'search', query: node.heading ?? node.title };
    }
    if (opened.size && lastRole(messages) !== 'user') {
      const sections = [...opened.values()].map((node) => openSectionResult(node, byId));
      messages.push({ role: 'user', content: `Opened sections:\n${JSON.stringify(sections)}` });
    }
  }

  // The answer turn must start a fresh assistant response. If the loop ended on
  // an assistant turn, nudge with a final user message so it isn't a prefill.
  if (lastRole(messages) === 'assistant') {
    messages.push({
      role: 'user',
      content:
        'Write the answer now. Begin directly with the answer itself — no preamble, no "based on…" opener, no headings. Link only to sections you opened, using their exact url.',
    });
  }

  const sources = sourcesFromNodes(opened, config.maxResults);
  yield { type: 'sources', sources };

  // Phase 2: streamed answer turn — no tools, so the model can only answer.
  for await (const event of tracedStream(
    stream,
    {
      apiKey,
      model: config.model,
      system: answerSystem(system, sources),
      messages,
      maxTokens: config.answerMaxTokens,
      signal,
    },
    telemetry,
  )) {
    if (event.type === 'text' && event.text) yield { type: 'token', text: event.text };
  }

  yield { type: 'done' };
}

function openSectionResult(node: DigestNode, byId: Map<string, Chunk>) {
  const base = {
    id: node.id,
    url: node.url,
    heading: node.heading,
    group: node.group,
    mode: node.mode,
    summary: node.summary,
    facts: node.facts.map((fact) => ({ kind: fact.kind, literal: fact.literal })),
  };
  // Reference sections carry dense literals; hand the model the source text so it
  // reads exact wording rather than trusting a paraphrase.
  if (node.mode === 'source-primary') {
    const text = byId.get(node.id)?.text ?? '';
    return { ...base, text: text.length > 1200 ? text.slice(0, 1200) + '…' : text };
  }
  return base;
}

function buildDigestSystemPrompt(digest: Digest): AnthropicTextBlock[] {
  return [
    {
      type: 'text',
      text: `You are the documentation assistant for this site. Answer the user's question using ONLY the documentation sections you open with the open_section tool.

You are given a map of the documentation below: every section, its id, and a short summary. Open the sections you need (open_section), reading their summary and exact facts, then write your answer. You may run up to a few opens. Open every section your answer draws on — you may only link to sections you opened.

Write a short, direct answer in Markdown:
- Start IMMEDIATELY with the substance. Your first sentence must answer the question. Never open with "Based on…", "Here is…", "Sure", a restatement of the question, or any summary/preamble.
- Keep it tight: one or two short paragraphs, plus a short bullet list only if it genuinely helps. This renders in a small search popover, so do NOT use headings (#, ##) or horizontal rules (---).
- For exact strings (flags, commands, identifiers, versions), quote the section's \`facts\` verbatim — never reword them.
- When you reference a section, link to it inline using its exact \`url\`, e.g. [autoscaling](/docs/concepts#kubernetes-autoscaling). Never invent a URL or anchor.
- If the documentation does not cover the question, say so plainly in one sentence and do not fabricate an answer.`,
    },
    {
      type: 'text',
      text: `<map>\n${digest.overview || renderNodeMap(digest.nodes)}\n</map>\n\n<summaries>\n${digest.nodes
        .map((node) => `- \`${node.id}\`${node.mode === 'source-primary' ? ' (reference)' : ''}: ${node.summary}`)
        .join('\n')}\n</summaries>`,
      cache_control: { type: 'ephemeral' },
    },
  ];
}

/** Fallback map when a digest predates the stored `overview`. */
function renderNodeMap(nodes: DigestNode[]): string {
  return nodes.map((node) => `- ${node.heading ?? node.title} — \`${node.id}\``).join('\n');
}

// ---------------------------------------------------------------------------
// Routed path: navigate a large digest by search instead of inlining it whole.
// ---------------------------------------------------------------------------

const SEARCH_SECTIONS_TOOL: AnthropicTool = {
  name: 'search_sections',
  description:
    'Search the documentation for sections relevant to a focused sub-query. Returns matching section ids with their group, heading, and a one-line summary. Use it to find the ids you then read with open_section.',
  input_schema: {
    type: 'object',
    properties: {
      query: { type: 'string', description: 'Focused keyword query or synonym expansion to search for.' },
    },
    required: ['query'],
  },
};

/** Compact group → page map: orientation only, so the prompt stays bounded. */
function routedDigestMap(nodes: DigestNode[]): string {
  const byGroup = new Map<string, Set<string>>();
  for (const node of nodes) {
    const group = node.group ?? 'Docs';
    if (!byGroup.has(group)) byGroup.set(group, new Set());
    byGroup.get(group)!.add(node.title);
  }
  const lines: string[] = [];
  for (const [group, pages] of byGroup) {
    lines.push(`## ${group}`);
    for (const page of pages) lines.push(`- ${page}`);
  }
  return lines.join('\n');
}

function routedDigestSystemPrompt(digest: Digest): AnthropicTextBlock[] {
  return [
    {
      type: 'text',
      text: `You are the documentation assistant for this site. Answer the user's question using ONLY documentation sections you retrieve.

The documentation is large, so it is not all shown here. Use search_sections to find relevant sections — each result includes a short summary you can answer from directly. When you need a section's exact facts (flags, commands, identifiers), open_section it. One or two focused searches is plenty: once the results cover the question, STOP searching and answer. Do not keep searching for a perfect match.

Write a short, direct answer in Markdown:
- Start IMMEDIATELY with the substance. Your first sentence must answer the question. Never open with "Based on…", "Here is…", "Sure", a restatement of the question, or any summary/preamble.
- Keep it tight: one or two short paragraphs, plus a short bullet list only if it genuinely helps. This renders in a small search popover, so do NOT use headings (#, ##) or horizontal rules (---).
- For exact strings (flags, commands, identifiers, versions), quote a section's \`facts\` verbatim — never reword them.
- When you reference a section, link to it inline using its exact \`url\` from your search results or open_section, e.g. [autoscaling](/docs/concepts#kubernetes-autoscaling). Never invent a URL or anchor.
- If the documentation does not cover the question, say so plainly in one sentence and do not fabricate an answer.`,
    },
    {
      type: 'text',
      text: `<domain_context>\n${digest.context || 'No digest context is available.'}\n</domain_context>\n\n<map>\n${routedDigestMap(digest.nodes)}\n</map>`,
      cache_control: { type: 'ephemeral' },
    },
  ];
}

/** Search the digest's nodes for a sub-query; returns distilled candidates. */
function searchSections(
  searchQuery: string,
  chunks: Chunk[],
  nodesById: Map<string, DigestNode>,
  digest: Digest,
  config: SearchLoopConfig,
) {
  return prefilter(chunks, searchQuery, digest.glossary, config.candidatePerSearch, config.perDocCap, digest.nodes)
    .map((candidate) => nodesById.get(candidate.id))
    .filter((node): node is DigestNode => node !== undefined)
    .map((node) => ({
      id: node.id,
      url: node.url,
      group: node.group,
      heading: node.heading,
      summary: node.summary,
      ...(node.mode === 'source-primary' ? { reference: true } : {}),
    }));
}

async function* routedDigestAnswerLoop({
  apiKey,
  query,
  chunks,
  digest,
  config,
  signal,
  call = callClaude,
  stream = streamClaude,
  telemetry = makeTelemetry(),
}: AnswerLoopArgs): AsyncGenerator<AgenticEvent> {
  const byId = new Map(chunks.map((chunk) => [chunk.id, chunk]));
  const nodesById = new Map(digest.nodes.map((node) => [node.id, node]));
  const opened = new Map<string, DigestNode>();
  const seen = new Map<string, DigestNode>(); // sections surfaced by search, in rank order
  const messages: AnthropicMessage[] = [{ role: 'user', content: `Query: ${query}` }];
  const system = routedDigestSystemPrompt(digest);

  const open = (id: string): DigestNode | null => {
    const node = nodesById.get(id);
    if (node) opened.set(id, node);
    return node ?? null;
  };

  // Phase 1: bounded loop of searches and section opens (non-streaming tool turns).
  for (let i = 0; i < config.maxIterations; i += 1) {
    const response = await tracedCall(
      call,
      {
        apiKey,
        model: config.model,
        system,
        messages,
        tools: [SEARCH_SECTIONS_TOOL, OPEN_SECTION_TOOL],
        toolChoice: { type: 'auto' },
        maxTokens: 1024,
        signal,
      },
      telemetry,
      i,
    );

    messages.push({ role: 'assistant', content: response.content });
    const toolResults: AnthropicToolResultBlock[] = [];

    for (const block of response.content) {
      if (block.type !== 'tool_use') continue;
      if (block.name === 'search_sections') {
        const searchQuery = normalizeToolQuery(block.input) || query;
        yield { type: 'search', query: searchQuery };
        const results = searchSections(searchQuery, chunks, nodesById, digest, config);
        for (const result of results) {
          const node = nodesById.get(result.id);
          if (node && !seen.has(node.id)) seen.set(node.id, node);
        }
        toolResults.push({
          type: 'tool_result',
          tool_use_id: block.id,
          content: JSON.stringify(results),
        });
      } else if (block.name === 'open_section') {
        const id = normalizeId(block.input);
        const node = open(id);
        toolResults.push({
          type: 'tool_result',
          tool_use_id: block.id,
          content: node
            ? JSON.stringify(openSectionResult(node, byId))
            : JSON.stringify({ error: `No section "${id}". Search first, then open an exact id from the results.` }),
        });
      }
    }

    if (!toolResults.length) break; // model is ready to answer
    messages.push({ role: 'user', content: toolResults });
  }

  // Fallback: if the model never searched or opened anything, ground on the best
  // keyword matches for the original query so the answer isn't empty.
  if (!opened.size && !seen.size) {
    for (const candidate of prefilter(chunks, query, digest.glossary, config.maxResults, config.perDocCap, digest.nodes)) {
      const node = nodesById.get(candidate.id);
      if (node && !seen.has(node.id)) seen.set(node.id, node);
    }
  }

  // Ground the answer in everything surfaced: opened sections (full facts) first,
  // then searched summaries, capped to maxResults.
  const grounded = [...new Map<string, DigestNode>([...opened, ...seen]).values()].slice(0, config.maxResults);
  const sources = sourcesFromNodes(new Map(grounded.map((node) => [node.id, node])), config.maxResults);
  yield { type: 'sources', sources };

  // Phase 2: a clean answer turn. Replaying the tool transcript keeps the model in
  // "let me search more" mode (and it tries to call tools that no longer exist), so
  // instead hand it just the question and the gathered sections — now it can only
  // write the final prose answer.
  const answerContext = grounded.map((node) => openSectionResult(node, byId));
  const answerMessages: AnthropicMessage[] = [
    {
      role: 'user',
      content: `Question: ${query}\n\nAnswer using only these documentation sections:\n${JSON.stringify(answerContext)}`,
    },
  ];
  for await (const event of tracedStream(
    stream,
    {
      apiKey,
      model: config.model,
      system: routedAnswerSystemPrompt(digest),
      messages: answerMessages,
      maxTokens: config.answerMaxTokens,
      signal,
    },
    telemetry,
  )) {
    if (event.type === 'text' && event.text) yield { type: 'token', text: event.text };
  }

  yield { type: 'done' };
}

/** Answer-only system prompt for the routed loop's final turn (no tools). */
function routedAnswerSystemPrompt(digest: Digest): AnthropicTextBlock[] {
  return [
    {
      type: 'text',
      text: `You are the documentation assistant for this site. Write the answer to the user's question using ONLY the documentation sections provided in the next message. You have no tools — produce the final prose answer now.

- Start IMMEDIATELY with the substance. Your first sentence must answer the question. Never open with "Based on…", "Here is…", "Sure", "Let me…", or any preamble or statement about searching, opening, or checking further.
- Keep it tight: one or two short paragraphs, plus a short bullet list only if it genuinely helps. This renders in a small search popover, so do NOT use headings (#, ##) or horizontal rules (---).
- For exact strings (flags, commands, identifiers, versions), quote a section's \`facts\` verbatim — never reword them.
- Link to sections inline using their exact \`url\` from the provided sections, e.g. [autoscaling](/docs/concepts#kubernetes-autoscaling). Never invent a URL or anchor.
- If the provided sections do not answer the question, say so plainly in one sentence and do not fabricate an answer.`,
    },
    {
      type: 'text',
      text: `<domain_context>\n${digest.context || ''}\n</domain_context>`,
      cache_control: { type: 'ephemeral' },
    },
  ];
}

function sourcesFromNodes(opened: Map<string, DigestNode>, maxResults: number): Source[] {
  const sources: Source[] = [];
  const urls = new Set<string>();
  for (const node of opened.values()) {
    if (urls.has(node.url)) continue;
    urls.add(node.url);
    sources.push({
      title: node.title,
      heading: node.heading ?? undefined,
      url: node.url,
      group: node.group ?? undefined,
      terms: node.terms,
    });
    if (sources.length >= maxResults) break;
  }
  return sources;
}

function normalizeId(input: unknown): string {
  return typeof (input as { id?: unknown })?.id === 'string' ? (input as { id: string }).id.trim() : '';
}

// ---------------------------------------------------------------------------
// Legacy path: original keyword-search loop, used for node-less graphs.
// ---------------------------------------------------------------------------

interface SeenCandidate {
  chunk: Chunk;
  snippet: string;
}

const SEARCH_TOOL: AnthropicTool = {
  name: 'search',
  description: 'Search the documentation heading chunks with a focused sub-query.',
  input_schema: {
    type: 'object',
    properties: {
      query: { type: 'string', description: 'Focused keyword query or synonym expansion to search for.' },
    },
    required: ['query'],
  },
};

async function* legacyAnswerLoop({
  apiKey,
  query,
  chunks,
  digest,
  config,
  signal,
  call = callClaude,
  stream = streamClaude,
  telemetry = makeTelemetry(),
}: AnswerLoopArgs): AsyncGenerator<AgenticEvent> {
  const byId = new Map(chunks.map((chunk) => [chunk.id, chunk]));
  const seen = new Map<string, SeenCandidate>();
  const messages: AnthropicMessage[] = [{ role: 'user', content: `Query: ${query}` }];
  const system = buildSystemPrompt(digest);

  // Phase 1: bounded, non-streaming search loop.
  for (let i = 0; i < config.maxIterations; i += 1) {
    const response = await tracedCall(
      call,
      {
        apiKey,
        model: config.model,
        system,
        messages,
        tools: [SEARCH_TOOL],
        toolChoice: { type: 'auto' },
        maxTokens: 1024,
        signal,
      },
      telemetry,
      i,
    );

    messages.push({ role: 'assistant', content: response.content });
    const toolResults: AnthropicToolResultBlock[] = [];

    for (const block of response.content) {
      if (block.type !== 'tool_use' || block.name !== 'search') continue;
      const searchQuery = normalizeToolQuery(block.input) || query;
      yield { type: 'search', query: searchQuery };
      const fresh = runSearchTool(searchQuery, chunks, byId, seen, digest, config);
      toolResults.push({
        type: 'tool_result',
        tool_use_id: block.id,
        content: JSON.stringify(fresh.map((candidate) => candidateForToolResult(candidate, byId))),
      });
    }

    if (!toolResults.length) break;
    messages.push({ role: 'user', content: toolResults });
  }

  if (!seen.size) {
    const fresh = runSearchTool(query, chunks, byId, seen, digest, config);
    yield { type: 'search', query };
    if (lastRole(messages) !== 'user') {
      messages.push({
        role: 'user',
        content: `Search results:\n${JSON.stringify(fresh.map((candidate) => candidateForToolResult(candidate, byId)))}`,
      });
    }
  }

  if (lastRole(messages) === 'assistant') {
    messages.push({
      role: 'user',
      content:
        'Write the answer now. Begin directly with the answer itself — no preamble, no "based on…" opener, no headings. Link only with the provided URLs.',
    });
  }

  const sources = sourcesFromSeen(seen, config.maxResults);
  yield { type: 'sources', sources };

  for await (const event of tracedStream(
    stream,
    {
      apiKey,
      model: config.model,
      system: answerSystem(system, sources),
      messages,
      maxTokens: config.answerMaxTokens,
      signal,
    },
    telemetry,
  )) {
    if (event.type === 'text' && event.text) yield { type: 'token', text: event.text };
  }

  yield { type: 'done' };
}

function runSearchTool(
  searchQuery: string,
  chunks: Chunk[],
  byId: Map<string, Chunk>,
  seen: Map<string, SeenCandidate>,
  digest: Digest,
  config: SearchLoopConfig,
): Candidate[] {
  return prefilter(chunks, searchQuery, digest.glossary, config.candidatePerSearch, config.perDocCap, digest.nodes)
    .filter((candidate) => !seen.has(candidate.id))
    .map((candidate) => {
      const chunk = byId.get(candidate.id);
      if (chunk) seen.set(candidate.id, { chunk, snippet: candidate.snippet });
      return candidate;
    })
    .filter((candidate) => byId.has(candidate.id));
}

function candidateForToolResult(candidate: Candidate, byId: Map<string, Chunk>) {
  return {
    id: candidate.id,
    docTitle: candidate.docTitle,
    heading: candidate.heading,
    url: byId.get(candidate.id)?.url ?? '',
    snippet: candidate.snippet,
  };
}

function buildSystemPrompt(digest: Digest): AnthropicTextBlock[] {
  return [
    {
      type: 'text',
      text: `You are the documentation assistant for this site. Answer the user's question using ONLY the documentation sections returned by the search tool.

You decide how many searches to run. Issue focused sub-queries with the search tool: vary terms, try synonyms, and decompose multi-part questions. When you have gathered enough context, stop calling the search tool and write your answer.

Write a short, direct answer in Markdown:
- Start IMMEDIATELY with the substance. Your first sentence must answer the question. Never open with "Based on…", "Here is…", "Sure", a restatement of the question, or any summary/preamble.
- Keep it tight: one or two short paragraphs, plus a short bullet list only if it genuinely helps. This renders in a small search popover, so do NOT use headings (#, ##) or horizontal rules (---).
- Ground every claim in the retrieved sections.
- When you reference a section, link to it inline using its exact \`url\` from the search results, for example: [autoscaling](/docs/concepts#kubernetes-autoscaling). Never invent a URL or anchor — only link to URLs that appear in the search results.
- If the documentation does not cover the question, say so plainly in one sentence and do not fabricate an answer.`,
    },
    {
      type: 'text',
      text: `<domain_context>\n${digest.context || 'No digest context is available.'}\n</domain_context>`,
      cache_control: { type: 'ephemeral' },
    },
  ];
}

/** Appends the grounding allow-list to the system prompt for the answer turn. */
function answerSystem(system: AnthropicTextBlock[], sources: Source[]): AnthropicTextBlock[] {
  if (!sources.length) return system;
  const list = sources.map((source) => `- ${source.url} (${source.heading ?? source.title})`).join('\n');
  return [
    ...system,
    {
      type: 'text',
      text: `You have finished gathering context. Write the answer now. Use only these URLs when linking:\n${list}`,
    },
  ];
}

function sourcesFromSeen(seen: Map<string, SeenCandidate>, maxResults: number): Source[] {
  const sources: Source[] = [];
  const urls = new Set<string>();
  for (const { chunk } of seen.values()) {
    if (urls.has(chunk.url)) continue;
    urls.add(chunk.url);
    sources.push({ title: chunk.docTitle, heading: chunk.heading, url: chunk.url, group: chunk.group });
    if (sources.length >= maxResults) break;
  }
  return sources;
}

function lastRole(messages: AnthropicMessage[]): AnthropicMessage['role'] | undefined {
  return messages[messages.length - 1]?.role;
}

function normalizeToolQuery(input: unknown): string {
  return typeof (input as { query?: unknown })?.query === 'string' ? (input as { query: string }).query.trim() : '';
}

export function toolUse(id: string, name: string, input: unknown): AnthropicResponse['content'][number] {
  return { type: 'tool_use', id, name, input };
}
