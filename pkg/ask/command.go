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
		if len(commandArgs) != 0 {
			return fmt.Errorf("tree takes no arguments")
		}
		return group.withOptions(options).runTree(ctx, stdout)
	case "ls":
		path := ""
		if len(commandArgs) > 1 {
			return fmt.Errorf("ls takes at most one path")
		}
		if len(commandArgs) == 1 {
			path = commandArgs[0]
		}
		return group.withOptions(options).runLs(ctx, path, stdout)
	case "head":
		if len(commandArgs) == 0 {
			return fmt.Errorf("head requires a path")
		}
		return group.withOptions(options).runHead(ctx, strings.Join(commandArgs, " "), stdout)
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
	case "glossary":
		return group.withOptions(options).runGlossary(ctx, commandArgs, stdout)
	case "sections":
		return group.withOptions(options).runSections(ctx, commandArgs, stdout)
	case "section":
		return group.withOptions(options).runSection(ctx, commandArgs, stdout)
	case "overview":
		if len(commandArgs) != 0 {
			return fmt.Errorf("overview takes no arguments")
		}
		return group.withOptions(options).runOverview(ctx, stdout)
	case "search":
		if len(commandArgs) == 0 {
			return fmt.Errorf("search requires a query")
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

func (group CommandGroup) runTree(ctx context.Context, stdout io.Writer) error {
	digest, err := group.loadDigest(ctx)
	if err != nil {
		return err
	}
	payload := map[string]any{"tree": RenderDigestTree(digest), "entries": digestLeafEntries(digest)}
	return group.writeOutput(stdout, payload, func(w io.Writer) {
		fmt.Fprintln(w, payload["tree"])
	})
}

func (group CommandGroup) runLs(ctx context.Context, path string, stdout io.Writer) error {
	digest, err := group.loadDigest(ctx)
	if err != nil {
		return err
	}
	entries := ListDigestPath(digest, path)
	return group.writeOutput(stdout, map[string]any{"entries": entries}, func(w io.Writer) {
		for _, entry := range entries {
			name := entry.Path
			if prefix := cleanDigestPath(path); prefix != "" {
				name = strings.TrimPrefix(strings.TrimPrefix(entry.Path, prefix), "/")
			}
			if entry.Kind == "dir" && !strings.HasSuffix(name, "/") {
				name += "/"
			}
			if entry.Title != "" && entry.Kind != "dir" {
				fmt.Fprintf(w, "%s\t%s\n", name, entry.Title)
			} else {
				fmt.Fprintln(w, name)
			}
		}
	})
}

func (group CommandGroup) runHead(ctx context.Context, path string, stdout io.Writer) error {
	digest, err := group.loadDigest(ctx)
	if err != nil {
		return err
	}
	head, ok := HeadDigestPath(digest, path)
	if !ok {
		return fmt.Errorf("no digest path matched %q", path)
	}
	return group.writeOutput(stdout, head, func(w io.Writer) {
		fmt.Fprintln(w, head.Title)
		if head.URL != "" {
			fmt.Fprintln(w, head.URL)
		}
		if strings.TrimSpace(head.Summary) != "" {
			fmt.Fprintf(w, "\n%s\n", head.Summary)
		}
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

func (group CommandGroup) runGlossary(ctx context.Context, args []string, stdout io.Writer) error {
	subcommand := "list"
	if len(args) > 0 {
		subcommand = args[0]
		args = args[1:]
	}
	switch subcommand {
	case "list":
		if len(args) != 0 {
			return fmt.Errorf("glossary list takes no arguments")
		}
		terms, err := group.listGlossary(ctx)
		if err != nil {
			return err
		}
		return group.writeOutput(stdout, map[string]any{"terms": terms}, func(w io.Writer) {
			for _, term := range terms {
				fmt.Fprintf(w, "%s\n", term.Term)
			}
		})
	case "get":
		if len(args) == 0 {
			return fmt.Errorf("glossary get requires a term")
		}
		entry, err := group.getGlossaryEntry(ctx, strings.Join(args, " "))
		if err != nil {
			return err
		}
		return group.writeOutput(stdout, entry, func(w io.Writer) {
			fmt.Fprintf(w, "%s\n%s\n", entry.Term, entry.Definition)
		})
	default:
		return fmt.Errorf("unknown glossary command %q", subcommand)
	}
}

func (group CommandGroup) runSections(ctx context.Context, args []string, stdout io.Writer) error {
	subcommand := "list"
	if len(args) > 0 && !strings.HasPrefix(args[0], "--") {
		subcommand = args[0]
		args = args[1:]
	}
	if subcommand != "list" {
		return fmt.Errorf("unknown sections command %q", subcommand)
	}
	groupName, rest, err := parseCommandValueFlag(args, "--group")
	if err != nil {
		return err
	}
	if len(rest) != 0 {
		return fmt.Errorf("sections list got unexpected arguments: %s", strings.Join(rest, " "))
	}
	sections, err := group.listSections(ctx, groupName)
	if err != nil {
		return err
	}
	return group.writeOutput(stdout, map[string]any{"sections": sections}, func(w io.Writer) {
		for _, section := range sections {
			fmt.Fprintf(w, "%s\t%s\t%s\n", section.ID, section.Title, section.URL)
		}
	})
}

func (group CommandGroup) runSection(ctx context.Context, args []string, stdout io.Writer) error {
	if len(args) == 0 || args[0] != "get" {
		return fmt.Errorf("usage: section get <id>")
	}
	if len(args) < 2 {
		return fmt.Errorf("section get requires an id")
	}
	node, err := group.getSection(ctx, strings.Join(args[1:], " "))
	if err != nil {
		return err
	}
	return group.writeOutput(stdout, node, func(w io.Writer) {
		fmt.Fprintf(w, "%s\n%s\n", node.Title, node.Summary)
	})
}

func (group CommandGroup) runOverview(ctx context.Context, stdout io.Writer) error {
	overview, err := group.overview(ctx)
	if err != nil {
		return err
	}
	return group.writeOutput(stdout, overview, func(w io.Writer) {
		if strings.TrimSpace(overview.Context) != "" {
			fmt.Fprintf(w, "%s\n\n", overview.Context)
		}
		fmt.Fprintln(w, overview.Overview)
	})
}

func (group CommandGroup) runSearch(ctx context.Context, query string, stdout io.Writer) error {
	response, err := group.search(ctx, query)
	if err != nil {
		return err
	}
	return group.writeOutput(stdout, response, func(w io.Writer) {
		for _, result := range response.Results {
			fmt.Fprintf(w, "%s\n%s\n%s\n\n", result.Title, result.URL, result.Snippet)
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

func (group CommandGroup) listGlossary(ctx context.Context) ([]GlossaryEntry, error) {
	if group.options.Endpoint != "" {
		return NewEndpointClient(group.options.Endpoint).ListGlossary(ctx)
	}
	digest, err := group.loadDigest(ctx)
	if err != nil {
		return nil, err
	}
	return ListGlossary(digest), nil
}

func (group CommandGroup) getGlossaryEntry(ctx context.Context, term string) (GlossaryEntry, error) {
	if group.options.Endpoint != "" {
		return NewEndpointClient(group.options.Endpoint).GetGlossaryEntry(ctx, term)
	}
	digest, err := group.loadDigest(ctx)
	if err != nil {
		return GlossaryEntry{}, err
	}
	entry, ok := GetGlossaryEntry(digest, term)
	if !ok {
		return GlossaryEntry{}, fmt.Errorf("no glossary entry matched %q", term)
	}
	return entry, nil
}

func (group CommandGroup) listSections(ctx context.Context, groupName string) ([]SectionSummary, error) {
	if group.options.Endpoint != "" {
		return NewEndpointClient(group.options.Endpoint).ListSections(ctx, groupName)
	}
	digest, err := group.loadDigest(ctx)
	if err != nil {
		return nil, err
	}
	return ListSectionSummaries(digest, groupName), nil
}

func (group CommandGroup) getSection(ctx context.Context, id string) (DigestNode, error) {
	if group.options.Endpoint != "" {
		return NewEndpointClient(group.options.Endpoint).GetSection(ctx, id)
	}
	digest, err := group.loadDigest(ctx)
	if err != nil {
		return DigestNode{}, err
	}
	node, ok := GetSection(digest, id)
	if !ok {
		return DigestNode{}, fmt.Errorf("no section matched %q", id)
	}
	return node, nil
}

func (group CommandGroup) overview(ctx context.Context) (Overview, error) {
	if group.options.Endpoint != "" {
		return NewEndpointClient(group.options.Endpoint).Overview(ctx)
	}
	digest, err := group.loadDigest(ctx)
	if err != nil {
		return Overview{}, err
	}
	return GetOverview(digest), nil
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

func parseCommandValueFlag(args []string, name string) (string, []string, error) {
	var value string
	var rest []string
	for i := 0; i < len(args); i++ {
		if args[i] != name {
			rest = append(rest, args[i])
			continue
		}
		if i+1 >= len(args) {
			return "", nil, fmt.Errorf("%s requires a value", name)
		}
		value = args[i+1]
		i++
	}
	return value, rest, nil
}

func writeCommandUsage(w io.Writer) {
	fmt.Fprintln(w, `Usage:
  ask [--digest-dir .hev-ask] [--endpoint URL] [--json] <command>

Commands:
  tree
  ls [path]
  head <path>
  cat <path>
  facts <path>
  grep <query>
  glossary list
  glossary get <term>
  sections list [--group GROUP]
  section get <id>
  overview
  search <query>
  answer <query>
  mcp`)
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
