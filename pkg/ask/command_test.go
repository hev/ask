package ask

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandGroupRunGrepJSON(t *testing.T) {
	path := writeCommandTestGraph(t)
	group := NewCommandGroup(CommandOptions{})

	var stdout, stderr bytes.Buffer
	err := group.Run(context.Background(), []string{"--digest-path", path, "--json", "--max-results", "1", "grep", "kg", "path"}, strings.NewReader(""), &stdout, &stderr)
	if err != nil {
		t.Fatalf("run failed: %v\nstderr: %s", err, stderr.String())
	}

	var response KeywordResponse
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("decode output %s: %v", stdout.String(), err)
	}
	if len(response.Results) != 1 || response.Results[0].URL != "/docs/api/cli#flags" {
		t.Fatalf("unexpected grep output: %#v", response)
	}
}

// Glossary entries are reached through cat on a _glossary/ path; the lookup is
// alias-aware, so "kg" resolves to its canonical term.
func TestCommandGroupCatGlossaryAlias(t *testing.T) {
	path := writeCommandTestGraph(t)
	group := NewCommandGroup(CommandOptions{DigestPath: path, JSONOutput: true})

	var stdout, stderr bytes.Buffer
	err := group.Run(context.Background(), []string{"cat", "_glossary/kg"}, strings.NewReader(""), &stdout, &stderr)
	if err != nil {
		t.Fatalf("run failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"term": "Knowledge digest"`) {
		t.Fatalf("unexpected glossary output: %s", stdout.String())
	}
}

func TestCommandGroupRunTreeCatAndFacts(t *testing.T) {
	path := writeCommandTestGraph(t)
	group := NewCommandGroup(CommandOptions{DigestDir: path})

	// The default depth maps the tree a couple of levels deep, collapsing the
	// deeper section into a count rather than printing its label.
	var stdout, stderr bytes.Buffer
	if err := group.Run(context.Background(), []string{"tree"}, strings.NewReader(""), &stdout, &stderr); err != nil {
		t.Fatalf("tree failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "api/") || !strings.Contains(stdout.String(), "cli/  (+1)") {
		t.Fatalf("unexpected tree output: %s", stdout.String())
	}
	if strings.Contains(stdout.String(), "CLI > Flags") {
		t.Fatalf("default tree should not descend to the leaf label: %s", stdout.String())
	}

	// --depth all expands every level.
	stdout.Reset()
	if err := group.Run(context.Background(), []string{"tree", "--depth", "all"}, strings.NewReader(""), &stdout, &stderr); err != nil {
		t.Fatalf("tree --depth all failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "CLI > Flags") {
		t.Fatalf("unexpected deep tree output: %s", stdout.String())
	}

	stdout.Reset()
	if err := group.Run(context.Background(), []string{"cat", "api/cli/flags"}, strings.NewReader(""), &stdout, &stderr); err != nil {
		t.Fatalf("cat failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Command flags configure digest paths and output.") {
		t.Fatalf("unexpected cat output: %s", stdout.String())
	}

	stdout.Reset()
	if err := group.Run(context.Background(), []string{"facts", "api/cli/flags"}, strings.NewReader(""), &stdout, &stderr); err != nil {
		t.Fatalf("facts failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "--digest-path") {
		t.Fatalf("unexpected facts output: %s", stdout.String())
	}
}

func writeCommandTestGraph(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".hev-ask")
	if err := WriteDigest(path, testDigest()); err != nil {
		t.Fatal(err)
	}
	return path
}
