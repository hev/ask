---
id: "digest-creation#verify-review-commit"
title: "Digest creation"
heading: "Verify, review, commit"
group: "Overview"
order: 74
url: "/docs/digest-creation#verify-review-commit"
anchor: "verify-review-commit"
terms: ["verify","review","commit","command","gate","builds","site","fails","section","anchor","missing","rendered","html","warns","coverage","fidelity","drift","after","tree","because","markdown","distilled","prose","grounded","facts","change","together","reviewable","diff","reviewed","pull","requests","digest","regenerates","only","content","changes","build","runs","runtime"]
hash: "3e04fbcdb9c42dcd79f213c0ffc5ff67e8eba56d67c1e18de255ae88432e8c8a"
mode: "agent-primary"
facts: [{"kind":"code","literal":"pnpm exec ask digest verify     # builds the site, checks every anchor resolves\ngit add .hev-ask","chunkId":"digest-creation#verify-review-commit"},{"kind":"code","literal":"ask digest verify","chunkId":"digest-creation#verify-review-commit"}]
sources: [{"chunkId":"digest-creation#verify-review-commit","url":"/docs/digest-creation#verify-review-commit","anchor":"verify-review-commit"}]
---

The verify command is the CI gate: it builds the site, fails when any section's anchor is missing from the rendered HTML, and warns on coverage or fidelity drift, after which you commit the tree. Because the tree is markdown, a section's distilled prose and grounded facts change together in one reviewable diff and are reviewed in pull requests; the digest regenerates only when content changes and a build runs, with the runtime logging a warning on hash mismatch as a rebuild cue while a stale digest degrades rather than breaks.
