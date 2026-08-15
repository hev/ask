package ask

import (
	"fmt"
	"sort"
	"strings"
)

// NormalizeVersion strips a leading v/V prefix and lowercases.
// "v20" → "20", "V1.2.3" → "1.2.3", "20.0.0-rc1" → "20.0.0-rc1".
func NormalizeVersion(v string) string {
	lower := strings.ToLower(strings.TrimSpace(v))
	return strings.TrimPrefix(lower, "v")
}

type versionEntry struct {
	version string // already normalized
	chunkID string
}

// collectVersionEntries groups all VersionRefs from a digest by normalized subject.
func collectVersionEntries(digest Digest) map[string][]versionEntry {
	out := map[string][]versionEntry{}
	for _, node := range digest.Nodes {
		for _, ref := range node.VersionRefs {
			subject := strings.ToLower(strings.TrimSpace(ref.Subject))
			if subject == "" {
				continue
			}
			out[subject] = append(out[subject], versionEntry{
				version: NormalizeVersion(ref.Version),
				chunkID: ref.ChunkID,
			})
		}
	}
	return out
}

// consensusVersion returns the single normalized version for a subject's entries,
// or "" if the entries are empty or contradictory (multiple distinct values).
func consensusVersion(entries []versionEntry) string {
	if len(entries) == 0 {
		return ""
	}
	seen := map[string]bool{}
	for _, e := range entries {
		seen[e.version] = true
	}
	if len(seen) != 1 {
		return ""
	}
	for v := range seen {
		return v
	}
	return ""
}

// DiffVersionRefs compares version refs between an existing committed digest and a
// freshly assembled one. It returns one VersionDrift per subject whose normalized
// version changed. Additions (subject absent from existing) and removals are not
// reported. If existing has no version data at all (a pre-HEV-30 digest), it
// returns an empty slice rather than treating everything as new.
func DiffVersionRefs(existing, updated Digest) []VersionDrift {
	hasExisting := false
	for _, node := range existing.Nodes {
		if len(node.VersionRefs) > 0 {
			hasExisting = true
			break
		}
	}
	if !hasExisting {
		return []VersionDrift{}
	}

	oldBySubject := collectVersionEntries(existing)
	newBySubject := collectVersionEntries(updated)

	var drifts []VersionDrift
	for subject, newEntries := range newBySubject {
		oldEntries, inOld := oldBySubject[subject]
		if !inOld {
			continue
		}
		oldVer := consensusVersion(oldEntries)
		newVer := consensusVersion(newEntries)
		if oldVer == "" || newVer == "" || oldVer == newVer {
			continue
		}
		sections := uniqueChunkIDs(newEntries)
		sort.Strings(sections)
		drifts = append(drifts, VersionDrift{
			Subject:  subject,
			Old:      oldVer,
			New:      newVer,
			Sections: sections,
		})
	}
	sort.Slice(drifts, func(i, j int) bool { return drifts[i].Subject < drifts[j].Subject })
	return drifts
}

// FindContradictions finds subjects in a single digest where two or more sections
// claim different normalized versions. A corpus that disagrees with itself is a
// content error that verify should catch and fail on.
func FindContradictions(digest Digest) []VersionContradiction {
	bySubject := collectVersionEntries(digest)

	var out []VersionContradiction
	for subject, entries := range bySubject {
		// group by normalized version
		sectionsByVersion := map[string][]string{}
		for _, e := range entries {
			if !stringSliceContains(sectionsByVersion[e.version], e.chunkID) {
				sectionsByVersion[e.version] = append(sectionsByVersion[e.version], e.chunkID)
			}
		}
		if len(sectionsByVersion) <= 1 {
			continue
		}
		var conflicts []VersionConflict
		for ver, sections := range sectionsByVersion {
			sorted := make([]string, len(sections))
			copy(sorted, sections)
			sort.Strings(sorted)
			conflicts = append(conflicts, VersionConflict{Version: ver, Sections: sorted})
		}
		sort.Slice(conflicts, func(i, j int) bool { return conflicts[i].Version < conflicts[j].Version })
		out = append(out, VersionContradiction{Subject: subject, Versions: conflicts})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Subject < out[j].Subject })
	return out
}

// FormatVersionDrift formats a drift entry as "subject: old -> new (section1, section2)".
func FormatVersionDrift(d VersionDrift) string {
	return fmt.Sprintf("%s: %s -> %s (%s)", d.Subject, d.Old, d.New, strings.Join(d.Sections, ", "))
}

// FormatVersionContradiction formats a contradiction as
// `"v1" (s1, s2) vs "v2" (s3)`.
func FormatVersionContradiction(c VersionContradiction) string {
	parts := make([]string, 0, len(c.Versions))
	for _, v := range c.Versions {
		parts = append(parts, fmt.Sprintf("%q (%s)", v.Version, strings.Join(v.Sections, ", ")))
	}
	return strings.Join(parts, " vs ")
}

func uniqueChunkIDs(entries []versionEntry) []string {
	seen := map[string]bool{}
	var out []string
	for _, e := range entries {
		if !seen[e.chunkID] {
			seen[e.chunkID] = true
			out = append(out, e.chunkID)
		}
	}
	return out
}

func stringSliceContains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
