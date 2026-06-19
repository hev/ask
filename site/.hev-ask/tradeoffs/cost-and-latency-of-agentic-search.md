---
id: "tradeoffs#cost-and-latency-of-agentic-search"
title: "Tradeoffs"
heading: "Cost and latency of agentic search"
group: "Overview"
order: 102
url: "/docs/tradeoffs#cost-and-latency-of-agentic-search"
anchor: "cost-and-latency-of-agentic-search"
terms: ["cost","latency","agentic","search","path","calls","model","meaning","real","small","worst","case","roughly","iteration","count","round","trips","seconds","keyword","staying","instant","knob","bounded","loop","submitted","query","default","domain","context","prompt","cached","across","rounds","offline","build","uses","stronger","hash","gate","means"]
hash: "64c96faf9220cd1746ff5fe83b61b585b03f23674074ae444bfec6a16b4411de"
mode: "agent-primary"
facts: [{"kind":"code","literal":"maxIterations","chunkId":"tradeoffs#cost-and-latency-of-agentic-search"}]
sources: [{"chunkId":"tradeoffs#cost-and-latency-of-agentic-search","url":"/docs/tradeoffs#cost-and-latency-of-agentic-search","anchor":"cost-and-latency-of-agentic-search"}]
---

The agentic path calls the model, meaning real if small cost and latency: worst-case latency is roughly the iteration count of small-model round-trips (a few seconds) with the keyword path staying instant and the iteration cap as the knob, and cost is one bounded loop per submitted query on the default small model with domain context prompt-cached across rounds. The offline build uses a stronger model but the hash gate means you pay for it only when content changes, and keyword-only is a first-class mode if you want no key in the loop.
