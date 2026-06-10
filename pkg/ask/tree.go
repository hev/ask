package ask

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	digestMetaFile      = "_meta.md"
	digestGlossaryDir   = "_glossary"
	metaOverviewHeading = "\n## Overview\n\n"
)

type DigestEntry struct {
	Path  string `json:"path"`
	Kind  string `json:"kind"`
	Title string `json:"title"`
	URL   string `json:"url,omitempty"`
}

type DigestHead struct {
	Path    string  `json:"path"`
	Title   string  `json:"title"`
	Summary string  `json:"summary"`
	URL     string  `json:"url,omitempty"`
	Heading *string `json:"heading,omitempty"`
	Group   *string `json:"group,omitempty"`
}

type DigestFacts struct {
	Path    string      `json:"path"`
	Facts   []Fact      `json:"facts"`
	Sources []SourceRef `json:"sources"`
	Terms   []string    `json:"terms"`
}

type DigestTreeIntegrity struct {
	Orphans []string
}

type DigestTreeFile struct {
	Path string
	Body string
}

type MigrateResult struct {
	From   string
	Path   string
	Chunks int
}

func isDigestTreePath(path string) bool {
	return strings.ToLower(filepath.Ext(path)) != ".json"
}

func LoadDigestFS(fsys fs.FS) (Digest, error) {
	return loadDigestTreeFS(fsys, ".")
}

func LoadDigestTree(path string) (Digest, error) {
	return loadDigestTreeFS(os.DirFS(path), ".")
}

func loadDigestTreeFS(fsys fs.FS, root string) (Digest, error) {
	metaPath := filepath.ToSlash(filepath.Join(root, digestMetaFile))
	metaRaw, err := fs.ReadFile(fsys, metaPath)
	if err != nil {
		return Digest{}, fmt.Errorf("read digest tree meta %q: %w", digestMetaFile, err)
	}
	metaDoc := ParseFrontmatter(string(metaRaw))
	context, overview := parseMetaBody(metaDoc.Body)
	digest := Digest{
		Version:     intField(metaDoc.Data, "version", 2),
		GeneratedAt: stringField(metaDoc.Data, "generatedAt", ""),
		ContentHash: stringField(metaDoc.Data, "contentHash", ""),
		Context:     context,
		Overview:    overview,
		Suggestions: stringSliceField(metaDoc.Data, "suggestions"),
		Glossary:    []GlossaryEntry{},
		Nodes:       []DigestNode{},
		Edges:       []DigestEdge{},
	}

	if err := fs.WalkDir(fsys, root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == digestMetaFile || filepath.Ext(rel) != ".md" {
			return nil
		}
		raw, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		doc := ParseFrontmatter(string(raw))
		if strings.HasPrefix(rel, digestGlossaryDir+"/") {
			entry := GlossaryEntry{
				Term:       stringField(doc.Data, "term", titleFromPath(rel)),
				Aliases:    stringSliceField(doc.Data, "aliases"),
				Definition: strings.TrimSpace(doc.Body),
			}
			if entry.Term != "" {
				digest.Glossary = append(digest.Glossary, entry)
			}
			return nil
		}
		node := nodeFromTreeFile(rel, doc)
		if node.ID != "" && node.URL != "" {
			digest.Nodes = append(digest.Nodes, node)
		}
		return nil
	}); err != nil {
		return Digest{}, fmt.Errorf("read digest tree: %w", err)
	}

	sort.Slice(digest.Glossary, func(i, j int) bool {
		return strings.ToLower(digest.Glossary[i].Term) < strings.ToLower(digest.Glossary[j].Term)
	})
	sort.Slice(digest.Nodes, func(i, j int) bool { return digest.Nodes[i].ID < digest.Nodes[j].ID })
	normalizeDigest(&digest)
	return digest, nil
}

func nodeFromTreeFile(rel string, doc FrontmatterDocument) DigestNode {
	pathKey := strings.TrimSuffix(filepath.ToSlash(rel), ".md")
	heading := optionalStringField(doc.Data, "heading")
	group := optionalStringField(doc.Data, "group")
	anchor := optionalStringField(doc.Data, "anchor")
	id := stringField(doc.Data, "id", "")
	if id == "" {
		id = pathKeyToSectionID(pathKey)
	}
	url := stringField(doc.Data, "url", "")
	body := strings.TrimSpace(doc.Body)
	summary := firstParagraph(body)
	if summary == "" {
		summary = stringField(doc.Data, "summary", "")
	}
	sources := sourceRefSliceField(doc.Data, "sources")
	if len(sources) == 0 && url != "" {
		sources = []SourceRef{{ChunkID: id, URL: url, Anchor: anchor}}
	}
	return DigestNode{
		ID:      id,
		Kind:    "section",
		Title:   stringField(doc.Data, "title", id),
		Heading: heading,
		Group:   group,
		URL:     url,
		Summary: summary,
		Hash:    stringField(doc.Data, "hash", ""),
		Facts:   factSliceField(doc.Data, "facts"),
		Sources: sources,
		Mode:    stringField(doc.Data, "mode", "agent-primary"),
		Terms:   stringSliceField(doc.Data, "terms"),
	}
}

func WriteDigestTree(path string, digest Digest) error {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return err
	}
	desired := map[string]bool{}
	for _, file := range DigestTreeFiles(digest) {
		if err := writeDigestTreeFile(filepath.Join(path, filepath.FromSlash(file.Path)), file.Body); err != nil {
			return err
		}
		desired[file.Path] = true
	}
	return removeOrphanDigestMarkdown(path, desired)
}

func DigestTreeFiles(digest Digest) []DigestTreeFile {
	files := []DigestTreeFile{{Path: digestMetaFile, Body: renderDigestMetaFile(digest)}}
	for _, entry := range digest.Glossary {
		files = append(files, DigestTreeFile{
			Path: filepath.ToSlash(filepath.Join(digestGlossaryDir, glossaryPath(entry)+".md")),
			Body: renderGlossaryFile(entry),
		})
	}
	for i, node := range digest.Nodes {
		files = append(files, DigestTreeFile{
			Path: NodePath(node) + ".md",
			Body: renderSectionFile(node, i),
		})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files
}

func WriteDigestArchive(w io.Writer, digest Digest) error {
	gzipWriter := gzip.NewWriter(w)
	tarWriter := tar.NewWriter(gzipWriter)
	for _, file := range DigestTreeFiles(digest) {
		body := []byte(file.Body)
		if err := tarWriter.WriteHeader(&tar.Header{
			Name:    file.Path,
			Mode:    0o644,
			Size:    int64(len(body)),
			ModTime: time.Unix(0, 0).UTC(),
		}); err != nil {
			return err
		}
		if _, err := tarWriter.Write(body); err != nil {
			return err
		}
	}
	if err := tarWriter.Close(); err != nil {
		return err
	}
	return gzipWriter.Close()
}

func ExtractDigestArchive(r io.Reader, path string) error {
	parent := filepath.Dir(path)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	tmp, err := os.MkdirTemp(parent, ".hev-ask-archive-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("read digest archive gzip: %w", err)
	}
	defer gzipReader.Close()
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read digest archive tar: %w", err)
		}
		if header.Typeflag != tar.TypeReg && header.Typeflag != tar.TypeRegA && header.Typeflag != 0 {
			continue
		}
		if header.FileInfo().IsDir() {
			continue
		}
		name := filepath.ToSlash(header.Name)
		if !safeArchivePath(name) {
			return fmt.Errorf("unsafe digest archive path %q", header.Name)
		}
		outPath := filepath.Join(tmp, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}
		out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(out, tarReader)
		closeErr := out.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
	}
	if _, err := LoadDigestTree(tmp); err != nil {
		return fmt.Errorf("validate digest archive: %w", err)
	}
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func safeArchivePath(path string) bool {
	cleaned := filepath.ToSlash(filepath.Clean(path))
	local := filepath.FromSlash(path)
	return cleaned == path &&
		path != "." &&
		!strings.HasPrefix(path, "/") &&
		!strings.HasPrefix(path, "../") &&
		!strings.Contains(path, "/../") &&
		!filepath.IsAbs(local) &&
		filepath.VolumeName(local) == ""
}

func MigrateLegacyDigest(siteRoot string, inputPath string, digestPath string) (MigrateResult, error) {
	if digestPath == "" {
		digestPath = ".hev-ask"
	}
	input := resolveSitePath(siteRoot, inputPath)
	output := resolveSitePath(siteRoot, digestPath)
	if !isDigestTreePath(output) {
		output = filepath.Dir(output)
	}
	digest, err := loadDigestJSON(input)
	if err != nil {
		return MigrateResult{}, err
	}
	if err := WriteDigestTree(output, digest); err != nil {
		return MigrateResult{}, err
	}
	return MigrateResult{From: input, Path: output, Chunks: len(digest.Nodes)}, nil
}

func renderDigestMetaFile(digest Digest) string {
	fields := []frontmatterField{
		{"version", digest.Version},
		{"generatedAt", firstNonEmpty(digest.GeneratedAt, time.Now().UTC().Format("2006-01-02T15:04:05.000Z"))},
		{"contentHash", digest.ContentHash},
		{"suggestions", digest.Suggestions},
	}
	body := strings.TrimSpace("## Context\n\n" + strings.TrimSpace(digest.Context) + metaOverviewHeading + strings.TrimSpace(digest.Overview))
	return markdownWithFrontmatter(fields, body)
}

func renderGlossaryFile(entry GlossaryEntry) string {
	fields := []frontmatterField{
		{"term", entry.Term},
		{"aliases", entry.Aliases},
	}
	return markdownWithFrontmatter(fields, strings.TrimSpace(entry.Definition))
}

func renderSectionFile(node DigestNode, order int) string {
	anchor := ""
	if len(node.Sources) > 0 && node.Sources[0].Anchor != nil {
		anchor = *node.Sources[0].Anchor
	}
	fields := []frontmatterField{
		{"id", node.ID},
		{"title", node.Title},
		{"heading", stringPtrValue(node.Heading)},
		{"group", stringPtrValue(node.Group)},
		{"order", order},
		{"url", node.URL},
		{"anchor", anchor},
		{"terms", node.Terms},
		{"hash", node.Hash},
		{"mode", node.Mode},
		{"facts", node.Facts},
		{"sources", node.Sources},
	}
	return markdownWithFrontmatter(fields, strings.TrimSpace(node.Summary))
}

type frontmatterField struct {
	key   string
	value any
}

func markdownWithFrontmatter(fields []frontmatterField, body string) string {
	var b strings.Builder
	b.WriteString("---\n")
	for _, field := range fields {
		b.WriteString(field.key)
		b.WriteString(": ")
		b.WriteString(formatFrontmatterValue(field.value))
		b.WriteByte('\n')
	}
	b.WriteString("---\n\n")
	b.WriteString(strings.TrimSpace(body))
	b.WriteByte('\n')
	return b.String()
}

func writeDigestTreeFile(path string, body string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(body), 0o644)
}

func formatFrontmatterValue(value any) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case string:
		if v == "" {
			return "null"
		}
		data, _ := json.Marshal(v)
		return string(data)
	case int:
		return strconv.Itoa(v)
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return "null"
		}
		return string(data)
	}
}

func removeOrphanDigestMarkdown(root string, desired map[string]bool) error {
	return filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == "shards" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.ToLower(filepath.Ext(path)) != ".md" {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !desired[rel] {
			return os.Remove(path)
		}
		return nil
	})
}

func CheckDigestTreeIntegrity(path string, digest Digest) (DigestTreeIntegrity, error) {
	if !isDigestTreePath(path) {
		return DigestTreeIntegrity{}, nil
	}
	stat, err := os.Stat(path)
	if err != nil {
		return DigestTreeIntegrity{}, err
	}
	if !stat.IsDir() {
		return DigestTreeIntegrity{}, fmt.Errorf("digest path %q is not a directory", path)
	}
	desired := map[string]bool{digestMetaFile: true}
	for _, entry := range digest.Glossary {
		desired[filepath.ToSlash(filepath.Join(digestGlossaryDir, glossaryPath(entry)+".md"))] = true
	}
	for _, node := range digest.Nodes {
		desired[NodePath(node)+".md"] = true
	}
	var orphans []string
	if err := filepath.WalkDir(path, func(file string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == "shards" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.ToLower(filepath.Ext(file)) != ".md" {
			return nil
		}
		rel, err := filepath.Rel(path, file)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !desired[rel] {
			orphans = append(orphans, rel)
		}
		return nil
	}); err != nil {
		return DigestTreeIntegrity{}, err
	}
	sort.Strings(orphans)
	return DigestTreeIntegrity{Orphans: orphans}, nil
}

func NodePath(node DigestNode) string {
	if strings.TrimSpace(node.ID) == "" {
		return ""
	}
	return strings.ReplaceAll(strings.Trim(strings.TrimSpace(node.ID), "/"), "#", "/")
}

func FindNodeByPath(digest Digest, path string) (DigestNode, string, bool) {
	key := cleanDigestPath(path)
	if key == "" {
		return DigestNode{}, "", false
	}
	keyAsPath := strings.ReplaceAll(key, "#", "/")
	for _, node := range digest.Nodes {
		nodePath := NodePath(node)
		if nodePath == key || nodePath == keyAsPath || node.ID == key {
			return node, nodePath, true
		}
	}
	return DigestNode{}, "", false
}

func ListDigestPath(digest Digest, path string) []DigestEntry {
	prefix := cleanDigestPath(path)
	entries := digestLeafEntries(digest)
	children := map[string]DigestEntry{}
	for _, entry := range entries {
		if prefix != "" {
			if entry.Path == prefix {
				children[entry.Path] = entry
				continue
			}
			if !strings.HasPrefix(entry.Path, prefix+"/") {
				continue
			}
		}
		rest := strings.TrimPrefix(entry.Path, prefix)
		rest = strings.TrimPrefix(rest, "/")
		if rest == "" {
			continue
		}
		segment := rest
		if slash := strings.IndexByte(rest, '/'); slash >= 0 {
			segment = rest[:slash]
		}
		childPath := segment
		if prefix != "" {
			childPath = prefix + "/" + segment
		}
		if slash := strings.IndexByte(rest, '/'); slash >= 0 {
			children[childPath+"/"] = DigestEntry{Path: childPath, Kind: "dir", Title: segment + "/"}
			continue
		}
		children[childPath] = entry
	}
	out := make([]DigestEntry, 0, len(children))
	for _, entry := range children {
		out = append(out, entry)
	}
	sortDigestEntries(out)
	return out
}

func RenderDigestTree(digest Digest) string {
	entries := digestLeafEntries(digest)
	sortDigestEntries(entries)
	var lines []string
	seenDirs := map[string]bool{}
	for _, entry := range entries {
		parts := strings.Split(entry.Path, "/")
		for i := 1; i < len(parts); i++ {
			dir := strings.Join(parts[:i], "/")
			if seenDirs[dir] {
				continue
			}
			seenDirs[dir] = true
			lines = append(lines, strings.Repeat("  ", i-1)+parts[i-1]+"/")
		}
		indent := strings.Repeat("  ", len(parts)-1)
		lines = append(lines, indent+parts[len(parts)-1]+"  "+entry.Title)
	}
	return strings.Join(lines, "\n")
}

func HeadDigestPath(digest Digest, path string) (DigestHead, bool) {
	key := cleanDigestPath(path)
	if key == "" || key == "_meta" {
		return DigestHead{Path: "_meta", Title: "Digest metadata", Summary: firstParagraph(digest.Context)}, true
	}
	if strings.HasPrefix(key, digestGlossaryDir+"/") {
		if entry, glossaryPath, ok := FindGlossaryByPath(digest, key); ok {
			return DigestHead{Path: glossaryPath, Title: entry.Term, Summary: entry.Definition}, true
		}
		return DigestHead{}, false
	}
	node, nodePath, ok := FindNodeByPath(digest, key)
	if !ok {
		return DigestHead{}, false
	}
	return DigestHead{Path: nodePath, Title: NodeLabel(node), Summary: node.Summary, URL: node.URL, Heading: node.Heading, Group: node.Group}, true
}

func FactsDigestPath(digest Digest, path string) (DigestFacts, bool) {
	node, nodePath, ok := FindNodeByPath(digest, path)
	if !ok {
		return DigestFacts{}, false
	}
	return DigestFacts{Path: nodePath, Facts: node.Facts, Sources: node.Sources, Terms: node.Terms}, true
}

func FindGlossaryByPath(digest Digest, path string) (GlossaryEntry, string, bool) {
	key := strings.TrimPrefix(cleanDigestPath(path), digestGlossaryDir+"/")
	key = strings.TrimPrefix(key, "glossary/")
	for _, entry := range digest.Glossary {
		entryPath := digestGlossaryDir + "/" + glossaryPath(entry)
		if key == glossaryPath(entry) || cleanDigestPath(path) == entryPath || normalizeLookup(key) == normalizeLookup(entry.Term) {
			return entry, entryPath, true
		}
		for _, alias := range entry.Aliases {
			if normalizeLookup(key) == normalizeLookup(alias) {
				return entry, entryPath, true
			}
		}
	}
	return GlossaryEntry{}, "", false
}

func NodeLabel(node DigestNode) string {
	if node.Heading != nil && *node.Heading != "" {
		return node.Title + " > " + *node.Heading
	}
	return node.Title
}

func digestLeafEntries(digest Digest) []DigestEntry {
	entries := []DigestEntry{{Path: "_meta", Kind: "meta", Title: "Digest metadata"}}
	if len(digest.Glossary) > 0 {
		for _, entry := range digest.Glossary {
			entries = append(entries, DigestEntry{Path: digestGlossaryDir + "/" + glossaryPath(entry), Kind: "glossary", Title: entry.Term})
		}
	} else {
		entries = append(entries, DigestEntry{Path: digestGlossaryDir, Kind: "dir", Title: digestGlossaryDir + "/"})
	}
	for _, node := range digest.Nodes {
		entries = append(entries, DigestEntry{Path: NodePath(node), Kind: "section", Title: NodeLabel(node), URL: node.URL})
	}
	return entries
}

func sortDigestEntries(entries []DigestEntry) {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Kind == "dir" && entries[j].Kind != "dir" {
			return true
		}
		if entries[i].Kind != "dir" && entries[j].Kind == "dir" {
			return false
		}
		return entries[i].Path < entries[j].Path
	})
}

func cleanDigestPath(path string) string {
	cleaned := strings.TrimSpace(strings.ReplaceAll(path, "\\", "/"))
	cleaned = strings.TrimPrefix(cleaned, "./")
	cleaned = strings.Trim(cleaned, "/")
	cleaned = strings.TrimSuffix(cleaned, ".md")
	if cleaned == "." {
		return ""
	}
	return cleaned
}

func pathKeyToSectionID(pathKey string) string {
	return strings.Replace(pathKey, "/", "#", strings.Count(pathKey, "/"))
}

func glossaryPath(entry GlossaryEntry) string {
	slug := githubSlug(entry.Term)
	if slug == "" {
		slug = "term"
	}
	return slug
}

func parseMetaBody(body string) (string, string) {
	trimmed := strings.TrimSpace(body)
	trimmed = strings.TrimPrefix(trimmed, "## Context")
	trimmed = strings.TrimSpace(trimmed)
	if at := strings.Index(trimmed, metaOverviewHeading); at >= 0 {
		return strings.TrimSpace(trimmed[:at]), strings.TrimSpace(trimmed[at+len(metaOverviewHeading):])
	}
	if at := strings.Index(trimmed, "\n## Overview\n"); at >= 0 {
		return strings.TrimSpace(trimmed[:at]), strings.TrimSpace(strings.TrimPrefix(trimmed[at:], "\n## Overview"))
	}
	return trimmed, ""
}

func firstParagraph(body string) string {
	for _, block := range strings.Split(strings.TrimSpace(body), "\n\n") {
		if trimmed := strings.TrimSpace(block); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func titleFromPath(path string) string {
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	base = strings.ReplaceAll(base, "-", " ")
	if base == "" {
		return path
	}
	return strings.Title(base)
}

func optionalStringField(data map[string]any, key string) *string {
	value := stringField(data, key, "")
	if value == "" {
		return nil
	}
	return &value
}

func stringPtrValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func intField(data map[string]any, key string, fallback int) int {
	value, ok := data[key]
	if !ok {
		return fallback
	}
	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}
	return fallback
}

func stringSliceField(data map[string]any, key string) []string {
	raw, ok := data[key]
	if !ok || raw == nil {
		return []string{}
	}
	values, ok := raw.([]any)
	if !ok {
		if one, ok := raw.(string); ok && one != "" {
			return []string{one}
		}
		return []string{}
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if text, ok := value.(string); ok && text != "" {
			out = append(out, text)
		}
	}
	return out
}

func factSliceField(data map[string]any, key string) []Fact {
	raw, ok := data[key].([]any)
	if !ok {
		return []Fact{}
	}
	out := make([]Fact, 0, len(raw))
	for _, item := range raw {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		literal := stringField(obj, "literal", "")
		if literal == "" {
			continue
		}
		out = append(out, Fact{
			Kind:    stringField(obj, "kind", "value"),
			Literal: literal,
			ChunkID: stringField(obj, "chunkId", ""),
		})
	}
	return out
}

func sourceRefSliceField(data map[string]any, key string) []SourceRef {
	raw, ok := data[key].([]any)
	if !ok {
		return []SourceRef{}
	}
	out := make([]SourceRef, 0, len(raw))
	for _, item := range raw {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		chunkID := stringField(obj, "chunkId", "")
		url := stringField(obj, "url", "")
		if chunkID == "" || url == "" {
			continue
		}
		out = append(out, SourceRef{
			ChunkID: chunkID,
			URL:     url,
			Anchor:  optionalStringField(obj, "anchor"),
		})
	}
	return out
}
