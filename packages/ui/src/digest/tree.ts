import { readFileSync, readdirSync, statSync } from 'node:fs';
import path from 'node:path';
import { parseFrontmatter } from './frontmatter.ts';
import { EMPTY_DIGEST, normalizeDigest, type Digest, type DigestNode, type Fact, type SourceRef } from './schema.ts';

const META_FILE = '_meta.md';
const GLOSSARY_DIR = '_glossary';
const META_OVERVIEW = '\n## Overview\n\n';

export interface DigestTreeFile {
  path: string;
  body: string;
}

export function readDigestArtifact(siteRoot: string, digestPath: string): Digest {
  const resolved = path.resolve(siteRoot, digestPath);
  try {
    const stat = statSync(resolved);
    if (stat.isDirectory()) {
      try {
        return readDigestTree(resolved);
      } catch {
        return readLegacyJson(path.join(resolved, 'digest.json'));
      }
    }
    return readLegacyJson(resolved);
  } catch {
    try {
      return readLegacyJson(path.join(resolved, 'digest.json'));
    } catch {
      return EMPTY_DIGEST;
    }
  }
}

export function readDigestTree(root: string): Digest {
  const meta = parseFrontmatter(readFileSync(path.join(root, META_FILE), 'utf8'));
  const { context, overview } = parseMetaBody(meta.body);
  const digest: Digest = {
    version: 2,
    generatedAt: stringField(meta.data, 'generatedAt'),
    contentHash: stringField(meta.data, 'contentHash'),
    context,
    glossary: [],
    overview,
    suggestions: stringArrayField(meta.data, 'suggestions'),
    nodes: [],
    edges: [],
  };

  for (const file of walk(root)) {
    const rel = path.relative(root, file).replace(/\\/g, '/');
    if (rel === META_FILE || !rel.endsWith('.md')) continue;
    const parsed = parseFrontmatter(readFileSync(file, 'utf8'));
    if (rel.startsWith(`${GLOSSARY_DIR}/`)) {
      const term = stringField(parsed.data, 'term') || titleFromPath(rel);
      if (!term) continue;
      digest.glossary.push({
        term,
        aliases: stringArrayField(parsed.data, 'aliases'),
        definition: parsed.body.trim(),
      });
      continue;
    }
    const node = nodeFromTreeFile(rel, parsed);
    if (node) digest.nodes.push(node);
  }

  digest.glossary.sort((a, b) => a.term.localeCompare(b.term));
  digest.nodes.sort((a, b) => a.id.localeCompare(b.id));
  return normalizeDigest(digest);
}

export function digestTreeFiles(digest: Digest): DigestTreeFile[] {
  const files: DigestTreeFile[] = [{ path: META_FILE, body: renderMetaFile(digest) }];
  for (const entry of digest.glossary) {
    files.push({
      path: `${GLOSSARY_DIR}/${glossaryPath(entry.term)}.md`,
      body: markdownWithFrontmatter(
        [
          ['term', entry.term],
          ['aliases', entry.aliases],
        ],
        entry.definition,
      ),
    });
  }
  digest.nodes.forEach((node, order) => {
    files.push({ path: `${nodePath(node)}.md`, body: renderNodeFile(node, order) });
  });
  return files.sort((a, b) => a.path.localeCompare(b.path));
}

function readLegacyJson(file: string): Digest {
  return normalizeDigest(JSON.parse(readFileSync(file, 'utf8')));
}

function renderMetaFile(digest: Digest): string {
  return markdownWithFrontmatter(
    [
      ['version', digest.version],
      ['generatedAt', digest.generatedAt || new Date().toISOString()],
      ['contentHash', digest.contentHash],
      ['suggestions', digest.suggestions],
    ],
    `## Context\n\n${digest.context.trim()}${META_OVERVIEW}${digest.overview.trim()}`.trim(),
  );
}

function renderNodeFile(node: DigestNode, order: number): string {
  return markdownWithFrontmatter(
    [
      ['id', node.id],
      ['title', node.title],
      ['heading', node.heading],
      ['group', node.group],
      ['order', order],
      ['url', node.url],
      ['anchor', node.sources[0]?.anchor ?? null],
      ['terms', node.terms],
      ['hash', node.hash ?? null],
      ['mode', node.mode],
      ['facts', node.facts],
      ['sources', node.sources],
    ],
    node.summary,
  );
}

function markdownWithFrontmatter(fields: Array<[string, unknown]>, body: string): string {
  const frontmatter = fields.map(([key, value]) => `${key}: ${formatFrontmatterValue(value)}`).join('\n');
  return `---\n${frontmatter}\n---\n\n${body.trim()}\n`;
}

function formatFrontmatterValue(value: unknown): string {
  if (value === null || value === undefined || value === '') return 'null';
  if (typeof value === 'number') return String(value);
  return JSON.stringify(value);
}

function nodePath(node: DigestNode): string {
  return node.id.trim().replace(/#/g, '/').replace(/^\/+|\/+$/g, '');
}

function glossaryPath(term: string): string {
  const slug = term
    .toLowerCase()
    .replace(/[^a-z0-9 _-]+/g, '')
    .trim()
    .replace(/\s+/g, '-');
  return slug || 'term';
}

function nodeFromTreeFile(rel: string, parsed: ReturnType<typeof parseFrontmatter>): DigestNode | null {
  const id = stringField(parsed.data, 'id') || rel.replace(/\.md$/, '').replace(/\//g, '#');
  const url = stringField(parsed.data, 'url');
  if (!id || !url) return null;
  const anchor = nullableStringField(parsed.data, 'anchor');
  const sources = sourceArrayField(parsed.data, 'sources');
  return {
    id,
    kind: 'section',
    title: stringField(parsed.data, 'title') || id,
    heading: nullableStringField(parsed.data, 'heading'),
    group: nullableStringField(parsed.data, 'group'),
    url,
    summary: firstParagraph(parsed.body) || stringField(parsed.data, 'summary'),
    hash: stringField(parsed.data, 'hash') || undefined,
    facts: factArrayField(parsed.data, 'facts'),
    sources: sources.length ? sources : [{ chunkId: id, url, anchor }],
    mode: stringField(parsed.data, 'mode') === 'source-primary' ? 'source-primary' : 'agent-primary',
    terms: stringArrayField(parsed.data, 'terms'),
  };
}

function parseMetaBody(body: string): { context: string; overview: string } {
  let trimmed = body.trim();
  if (trimmed.startsWith('## Context')) trimmed = trimmed.slice('## Context'.length).trim();
  const at = trimmed.indexOf(META_OVERVIEW);
  if (at >= 0) {
    return { context: trimmed.slice(0, at).trim(), overview: trimmed.slice(at + META_OVERVIEW.length).trim() };
  }
  return { context: trimmed, overview: '' };
}

function walk(dir: string): string[] {
  const out: string[] = [];
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const file = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      if (entry.name === 'shards') continue;
      out.push(...walk(file));
    } else {
      out.push(file);
    }
  }
  return out;
}

function firstParagraph(body: string): string {
  return body
    .trim()
    .split(/\n{2,}/)
    .map((part) => part.trim())
    .find(Boolean) ?? '';
}

function titleFromPath(value: string): string {
  return path.basename(value, '.md').replace(/-/g, ' ');
}

function stringField(data: Record<string, unknown>, key: string): string {
  const value = data[key];
  return typeof value === 'string' ? value : '';
}

function nullableStringField(data: Record<string, unknown>, key: string): string | null {
  return stringField(data, key) || null;
}

function stringArrayField(data: Record<string, unknown>, key: string): string[] {
  const value = data[key];
  if (Array.isArray(value)) return value.filter((item): item is string => typeof item === 'string' && item.length > 0);
  return typeof value === 'string' && value ? [value] : [];
}

function factArrayField(data: Record<string, unknown>, key: string): Fact[] {
  const value = data[key];
  if (!Array.isArray(value)) return [];
  return value
    .map((item) => {
      if (!item || typeof item !== 'object') return null;
      const maybe = item as Partial<Fact>;
      if (typeof maybe.literal !== 'string' || !maybe.literal) return null;
      return {
        kind: maybe.kind ?? 'value',
        literal: maybe.literal,
        chunkId: typeof maybe.chunkId === 'string' ? maybe.chunkId : '',
      } satisfies Fact;
    })
    .filter((item): item is Fact => item !== null);
}

function sourceArrayField(data: Record<string, unknown>, key: string): SourceRef[] {
  const value = data[key];
  if (!Array.isArray(value)) return [];
  return value
    .map((item) => {
      if (!item || typeof item !== 'object') return null;
      const maybe = item as Partial<SourceRef>;
      if (typeof maybe.chunkId !== 'string' || typeof maybe.url !== 'string') return null;
      return {
        chunkId: maybe.chunkId,
        url: maybe.url,
        anchor: typeof maybe.anchor === 'string' ? maybe.anchor : null,
      } satisfies SourceRef;
    })
    .filter((item): item is SourceRef => item !== null);
}
