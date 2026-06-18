---
id: "digest-creation#verify-review-commit"
title: "Digest creation"
heading: "Verify, review, commit"
group: "Overview"
order: 74
url: "/docs/digest-creation#verify-review-commit"
anchor: "verify-review-commit"
terms: ["verify","review","commit","step","gate","builds","site","fails","section","anchor","missing","rendered","html","warning","coverage","fidelity","drift","after","tree","because","markdown","prose","facts","change","together","reviewable","diff","digest","only","regenerates","content","changes","build","runs","runtime","logs","line","hash","mismatch","rebuild"]
hash: "3e04fbcdb9c42dcd79f213c0ffc5ff67e8eba56d67c1e18de255ae88432e8c8a"
mode: "agent-primary"
facts: [{"kind":"code","literal":"pnpm exec ask digest verify     # builds the site, checks every anchor resolves\ngit add .hev-ask","chunkId":"digest-creation#verify-review-commit"},{"kind":"code","literal":"ask digest verify","chunkId":"digest-creation#verify-review-commit"}]
sources: [{"chunkId":"digest-creation#verify-review-commit","url":"/docs/digest-creation#verify-review-commit","anchor":"verify-review-commit"}]
---

The verify step is the CI gate: it builds the site and fails when any section's anchor is missing from the rendered HTML, warning on coverage or fidelity drift, after which you commit the tree. Because the tree is markdown, a section's prose and facts change together in one reviewable diff, the digest only regenerates when content changes and a build runs, and the runtime logs a one-line warning on a hash mismatch as a cue to rebuild while a stale digest keeps serving.
