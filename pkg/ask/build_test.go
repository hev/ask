package ask

import (
	"os"
	"path/filepath"
	"testing"
)

func writeParityFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	mustWrite := func(path string, body string) {
		t.Helper()
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	mustWrite("src/content/docs/index.mdx", "---\n"+
		"title: \"Intro\"\n"+
		"description: \"Start here.\"\n"+
		"group: \"Overview\"\n"+
		"---\n"+
		"import X from \"./x\";\n\n"+
		"# Page Title\n\n"+
		"Intro with [a link](/x), `code`, and <Badge />.\n\n"+
		"## Install `ask`!\n\n"+
		"Use `ask search` and `--digest-path`.\n\n"+
		"### Install `ask`! ###\n\n"+
		"Duplicate heading with v1.2.3 and `claude-haiku-4-5`.\n\n"+
		"#### Too Deep\n\n"+
		"Stays in previous section.\n\n"+
		"## Tables & Lists\n\n"+
		"| Flag | Meaning |\n"+
		"|---|---|\n"+
		"| `--json` | machine output |\n\n"+
		"- item with _emphasis_\n")
	mustWrite("src/content/docs/api/config.mdx", "---\n"+
		"title: \"Config\"\n"+
		"description: \"API details.\"\n"+
		"group: \"API\"\n"+
		"---\n\n"+
		"## Options\n\n"+
		"Set `endpoint` to `/api/ask`.\n")
	return root
}

func TestBuildCorpusMatchesTypeScriptFixture(t *testing.T) {
	root := writeParityFixture(t)
	corpus, err := BuildCorpus(BuildOptions{SiteRoot: root, Collections: []string{"docs"}, BasePath: "/docs/", ChunkHeadingDepth: 3})
	if err != nil {
		t.Fatal(err)
	}
	if corpus.ContentHash != "2409c203fad949741c1b2121f36b8a2c7990f8d32579734903636874095257fd" {
		t.Fatalf("unexpected content hash: %s", corpus.ContentHash)
	}
	gotIDs := make([]string, len(corpus.Chunks))
	for i, chunk := range corpus.Chunks {
		gotIDs[i] = chunk.ID
	}
	wantIDs := []string{
		"api/config",
		"api/config#options",
		"index",
		"index#install-ask",
		"index#install-ask-1",
		"index#tables--lists",
	}
	for i, want := range wantIDs {
		if gotIDs[i] != want {
			t.Fatalf("chunk ids mismatch at %d: got %#v want %#v", i, gotIDs, wantIDs)
		}
	}
	wantText := "Install ask! ### Duplicate heading with v1.2.3 and claude-haiku-4-5. Too Deep Stays in previous section."
	if corpus.Chunks[4].Text != wantText {
		t.Fatalf("unexpected cleaned text:\n%s", corpus.Chunks[4].Text)
	}
}

func TestAssembleDigestDerivesNodesFactsAndOverview(t *testing.T) {
	root := writeParityFixture(t)
	corpus, err := BuildCorpus(BuildOptions{SiteRoot: root, Collections: []string{"docs"}, BasePath: "/docs/", ChunkHeadingDepth: 3})
	if err != nil {
		t.Fatal(err)
	}
	digest := AssembleDigest(EmittedDistillation{
		Context: "Fixture docs.",
		Glossary: []GlossaryEntry{
			{Term: "ask CLI", Aliases: []string{"ask"}, Definition: "Reads docs."},
		},
		Summaries: []SectionSummaryIn{{ID: "api/config#options", Summary: "Options configure the endpoint."}},
		Suggestions: []string{
			"How do I configure the endpoint?",
		},
	}, corpus)
	if digest.ContentHash != corpus.ContentHash || len(digest.Nodes) != len(corpus.Chunks) {
		t.Fatalf("unexpected digest shape: %#v", digest)
	}
	node, ok := GetSection(digest, "api/config#options")
	if !ok {
		t.Fatal("missing api/config#options node")
	}
	if node.Mode != "source-primary" {
		t.Fatalf("API node should be source-primary, got %q", node.Mode)
	}
	if len(node.Facts) != 2 || node.Facts[0].Literal != "endpoint" || node.Facts[1].Literal != "/api/ask" {
		t.Fatalf("unexpected facts: %#v", node.Facts)
	}
	if digest.Overview == "" || digest.Suggestions[0] == "" {
		t.Fatalf("expected overview and suggestions: %#v", digest)
	}
}
