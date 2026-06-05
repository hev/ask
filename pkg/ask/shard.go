package ask

// Sharded digest builds. A corpus too large for one model context is split
// into shards along slug-prefix boundaries; each shard is distilled in its own
// fresh context (by the build-digest skill), and the shards are merged back
// into a single distillation before assembly. Planning and merging here are
// deterministic and keyless — the model never runs in this file.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	// DefaultShardBytes is ~50k tokens of section text per shard — comfortable
	// for one fresh distillation context.
	DefaultShardBytes = 200_000
	// GlossaryCap bounds the merged glossary so the runtime prompt and
	// prefilter stay bounded no matter how many shards contribute terms.
	GlossaryCap = 75
)

type PlannedShard struct {
	// ID is filesystem-safe, derived from the slug prefix(es) this shard covers.
	ID string
	// Prefixes are the slug prefixes covered: one, or a merged lexicographic run.
	Prefixes []string
	Sections []CorpusSection
	Bytes    int
	// ShardHash covers this shard's section ids + text. A distillation is valid
	// only for this exact hash.
	ShardHash string
}

type ShardManifestEntry struct {
	ID        string   `json:"id"`
	Prefixes  []string `json:"prefixes"`
	Sections  int      `json:"sections"`
	Bytes     int      `json:"bytes"`
	ShardHash string   `json:"shardHash"`
	// Distilled is a snapshot at manifest-write time: a distillation for this
	// exact shardHash exists on disk.
	Distilled bool `json:"distilled"`
}

type ShardManifest struct {
	Version     int                  `json:"version"`
	ContentHash string               `json:"contentHash"`
	DigestPath  string               `json:"digestPath"`
	UpToDate    bool                 `json:"upToDate"`
	ShardBytes  int                  `json:"shardBytes"`
	Shards      []ShardManifestEntry `json:"shards"`
}

// ShardInput is one shard's model-input payload (`input-<id>.json`). The
// distiller must copy ShardHash into its distill file — that is the proof the
// distillation matches this exact content.
type ShardInput struct {
	ShardID     string          `json:"shardId"`
	ShardHash   string          `json:"shardHash"`
	ContentHash string          `json:"contentHash"`
	Sections    []CorpusSection `json:"sections"`
}

// ShardDistillation is the skill-authored `distill-<id>.json`: the shard's
// distilled fields plus the hash of the input it was distilled from and a
// short gist used by the final global synthesis pass.
type ShardDistillation struct {
	ShardHash string `json:"shardHash"`
	Notes     string `json:"notes"`
	EmittedDistillation
}

// PlanShards groups the corpus sections into shards by slug path prefix —
// `workers/observability/...` style — so shard membership is stable: editing a
// doc re-hashes only the shard that owns its prefix. Oversize groups split one
// path segment deeper (a single doc never splits); small adjacent groups are
// greedily merged up to maxBytes.
func PlanShards(sections []CorpusSection, maxBytes int) []PlannedShard {
	if maxBytes <= 0 {
		maxBytes = DefaultShardBytes
	}
	sorted := append([]CorpusSection{}, sections...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID < sorted[j].ID })
	leaves := splitByPrefix(sorted, 0, maxBytes)

	// Greedy run-merge of adjacent leaves. Deterministic; membership only
	// shifts when a neighbor's size crosses the cap.
	type run struct {
		prefixes []string
		sections []CorpusSection
		bytes    int
	}
	var merged []run
	for _, leaf := range leaves {
		if n := len(merged); n > 0 && merged[n-1].bytes+leaf.bytes <= maxBytes {
			merged[n-1].prefixes = append(merged[n-1].prefixes, leaf.prefix)
			merged[n-1].sections = append(merged[n-1].sections, leaf.sections...)
			merged[n-1].bytes += leaf.bytes
			continue
		}
		merged = append(merged, run{prefixes: []string{leaf.prefix}, sections: leaf.sections, bytes: leaf.bytes})
	}

	used := map[string]bool{}
	shards := make([]PlannedShard, 0, len(merged))
	for _, shard := range merged {
		shards = append(shards, PlannedShard{
			ID:        uniqueShardID(shardIDFor(shard.prefixes), used),
			Prefixes:  shard.prefixes,
			Sections:  shard.sections,
			Bytes:     shard.bytes,
			ShardHash: ShardHash(shard.sections),
		})
	}
	return shards
}

type prefixLeaf struct {
	prefix   string
	sections []CorpusSection
	bytes    int
}

func splitByPrefix(sections []CorpusSection, depth int, maxBytes int) []prefixLeaf {
	groups := map[string][]CorpusSection{}
	for _, section := range sections {
		key := prefixAtDepth(slugOf(section.ID), depth)
		groups[key] = append(groups[key], section)
	}
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var leaves []prefixLeaf
	for _, prefix := range keys {
		members := groups[prefix]
		bytes := sectionBytes(members)
		if bytes <= maxBytes {
			leaves = append(leaves, prefixLeaf{prefix: prefix, sections: members, bytes: bytes})
			continue
		}
		// Can we make progress one segment deeper? A single doc (every slug
		// equal to the prefix) cannot split — it ships as one oversize leaf.
		deeper := map[string]bool{}
		for _, member := range members {
			deeper[prefixAtDepth(slugOf(member.ID), depth+1)] = true
		}
		if len(deeper) <= 1 {
			leaves = append(leaves, prefixLeaf{prefix: prefix, sections: members, bytes: bytes})
			continue
		}
		leaves = append(leaves, splitByPrefix(members, depth+1, maxBytes)...)
	}
	return leaves
}

func slugOf(sectionID string) string {
	if i := strings.IndexByte(sectionID, '#'); i >= 0 {
		return sectionID[:i]
	}
	return sectionID
}

func prefixAtDepth(slug string, depth int) string {
	parts := strings.Split(slug, "/")
	if len(parts) > depth+1 {
		parts = parts[:depth+1]
	}
	return strings.Join(parts, "/")
}

func sectionBytes(sections []CorpusSection) int {
	total := 0
	for _, section := range sections {
		total += len(section.Text)
	}
	return total
}

// ShardHash hashes a shard's section ids + text.
func ShardHash(sections []CorpusSection) string {
	hash := sha256.New()
	for _, section := range sections {
		hash.Write([]byte(section.ID + " " + section.Text + " "))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

var shardIDUnsafe = regexp.MustCompile(`[^a-zA-Z0-9_.-]`)

func shardIDFor(prefixes []string) string {
	first := sanitizePrefix(prefixes[0])
	if len(prefixes) == 1 {
		return first
	}
	return first + "--" + sanitizePrefix(prefixes[len(prefixes)-1])
}

func sanitizePrefix(prefix string) string {
	cleaned := shardIDUnsafe.ReplaceAllString(strings.ReplaceAll(prefix, "/", "__"), "-")
	if cleaned == "" {
		return "root"
	}
	return cleaned
}

func uniqueShardID(base string, used map[string]bool) string {
	id := base
	for n := 2; used[id]; n++ {
		id = fmt.Sprintf("%s-%d", base, n)
	}
	used[id] = true
	return id
}

// ---------------------------------------------------------------------------
// Merging shard distillations back into one.
// ---------------------------------------------------------------------------

type shardPart struct {
	shardID string
	part    EmittedDistillation
}

type mergeResult struct {
	merged EmittedDistillation
	// missing are section ids with no summary in any shard — they fall back to
	// excerpts at assembly.
	missing []string
	// glossaryDropped counts entries dropped by the cap.
	glossaryDropped int
}

// mergeDistillations merges per-shard distillations with the global synthesis.
// Context and suggestions come from the global part only; summaries union
// across shards (duplicate or unknown ids are hard errors naming the shard);
// glossaries are deduped case-insensitively and capped.
func mergeDistillations(globalPart EmittedDistillation, shardParts []shardPart, validIDs map[string]bool) (mergeResult, error) {
	ownerByID := map[string]string{}
	summaryByID := map[string]string{}
	for _, shard := range shardParts {
		for _, entry := range shard.part.Summaries {
			if !validIDs[entry.ID] {
				return mergeResult{}, fmt.Errorf(
					"distill shard %q summarizes unknown section id %q — it is stale; re-run `ask digest corpus --shards-dir …` and re-distil that shard",
					shard.shardID, entry.ID,
				)
			}
			if owner, taken := ownerByID[entry.ID]; taken {
				return mergeResult{}, fmt.Errorf("section id %q is summarized by both %q and %q — shards must not overlap", entry.ID, owner, shard.shardID)
			}
			ownerByID[entry.ID] = shard.shardID
			summaryByID[entry.ID] = entry.Summary
		}
	}

	var missing []string
	for id := range validIDs {
		if _, ok := summaryByID[id]; !ok {
			missing = append(missing, id)
		}
	}
	sort.Strings(missing)

	glossary := []GlossaryEntry{}
	seenTerms := map[string]int{}
	glossaryDropped := 0
	allEntries := append([]GlossaryEntry{}, globalPart.Glossary...)
	for _, shard := range shardParts {
		allEntries = append(allEntries, shard.part.Glossary...)
	}
	for _, entry := range allEntries {
		key := strings.ToLower(strings.TrimSpace(entry.Term))
		if key == "" {
			continue
		}
		if at, ok := seenTerms[key]; ok {
			// First definition wins; aliases union.
			existing := &glossary[at]
			for _, alias := range entry.Aliases {
				duplicate := false
				for _, have := range existing.Aliases {
					if strings.EqualFold(have, alias) {
						duplicate = true
						break
					}
				}
				if !duplicate {
					existing.Aliases = append(existing.Aliases, alias)
				}
			}
			continue
		}
		if len(glossary) >= GlossaryCap {
			glossaryDropped++
			continue
		}
		seenTerms[key] = len(glossary)
		aliases := entry.Aliases
		if aliases == nil {
			aliases = []string{}
		}
		glossary = append(glossary, GlossaryEntry{Term: entry.Term, Aliases: aliases, Definition: entry.Definition})
	}

	summaries := make([]SectionSummaryIn, 0, len(summaryByID))
	for id, summary := range summaryByID {
		summaries = append(summaries, SectionSummaryIn{ID: id, Summary: summary})
	}
	sort.Slice(summaries, func(i, j int) bool { return summaries[i].ID < summaries[j].ID })

	suggestions := globalPart.Suggestions
	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	return mergeResult{
		merged: EmittedDistillation{
			Context:     globalPart.Context,
			Glossary:    glossary,
			Summaries:   summaries,
			Suggestions: suggestions,
		},
		missing:         missing,
		glossaryDropped: glossaryDropped,
	}, nil
}

// ---------------------------------------------------------------------------
// File IO: sharded corpus, shard-merging assemble, status.
// ---------------------------------------------------------------------------

type CorpusShardsResult struct {
	Dir      string
	Sections int
	Shards   int
	Pending  int
	UpToDate bool
}

// WriteCorpusShards splits the corpus into prefix-stable shards and writes one
// `input-<id>.json` per shard plus `manifest.json`. Re-running is the resume
// mechanism: a shard whose hash is unchanged keeps its existing distillation;
// a changed shard is marked pending again. Orphan shard files are removed.
func WriteCorpusShards(options BuildOptions, shardsDir string, shardBytes int) (CorpusShardsResult, error) {
	normalizeBuildOptions(&options)
	if shardBytes <= 0 {
		shardBytes = DefaultShardBytes
	}
	corpus, err := BuildCorpus(options)
	if err != nil {
		return CorpusShardsResult{}, err
	}
	committed, err := LoadDigest(resolveSitePath(options.SiteRoot, options.DigestPath))
	upToDate := err == nil && committed.Version == 2 && committed.ContentHash == corpus.ContentHash && len(committed.Nodes) > 0

	dir := resolveSitePath(options.SiteRoot, shardsDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return CorpusShardsResult{}, err
	}
	planned := PlanShards(CorpusSections(corpus), shardBytes)

	entries := make([]ShardManifestEntry, 0, len(planned))
	pending := 0
	for _, shard := range planned {
		input := ShardInput{ShardID: shard.ID, ShardHash: shard.ShardHash, ContentHash: corpus.ContentHash, Sections: shard.Sections}
		data, err := json.MarshalIndent(input, "", "  ")
		if err != nil {
			return CorpusShardsResult{}, err
		}
		if err := os.WriteFile(filepath.Join(dir, "input-"+shard.ID+".json"), append(data, '\n'), 0o644); err != nil {
			return CorpusShardsResult{}, err
		}
		distill, _ := readShardDistillation(dir, shard.ID)
		distilled := distill != nil && distill.ShardHash == shard.ShardHash
		if !distilled {
			pending++
		}
		entries = append(entries, ShardManifestEntry{
			ID:        shard.ID,
			Prefixes:  shard.Prefixes,
			Sections:  len(shard.Sections),
			Bytes:     shard.Bytes,
			ShardHash: shard.ShardHash,
			Distilled: distilled,
		})
	}

	manifest := ShardManifest{
		Version:     1,
		ContentHash: corpus.ContentHash,
		DigestPath:  options.DigestPath,
		UpToDate:    upToDate,
		ShardBytes:  shardBytes,
		Shards:      entries,
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return CorpusShardsResult{}, err
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), append(data, '\n'), 0o644); err != nil {
		return CorpusShardsResult{}, err
	}
	if err := removeOrphanShardFiles(dir, planned); err != nil {
		return CorpusShardsResult{}, err
	}

	return CorpusShardsResult{Dir: dir, Sections: len(corpus.Chunks), Shards: len(entries), Pending: pending, UpToDate: upToDate}, nil
}

type AssembleFromShardsResult struct {
	BuildResult
	Shards int
	// SkippedShards were not (or stalely) distilled — their sections fell back
	// to excerpts.
	SkippedShards []string
	// Missing are section ids that fell back to deterministic excerpts.
	Missing         []string
	GlossaryDropped int
}

// AssembleFromShards merges every current shard distillation plus the global
// synthesis (`global.json`) and writes the committed digest. Shards that are
// pending or stale are skipped with a warning — their sections fall back to
// excerpts — so the digest is usable after every distillation wave.
func AssembleFromShards(options BuildOptions, shardsDir string) (AssembleFromShardsResult, error) {
	normalizeBuildOptions(&options)
	dir := resolveSitePath(options.SiteRoot, shardsDir)
	manifest, err := readShardManifest(dir)
	if err != nil {
		return AssembleFromShardsResult{}, fmt.Errorf("no shard manifest in %s — run `ask digest corpus --shards-dir %s` first", shardsDir, shardsDir)
	}

	corpus, err := BuildCorpus(options)
	if err != nil {
		return AssembleFromShardsResult{}, err
	}
	if corpus.ContentHash != manifest.ContentHash {
		return AssembleFromShardsResult{}, fmt.Errorf(
			"content changed since the corpus was sharded — re-run `ask digest corpus --shards-dir %s` (and re-distil any shards it marks pending)", shardsDir,
		)
	}

	globalPart, err := readGlobalDistillation(dir, shardsDir)
	if err != nil {
		return AssembleFromShardsResult{}, err
	}

	var parts []shardPart
	var skipped []string
	for _, entry := range manifest.Shards {
		distill, _ := readShardDistillation(dir, entry.ID)
		if distill == nil || distill.ShardHash != entry.ShardHash {
			skipped = append(skipped, entry.ID)
			continue
		}
		part := distill.EmittedDistillation
		normalizeDistillation(&part)
		parts = append(parts, shardPart{shardID: entry.ID, part: part})
	}

	validIDs := make(map[string]bool, len(corpus.Chunks))
	for _, chunk := range corpus.Chunks {
		validIDs[chunk.ID] = true
	}
	merge, err := mergeDistillations(globalPart, parts, validIDs)
	if err != nil {
		return AssembleFromShardsResult{}, err
	}

	graph := AssembleDigest(merge.merged, corpus)
	out := resolveSitePath(options.SiteRoot, options.DigestPath)
	if err := WriteDigest(out, graph); err != nil {
		return AssembleFromShardsResult{}, err
	}
	return AssembleFromShardsResult{
		BuildResult:     BuildResult{Status: "built", Path: out, ContentHash: corpus.ContentHash, Chunks: len(corpus.Chunks)},
		Shards:          len(manifest.Shards),
		SkippedShards:   skipped,
		Missing:         merge.missing,
		GlossaryDropped: merge.glossaryDropped,
	}, nil
}

type ShardState string

const (
	ShardDistilled ShardState = "distilled"
	ShardPending   ShardState = "pending"
	ShardStale     ShardState = "stale"
)

type ShardStatusEntry struct {
	ShardManifestEntry
	State ShardState
}

type ShardStatusResult struct {
	Dir         string
	ContentHash string
	UpToDate    bool
	HasGlobal   bool
	Shards      []ShardStatusEntry
}

// ShardStatus is a disk-only report of shard coverage — what is distilled,
// pending, or stale.
func ShardStatus(siteRoot string, shardsDir string) (ShardStatusResult, error) {
	if siteRoot == "" {
		siteRoot = "."
	}
	dir := resolveSitePath(siteRoot, shardsDir)
	manifest, err := readShardManifest(dir)
	if err != nil {
		return ShardStatusResult{}, fmt.Errorf("no shard manifest in %s — run `ask digest corpus --shards-dir %s` first", shardsDir, shardsDir)
	}
	shards := make([]ShardStatusEntry, 0, len(manifest.Shards))
	for _, entry := range manifest.Shards {
		distill, _ := readShardDistillation(dir, entry.ID)
		state := ShardPending
		if distill != nil {
			if distill.ShardHash == entry.ShardHash {
				state = ShardDistilled
			} else {
				state = ShardStale
			}
		}
		entry.Distilled = state == ShardDistilled
		shards = append(shards, ShardStatusEntry{ShardManifestEntry: entry, State: state})
	}
	hasGlobal := false
	if data, err := os.ReadFile(filepath.Join(dir, "global.json")); err == nil {
		hasGlobal = json.Valid(data)
	}
	return ShardStatusResult{Dir: dir, ContentHash: manifest.ContentHash, UpToDate: manifest.UpToDate, HasGlobal: hasGlobal, Shards: shards}, nil
}

func readShardManifest(dir string) (ShardManifest, error) {
	data, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		return ShardManifest{}, err
	}
	var manifest ShardManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return ShardManifest{}, err
	}
	if manifest.Version != 1 {
		return ShardManifest{}, fmt.Errorf("unsupported shard manifest version %d", manifest.Version)
	}
	return manifest, nil
}

func readShardDistillation(dir string, shardID string) (*ShardDistillation, error) {
	data, err := os.ReadFile(filepath.Join(dir, "distill-"+shardID+".json"))
	if err != nil {
		return nil, err
	}
	var distill ShardDistillation
	if err := json.Unmarshal(data, &distill); err != nil {
		return nil, fmt.Errorf("parse distill-%s.json: %w", shardID, err)
	}
	if distill.ShardHash == "" {
		return nil, fmt.Errorf("distill-%s.json is missing shardHash", shardID)
	}
	return &distill, nil
}

func readGlobalDistillation(dir string, shardsDir string) (EmittedDistillation, error) {
	data, err := os.ReadFile(filepath.Join(dir, "global.json"))
	if err != nil {
		return EmittedDistillation{}, fmt.Errorf(
			"could not read %s/global.json — the global synthesis ({context, suggestions, glossary?}) is required to assemble a sharded digest", shardsDir,
		)
	}
	var globalPart EmittedDistillation
	if err := json.Unmarshal(data, &globalPart); err != nil {
		return EmittedDistillation{}, fmt.Errorf("parse %s/global.json: %w", shardsDir, err)
	}
	normalizeDistillation(&globalPart)
	return globalPart, nil
}

var shardFileName = regexp.MustCompile(`^(input|distill)-(.+)\.json$`)

func removeOrphanShardFiles(dir string, planned []PlannedShard) error {
	valid := map[string]bool{}
	for _, shard := range planned {
		valid[shard.ID] = true
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		match := shardFileName.FindStringSubmatch(entry.Name())
		if match != nil && !valid[match[2]] {
			if err := os.Remove(filepath.Join(dir, entry.Name())); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}
	return nil
}
