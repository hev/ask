package ask

import (
	"regexp"
	"strings"
)

// VersionRef is an extraction artifact; it becomes a persistence type in HEV-30.
// "product" Kind is reserved but has no recognition rule yet — a follow-on story
// will define product-name heuristics. Until then Kind will never equal "product".
// Note: ExtractFacts (via versionRE) and ExtractVersionRefs both fire on bare
// semver tokens in the same chunk. HEV-30 decides whether to deduplicate them at
// the digest layer.
type VersionRef struct {
	Subject    string // what the version applies to: "node", "@hevmind/ask", "python"
	Raw        string // the literal text matched, e.g. "Node 20"
	Normalized string // canonical form, e.g. "20"
	Kind       string // "runtime" | "package" | "product" | "unknown"
	ChunkID    string // chunk the reference came from
}

var (
	// npm package pins: @scope/pkg@ver or pkg@ver
	npmPinRE = regexp.MustCompile(`(@[a-zA-Z0-9_-]+/[a-zA-Z0-9_.-]+|[a-zA-Z][a-zA-Z0-9_.-]*)@(v?\d+(?:\.\d+)*)`)

	// runtime keyword patterns in prose
	nodeRE   = regexp.MustCompile(`(?i)\bNode(?:\.js)?\s*(?:>=|<=|>|<|~=|==|!=)?\s*(v?\d+(?:\.\d+)*)`)
	pythonRE = regexp.MustCompile(`(?i)\bPython\s*(?:>=|<=|>|<|~=|==|!=)?\s*(v?\d+(?:\.\d+)+)`)
	rubyRE   = regexp.MustCompile(`(?i)\bRuby\s*(?:>=|<=|>|<|~=|==|!=)?\s*(v?\d+(?:\.\d+)+)`)

	// go directive inside a fenced code block: "go 1.23"
	goDirectiveRE = regexp.MustCompile(`(?m)^go\s+(\d+(?:\.\d+)+)\s*$`)

	// bare semver: v2.1.3, 1.4.0, v2
	bareSemverRE = regexp.MustCompile(`\bv?(\d+(?:\.\d+)+|\d+)\b`)

	// exclusion detectors
	isoDateRE    = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	urlFragmentRE = regexp.MustCompile(`[/#]v?\d`)
)

// ExtractVersionRefs finds version references in a chunk of markdown and returns
// them as structured values. Pure: same input always yields the same output.
// Results are in document order (fenced blocks first, then prose), deduplicated by Raw.
func ExtractVersionRefs(chunkID string, raw string) []VersionRef {
	seen := map[string]bool{}
	refs := []VersionRef{}
	push := func(r VersionRef) {
		r.Raw = strings.TrimSpace(r.Raw)
		if r.Raw == "" || seen[r.Raw] {
			return
		}
		seen[r.Raw] = true
		r.ChunkID = chunkID
		refs = append(refs, r)
	}

	// Strip ISO dates from consideration by replacing with spaces to avoid false positives
	stripped := isoDateRE.ReplaceAllStringFunc(raw, func(s string) string {
		return strings.Repeat(" ", len(s))
	})

	// --- Fenced code blocks first ---
	for _, match := range fenceRE.FindAllStringSubmatch(stripped, -1) {
		block := match[1]
		for _, m := range goDirectiveRE.FindAllStringSubmatch(block, -1) {
			push(VersionRef{Subject: "go", Raw: "go " + m[1], Normalized: m[1], Kind: "runtime"})
		}
	}

	// Strip fenced blocks before prose scan
	prose := fenceRE.ReplaceAllString(stripped, " ")

	// --- npm package pins ---
	for _, m := range npmPinRE.FindAllStringSubmatch(prose, -1) {
		fullMatch := m[0]
		pkg := m[1]
		ver := strings.TrimPrefix(m[2], "v")
		if shouldExclude(prose, fullMatch) {
			continue
		}
		push(VersionRef{Subject: pkg, Raw: fullMatch, Normalized: ver, Kind: "package"})
	}

	// --- Runtime keyword patterns ---
	for _, m := range nodeRE.FindAllStringSubmatch(prose, -1) {
		fullMatch := m[0]
		ver := strings.TrimPrefix(m[1], "v")
		push(VersionRef{Subject: "node", Raw: fullMatch, Normalized: ver, Kind: "runtime"})
	}
	for _, m := range pythonRE.FindAllStringSubmatch(prose, -1) {
		fullMatch := m[0]
		ver := strings.TrimPrefix(m[1], "v")
		// strip trailing + from "3.11+"
		ver = strings.TrimSuffix(ver, "+")
		push(VersionRef{Subject: "python", Raw: fullMatch, Normalized: ver, Kind: "runtime"})
	}
	for _, m := range rubyRE.FindAllStringSubmatch(prose, -1) {
		fullMatch := m[0]
		ver := strings.TrimPrefix(m[1], "v")
		push(VersionRef{Subject: "ruby", Raw: fullMatch, Normalized: ver, Kind: "runtime"})
	}

	// --- Bare semver (last, so runtime/package matches take priority) ---
	// Remove already-matched text to avoid double-counting
	bareText := npmPinRE.ReplaceAllString(prose, " ")
	bareText = nodeRE.ReplaceAllString(bareText, " ")
	bareText = pythonRE.ReplaceAllString(bareText, " ")
	bareText = rubyRE.ReplaceAllString(bareText, " ")

	for _, loc := range bareSemverRE.FindAllStringIndex(bareText, -1) {
		start, end := loc[0], loc[1]
		fullMatch := bareText[start:end]

		// skip bare integers with no dot (e.g. plain "20" in "Node 20" would have
		// been consumed by nodeRE already; remaining bare integers like "4096" are sizes)
		if !strings.Contains(fullMatch, ".") && !strings.HasPrefix(fullMatch, "v") {
			continue
		}

		if shouldExcludeAt(bareText, start, end) {
			continue
		}

		norm := strings.TrimPrefix(fullMatch, "v")
		push(VersionRef{Subject: "", Raw: fullMatch, Normalized: norm, Kind: "unknown"})
	}

	return refs
}

// shouldExclude checks whether the matched string at any position in text should be excluded.
func shouldExclude(text, match string) bool {
	idx := strings.Index(text, match)
	if idx < 0 {
		return false
	}
	return shouldExcludeAt(text, idx, idx+len(match))
}

// shouldExcludeAt applies post-filter exclusion rules given byte offsets in text.
func shouldExcludeAt(text string, start, end int) bool {
	// Check character immediately before the match
	if start > 0 {
		prev := text[start-1]
		if prev == '/' || prev == '#' || prev == ':' {
			return true
		}
	}
	// Check character immediately after the match
	if end < len(text) {
		next := text[end]
		if next == 'x' || next == 'X' || next == '%' {
			return true
		}
		if next == '/' {
			return true
		}
	}
	// Check for size suffixes: " MB", " KB", " GB", " TB"
	if end+3 <= len(text) {
		suffix := strings.ToUpper(text[end : end+3])
		if suffix == " MB" || suffix == " KB" || suffix == " GB" || suffix == " TB" {
			return true
		}
	}
	// Check URL/anchor context: preceded by slash or hash within 20 chars
	contextStart := start - 20
	if contextStart < 0 {
		contextStart = 0
	}
	context := text[contextStart:end]
	if urlFragmentRE.MatchString(context) {
		return true
	}
	return false
}
