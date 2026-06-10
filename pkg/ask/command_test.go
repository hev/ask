package ask

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandGroupRunSearchJSON(t *testing.T) {
	path := writeCommandTestGraph(t)
	group := NewCommandGroup(CommandOptions{})

	var stdout, stderr bytes.Buffer
	err := group.Run(context.Background(), []string{"--digest-path", path, "--json", "--max-results", "1", "search", "kg", "path"}, strings.NewReader(""), &stdout, &stderr)
	if err != nil {
		t.Fatalf("run failed: %v\nstderr: %s", err, stderr.String())
	}

	var response KeywordResponse
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("decode output %s: %v", stdout.String(), err)
	}
	if len(response.Results) != 1 || response.Results[0].URL != "/docs/api/cli#flags" {
		t.Fatalf("unexpected search output: %#v", response)
	}
}

func TestCommandGroupRunGlossaryAlias(t *testing.T) {
	path := writeCommandTestGraph(t)
	group := NewCommandGroup(CommandOptions{DigestPath: path, JSONOutput: true})

	var stdout, stderr bytes.Buffer
	err := group.Run(context.Background(), []string{"glossary", "get", "kg"}, strings.NewReader(""), &stdout, &stderr)
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

	var stdout, stderr bytes.Buffer
	if err := group.Run(context.Background(), []string{"tree"}, strings.NewReader(""), &stdout, &stderr); err != nil {
		t.Fatalf("tree failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "api/") || !strings.Contains(stdout.String(), "CLI > Flags") {
		t.Fatalf("unexpected tree output: %s", stdout.String())
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
