package ask

import (
	"testing"
)

// TestNormalizeVersion validates the normalization contract used by all comparisons.
func TestNormalizeVersion(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"v20", "20"},
		{"V1.2.3", "1.2.3"},
		{"20", "20"},
		{"20.0.0-rc1", "20.0.0-rc1"},
		{"v20.0.0-rc1", "20.0.0-rc1"},
		{"", ""},
		{"  v3  ", "3"},
	}
	for _, tc := range cases {
		if got := NormalizeVersion(tc.input); got != tc.want {
			t.Errorf("NormalizeVersion(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// TestDiffVersionRefs_NoDrift ensures identical digests produce no drift.
func TestDiffVersionRefs_NoDrift(t *testing.T) {
	d := digestWithVersionRefs([]testNodeVersions{
		{id: "quickstart", refs: []VersionRef{{Subject: "node", Version: "20", ChunkID: "quickstart"}}},
	})
	drifts := DiffVersionRefs(d, d)
	if len(drifts) != 0 {
		t.Fatalf("expected no drift on identical digest, got %v", drifts)
	}
}

// TestDiffVersionRefs_DriftPresent checks that a version change is detected with correct fields.
func TestDiffVersionRefs_DriftPresent(t *testing.T) {
	old := digestWithVersionRefs([]testNodeVersions{
		{id: "quickstart", refs: []VersionRef{{Subject: "node", Version: "20", ChunkID: "quickstart"}}},
		{id: "api/cli", refs: []VersionRef{{Subject: "node", Version: "20", ChunkID: "api/cli"}}},
	})
	updated := digestWithVersionRefs([]testNodeVersions{
		{id: "quickstart", refs: []VersionRef{{Subject: "node", Version: "24", ChunkID: "quickstart"}}},
		{id: "api/cli", refs: []VersionRef{{Subject: "node", Version: "24", ChunkID: "api/cli"}}},
	})
	drifts := DiffVersionRefs(old, updated)
	if len(drifts) != 1 {
		t.Fatalf("expected 1 drift, got %d: %v", len(drifts), drifts)
	}
	d := drifts[0]
	if d.Subject != "node" {
		t.Errorf("Subject = %q, want %q", d.Subject, "node")
	}
	if d.Old != "20" {
		t.Errorf("Old = %q, want %q", d.Old, "20")
	}
	if d.New != "24" {
		t.Errorf("New = %q, want %q", d.New, "24")
	}
	if len(d.Sections) != 2 || d.Sections[0] != "api/cli" || d.Sections[1] != "quickstart" {
		t.Errorf("Sections = %v, want [api/cli quickstart]", d.Sections)
	}
	want := "node: 20 -> 24 (api/cli, quickstart)"
	if got := FormatVersionDrift(d); got != want {
		t.Errorf("FormatVersionDrift = %q, want %q", got, want)
	}
}

// TestDiffVersionRefs_NoVersionDataInExisting checks AC3: old digest with no version refs
// produces no drift even when the new digest has version refs.
func TestDiffVersionRefs_NoVersionDataInExisting(t *testing.T) {
	old := digestWithVersionRefs([]testNodeVersions{
		{id: "quickstart", refs: nil},
	})
	updated := digestWithVersionRefs([]testNodeVersions{
		{id: "quickstart", refs: []VersionRef{{Subject: "node", Version: "20", ChunkID: "quickstart"}}},
	})
	drifts := DiffVersionRefs(old, updated)
	if len(drifts) != 0 {
		t.Fatalf("expected no drift when existing has no version data, got %v", drifts)
	}
}

// TestDiffVersionRefs_SubjectRenamedOrRemoved checks that additions and removals are not drift.
func TestDiffVersionRefs_SubjectRenamedOrRemoved(t *testing.T) {
	old := digestWithVersionRefs([]testNodeVersions{
		{id: "a", refs: []VersionRef{{Subject: "node", Version: "20", ChunkID: "a"}}},
	})
	updated := digestWithVersionRefs([]testNodeVersions{
		// node is gone, python is new
		{id: "a", refs: []VersionRef{{Subject: "python", Version: "3.12", ChunkID: "a"}}},
	})
	drifts := DiffVersionRefs(old, updated)
	if len(drifts) != 0 {
		t.Fatalf("expected no drift for additions/removals, got %v", drifts)
	}
}

// TestFindContradictions_ConflictFound checks that two sections claiming different
// versions for the same subject are reported as a contradiction.
func TestFindContradictions_ConflictFound(t *testing.T) {
	d := digestWithVersionRefs([]testNodeVersions{
		{id: "quickstart", refs: []VersionRef{{Subject: "node", Version: "20", ChunkID: "quickstart"}}},
		{id: "api/cli", refs: []VersionRef{{Subject: "node", Version: "24", ChunkID: "api/cli"}}},
	})
	contradictions := FindContradictions(d)
	if len(contradictions) != 1 {
		t.Fatalf("expected 1 contradiction, got %d: %v", len(contradictions), contradictions)
	}
	c := contradictions[0]
	if c.Subject != "node" {
		t.Errorf("Subject = %q, want %q", c.Subject, "node")
	}
	if len(c.Versions) != 2 {
		t.Fatalf("expected 2 version conflicts, got %d", len(c.Versions))
	}
	// sorted by version: "20" < "24"
	if c.Versions[0].Version != "20" || len(c.Versions[0].Sections) != 1 || c.Versions[0].Sections[0] != "quickstart" {
		t.Errorf("first conflict = %+v, want {Version:20 Sections:[quickstart]}", c.Versions[0])
	}
	if c.Versions[1].Version != "24" || len(c.Versions[1].Sections) != 1 || c.Versions[1].Sections[0] != "api/cli" {
		t.Errorf("second conflict = %+v, want {Version:24 Sections:[api/cli]}", c.Versions[1])
	}
}

// TestFindContradictions_NoConflict ensures a consistent corpus is clean.
func TestFindContradictions_NoConflict(t *testing.T) {
	d := digestWithVersionRefs([]testNodeVersions{
		{id: "quickstart", refs: []VersionRef{{Subject: "node", Version: "20", ChunkID: "quickstart"}}},
		{id: "api/cli", refs: []VersionRef{{Subject: "node", Version: "20", ChunkID: "api/cli"}}},
	})
	if got := FindContradictions(d); len(got) != 0 {
		t.Fatalf("expected no contradictions, got %v", got)
	}
}

// TestFindContradictions_NormalizationUnifiesVariants ensures v20 and 20 are the same.
func TestFindContradictions_NormalizationUnifiesVariants(t *testing.T) {
	d := digestWithVersionRefs([]testNodeVersions{
		{id: "a", refs: []VersionRef{{Subject: "node", Version: "v20", ChunkID: "a"}}},
		{id: "b", refs: []VersionRef{{Subject: "node", Version: "20", ChunkID: "b"}}},
	})
	if got := FindContradictions(d); len(got) != 0 {
		t.Fatalf("v20 and 20 should not be a contradiction, got %v", got)
	}
}

// TestVerifyResult_ContradictionsPopulated verifies the VerifyResult plumbing
// without calling os.Exit — tests the struct field, not the CLI exit code.
func TestVerifyResult_ContradictionsPopulated(t *testing.T) {
	d := digestWithVersionRefs([]testNodeVersions{
		{id: "quickstart", refs: []VersionRef{{Subject: "node", Version: "20", ChunkID: "quickstart"}}},
		{id: "api/cli", refs: []VersionRef{{Subject: "node", Version: "24", ChunkID: "api/cli"}}},
	})
	contradictions := FindContradictions(d)
	if len(contradictions) == 0 {
		t.Fatal("expected FindContradictions to return contradictions for conflicting corpus")
	}
	result := VerifyResult{Contradictions: contradictions}
	if len(result.Contradictions) != 1 {
		t.Fatalf("VerifyResult.Contradictions should hold the conflicts, got %d", len(result.Contradictions))
	}
}

// testNodeVersions is a helper type for building test digests.
type testNodeVersions struct {
	id   string
	refs []VersionRef
}

func digestWithVersionRefs(nodes []testNodeVersions) Digest {
	digestNodes := make([]DigestNode, 0, len(nodes))
	for _, n := range nodes {
		digestNodes = append(digestNodes, DigestNode{
			ID:          n.id,
			Kind:        "section",
			URL:         "/docs/" + n.id,
			VersionRefs: n.refs,
		})
	}
	return Digest{
		Version: 2,
		Nodes:   digestNodes,
	}
}
