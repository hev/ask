import { defineCollection, z } from "astro:content";
import { glob } from "astro/loaders";

const docs = defineCollection({
	loader: glob({ pattern: "**/*.{md,mdx}", base: "./src/content/docs" }),
	schema: z.object({
		title: z.string(),
		description: z.string(),
		group: z.string(),
		order: z.number(),
	}),
});

// Public slide decks. Each entry is one markdown file; reveal.js splits slides
// on a line of `---` (vertical sub-slides on `--`, speaker notes after `Note:`).
// Body is rendered full-screen by src/layouts/Deck.astro; the reveal runtime +
// theme are vendored under public/reveal/.
const decks = defineCollection({
	loader: glob({ pattern: "**/*.md", base: "./src/content/decks" }),
	schema: z.object({
		title: z.string(),
		description: z.string().optional(),
		pubDate: z.coerce.date(),
		draft: z.boolean().default(false),
		// Reachable by direct link, kept out of the listing.
		unlisted: z.boolean().default(false),
	}),
});

export const collections = { docs, decks };
