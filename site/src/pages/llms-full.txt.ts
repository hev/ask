import type { APIRoute } from "astro";
import { digestData } from "../lib/digest-data";

const SITE = "https://hevask.com";

// The full LLM text is the committed ask digest, concatenated — NOT the raw doc
// bodies. Each section is the distilled summary plus its verbatim facts and a
// source deep link, so an LLM gets the navigable, distilled corpus (the same
// artifact the overlay and a coding agent read) instead of a raw page dump.
export const GET: APIRoute = async () => {
	const d = digestData;
	const abs = (url: string) => (url?.startsWith("http") ? url : `${SITE}${url}`);
	const parts: string[] = [];

	parts.push("# hev ask — full docs (distilled)\n\n");
	parts.push(
		`> The committed ask digest, concatenated: one distilled section per heading, with its source deep link and verbatim facts. Index at ${SITE}/llms.txt.\n\n`,
	);
	if (d.overview?.trim()) parts.push(`${d.overview.trim()}\n\n`);
	if (d.context?.trim()) parts.push(`${d.context.trim()}\n\n`);

	for (const node of d.nodes ?? []) {
		parts.push(`---\n\n## ${node.title}\n\n`);
		parts.push(`Source: ${abs(node.url)}\n\n`);
		if (node.summary?.trim()) parts.push(`${node.summary.trim()}\n\n`);
		if (node.facts?.length) {
			parts.push("Facts:\n\n");
			for (const fact of node.facts) {
				parts.push("```\n" + fact.literal.trim() + "\n```\n\n");
			}
		}
	}

	if (d.glossary?.length) {
		parts.push("---\n\n## Glossary\n\n");
		for (const term of d.glossary) {
			const aliases = term.aliases?.length
				? ` (aliases: ${term.aliases.join(", ")})`
				: "";
			parts.push(`- **${term.term}**${aliases}: ${term.definition}\n`);
		}
		parts.push("\n");
	}

	return new Response(parts.join(""), {
		headers: { "Content-Type": "text/plain; charset=utf-8" },
	});
};
