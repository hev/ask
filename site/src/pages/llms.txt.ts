import type { APIRoute } from "astro";
import { docsNav, getAllDocs, getDocHref } from "../lib/docs";
import { digestData } from "../lib/digest-data";

const SITE = "https://hevask.com";

// The map is the curated docs nav (human-ordered, page-level — the right shape
// for an llms.txt index). The orientation blurb and glossary come from the
// committed ask digest, so the index and the distilled full text at
// /llms-full.txt are sourced from the same artifact.
export const GET: APIRoute = async () => {
	const all = await getAllDocs();
	const byId = new Map(all.map((entry) => [entry.id, entry]));
	const orientation =
		digestData.context?.trim() ||
		"hev ask is an Astro integration that adds instant keyword search over heading anchors, plus an optional Claude-powered agentic search loop on Enter. The corpus is your content collection; an offline-built, committed ask digest gives the loop domain context and a glossary.";

	const lines: string[] = [];
	lines.push("# hev ask");
	lines.push("");
	lines.push("> A ⌘K search overlay for docs sites, built on a committed ask digest.");
	lines.push("");
	lines.push(orientation);
	lines.push("");
	lines.push(`The full distilled docs are at ${SITE}/llms-full.txt.`);
	lines.push("");

	for (const group of docsNav) {
		lines.push(`## ${group.label}`);
		for (const id of group.items) {
			const entry = byId.get(id);
			if (!entry) continue;
			const url = `${SITE}${getDocHref(id)}`;
			lines.push(`- [${entry.data.title}](${url}): ${entry.data.description}`);
		}
		lines.push("");
	}

	if (digestData.glossary?.length) {
		lines.push("## Glossary");
		for (const term of digestData.glossary) {
			const aliases = term.aliases?.length
				? ` (aliases: ${term.aliases.join(", ")})`
				: "";
			lines.push(`- **${term.term}**${aliases}: ${term.definition}`);
		}
		lines.push("");
	}

	return new Response(lines.join("\n"), {
		headers: { "Content-Type": "text/plain; charset=utf-8" },
	});
};
