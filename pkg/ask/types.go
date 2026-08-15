package ask

type GlossaryEntry struct {
	Term       string   `json:"term"`
	Aliases    []string `json:"aliases"`
	Definition string   `json:"definition"`
}

type Fact struct {
	Kind    string `json:"kind"`
	Literal string `json:"literal"`
	ChunkID string `json:"chunkId"`
}

type SourceRef struct {
	ChunkID string  `json:"chunkId"`
	URL     string  `json:"url"`
	Anchor  *string `json:"anchor"`
}

type VersionRef struct {
	Subject string `json:"subject"`
	Version string `json:"version"`
	ChunkID string `json:"chunkId"`
}

type VersionDrift struct {
	Subject  string
	Old      string
	New      string
	Sections []string
}

type VersionConflict struct {
	Version  string
	Sections []string
}

type VersionContradiction struct {
	Subject  string
	Versions []VersionConflict
}

type DigestNode struct {
	ID          string       `json:"id"`
	Kind        string       `json:"kind"`
	Title       string       `json:"title"`
	Heading     *string      `json:"heading"`
	Group       *string      `json:"group"`
	URL         string       `json:"url"`
	Summary     string       `json:"summary"`
	Hash        string       `json:"hash,omitempty"`
	Facts       []Fact       `json:"facts"`
	Sources     []SourceRef  `json:"sources"`
	Mode        string       `json:"mode"`
	Terms       []string     `json:"terms"`
	VersionRefs []VersionRef `json:"versionRefs,omitempty"`
}

type DigestEdge struct {
	Rel  string `json:"rel"`
	From string `json:"from"`
	To   string `json:"to"`
}

type Digest struct {
	Version     int             `json:"version"`
	GeneratedAt string          `json:"generatedAt"`
	ContentHash string          `json:"contentHash"`
	Context     string          `json:"context"`
	Glossary    []GlossaryEntry `json:"glossary"`
	Overview    string          `json:"overview"`
	Suggestions []string        `json:"suggestions"`
	Nodes       []DigestNode    `json:"nodes"`
	Edges       []DigestEdge    `json:"edges"`
}

type SectionSummary struct {
	ID      string  `json:"id"`
	Title   string  `json:"title"`
	Heading *string `json:"heading"`
	Group   *string `json:"group"`
	URL     string  `json:"url"`
}

type KeywordResult struct {
	Title   string  `json:"title"`
	Heading *string `json:"heading,omitempty"`
	URL     string  `json:"url"`
	Group   *string `json:"group,omitempty"`
	Snippet string  `json:"snippet"`
}

type KeywordResponse struct {
	Results []KeywordResult `json:"results"`
	Query   string          `json:"query"`
	Model   string          `json:"model,omitempty"`
	Mode    string          `json:"mode"`
	Warning string          `json:"warning,omitempty"`
}
