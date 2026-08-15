package ask

import (
	"reflect"
	"testing"
)

func TestExtractVersionRefs_MustRecognise(t *testing.T) {
	cases := []struct {
		name       string
		input      string
		wantSubj   string
		wantKind   string
		wantNorm   string
	}{
		{
			name:     "bare semver with v prefix",
			input:    "`v2.1.3`",
			wantSubj: "",
			wantKind: "unknown",
			wantNorm: "2.1.3",
		},
		{
			name:     "bare semver no prefix",
			input:    "`1.4.0`",
			wantSubj: "",
			wantKind: "unknown",
			wantNorm: "1.4.0",
		},
		{
			name:     "bare major-only with v",
			input:    "`v2`",
			wantSubj: "",
			wantKind: "unknown",
			wantNorm: "2",
		},
		{
			name:     "npm scoped package pin",
			input:    "`@hevmind/ask@1.4.0`",
			wantSubj: "@hevmind/ask",
			wantKind: "package",
			wantNorm: "1.4.0",
		},
		{
			name:     "npm unscoped package pin major-only",
			input:    "`astro@5`",
			wantSubj: "astro",
			wantKind: "package",
			wantNorm: "5",
		},
		{
			name:     "node runtime prose",
			input:    "Requires Node 20",
			wantSubj: "node",
			wantKind: "runtime",
			wantNorm: "20",
		},
		{
			name:     "node.js with operator",
			input:    "Node.js >= 18",
			wantSubj: "node",
			wantKind: "runtime",
			wantNorm: "18",
		},
		{
			name:     "python with plus suffix",
			input:    "Python 3.11+",
			wantSubj: "python",
			wantKind: "runtime",
			wantNorm: "3.11",
		},
		{
			name:     "go directive in fenced block",
			input:    "```go\ngo 1.23\n```",
			wantSubj: "go",
			wantKind: "runtime",
			wantNorm: "1.23",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			refs := ExtractVersionRefs("chunk1", tc.input)
			if len(refs) == 0 {
				t.Fatalf("expected at least one VersionRef, got none for input: %q", tc.input)
			}
			found := false
			for _, r := range refs {
				if r.Normalized == tc.wantNorm && r.Kind == tc.wantKind {
					if tc.wantSubj == "" || r.Subject == tc.wantSubj {
						found = true
						break
					}
				}
			}
			if !found {
				t.Fatalf("no ref matched subject=%q kind=%q norm=%q in %+v", tc.wantSubj, tc.wantKind, tc.wantNorm, refs)
			}
		})
	}
}

func TestExtractVersionRefs_MustNotMatch(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"ISO date", "2026-06-19"},
		{"date in prose", "June 19 2026"},
		{"URL path fragment", "/docs/v2/quickstart"},
		{"heading anchor", "#v2-migration"},
		{"multiplier decimal", "a 2.5x speedup"},
		{"percentage decimal", "99.9% uptime"},
		{"port number", "localhost:8080"},
		{"size value", "4096 MB"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			refs := ExtractVersionRefs("chunk1", tc.input)
			if len(refs) != 0 {
				t.Fatalf("expected no VersionRefs for %q, got %+v", tc.input, refs)
			}
		})
	}
}

func TestExtractVersionRefs_ChunkIDPropagated(t *testing.T) {
	refs := ExtractVersionRefs("my-chunk-id", "Node 20")
	if len(refs) == 0 {
		t.Fatal("expected refs")
	}
	for _, r := range refs {
		if r.ChunkID != "my-chunk-id" {
			t.Fatalf("ChunkID = %q, want %q", r.ChunkID, "my-chunk-id")
		}
	}
}

func TestExtractVersionRefs_Deduplication(t *testing.T) {
	refs := ExtractVersionRefs("c", "Use v1.4.0 and also v1.4.0 again.")
	count := 0
	for _, r := range refs {
		if r.Normalized == "1.4.0" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected 1 result for v1.4.0, got %d", count)
	}
}

func TestExtractVersionRefs_StableOrder(t *testing.T) {
	input := "Node 20 and Python 3.11+ and v2.1.3"
	a := ExtractVersionRefs("c", input)
	b := ExtractVersionRefs("c", input)
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("unstable order:\nfirst:  %+v\nsecond: %+v", a, b)
	}
}

func TestExtractVersionRefs_EmptyInput(t *testing.T) {
	refs := ExtractVersionRefs("c", "")
	if refs == nil {
		t.Fatal("expected non-nil empty slice, got nil")
	}
	if len(refs) != 0 {
		t.Fatalf("expected empty slice, got %+v", refs)
	}
}

func TestExtractVersionRefs_MultipleVersionsInOneChunk(t *testing.T) {
	input := "Requires Node 20 and @hevmind/ask@1.4.0.\n```go\ngo 1.23\n```"
	refs := ExtractVersionRefs("multi", input)

	findKind := func(kind, subj string) bool {
		for _, r := range refs {
			if r.Kind == kind && r.Subject == subj {
				return true
			}
		}
		return false
	}

	if !findKind("runtime", "node") {
		t.Fatalf("missing node runtime ref in %+v", refs)
	}
	if !findKind("package", "@hevmind/ask") {
		t.Fatalf("missing @hevmind/ask package ref in %+v", refs)
	}
	if !findKind("runtime", "go") {
		t.Fatalf("missing go directive ref in %+v", refs)
	}
}
