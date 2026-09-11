// Minimal, dependency-free, streaming-safe Markdown renderer for the AI answer.
//
// The renderer is fed the *whole* accumulated answer on every token, so partial
// syntax (an unterminated `[label](`) simply fails to match and renders as
// literal text until the closing token arrives — never a broken DOM. Model
// output is untrusted: the source text is HTML-escaped first, then our own
// trusted tags are injected. Links are validated against the grounding source
// set so a hallucinated URL degrades to plain text — and, when the source
// carries distinctive `terms`, a link whose surrounding text shares none of
// them (a misattributed citation) also degrades to plain text.

export interface Source {
  title: string;
  heading?: string;
  url: string;
  group?: string;
  /** Distinctive tokens of the cited section, for the link-support check. */
  terms?: string[];
}

const LINK_RE = /\[([^\]]+)\]\(([^)\s]+)\)/g;
const LIST_ITEM_RE = /^\s*(?:[-*+]|\d+\.)\s+/;

export function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

/** Breadcrumb for a source, e.g. "Concepts › Kubernetes autoscaling". */
export function sourceBreadcrumb(source: Source): string {
  return [source.group, source.heading ?? source.title].filter(Boolean).join(' › ');
}

function tokenSet(text: string): Set<string> {
  // Strip link *targets* (keep labels) so a URL/anchor slug can't leak its own
  // terms into the support check — otherwise a link to #autoscaling would always
  // "support" itself via the word "autoscaling" in its href.
  const withoutUrls = text.replace(LINK_RE, (_m, label: string) => label);
  return new Set(withoutUrls.toLowerCase().match(/[a-z0-9]+/g) ?? []);
}

/**
 * A cited link survives only if the text around it shares a distinctive term
 * with the cited section. Lenient by design: a source with no `terms` (e.g. a
 * legacy/degraded digest) is never degraded on this basis.
 */
function supportsClaim(unitTokens: Set<string>, source: Source): boolean {
  if (!source.terms || source.terms.length === 0) return true;
  for (const term of source.terms) if (unitTokens.has(term)) return true;
  return false;
}

// Root-relative links or explicit web URLs only, even if source metadata is unsafe.
function safeUrl(url: string): boolean {
  if (/[\\\u0000-\u0020\u007f]/.test(url)) return false;
  if (url.startsWith('/') && !url.startsWith('//')) return true;
  try {
    const parsed = new URL(url);
    return parsed.protocol === 'https:' || parsed.protocol === 'http:';
  } catch { return false; }
}

const FENCE_RE = /^ {0,3}(`{3,}|~{3,})(.*)$/;
const RULE_RE = /^(?:-{3,}|\*{3,}|_{3,})$/;

export function renderMarkdown(md: string, sources: Source[] = []): string {
  const urlMap = new Map(sources.map((source) => [source.url, source]));
  // Retain line endings in code, including blank lines and the last partial line.
  const lines = md.match(/[^\n]*(?:\n|$)/g)?.filter(Boolean) ?? [];
  const text = (i: number) => lines[i].replace(/\r?\n$/, '');
  const fenceAt = (i: number) => {
    const match = text(i).match(FENCE_RE);
    // Backtick fence info strings cannot themselves contain a backtick.
    return match && !(match[1][0] === '`' && match[2].includes('`')) ? match : null;
  };
  const inline = (value: string) => renderInline(value, urlMap, tokenSet(value));
  const out: string[] = [];
  let i = 0;
  while (i < lines.length) {
    const line = text(i);
    if (!line.trim()) { i++; continue; }
    const fence = fenceAt(i);
    if (fence) {
      const closing = new RegExp(`^ {0,3}${fence[1][0]}{${fence[1].length},}[ \\t]*$`);
      let code = '';
      i++;
      while (i < lines.length && !closing.test(text(i))) code += lines[i++];
      if (i < lines.length) i++;
      // Never feed code or the optional language label through inline parsing.
      out.push(`<pre><code>${escapeHtml(code)}</code></pre>`);
      continue;
    }
    if (RULE_RE.test(line.trim())) { i++; continue; }
    const heading = line.match(/^#{1,6}\s+(.+)$/);
    if (heading) { out.push(`<p><strong>${inline(heading[1])}</strong></p>`); i++; continue; }
    if (LIST_ITEM_RE.test(line)) {
      const ordered = /^\s*\d+\./.test(line);
      const tag = ordered ? 'ol' : 'ul';
      const items: string[] = [];
      while (i < lines.length && LIST_ITEM_RE.test(text(i)) &&
        /^\s*\d+\./.test(text(i)) === ordered && !RULE_RE.test(text(i).trim())) {
        let item = text(i++).replace(LIST_ITEM_RE, '');
        while (i < lines.length && /^\s+\S/.test(text(i)) &&
          !LIST_ITEM_RE.test(text(i)) && !fenceAt(i)) item += ' ' + text(i++).trim();
        items.push(`<li>${inline(item)}</li>`);
      }
      out.push(`<${tag}>${items.join('')}</${tag}>`);
      continue;
    }
    const paragraph = [line];
    i++;
    while (i < lines.length && text(i).trim() && !fenceAt(i) &&
      !LIST_ITEM_RE.test(text(i)) && !RULE_RE.test(text(i).trim()) && !/^#{1,6}\s/.test(text(i))) {
      paragraph.push(text(i++));
    }
    out.push(`<p>${inline(paragraph.join(' '))}</p>`);
  }
  return out.join('');
}

function renderInline(raw: string, urlMap: Map<string, Source>, unitTokens: Set<string>, links = true): string {
  // Tokenize before escaping/injecting tags. Code is opaque, and link labels
  // cannot create nested anchors or cause a later regex to rewrite an href.
  const tokens = /`([^`]+)`|\[([^\]]+)\]\(([^)\s]+)\)|\*\*([^*]+)\*\*/g;
  let out = '', end = 0;
  for (const match of raw.matchAll(tokens)) {
    out += escapeHtml(raw.slice(end, match.index));
    if (match[1] !== undefined) out += `<code>${escapeHtml(match[1])}</code>`;
    else if (match[2] !== undefined) {
      const label = renderInline(match[2], urlMap, unitTokens, false);
      const url = match[3];
      const source = urlMap.get(url);
      out += links && source && safeUrl(url) && supportsClaim(unitTokens, source)
        ? `<a class="as-answer-link" href="${escapeHtml(url)}" title="${escapeHtml(sourceBreadcrumb(source))}">${label}</a>`
        : label;
    } else out += `<strong>${renderInline(match[4], urlMap, unitTokens, links)}</strong>`;
    end = match.index + match[0].length;
  }
  return out + escapeHtml(raw.slice(end));
}
