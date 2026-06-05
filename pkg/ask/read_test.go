package ask

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func testDigest() Digest {
	apiGroup := "API"
	overviewGroup := "Overview"
	flagsHeading := "Flags"
	introHeading := "Introduction"
	return Digest{
		Version: 2,
		Context: "Docs orientation.",
		Glossary: []GlossaryEntry{
			{Term: "Knowledge digest", Aliases: []string{"kg", "shadow site"}, Definition: "Committed docs digest."},
		},
		Overview: "## API\n- Flags - `api/cli#flags`",
		Nodes: []DigestNode{
			{
				ID:      "api/cli#flags",
				Kind:    "section",
				Title:   "CLI",
				Heading: &flagsHeading,
				Group:   &apiGroup,
				URL:     "/docs/api/cli#flags",
				Summary: "Command flags configure digest paths and output.",
				Facts:   []Fact{{Kind: "flag", Literal: "--digest-path", ChunkID: "api/cli#flags"}},
				Mode:    "source-primary",
				Terms:   []string{"flags", "digest", "paths"},
			},
			{
				ID:      "index#intro",
				Kind:    "section",
				Title:   "Intro",
				Heading: &introHeading,
				Group:   &overviewGroup,
				URL:     "/docs#intro",
				Summary: "The overlay and CLI read the same digest.",
				Mode:    "agent-primary",
				Terms:   []string{"overlay", "cli"},
			},
		},
	}
}

func TestLoadDigestNormalizesSlices(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "digest.json")
	if err := os.WriteFile(path, []byte(`{"version":2,"nodes":[{"id":"x","url":"/x"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	digest, err := LoadDigest(path)
	if err != nil {
		t.Fatal(err)
	}
	if digest.Glossary == nil || digest.Nodes[0].Facts == nil || digest.Nodes[0].Terms == nil {
		t.Fatalf("expected nil slices to normalize to empty slices: %#v", digest)
	}
	if digest.Nodes[0].Kind != "section" {
		t.Fatalf("expected default node kind, got %q", digest.Nodes[0].Kind)
	}
}

func TestReadHelpers(t *testing.T) {
	digest := testDigest()
	if entry, ok := GetGlossaryEntry(digest, "KG"); !ok || entry.Term != "Knowledge digest" {
		t.Fatalf("expected alias lookup to resolve term, got %#v %v", entry, ok)
	}
	sections := ListSectionSummaries(digest, "api")
	if len(sections) != 1 || sections[0].ID != "api/cli#flags" {
		t.Fatalf("unexpected filtered sections: %#v", sections)
	}
	node, ok := GetSection(digest, "api%2Fcli%23flags")
	if !ok || node.URL != "/docs/api/cli#flags" {
		t.Fatalf("expected encoded section lookup to resolve, got %#v %v", node, ok)
	}
	overview := GetOverview(digest)
	if overview.Context != "Docs orientation." || overview.Overview == "" {
		t.Fatalf("unexpected overview: %#v", overview)
	}
}

func TestSearchDigestUsesGlossaryAndFacts(t *testing.T) {
	response := SearchDigest(testDigest(), "kg path", SearchOptions{MaxResults: 4})
	if len(response.Results) == 0 {
		t.Fatal("expected search results")
	}
	data, _ := json.Marshal(response.Results[0])
	if response.Results[0].URL != "/docs/api/cli#flags" {
		t.Fatalf("expected flags section first, got %s", data)
	}
}
