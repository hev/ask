package ask

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func LoadDigest(path string) (Digest, error) {
	stat, err := os.Stat(path)
	if err == nil && stat.IsDir() {
		digest, treeErr := LoadDigestTree(path)
		if treeErr == nil {
			return digest, nil
		}
		legacyPath := filepath.Join(path, "digest.json")
		if _, legacyErr := os.Stat(legacyPath); legacyErr == nil {
			return loadDigestJSON(legacyPath)
		}
		return Digest{}, treeErr
	}
	if err != nil && isDigestTreePath(path) {
		legacyPath := filepath.Join(path, "digest.json")
		if _, legacyErr := os.Stat(legacyPath); legacyErr == nil {
			return loadDigestJSON(legacyPath)
		}
	}
	return loadDigestJSON(path)
}

func loadDigestJSON(path string) (Digest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Digest{}, fmt.Errorf("read digest %q: %w", path, err)
	}
	var digest Digest
	if err := json.Unmarshal(data, &digest); err != nil {
		return Digest{}, fmt.Errorf("parse digest %q: %w", path, err)
	}
	normalizeDigest(&digest)
	return digest, nil
}

func normalizeDigest(digest *Digest) {
	if digest.Glossary == nil {
		digest.Glossary = []GlossaryEntry{}
	}
	if digest.Suggestions == nil {
		digest.Suggestions = []string{}
	}
	if digest.Nodes == nil {
		digest.Nodes = []DigestNode{}
	}
	if digest.Edges == nil {
		digest.Edges = []DigestEdge{}
	}
	for i := range digest.Glossary {
		if digest.Glossary[i].Aliases == nil {
			digest.Glossary[i].Aliases = []string{}
		}
	}
	for i := range digest.Nodes {
		node := &digest.Nodes[i]
		if node.Kind == "" {
			node.Kind = "section"
		}
		if node.Facts == nil {
			node.Facts = []Fact{}
		}
		if node.Sources == nil {
			node.Sources = []SourceRef{}
		}
		if node.Terms == nil {
			node.Terms = []string{}
		}
		if node.Mode == "" {
			node.Mode = "agent-primary"
		}
	}
}
