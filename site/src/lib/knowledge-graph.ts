import kgRaw from "../../.hev-ask/digest.json?raw";

export interface GlossaryTerm {
	term: string;
	aliases?: string[];
	definition: string;
}

export interface KgFact {
	kind: "code" | "flag" | "value" | string;
	literal: string;
	chunkId: string;
}

export interface KgNode {
	id: string;
	kind: string;
	title: string;
	heading: string | null;
	group: string | null;
	url: string;
	summary: string;
	facts: KgFact[];
	mode: "agent-primary" | "source-primary";
	terms: string[];
}

export interface KnowledgeGraph {
	version: number;
	generatedAt: string;
	contentHash: string;
	context: string;
	glossary: GlossaryTerm[];
	overview: string;
	nodes: KgNode[];
}

export const knowledgeGraph = JSON.parse(kgRaw) as KnowledgeGraph;

export const knowledgeGraphHref = "/digest";

/** Nodes grouped by their docs group, preserving id order within a group. */
export function nodesByGroup(): { group: string; nodes: KgNode[] }[] {
	const map = new Map<string, KgNode[]>();
	for (const node of knowledgeGraph.nodes ?? []) {
		const group = node.group ?? "Docs";
		if (!map.has(group)) map.set(group, []);
		map.get(group)!.push(node);
	}
	return [...map.entries()]
		.sort((a, b) => a[0].localeCompare(b[0]))
		.map(([group, nodes]) => ({ group, nodes }));
}

export function getKnowledgeGraphRawJson() {
	return JSON.stringify(knowledgeGraph, null, 2);
}

function escapeHtml(text: string) {
	return text
		.replace(/&/g, "&amp;")
		.replace(/</g, "&lt;")
		.replace(/>/g, "&gt;");
}

function renderInline(text: string) {
	return escapeHtml(text)
		.replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>")
		.replace(/`([^`]+)`/g, "<code>$1</code>");
}

/**
 * Render the KG's constrained markdown (the `context` field) to HTML:
 * `##`/`###` headings, `-` bullet lists, paragraphs, `**bold**`, `code`.
 */
export function renderKgMarkdown(md: string): string {
	const out: string[] = [];
	let list: string[] | null = null;
	const flushList = () => {
		if (list) {
			out.push(`<ul>${list.join("")}</ul>`);
			list = null;
		}
	};
	for (const block of md.split(/\n{2,}/)) {
		for (const line of block.split("\n")) {
			const trimmed = line.trim();
			if (!trimmed) continue;
			const heading = trimmed.match(/^(#{2,4})\s+(.*)$/);
			const bullet = trimmed.match(/^-\s+(.*)$/);
			if (bullet) {
				(list ??= []).push(`<li>${renderInline(bullet[1])}</li>`);
			} else if (heading) {
				flushList();
				const level = Math.min(heading[1].length + 1, 6);
				out.push(`<h${level}>${renderInline(heading[2])}</h${level}>`);
			} else {
				flushList();
				out.push(`<p>${renderInline(trimmed)}</p>`);
			}
		}
		flushList();
	}
	return out.join("\n");
}

export function formatGeneratedAt() {
	const date = new Date(knowledgeGraph.generatedAt);
	return Number.isNaN(date.getTime()) ? knowledgeGraph.generatedAt || "—" : date.toISOString();
}
