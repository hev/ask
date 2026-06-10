package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	askpkg "github.com/hev/ask/pkg/ask"
)

func writeTestDigest(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".hev-ask")
	apiGroup := "API"
	flagsHeading := "Flags"
	digest := askpkg.Digest{
		Version: 2,
		Context: "Docs orientation.",
		Glossary: []askpkg.GlossaryEntry{
			{Term: "Knowledge digest", Aliases: []string{"kg"}, Definition: "Committed docs digest."},
		},
		Overview: "## API\n- Flags - `api/cli#flags`",
		Nodes: []askpkg.DigestNode{
			{
				ID:      "api/cli#flags",
				Kind:    "section",
				Title:   "CLI",
				Heading: &flagsHeading,
				Group:   &apiGroup,
				URL:     "/docs/api/cli#flags",
				Summary: "Command flags configure digest paths.",
				Facts:   []askpkg.Fact{{Kind: "flag", Literal: "--digest-path", ChunkID: "api/cli#flags"}},
				Mode:    "source-primary",
				Terms:   []string{"flags", "digest"},
			},
		},
	}
	if err := askpkg.WriteDigest(path, digest); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunGlossaryGetJSON(t *testing.T) {
	path := writeTestDigest(t)
	var stdout, stderr bytes.Buffer
	err := run(context.Background(), []string{"--digest-path", path, "--json", "glossary", "get", "kg"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"term": "Knowledge digest"`) {
		t.Fatalf("unexpected output: %s", stdout.String())
	}
}

func TestRunSectionsListGroupJSON(t *testing.T) {
	path := writeTestDigest(t)
	var stdout, stderr bytes.Buffer
	err := run(context.Background(), []string{"--digest-path", path, "--json", "sections", "list", "--group", "api"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"id": "api/cli#flags"`) {
		t.Fatalf("unexpected output: %s", stdout.String())
	}
}

func TestRunSearchJSON(t *testing.T) {
	path := writeTestDigest(t)
	var stdout, stderr bytes.Buffer
	err := run(context.Background(), []string{"--digest-path", path, "--json", "search", "kg", "path"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"url": "/docs/api/cli#flags"`) {
		t.Fatalf("unexpected output: %s", stdout.String())
	}
}

func TestRunTreeCatFacts(t *testing.T) {
	path := writeTestDigest(t)
	var stdout, stderr bytes.Buffer
	if err := run(context.Background(), []string{"--digest-dir", path, "tree"}, &stdout, &stderr); err != nil {
		t.Fatalf("tree failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "api/cli/flags") {
		t.Fatalf("unexpected tree output: %s", stdout.String())
	}
	stdout.Reset()
	if err := run(context.Background(), []string{"--digest-dir", path, "cat", "api/cli/flags"}, &stdout, &stderr); err != nil {
		t.Fatalf("cat failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Command flags configure digest paths.") {
		t.Fatalf("unexpected cat output: %s", stdout.String())
	}
	stdout.Reset()
	if err := run(context.Background(), []string{"--digest-dir", path, "facts", "api/cli/flags"}, &stdout, &stderr); err != nil {
		t.Fatalf("facts failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "--digest-path") {
		t.Fatalf("unexpected facts output: %s", stdout.String())
	}
}

func TestRunDigestCorpusAndAssemble(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "src/content/docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "src/content/docs/index.mdx"), []byte("---\ntitle: \"Intro\"\ndescription: \"Start.\"\ngroup: \"Overview\"\n---\n\n## Install\n\nUse `ask`.\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(original); err != nil {
			t.Fatal(err)
		}
	}()

	var stdout, stderr bytes.Buffer
	if err := run(context.Background(), []string{"digest", "corpus"}, &stdout, &stderr); err != nil {
		t.Fatalf("corpus failed: %v\nstderr: %s", err, stderr.String())
	}
	input, err := os.ReadFile(filepath.Join(dir, ".hev-ask/digest-input.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(input), `"id": "index#install"`) {
		t.Fatalf("unexpected corpus: %s", input)
	}

	distill := `{"context":"ctx","glossary":[],"summaries":[{"id":"index#install","summary":"Install ask."}],"suggestions":["How do I install?"]}`
	if err := os.WriteFile(filepath.Join(dir, ".hev-ask/digest-distill.json"), []byte(distill), 0o600); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	if err := run(context.Background(), []string{"digest", "assemble"}, &stdout, &stderr); err != nil {
		t.Fatalf("assemble failed: %v\nstderr: %s", err, stderr.String())
	}
	digest, err := askpkg.LoadDigest(filepath.Join(dir, ".hev-ask"))
	if err != nil {
		t.Fatal(err)
	}
	if len(digest.Nodes) != 2 || digest.Nodes[1].Summary != "Install ask." {
		t.Fatalf("unexpected digest: %#v", digest)
	}
}

func TestRunDigestVerifySkipBuild(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "src/content/docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "src/content/docs/index.mdx"), []byte("---\ntitle: \"Intro\"\ndescription: \"Start.\"\ngroup: \"Overview\"\n---\n\n## Install\n\nUse `ask`.\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "dist/docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "dist/docs/index.html"), []byte(`<h2 id="install">Install</h2>`), 0o600); err != nil {
		t.Fatal(err)
	}
	digest := askpkg.AssembleDigest(askpkg.EmittedDistillation{
		Context:   "ctx",
		Glossary:  []askpkg.GlossaryEntry{},
		Summaries: []askpkg.SectionSummaryIn{{ID: "index#install", Summary: "Install ask."}},
	}, mustCorpus(t, dir))
	if err := askpkg.WriteDigest(filepath.Join(dir, ".hev-ask"), digest); err != nil {
		t.Fatal(err)
	}

	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(original); err != nil {
			t.Fatal(err)
		}
	}()

	var stdout, stderr bytes.Buffer
	if err := run(context.Background(), []string{"digest", "verify", "--skip-build"}, &stdout, &stderr); err != nil {
		t.Fatalf("verify failed: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "verified 1 anchors") {
		t.Fatalf("unexpected verify output: %s", stdout.String())
	}
}

func TestRunDigestBuildSkipsCurrentGraph(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "src/content/docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "src/content/docs/index.mdx"), []byte("---\ntitle: \"Intro\"\ndescription: \"Start.\"\ngroup: \"Overview\"\n---\n\n## Install\n\nUse `ask`.\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	corpus := mustCorpus(t, dir)
	digest := askpkg.AssembleDigest(askpkg.EmittedDistillation{Context: "ctx", Glossary: []askpkg.GlossaryEntry{}}, corpus)
	if err := askpkg.WriteDigest(filepath.Join(dir, ".hev-ask"), digest); err != nil {
		t.Fatal(err)
	}

	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(original); err != nil {
			t.Fatal(err)
		}
	}()

	var stdout, stderr bytes.Buffer
	if err := run(context.Background(), []string{"digest", "build"}, &stdout, &stderr); err != nil {
		t.Fatalf("build failed: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "digest:skipped") {
		t.Fatalf("unexpected build output: %s", stdout.String())
	}
}

func mustCorpus(t *testing.T, root string) askpkg.CorpusBuild {
	t.Helper()
	corpus, err := askpkg.BuildCorpus(askpkg.BuildOptions{SiteRoot: root, Collections: []string{"docs"}, BasePath: "/docs/", ChunkHeadingDepth: 3})
	if err != nil {
		t.Fatal(err)
	}
	return corpus
}
