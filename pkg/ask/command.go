package ask

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type CommandOptions struct {
	DigestDir  string
	DigestPath string
	Endpoint   string
	JSONOutput bool
	MaxResults int
}

type CommandGroup struct {
	options CommandOptions
}

func NewCommandGroup(options CommandOptions) CommandGroup {
	if options.DigestDir != "" {
		options.DigestPath = options.DigestDir
	}
	if options.DigestPath == "" {
		options.DigestPath = ".hev-ask"
	}
	options.DigestDir = options.DigestPath
	if options.MaxResults <= 0 {
		options.MaxResults = 8
	}
	return CommandGroup{options: options}
}

func (group CommandGroup) Run(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	options, rest, err := parseCommandFlags(group.options, args)
	if err != nil {
		return err
	}
	if len(rest) == 0 {
		writeCommandUsage(stderr)
		return errors.New("missing command")
	}

	command := rest[0]
	commandArgs := rest[1:]
	switch command {
	case "tree":
		path, depth, err := parseTreeArgs(commandArgs)
		if err != nil {
			return err
		}
		return group.withOptions(options).runTree(ctx, path, depth, stdout)
	case "cat":
		if len(commandArgs) == 0 {
			return fmt.Errorf("cat requires a path")
		}
		return group.withOptions(options).runCat(ctx, strings.Join(commandArgs, " "), stdout)
	case "facts":
		if len(commandArgs) == 0 {
			return fmt.Errorf("facts requires a path")
		}
		return group.withOptions(options).runFacts(ctx, strings.Join(commandArgs, " "), stdout)
	case "grep":
		if len(commandArgs) == 0 {
			return fmt.Errorf("grep requires a query")
		}
		return group.withOptions(options).runGrep(ctx, strings.Join(commandArgs, " "), stdout)
	case "answer":
		if len(commandArgs) == 0 {
			return fmt.Errorf("answer requires a query")
		}
		return group.withOptions(options).runAnswer(ctx, strings.Join(commandArgs, " "), stdout)
	case "mcp":
		if len(commandArgs) != 0 {
			return fmt.Errorf("mcp takes no arguments")
		}
		return ServeMCP(ctx, MCPOptions{
			DigestPath: options.DigestPath,
			Endpoint:   options.Endpoint,
			MaxResults: options.MaxResults,
		}, stdin, stdout)
	case "help", "-h", "--help":
		writeCommandUsage(stdout)
		return nil
	default:
		writeCommandUsage(stderr)
		return fmt.Errorf("unknown command %q", command)
	}
}

func (group CommandGroup) withOptions(options CommandOptions) CommandGroup {
	return CommandGroup{options: options}
}

func (group CommandGroup) runTree(ctx context.Context, path string, depth int, stdout io.Writer) error {
	digest, err := group.loadDigest(ctx)
	if err != nil {
		return err
	}
	rendered, entries, err := RenderDigestTreeDepth(digest, path, depth)
	if err != nil {
		return err
	}
	payload := map[string]any{"path": cleanDigestPath(path), "depth": depth, "tree": rendered, "entries": entries}
	return group.writeOutput(stdout, payload, func(w io.Writer) {
		fmt.Fprintln(w, rendered)
	})
}

func (group CommandGroup) runCat(ctx context.Context, path string, stdout io.Writer) error {
	digest, err := group.loadDigest(ctx)
	if err != nil {
		return err
	}
	key := cleanDigestPath(path)
	if key == "" || key == "_meta" {
		overview := GetOverview(digest)
		return group.writeOutput(stdout, overview, func(w io.Writer) { writeOverviewHuman(w, overview) })
	}
	if strings.HasPrefix(key, digestGlossaryDir+"/") || strings.HasPrefix(key, "glossary/") {
		entry, _, ok := FindGlossaryByPath(digest, key)
		if !ok {
			return fmt.Errorf("no glossary entry matched %q", path)
		}
		return group.writeOutput(stdout, entry, func(w io.Writer) { writeGlossaryEntryHuman(w, entry) })
	}
	node, _, ok := FindNodeByPath(digest, path)
	if !ok {
		return fmt.Errorf("no digest path matched %q", path)
	}
	return group.writeOutput(stdout, node, func(w io.Writer) { writeSectionHuman(w, node) })
}

func (group CommandGroup) runFacts(ctx context.Context, path string, stdout io.Writer) error {
	digest, err := group.loadDigest(ctx)
	if err != nil {
		return err
	}
	facts, ok := FactsDigestPath(digest, path)
	if !ok {
		return fmt.Errorf("no digest path matched %q", path)
	}
	return group.writeOutput(stdout, facts, func(w io.Writer) {
		if len(facts.Facts) > 0 {
			fmt.Fprintln(w, "Facts:")
			for _, fact := range facts.Facts {
				fmt.Fprintf(w, "- %s\t%s\n", fact.Kind, fact.Literal)
			}
		}
		if len(facts.Sources) > 0 {
			if len(facts.Facts) > 0 {
				fmt.Fprintln(w)
			}
			fmt.Fprintln(w, "Sources:")
			for _, source := range facts.Sources {
				fmt.Fprintf(w, "- %s\n", source.URL)
			}
		}
	})
}

func (group CommandGroup) runGrep(ctx context.Context, query string, stdout io.Writer) error {
	response, err := group.search(ctx, query)
	if err != nil {
		return err
	}
	return group.writeOutput(stdout, response, func(w io.Writer) {
		for _, result := range response.Results {
			fmt.Fprintf(w, "%s\t%s\t%s\n", result.URL, result.Title, strings.Join(strings.Fields(result.Snippet), " "))
		}
	})
}

func (group CommandGroup) runAnswer(ctx context.Context, query string, stdout io.Writer) error {
	if group.options.Endpoint == "" {
		return errors.New("answer requires --endpoint for the remote SSE answer path; without --endpoint, use search for keyless local retrieval")
	}
	encoder := json.NewEncoder(stdout)
	return NewEndpointClient(group.options.Endpoint).StreamAnswer(ctx, query, func(event AnswerEvent) error {
		if group.options.JSONOutput {
			return encoder.Encode(event)
		}
		switch event.Event {
		case "token":
			var payload struct {
				Text string `json:"text"`
			}
			if err := json.Unmarshal(event.Data, &payload); err != nil {
				return err
			}
			_, err := io.WriteString(stdout, payload.Text)
			return err
		case "keyword":
			var response KeywordResponse
			if err := json.Unmarshal(event.Data, &response); err != nil {
				return err
			}
			return group.writeOutput(stdout, response, func(w io.Writer) {
				for _, result := range response.Results {
					fmt.Fprintf(w, "%s\n%s\n%s\n\n", result.Title, result.URL, result.Snippet)
				}
			})
		case "done":
			_, err := io.WriteString(stdout, "\n")
			return err
		}
		return nil
	})
}

func (group CommandGroup) search(ctx context.Context, query string) (KeywordResponse, error) {
	if group.options.Endpoint != "" {
		return NewEndpointClient(group.options.Endpoint).Search(ctx, query)
	}
	digest, err := group.loadDigest(ctx)
	if err != nil {
		return KeywordResponse{}, err
	}
	return SearchDigest(digest, query, SearchOptions{MaxResults: group.options.MaxResults}), nil
}

func (group CommandGroup) loadDigest(ctx context.Context) (Digest, error) {
	if group.options.Endpoint != "" {
		return NewEndpointClient(group.options.Endpoint).Digest(ctx)
	}
	return LoadDigest(group.options.DigestPath)
}

func (group CommandGroup) writeOutput(stdout io.Writer, value any, human func(io.Writer)) error {
	if !group.options.JSONOutput {
		human(stdout)
		return nil
	}
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func parseCommandFlags(defaults CommandOptions, args []string) (CommandOptions, []string, error) {
	options := NewCommandGroup(defaults).options
	var rest []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--digest-dir":
			if i+1 >= len(args) {
				return options, nil, fmt.Errorf("--digest-dir requires a value")
			}
			options.DigestPath = args[i+1]
			options.DigestDir = options.DigestPath
			i++
		case "--digest-path":
			if i+1 >= len(args) {
				return options, nil, fmt.Errorf("--digest-path requires a value")
			}
			options.DigestPath = args[i+1]
			options.DigestDir = options.DigestPath
			i++
		case "--endpoint":
			if i+1 >= len(args) {
				return options, nil, fmt.Errorf("--endpoint requires a value")
			}
			options.Endpoint = args[i+1]
			i++
		case "--json":
			options.JSONOutput = true
		case "--max-results":
			if i+1 >= len(args) {
				return options, nil, fmt.Errorf("--max-results requires a value")
			}
			var value int
			if _, err := fmt.Sscanf(args[i+1], "%d", &value); err != nil || value <= 0 {
				return options, nil, fmt.Errorf("--max-results must be a positive integer")
			}
			options.MaxResults = value
			i++
		default:
			rest = append(rest, args[i])
		}
	}
	return options, rest, nil
}

// parseTreeArgs reads the optional path argument and --depth flag for `tree`.
// Depth defaults to 2 (a couple of levels); "all" or "0" means unlimited.
func parseTreeArgs(args []string) (string, int, error) {
	depth := 2
	path := ""
	pathSet := false
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--depth":
			if i+1 >= len(args) {
				return "", 0, fmt.Errorf("--depth requires a value")
			}
			value := args[i+1]
			i++
			if value == "all" || value == "0" {
				depth = 0
				continue
			}
			var n int
			if _, err := fmt.Sscanf(value, "%d", &n); err != nil || n < 0 {
				return "", 0, fmt.Errorf(`--depth must be a non-negative integer or "all"`)
			}
			depth = n
		case strings.HasPrefix(args[i], "--"):
			return "", 0, fmt.Errorf("unknown tree flag %q", args[i])
		case pathSet:
			return "", 0, fmt.Errorf("tree takes at most one path")
		default:
			path = args[i]
			pathSet = true
		}
	}
	return path, depth, nil
}

func writeCommandUsage(w io.Writer) {
	fmt.Fprintln(w, `Usage:
  ask [--digest-dir .hev-ask] [--endpoint URL] [--json] <command>

Commands:
  tree [path] [--depth N|all]   map the digest (default: 2 levels deep)
  cat <path>                    read a section, _glossary/<term>, or _meta (overview)
  facts <path>                  grounded literals + sources for a section
  grep <query>                  keyword search over the digest
  answer <query>                synthesized, cited reply (requires --endpoint)
  mcp                           serve the digest to an agent over stdio`)
}

func writeGlossaryEntryHuman(w io.Writer, entry GlossaryEntry) {
	fmt.Fprintf(w, "%s\n", entry.Term)
	if len(entry.Aliases) > 0 {
		fmt.Fprintf(w, "Aliases: %s\n", strings.Join(entry.Aliases, ", "))
	}
	fmt.Fprintf(w, "%s\n", entry.Definition)
}

func writeOverviewHuman(w io.Writer, overview Overview) {
	if strings.TrimSpace(overview.Context) != "" {
		fmt.Fprintf(w, "%s\n\n", overview.Context)
	}
	fmt.Fprintln(w, overview.Overview)
}

func writeSectionHuman(w io.Writer, node DigestNode) {
	fmt.Fprintf(w, "%s\n", NodeLabel(node))
	if node.URL != "" {
		fmt.Fprintf(w, "%s\n", node.URL)
	}
	if strings.TrimSpace(node.Summary) != "" {
		fmt.Fprintf(w, "\n%s\n", node.Summary)
	}
}
