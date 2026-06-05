package ask

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildDigestSkipsCurrentGraphWithoutAPIKey(t *testing.T) {
	root := writeParityFixture(t)
	options := BuildOptions{SiteRoot: root, Collections: []string{"docs"}, BasePath: "/docs/", ChunkHeadingDepth: 3}
	corpus, err := BuildCorpus(options)
	if err != nil {
		t.Fatal(err)
	}
	digest := AssembleDigest(EmittedDistillation{Context: "ctx", Glossary: []GlossaryEntry{}, Summaries: []SectionSummaryIn{}}, corpus)
	if err := WriteDigest(filepath.Join(root, ".hev-ask/digest.json"), digest); err != nil {
		t.Fatal(err)
	}
	result, err := BuildDigest(BuildDigestOptions{BuildOptions: options})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "skipped" || result.Chunks != len(corpus.Chunks) {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestBuildDigestCallsAnthropicAndWritesGraph(t *testing.T) {
	root := writeParityFixture(t)
	options := BuildOptions{SiteRoot: root, Collections: []string{"docs"}, BasePath: "/docs/", ChunkHeadingDepth: 3}
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Header.Get("x-api-key") != "test-key" {
			t.Fatalf("missing api key header")
		}
		if request.Header.Get("anthropic-version") != anthropicVersion {
			t.Fatalf("unexpected anthropic version")
		}
		var body map[string]any
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["model"] != "test-model" {
			t.Fatalf("unexpected model: %#v", body["model"])
		}
		encoded, _ := json.Marshal(body)
		if !strings.Contains(string(encoded), "api/config#options") {
			t.Fatalf("request did not include corpus: %s", encoded)
		}
		payload := map[string]any{
			"content": []map[string]any{
				{
					"type": "tool_use",
					"name": "emit_digest",
					"input": map[string]any{
						"context": "Fixture docs.",
						"glossary": []map[string]any{
							{"term": "ask", "aliases": []string{"cli"}, "definition": "A command."},
						},
						"summaries": []map[string]any{
							{"id": "api/config#options", "summary": "Configures the endpoint."},
							{"id": "index#install-ask", "summary": "Installs ask."},
						},
						"suggestions": []string{"How do I configure it?"},
					},
				},
			},
		}
		data, _ := json.Marshal(payload)
		return response(http.StatusOK, "application/json", data), nil
	})}
	result, err := BuildDigest(BuildDigestOptions{
		BuildOptions: options,
		DigestModel:  "test-model",
		APIKey:       "test-key",
		HTTPClient:   client,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "built" {
		t.Fatalf("unexpected result: %#v", result)
	}
	digest, err := LoadDigest(filepath.Join(root, ".hev-ask/digest.json"))
	if err != nil {
		t.Fatal(err)
	}
	node, ok := GetSection(digest, "api/config#options")
	if !ok || node.Summary != "Configures the endpoint." {
		t.Fatalf("unexpected digest node: %#v %v", node, ok)
	}
}

func TestBuildDigestRequiresKeyForFreshGraph(t *testing.T) {
	root := writeParityFixture(t)
	t.Setenv("ANTHROPIC_API_KEY", "")
	_, err := BuildDigest(BuildDigestOptions{
		BuildOptions: BuildOptions{SiteRoot: root, Collections: []string{"docs"}, BasePath: "/docs/", ChunkHeadingDepth: 3},
	})
	if err == nil || !strings.Contains(err.Error(), "ANTHROPIC_API_KEY") {
		t.Fatalf("expected key error, got %v", err)
	}
}
