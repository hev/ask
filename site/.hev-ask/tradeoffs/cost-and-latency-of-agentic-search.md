---
id: "tradeoffs#cost-and-latency-of-agentic-search"
title: "Tradeoffs"
heading: "Cost and latency of agentic search"
group: "Overview"
order: 102
url: "/docs/tradeoffs#cost-and-latency-of-agentic-search"
anchor: "cost-and-latency-of-agentic-search"
terms: ["cost","latency","agentic","search","path","calls","model","meaning","real","small","money","worst","case","roughly","iteration","worth","fast","round","trips","while","keyword","stays","instant","lane","bounded","loop","submitted","query","default","domain","context","prompt","cached","across","rounds","offline","build","uses","stronger","hash"]
hash: "64c96faf9220cd1746ff5fe83b61b585b03f23674074ae444bfec6a16b4411de"
mode: "agent-primary"
facts: [{"kind":"code","literal":"maxIterations","chunkId":"tradeoffs#cost-and-latency-of-agentic-search"}]
sources: [{"chunkId":"tradeoffs#cost-and-latency-of-agentic-search","url":"/docs/tradeoffs#cost-and-latency-of-agentic-search","anchor":"cost-and-latency-of-agentic-search"}]
---

The agentic path calls the model, meaning real if small money and latency: worst-case latency is roughly the iteration cap's worth of fast-model round-trips while the keyword path stays the instant lane, and cost is one bounded loop per submitted query on the default fast model with domain context prompt-cached across rounds. The offline build uses a stronger model but the hash gate means you pay for it only when content changes, and you can run keyword-only as a first-class mode if you want no key in the loop.
