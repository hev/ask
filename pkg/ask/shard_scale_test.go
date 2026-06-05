package ask

// Cloudflare-docs-scale benchmark for the sharded digest pipeline. The target
// shape is ~100 product dirs / ~5k MDX files / ~30MB of content — the size of
// developers.cloudflare.com's docs repo. Distillation is stubbed (excerpt
// summaries), so this proves the deterministic pipeline — shard, manifest,
// merge, assemble — fits the benchmark corpus without any model spend.
//
// Heavy by design; run explicitly with:
//
//	ASK_SCALE_BENCH=1 go test ./pkg/ask -run TestShardPipelineAtCloudflareScale -v

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestShardPipelineAtCloudflareScale(t *testing.T) {
	if os.Getenv("ASK_SCALE_BENCH") == "" {
		t.Skip("set ASK_SCALE_BENCH=1 to run the Cloudflare-scale benchmark")
	}

	root := t.TempDir()
	const products = 100
	const filesPerProduct = 50

	written := 0
	for p := 0; p < products; p++ {
		for f := 0; f < filesPerProduct; f++ {
			dir := filepath.Join(root, "src/content/docs", fmt.Sprintf("product-%02d", p))
			if f%10 != 0 {
				dir = filepath.Join(dir, fmt.Sprintf("area-%d", f/10))
			}
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			body := syntheticDoc(p, f)
			written += len(body)
			if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("page-%02d.mdx", f)), []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	t.Logf("synthetic corpus: %d files, %.1f MB", products*filesPerProduct, float64(written)/1e6)

	options := BuildOptions{SiteRoot: root, Collections: []string{"docs"}, BasePath: "/", ChunkHeadingDepth: 3}

	start := time.Now()
	result, err := WriteCorpusShards(options, ".hev-ask/shards", DefaultShardBytes)
	if err != nil {
		t.Fatal(err)
	}
	corpusDuration := time.Since(start)
	t.Logf("corpus+shard: %d sections, %d shards in %s (%s heap)", result.Sections, result.Shards, corpusDuration, heapMB())

	// Stub-distil every shard from its input file — stands in for the skill's
	// per-shard agents.
	dir := filepath.Join(root, ".hev-ask/shards")
	manifest, err := readShardManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	start = time.Now()
	for _, entry := range manifest.Shards {
		data, err := os.ReadFile(filepath.Join(dir, "input-"+entry.ID+".json"))
		if err != nil {
			t.Fatal(err)
		}
		var input ShardInput
		if err := json.Unmarshal(data, &input); err != nil {
			t.Fatal(err)
		}
		distill := ShardDistillation{ShardHash: input.ShardHash, Notes: "Covers " + entry.ID}
		for _, sec := range input.Sections {
			distill.Summaries = append(distill.Summaries, SectionSummaryIn{ID: sec.ID, Summary: textExcerpt(sec.Text, 160)})
		}
		out, err := json.Marshal(distill)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "distill-"+entry.ID+".json"), out, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	globalJSON := `{"context":"Synthetic platform docs.","suggestions":["How do I get started?","How do I deploy?"],"glossary":[]}`
	if err := os.WriteFile(filepath.Join(dir, "global.json"), []byte(globalJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Logf("stub distillation: %d shards in %s", len(manifest.Shards), time.Since(start))

	start = time.Now()
	assembled, err := AssembleFromShards(options, ".hev-ask/shards")
	if err != nil {
		t.Fatal(err)
	}
	assembleDuration := time.Since(start)
	if len(assembled.SkippedShards) != 0 || len(assembled.Missing) != 0 {
		t.Fatalf("expected full coverage: skipped=%d missing=%d", len(assembled.SkippedShards), len(assembled.Missing))
	}
	info, err := os.Stat(filepath.Join(root, ".hev-ask/digest.json"))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("assemble: %d chunks from %d shards in %s; digest.json %.1f MB (%s heap)",
		assembled.Chunks, assembled.Shards, assembleDuration, float64(info.Size())/1e6, heapMB())

	// Refresh: a single-doc edit must re-pend exactly one shard.
	page := filepath.Join(root, "src/content/docs/product-42/area-2/page-25.mdx")
	original, err := os.ReadFile(page)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(page, append(original, []byte("\nEdited paragraph.\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	start = time.Now()
	result, err = WriteCorpusShards(options, ".hev-ask/shards", DefaultShardBytes)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("incremental re-corpus after 1-doc edit: %d pending of %d shards in %s", result.Pending, result.Shards, time.Since(start))
	if result.Pending != 1 {
		t.Fatalf("expected exactly 1 pending shard after a one-doc edit, got %d", result.Pending)
	}
}

func syntheticDoc(p int, f int) string {
	page := fmt.Sprintf("Product %02d Page %02d", p, f)
	body := fmt.Sprintf("---\ntitle: \"%s\"\ndescription: \"Reference for %s.\"\ngroup: \"Product %02d\"\n---\n\n", page, page, p)
	body += fmt.Sprintf("This page explains how %s works, when to reach for it, and the limits that apply to it in production deployments.\n\n", page)
	for s := 0; s < 4; s++ {
		body += fmt.Sprintf("## Topic %d on page %d\n\n", s, f)
		for para := 0; para < 5; para++ {
			body += fmt.Sprintf(
				"The `--flag-%d-%d` option controls behavior %d for this feature. Combine it with the binding and the platform schedules work across regions, retries transient failures, and reports metrics to the dashboard so operators can audit throughput, latency, and error budgets over time. ",
				s, para, para)
		}
		body += fmt.Sprintf("\n\n```sh\nnpx tool deploy --product %d --topic %d\n```\n\n", p, s)
	}
	return body
}

func heapMB() string {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	return fmt.Sprintf("%.0fMB", float64(stats.HeapAlloc)/1e6)
}
