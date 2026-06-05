import type { GlossaryEntry, KnowledgeNode } from '../kg/schema';
import { expandQueryTerms } from '../kg/expand.ts';
import { tokenize } from './chunk.ts';
import type { Chunk } from './chunk';

export interface Candidate {
  id: string;
  docTitle: string;
  group?: string;
  heading?: string;
  snippet: string;
}

/**
 * Distinctive tokens the digest carries for a section: its `terms`, the
 * tokens of its distilled `summary`, and the tokens of its verbatim `facts`. A
 * query term hitting any of these means the graph considers that section central
 * to the term, so it earns a ranking boost over an incidental body mention.
 */
function nodeSignal(nodes: KnowledgeNode[] | undefined): Map<string, Set<string>> {
  const signal = new Map<string, Set<string>>();
  if (!nodes) return signal;
  for (const node of nodes) {
    const tokens = new Set<string>(node.terms);
    for (const token of tokenize(node.summary)) tokens.add(token);
    for (const fact of node.facts) for (const token of tokenize(fact.literal)) tokens.add(token);
    signal.set(node.id, tokens);
  }
  return signal;
}

/**
 * Keyword prefilter over heading chunks. Query terms are widened through the
 * digest glossary, scored by token overlap — boosted by the digest's
 * distilled per-section signal when `nodes` are present — then capped per
 * document so one page cannot crowd out the rest of the result set. With no
 * nodes it degrades to plain token overlap over the raw chunk text.
 */
export function prefilter(
  chunks: Chunk[],
  query: string,
  glossary: GlossaryEntry[],
  pool: number,
  perDocCap: number,
  nodes?: KnowledgeNode[],
): Candidate[] {
  const terms = expandQueryTerms(query, glossary);
  if (!terms.length) return [];

  const signal = nodeSignal(nodes);
  const scored = chunks
    .map((chunk) => {
      const boost = signal.get(chunk.id);
      let score = 0;
      for (const term of terms) {
        if (chunk.tokens.has(term)) score += 1;
        if (boost?.has(term)) score += 1;
      }
      return { chunk, score };
    })
    .filter((candidate) => candidate.score > 0)
    .sort((a, b) => b.score - a.score || a.chunk.id.localeCompare(b.chunk.id));

  const perDoc = new Map<string, number>();
  const capped = [];
  for (const item of scored) {
    const count = perDoc.get(item.chunk.docSlug) ?? 0;
    if (count >= perDocCap) continue;
    perDoc.set(item.chunk.docSlug, count + 1);
    capped.push(item);
    if (capped.length >= pool) break;
  }

  return capped.map(({ chunk }) => ({
    id: chunk.id,
    docTitle: chunk.docTitle,
    group: chunk.group,
    heading: chunk.heading,
    snippet: excerpt(chunk.text, terms),
  }));
}

function excerpt(text: string, terms: string[], radius = 200): string {
  const lower = text.toLowerCase();
  let pos = -1;
  for (const term of terms) {
    const i = lower.indexOf(term);
    if (i !== -1 && (pos === -1 || i < pos)) pos = i;
  }
  const start = pos === -1 ? 0 : Math.max(0, pos - 40);
  const slice = text.slice(start, start + radius).replace(/\s+/g, ' ').trim();
  return (start > 0 ? '...' : '') + slice + (start + radius < text.length ? '...' : '');
}
