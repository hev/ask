package ask

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type VerifyOptions struct {
	BuildOptions
	BuildCommand string
	DistDir      string
	SkipBuild    bool
}

type VerifyResult struct {
	Checked    int
	Missing    []MissingAnchor
	Dropped    []DroppedLiteral
	Uncovered  []string
	TreeErrors []string
}

type MissingAnchor struct {
	URL      string
	File     string
	AnchorID string
}

type DroppedLiteral struct {
	ID      string
	Literal string
}

func VerifyAnchors(options VerifyOptions) (VerifyResult, error) {
	normalizeBuildOptions(&options.BuildOptions)
	if options.DistDir == "" {
		options.DistDir = "dist"
	}
	if !options.SkipBuild {
		command := options.BuildCommand
		if command == "" {
			command = "pnpm build"
		}
		cmd := exec.Command("sh", "-c", command)
		cmd.Dir = options.SiteRoot
		output, err := cmd.CombinedOutput()
		if err != nil {
			return VerifyResult{}, fmt.Errorf("build command failed: %w\n%s", err, strings.TrimSpace(string(output)))
		}
	}

	corpus, err := BuildCorpus(options.BuildOptions)
	if err != nil {
		return VerifyResult{}, err
	}
	result := VerifyResult{}
	digestPath := resolveSitePath(options.SiteRoot, options.DigestPath)
	if stat, statErr := os.Stat(digestPath); statErr == nil && stat.IsDir() && isDigestTreePath(digestPath) {
		digest, err := LoadDigest(digestPath)
		if err != nil {
			return VerifyResult{}, fmt.Errorf("digest tree integrity failed: %w", err)
		}
		integrity, err := CheckDigestTreeIntegrity(digestPath, digest)
		if err != nil {
			return VerifyResult{}, fmt.Errorf("digest tree integrity failed: %w", err)
		}
		for _, orphan := range integrity.Orphans {
			result.TreeErrors = append(result.TreeErrors, "orphan digest file: "+orphan)
		}
		if digest.ContentHash != "" && digest.ContentHash != corpus.ContentHash {
			result.TreeErrors = append(result.TreeErrors, "meta contentHash does not match the current corpus; run `ask digest build`")
		}
	}
	distDir := resolveSitePath(options.SiteRoot, options.DistDir)
	for _, chunk := range corpus.Chunks {
		if chunk.AnchorID == "" {
			continue
		}
		result.Checked++
		files := htmlFilesForURL(distDir, chunk.URL)
		if !findHTMLWithID(files, chunk.AnchorID) {
			result.Missing = append(result.Missing, MissingAnchor{URL: chunk.URL, File: files[0], AnchorID: chunk.AnchorID})
		}
	}

	dropped, uncovered := verifyFidelity(options.BuildOptions, corpus.Chunks)
	result.Dropped = dropped
	result.Uncovered = uncovered
	return result, nil
}

func verifyFidelity(options BuildOptions, chunks []Chunk) ([]DroppedLiteral, []string) {
	digest, err := LoadDigest(resolveSitePath(options.SiteRoot, options.DigestPath))
	if err != nil || len(digest.Nodes) == 0 {
		uncovered := make([]string, 0, len(chunks))
		for _, chunk := range chunks {
			uncovered = append(uncovered, chunk.ID)
		}
		return []DroppedLiteral{}, uncovered
	}
	byID := map[string]DigestNode{}
	for _, node := range digest.Nodes {
		byID[node.ID] = node
	}
	dropped := []DroppedLiteral{}
	uncovered := []string{}
	for _, chunk := range chunks {
		node, ok := byID[chunk.ID]
		if !ok {
			uncovered = append(uncovered, chunk.ID)
			continue
		}
		if node.Mode == "source-primary" {
			continue
		}
		carried := map[string]bool{}
		for _, fact := range node.Facts {
			carried[fact.Literal] = true
		}
		for _, fact := range ExtractFacts(chunk.ID, chunk.Raw) {
			if !carried[fact.Literal] {
				dropped = append(dropped, DroppedLiteral{ID: chunk.ID, Literal: fact.Literal})
			}
		}
	}
	return dropped, uncovered
}

func htmlFilesForURL(distDir string, url string) []string {
	pathname := strings.TrimSuffix(strings.TrimPrefix(strings.Split(url, "#")[0], "/"), "/")
	return []string{
		filepath.Join(distDir, pathname, "index.html"),
		filepath.Join(distDir, "client", pathname, "index.html"),
	}
}

func findHTMLWithID(files []string, id string) bool {
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		if hasID(string(data), id) {
			return true
		}
	}
	return false
}

func hasID(html string, id string) bool {
	return strings.Contains(html, ` id="`+id+`"`) || strings.Contains(html, ` id='`+id+`'`)
}
