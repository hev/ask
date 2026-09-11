import assert from 'node:assert/strict';
import test from 'node:test';
import { readFileSync } from 'node:fs';
import { renderMarkdown, escapeHtml, type Source } from '../src/components/markdown.ts';

// Preserved public-doc answers captured on 2026-09-10; never call a provider.
const captures: { name: string; markdown: string; sources: Source[]; codeBlocks: string[] }[] =
  JSON.parse(readFileSync(new URL('./fixtures/captured-answers.json', import.meta.url), 'utf8'));
for (const fixture of captures) test(`replays ${fixture.name} with literal code and actual citations`, () => {
  const html = renderMarkdown(fixture.markdown, fixture.sources);
  assert.deepEqual([...html.matchAll(/<pre><code>([\s\S]*?)<\/code><\/pre>/g)].map(m => m[1]),
    fixture.codeBlocks.map(escapeHtml));
  assert.ok(html.includes('class="as-answer-link"'));
  if (fixture.name === 'owner-pro') assert.match(html, /<ul><li>/);
});

test('fences preserve code bytes, blank lines, CRLF, indentation and a partial last line', () => {
  for (const fence of ['```sh', '~~~sh']) {
    const body = '  echo "<&>"\r\n\r\nprintf `date`';
    assert.equal(renderMarkdown(fence + '\r\n' + body), `<pre><code>${escapeHtml(body)}</code></pre>`);
    assert.equal(renderMarkdown(fence + '\n' + body + '\n' + fence.slice(0,3)),
      `<pre><code>${escapeHtml(body + '\n')}</code></pre>`);
  }
});

test('fences do not interpret markup, shorter fences or unsafe language labels', () => {
  const body = '[safe](/docs) **bold** <img src=x onerror=alert(1)>\n```\n';
  assert.equal(renderMarkdown('````<script>\n' + body + '````'), `<pre><code>${escapeHtml(body)}</code></pre>`);
  const html = renderMarkdown('```sh\n[safe](/docs)\n```', [{title:'Docs',url:'/docs'}]);
  assert.ok(!html.includes('<a '));
  assert.equal(renderMarkdown('```sh\necho one\necho two\n```\nDone.'),
    '<pre><code>echo one\necho two\n</code></pre><p>Done.</p>');
});

test('paragraph-adjacent lists and indented continuations remain structural', () => {
  assert.equal(renderMarkdown('Kinds:\n- one\n  continued\n- two\nNext paragraph.\n1. first\n2. second'),
    '<p>Kinds:</p><ul><li>one continued</li><li>two</li></ul><p>Next paragraph.</p><ol><li>first</li><li>second</li></ol>');
  assert.equal(renderMarkdown('+ one\n+ two\n\nDone.'), '<ul><li>one</li><li>two</li></ul><p>Done.</p>');
});

test('unsafe source URLs never become answer anchors; escaped URLs still match raw allow-list', () => {
  for (const url of ['javascript:alert', 'data:text/html,evil', '//evil.example', '/\\evil', '/docs\u0001']) {
    assert.ok(!renderMarkdown(`[open](${url})`, [{title:'Unsafe',url}]).includes('<a '));
  }
  const url = '/docs?a=1&b=2';
  assert.match(renderMarkdown(`[open](${url})`, [{title:'Docs',url}]), /href="\/docs\?a=1&amp;b=2"/);
  assert.equal(renderMarkdown('`[open](/docs) **literal**`', [{title:'Docs',url:'/docs'}]),
    '<p><code>[open](/docs) **literal**</code></p>');
});

test('every partial prefix stays escaped and emits balanced trusted block tags', () => {
  const input = 'Intro:\n- safe\n```sh\n<img src=x onerror=alert(1)>\n\n[link](javascript:alert)\n```\nEnd.';
  for (let i=0; i<=input.length; i++) {
    const html = renderMarkdown(input.slice(0,i));
    assert.doesNotMatch(html, /<(?:img|script|iframe)\b/);
    for (const tag of ['pre','code','p','ul','li']) {
      assert.equal((html.match(new RegExp(`<${tag}>`, 'g')) ?? []).length,
        (html.match(new RegExp(`</${tag}>`, 'g')) ?? []).length, `${i}: ${tag}`);
    }
  }
});
