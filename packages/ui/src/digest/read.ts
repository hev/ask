import type { GlossaryEntry, Digest, DigestNode } from './schema.ts';

export interface SectionSummary {
  id: string;
  title: string;
  heading: string | null;
  group: string | null;
  url: string;
}

export function listGlossary(digest: Digest): GlossaryEntry[] {
  return digest.glossary;
}

export function getGlossaryEntry(digest: Digest, term: string): GlossaryEntry | null {
  const needle = normalizeLookup(term);
  if (!needle) return null;
  return (
    digest.glossary.find((entry) => {
      if (normalizeLookup(entry.term) === needle) return true;
      return entry.aliases.some((alias) => normalizeLookup(alias) === needle);
    }) ?? null
  );
}

export function listSectionSummaries(digest: Digest, group?: string | null): SectionSummary[] {
  const wantedGroup = group ? normalizeLookup(group) : '';
  return digest.nodes
    .filter((node) => !wantedGroup || normalizeLookup(node.group ?? '') === wantedGroup)
    .map(sectionSummary);
}

export function getSection(digest: Digest, id: string): DigestNode | null {
  const needle = decodePathValue(id).trim();
  if (!needle) return null;
  return digest.nodes.find((node) => node.id === needle) ?? null;
}

export function getOverview(digest: Digest): { overview: string; context: string } {
  return { overview: digest.overview, context: digest.context };
}

export function sectionSummary(node: DigestNode): SectionSummary {
  return {
    id: node.id,
    title: node.title,
    heading: node.heading,
    group: node.group,
    url: node.url,
  };
}

export function decodePathValue(value: string): string {
  try {
    return decodeURIComponent(value);
  } catch {
    return value;
  }
}

function normalizeLookup(value: string): string {
  return decodePathValue(value).trim().toLowerCase();
}
