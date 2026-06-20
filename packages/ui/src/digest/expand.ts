import { tokenize } from '../search/chunk.ts';
import type { GlossaryEntry } from './schema';

export function expandQueryTerms(query: string, glossary: GlossaryEntry[], cap = 24): string[] {
  const queryTokens = new Set(tokenize(query));
  if (!queryTokens.size) return [];
  const terms = new Set(queryTokens);

  for (const entry of glossary) {
    if (terms.size >= cap) break;
    // Expand only when the query contains a full glossary phrase — the term or one
    // of its aliases, every token present. Matching on any shared token (e.g. the
    // ubiquitous "authentication") drags in every entry that merely mentions it,
    // which floods results once the glossary is large.
    const phrases = [entry.term, ...entry.aliases].map((phrase) => tokenize(phrase)).filter((tokens) => tokens.length);
    const matched = phrases.some((phrase) => phrase.every((token) => queryTokens.has(token)));
    if (!matched) continue;
    for (const phrase of phrases) {
      for (const token of phrase) {
        if (token.length < 3) continue; // skip noisy short tokens like "ad", "id"
        terms.add(token);
        if (terms.size >= cap) break;
      }
      if (terms.size >= cap) break;
    }
  }

  return [...terms];
}
