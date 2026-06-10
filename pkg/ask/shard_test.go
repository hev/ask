package ask

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func section(id string, text string) CorpusSection {
	return CorpusSection{ID: id, URL: "/" + slugOf(id), Title: id, Text: text}
}

func TestPlanShardsGroupsByPrefixAndMergesSmallRuns(t *testing.T) {
	sections := []CorpusSection{
		section("alpha#intro", strings.Repeat("a", 40)),
		section("beta#intro", strings.Repeat("b", 40)),
		section("workers/start#intro", strings.Repeat("w", 120)),
		section("workers/start#deploy", strings.Repeat("w", 120)),
	}
	shards := PlanShards(sections, 100)
	if len(shards) != 2 {
		t.Fatalf("expected 2 shards, got %d: %+v", len(shards), shards)
	}
	// alpha+beta merge into one run under the cap; workers stands alone (oversize leaf).
	if got := shards[0].ID; got != "alpha--beta" {
		t.Fatalf("merged shard id: %s", got)
	}
	if got := len(shards[0].Sections); got != 2 {
		t.Fatalf("merged shard sections: %d", got)
	}
	if got := shards[1].ID; got != "workers" {
		t.Fatalf("oversize shard id: %s", got)
	}
	if shards[1].Bytes <= 100 {
		t.Fatalf("workers shard should be oversize, got %d bytes", shards[1].Bytes)
	}
}

func TestPlanShardsSplitsOversizeGroupsOneSegmentDeeper(t *testing.T) {
	sections := []CorpusSection{
		section("workers/ai#intro", strings.Repeat("x", 90)),
		section("workers/kv#intro", strings.Repeat("x", 90)),
		section("workers/queues#intro", strings.Repeat("x", 90)),
	}
	shards := PlanShards(sections, 100)
	if len(shards) != 3 {
		t.Fatalf("expected 3 shards after deep split, got %d", len(shards))
	}
	if shards[0].ID != "workers__ai" {
		t.Fatalf("deep shard id: %s", shards[0].ID)
	}
}

func TestPlanShardsIsStableUnderSingleDocEdits(t *testing.T) {
	var sections []CorpusSection
	for product := 0; product < 6; product++ {
		for page := 0; page < 4; page++ {
			id := fmt.Sprintf("p%d/page%d#body", product, page)
			sections = append(sections, section(id, strings.Repeat("t", 500)))
		}
	}
	before := PlanShards(sections, 2_500)

	// Edit one doc's text (same size class) — membership must not move and only
	// the owning shard's hash may change.
	edited := append([]CorpusSection{}, sections...)
	edited[5] = section(edited[5].ID, strings.Repeat("u", 500)) // p1/page1
	after := PlanShards(edited, 2_500)

	if len(before) != len(after) {
		t.Fatalf("shard count changed: %d -> %d", len(before), len(after))
	}
	changed := 0
	for i := range before {
		if before[i].ID != after[i].ID {
			t.Fatalf("shard %d id changed: %s -> %s", i, before[i].ID, after[i].ID)
		}
		if before[i].ShardHash != after[i].ShardHash {
			changed++
		}
	}
	if changed != 1 {
		t.Fatalf("expected exactly 1 shard hash to change, got %d", changed)
	}
}

func TestMergeDistillationsValidatesAndMerges(t *testing.T) {
	validIDs := map[string]bool{"a#x": true, "a#y": true, "b#z": true}
	globalPart := EmittedDistillation{
		Context:     "Site context.",
		Glossary:    []GlossaryEntry{{Term: "Overlay", Aliases: []string{"popover"}, Definition: "The UI."}},
		Suggestions: []string{"q1", "q2", "q3", "q4", "q5", "q6"},
	}
	parts := []shardPart{
		{shardID: "a", part: EmittedDistillation{
			Glossary:  []GlossaryEntry{{Term: "overlay", Aliases: []string{"search overlay"}, Definition: "Duplicate term, different case."}},
			Summaries: []SectionSummaryIn{{ID: "a#x", Summary: "X."}, {ID: "a#y", Summary: "Y."}},
		}},
		{shardID: "b", part: EmittedDistillation{}},
	}
	result, err := mergeDistillations(globalPart, parts, validIDs)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(result.merged.Summaries); got != 2 {
		t.Fatalf("merged summaries: %d", got)
	}
	if got := result.missing; len(got) != 1 || got[0] != "b#z" {
		t.Fatalf("missing ids: %v", got)
	}
	if got := len(result.merged.Glossary); got != 1 {
		t.Fatalf("glossary should dedupe case-insensitively, got %d entries", got)
	}
	aliases := result.merged.Glossary[0].Aliases
	if len(aliases) != 2 {
		t.Fatalf("aliases should union: %v", aliases)
	}
	if got := len(result.merged.Suggestions); got != 5 {
		t.Fatalf("suggestions should cap at 5, got %d", got)
	}
	if result.merged.Context != "Site context." {
		t.Fatalf("context must come from the global part: %q", result.merged.Context)
	}
}

func TestMergeDistillationsRejectsUnknownAndDuplicateIDs(t *testing.T) {
	validIDs := map[string]bool{"a#x": true}
	_, err := mergeDistillations(EmittedDistillation{}, []shardPart{
		{shardID: "a", part: EmittedDistillation{Summaries: []SectionSummaryIn{{ID: "ghost#id", Summary: "S."}}}},
	}, validIDs)
	if err == nil || !strings.Contains(err.Error(), "unknown section id") {
		t.Fatalf("expected unknown-id error, got %v", err)
	}

	_, err = mergeDistillations(EmittedDistillation{}, []shardPart{
		{shardID: "a", part: EmittedDistillation{Summaries: []SectionSummaryIn{{ID: "a#x", Summary: "S."}}}},
		{shardID: "b", part: EmittedDistillation{Summaries: []SectionSummaryIn{{ID: "a#x", Summary: "S2."}}}},
	}, validIDs)
	if err == nil || !strings.Contains(err.Error(), "must not overlap") {
		t.Fatalf("expected overlap error, got %v", err)
	}
}

func TestMergeDistillationsCapsGlossary(t *testing.T) {
	var entries []GlossaryEntry
	for i := 0; i < GlossaryCap+10; i++ {
		entries = append(entries, GlossaryEntry{Term: fmt.Sprintf("term-%d", i), Aliases: []string{}, Definition: "D."})
	}
	result, err := mergeDistillations(EmittedDistillation{Glossary: entries}, nil, map[string]bool{})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(result.merged.Glossary); got != GlossaryCap {
		t.Fatalf("glossary cap: %d", got)
	}
	if result.glossaryDropped != 10 {
		t.Fatalf("glossaryDropped: %d", result.glossaryDropped)
	}
}

func writeShardDistill(t *testing.T, dir string, shardID string, distill ShardDistillation) {
	t.Helper()
	data, err := json.Marshal(distill)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "distill-"+shardID+".json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestShardedCorpusRoundTrip(t *testing.T) {
	root := writeParityFixture(t)
	options := BuildOptions{SiteRoot: root, Collections: []string{"docs"}, BasePath: "/docs/", ChunkHeadingDepth: 3}

	// Tiny cap forces multiple shards even on the small fixture.
	result, err := WriteCorpusShards(options, ".hev-ask/shards", 200)
	if err != nil {
		t.Fatal(err)
	}
	if result.Shards < 2 {
		t.Fatalf("expected multiple shards, got %d", result.Shards)
	}
	if result.Pending != result.Shards {
		t.Fatalf("all shards should start pending: %d/%d", result.Pending, result.Shards)
	}

	dir := filepath.Join(root, ".hev-ask/shards")
	manifest, err := readShardManifest(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Distil every shard from its input file (stub summaries) + write global.json.
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
			distill.Summaries = append(distill.Summaries, SectionSummaryIn{ID: sec.ID, Summary: "Stub for " + sec.ID})
		}
		writeShardDistill(t, dir, entry.ID, distill)
	}
	globalJSON := `{"context":"Fixture docs.","suggestions":["How do I install?"],"glossary":[]}`
	if err := os.WriteFile(filepath.Join(dir, "global.json"), []byte(globalJSON), 0o600); err != nil {
		t.Fatal(err)
	}

	assembled, err := AssembleFromShards(options, ".hev-ask/shards")
	if err != nil {
		t.Fatal(err)
	}
	if len(assembled.SkippedShards) != 0 || len(assembled.Missing) != 0 {
		t.Fatalf("expected full coverage, skipped=%v missing=%v", assembled.SkippedShards, assembled.Missing)
	}
	graph, err := LoadDigest(filepath.Join(root, ".hev-ask"))
	if err != nil {
		t.Fatal(err)
	}
	if graph.Context != "Fixture docs." {
		t.Fatalf("context: %q", graph.Context)
	}
	for _, node := range graph.Nodes {
		if !strings.HasPrefix(node.Summary, "Stub for ") {
			t.Fatalf("node %s did not get its shard summary: %q", node.ID, node.Summary)
		}
	}

	// Resume: re-running corpus keeps every distillation.
	result, err = WriteCorpusShards(options, ".hev-ask/shards", 200)
	if err != nil {
		t.Fatal(err)
	}
	if result.Pending != 0 {
		t.Fatalf("expected 0 pending after distillation, got %d", result.Pending)
	}

	// Status agrees and sees global.json.
	status, err := ShardStatus(root, ".hev-ask/shards")
	if err != nil {
		t.Fatal(err)
	}
	if !status.HasGlobal {
		t.Fatal("status should report global.json present")
	}
	for _, shard := range status.Shards {
		if shard.State != ShardDistilled {
			t.Fatalf("shard %s state: %s", shard.ID, shard.State)
		}
	}

	// Content edit: only the owning shard goes pending again, and assemble
	// refuses the stale manifest until corpus is re-run.
	configPath := filepath.Join(root, "src/content/docs/api/config.mdx")
	original, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, append(original, []byte("\nNew paragraph.\n")...), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := AssembleFromShards(options, ".hev-ask/shards"); err == nil || !strings.Contains(err.Error(), "content changed") {
		t.Fatalf("expected content-changed error, got %v", err)
	}
	result, err = WriteCorpusShards(options, ".hev-ask/shards", 200)
	if err != nil {
		t.Fatal(err)
	}
	if result.Pending != 1 {
		t.Fatalf("expected exactly 1 pending shard after one-doc edit, got %d", result.Pending)
	}

	// Partial assemble still works: the stale shard's sections fall back.
	assembled, err = AssembleFromShards(options, ".hev-ask/shards")
	if err != nil {
		t.Fatal(err)
	}
	if len(assembled.SkippedShards) != 1 {
		t.Fatalf("expected 1 skipped shard, got %v", assembled.SkippedShards)
	}
	if len(assembled.Missing) == 0 {
		t.Fatal("expected excerpt fallbacks for the stale shard")
	}
}

func TestWriteCorpusShardsRemovesOrphans(t *testing.T) {
	root := writeParityFixture(t)
	options := BuildOptions{SiteRoot: root, Collections: []string{"docs"}, BasePath: "/docs/", ChunkHeadingDepth: 3}
	dir := filepath.Join(root, ".hev-ask/shards")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"input-ghost.json", "distill-ghost.json"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("{}"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := WriteCorpusShards(options, ".hev-ask/shards", 200); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"input-ghost.json", "distill-ghost.json"} {
		if _, err := os.Stat(filepath.Join(dir, name)); !os.IsNotExist(err) {
			t.Fatalf("orphan %s should be removed", name)
		}
	}
}
