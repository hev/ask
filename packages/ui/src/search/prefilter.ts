import type { GlossaryEntry, DigestNode } from '../digest/schema';
import { expandQueryTerms } from '../digest/expand.ts';
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
 * query term hitting any of these means the digest considers that section central
 * to the term, so it earns a ranking boost over an incidental body mention.
 */
function nodeSignal(nodes: DigestNode[] | undefined): Map<string, Set<string>> {
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
  nodes?: DigestNode[],
): Candidate[] {
  const terms = expandQueryTerms(query, glossary);
  if (!terms.length) return [];

  const signal = nodeSignal(nodes);
  // Inverse document frequency: down-weight terms common across the corpus
  // (stopwords, ubiquitous words like "authentication" or "setup") so a rare,
  // on-topic term ("oidc") dominates ranking. Without it, plain overlap buries
  // the specific section under pages that merely share several common words —
  // which degrades badly as the corpus grows (hundreds → thousands of sections).
  const df = new Map<string, number>();
  for (const chunk of chunks) for (const token of chunk.tokens) df.set(token, (df.get(token) ?? 0) + 1);
  const total = chunks.length;
  const weights = new Map(terms.map((term) => [term, Math.log(1 + total / (1 + (df.get(term) ?? 0)))]));

  const scored = chunks
    .map((chunk) => {
      const boost = signal.get(chunk.id);
      let score = 0;
      for (const term of terms) {
        const weight = weights.get(term) ?? 0;
        if (chunk.tokens.has(term)) score += weight;
        if (boost?.has(term)) score += weight;
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
