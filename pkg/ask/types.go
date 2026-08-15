package ask

type GlossaryEntry struct {
	Term       string   `json:"term"`
	Aliases    []string `json:"aliases"`
	Definition string   `json:"definition"`
}

type VersionRef struct {
	Literal string `json:"literal"`
	ChunkID string `json:"chunkId"`
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
	VersionRefs []VersionRef `json:"versionRefs"`
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
	VersionRefs []VersionRef    `json:"versionRefs"`
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
