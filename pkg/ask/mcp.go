package ask

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const mcpProtocolVersion = "2025-06-18"

type MCPOptions struct {
	DigestPath string
	Endpoint   string
	MaxResults int
	CacheDir   string
}

type rpcRequest struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      json.RawMessage  `json:"id"`
	Result  any              `json:"result,omitempty"`
	Error   *rpcErrorPayload `json:"error,omitempty"`
}

type rpcErrorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type mcpContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type mcpToolResult struct {
	Content           []mcpContent `json:"content"`
	StructuredContent any          `json:"structuredContent,omitempty"`
	IsError           bool         `json:"isError,omitempty"`
}

type mcpServer struct {
	options MCPOptions
}

func ServeMCP(ctx context.Context, options MCPOptions, in io.Reader, out io.Writer) error {
	if options.DigestPath == "" {
		options.DigestPath = ".hev-ask"
	}
	if options.MaxResults <= 0 {
		options.MaxResults = 8
	}

	server := mcpServer{options: options}
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if err := server.handleLine(ctx, []byte(line), out); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return ctx.Err()
}

func (server mcpServer) handleLine(ctx context.Context, line []byte, out io.Writer) error {
	var request rpcRequest
	if err := json.Unmarshal(line, &request); err != nil {
		return writeRPCError(out, nil, -32700, "parse error")
	}
	if request.Method == "" {
		if request.ID == nil {
			return nil
		}
		return writeRPCError(out, request.ID, -32600, "invalid request")
	}

	result, rpcErr := server.handleRequest(ctx, request)
	if request.ID == nil || strings.HasPrefix(request.Method, "notifications/") {
		return nil
	}
	if rpcErr != nil {
		return writeRPCError(out, request.ID, rpcErr.Code, rpcErr.Message)
	}
	return writeRPCResult(out, request.ID, result)
}

func (server mcpServer) handleRequest(ctx context.Context, request rpcRequest) (any, *rpcErrorPayload) {
	switch request.Method {
	case "initialize":
		return map[string]any{
			"protocolVersion": mcpProtocolVersion,
			"capabilities": map[string]any{
				"tools": map[string]any{},
			},
			"instructions": mcpInstructions(),
			"serverInfo": map[string]any{
				"name":    "hev-ask",
				"version": "0.0.1",
			},
		}, nil
	case "notifications/initialized":
		return nil, nil
	case "tools/list":
		return map[string]any{"tools": mcpTools()}, nil
	case "tools/call":
		var params struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		}
		if err := decodeObject(request.Params, &params); err != nil {
			return nil, &rpcErrorPayload{Code: -32602, Message: err.Error()}
		}
		if strings.TrimSpace(params.Name) == "" {
			return nil, &rpcErrorPayload{Code: -32602, Message: "tools/call requires a tool name"}
		}
		result, err := server.callTool(ctx, params.Name, params.Arguments)
		if err != nil {
			return mcpErrorResult(err), nil
		}
		return result, nil
	default:
		return nil, &rpcErrorPayload{Code: -32601, Message: "method not found"}
	}
}

func mcpTools() []map[string]any {
	return []map[string]any{
		{
			"name":        "fetch_docs",
			"description": "Materialize the hev ask digest tree to local disk and return its title tree.",
			"inputSchema": objectSchema(map[string]any{
				"force": map[string]any{"type": "boolean", "description": "Rebuild the local cache even when the content hash is unchanged."},
			}, nil),
		},
	}
}

func objectSchema(properties map[string]any, required []string) map[string]any {
	if properties == nil {
		properties = map[string]any{}
	}
	schema := map[string]any{
		"type":                 "object",
		"properties":           properties,
		"additionalProperties": false,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func stringSchema(description string) map[string]any {
	return map[string]any{"type": "string", "description": description}
}

func (server mcpServer) callTool(ctx context.Context, name string, arguments json.RawMessage) (mcpToolResult, error) {
	switch name {
	case "fetch_docs":
		var args struct {
			Force bool `json:"force"`
		}
		if err := decodeObject(arguments, &args); err != nil {
			return mcpToolResult{}, err
		}
		payload, err := server.fetchDocs(ctx, args.Force)
		if err != nil {
			return mcpToolResult{}, err
		}
		return mcpStructuredResult(payload, payload.Tree), nil
	default:
		return mcpToolResult{}, fmt.Errorf("unknown tool %q", name)
	}
}

type fetchDocsResult struct {
	Path        string `json:"path"`
	ContentHash string `json:"contentHash"`
	Sections    int    `json:"sections"`
	Tree        string `json:"tree"`
	UpToDate    bool   `json:"upToDate"`
}

func (server mcpServer) fetchDocs(ctx context.Context, force bool) (fetchDocsResult, error) {
	if server.options.Endpoint != "" {
		return server.fetchRemoteDocs(ctx, force)
	}
	digest, err := server.loadDigest(ctx)
	if err != nil {
		return fetchDocsResult{}, err
	}
	cachePath, err := server.cachePath()
	if err != nil {
		return fetchDocsResult{}, err
	}
	upToDate := false
	if !force {
		if cached, err := LoadDigest(cachePath); err == nil && cached.ContentHash != "" && cached.ContentHash == digest.ContentHash {
			upToDate = true
		}
	}
	if !upToDate {
		if err := WriteDigestTree(cachePath, digest); err != nil {
			return fetchDocsResult{}, err
		}
	}
	return fetchDocsResult{
		Path:        cachePath,
		ContentHash: digest.ContentHash,
		Sections:    len(digest.Nodes),
		Tree:        RenderDigestTree(digest),
		UpToDate:    upToDate,
	}, nil
}

func (server mcpServer) fetchRemoteDocs(ctx context.Context, force bool) (fetchDocsResult, error) {
	cachePath, err := server.cachePath()
	if err != nil {
		return fetchDocsResult{}, err
	}
	client := NewEndpointClient(server.options.Endpoint)
	remoteHash, hashErr := client.ArchiveContentHash(ctx)
	if hashErr == nil && !force && strings.TrimSpace(remoteHash) != "" {
		if cached, err := LoadDigest(cachePath); err == nil && cached.ContentHash == remoteHash {
			return fetchDocsResult{
				Path:        cachePath,
				ContentHash: cached.ContentHash,
				Sections:    len(cached.Nodes),
				Tree:        RenderDigestTree(cached),
				UpToDate:    true,
			}, nil
		}
	}
	digest, err := client.DownloadArchive(ctx, cachePath)
	if err != nil {
		// Older deployed sites may not expose /archive yet. Keep cross-site
		// hydrate usable by falling back to keyless JSON reads.
		if hashErr != nil || strings.Contains(err.Error(), "404") {
			digest, err = client.Digest(ctx)
			if err != nil {
				return fetchDocsResult{}, err
			}
			if err := WriteDigestTree(cachePath, digest); err != nil {
				return fetchDocsResult{}, err
			}
		} else {
			return fetchDocsResult{}, err
		}
	}
	return fetchDocsResult{
		Path:        cachePath,
		ContentHash: digest.ContentHash,
		Sections:    len(digest.Nodes),
		Tree:        RenderDigestTree(digest),
		UpToDate:    false,
	}, nil
}

func (server mcpServer) loadDigest(ctx context.Context) (Digest, error) {
	if server.options.Endpoint != "" {
		return NewEndpointClient(server.options.Endpoint).Digest(ctx)
	}
	return LoadDigest(server.options.DigestPath)
}

func (server mcpServer) cachePath() (string, error) {
	base := server.options.CacheDir
	if base == "" {
		cache, err := os.UserCacheDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(cache, "hev-ask")
	}
	return filepath.Join(base, server.cacheKey()), nil
}

func (server mcpServer) cacheKey() string {
	if server.options.Endpoint != "" {
		parsed, err := url.Parse(server.options.Endpoint)
		if err == nil && parsed.Host != "" {
			return sanitizeCacheKey(parsed.Host)
		}
		return "endpoint-" + SHA256Hex(server.options.Endpoint)[:12]
	}
	abs, err := filepath.Abs(server.options.DigestPath)
	if err != nil {
		abs = server.options.DigestPath
	}
	return "local-" + SHA256Hex(abs)[:12]
}

func sanitizeCacheKey(value string) string {
	replacer := strings.NewReplacer("/", "-", "\\", "-", ":", "-", " ", "-")
	cleaned := strings.Trim(replacer.Replace(value), "-")
	if cleaned == "" {
		return "docs"
	}
	return cleaned
}

func mcpInstructions() string {
	return "This server gives you a local hev ask digest tree. Call fetch_docs first. Read the returned title tree, then use filesystem tools on the returned path: ls/tree for titles, cat only for relevant section files, grep for specifics, and _glossary/ for aliases. Answer from those files and cite claims with each section file's url plus anchor frontmatter as a deep link. Do not read every file unless the task explicitly requires a full audit."
}

func (server mcpServer) listGlossary(ctx context.Context) ([]GlossaryEntry, error) {
	if server.options.Endpoint != "" {
		return NewEndpointClient(server.options.Endpoint).ListGlossary(ctx)
	}
	digest, err := LoadDigest(server.options.DigestPath)
	if err != nil {
		return nil, err
	}
	return ListGlossary(digest), nil
}

func (server mcpServer) getGlossaryEntry(ctx context.Context, term string) (GlossaryEntry, error) {
	if server.options.Endpoint != "" {
		return NewEndpointClient(server.options.Endpoint).GetGlossaryEntry(ctx, term)
	}
	digest, err := LoadDigest(server.options.DigestPath)
	if err != nil {
		return GlossaryEntry{}, err
	}
	entry, ok := GetGlossaryEntry(digest, term)
	if !ok {
		return GlossaryEntry{}, fmt.Errorf("no glossary entry matched %q", term)
	}
	return entry, nil
}

func (server mcpServer) listSections(ctx context.Context, group string) ([]SectionSummary, error) {
	if server.options.Endpoint != "" {
		return NewEndpointClient(server.options.Endpoint).ListSections(ctx, group)
	}
	digest, err := LoadDigest(server.options.DigestPath)
	if err != nil {
		return nil, err
	}
	return ListSectionSummaries(digest, group), nil
}

func (server mcpServer) getSection(ctx context.Context, id string) (DigestNode, error) {
	if server.options.Endpoint != "" {
		return NewEndpointClient(server.options.Endpoint).GetSection(ctx, id)
	}
	digest, err := LoadDigest(server.options.DigestPath)
	if err != nil {
		return DigestNode{}, err
	}
	node, ok := GetSection(digest, id)
	if !ok {
		return DigestNode{}, fmt.Errorf("no section matched %q", id)
	}
	return node, nil
}

func (server mcpServer) overview(ctx context.Context) (Overview, error) {
	if server.options.Endpoint != "" {
		return NewEndpointClient(server.options.Endpoint).Overview(ctx)
	}
	digest, err := LoadDigest(server.options.DigestPath)
	if err != nil {
		return Overview{}, err
	}
	return GetOverview(digest), nil
}

func (server mcpServer) search(ctx context.Context, query string, maxResults int) (KeywordResponse, error) {
	if server.options.Endpoint != "" {
		return NewEndpointClient(server.options.Endpoint).Search(ctx, query)
	}
	digest, err := LoadDigest(server.options.DigestPath)
	if err != nil {
		return KeywordResponse{}, err
	}
	if maxResults <= 0 {
		maxResults = server.options.MaxResults
	}
	return SearchDigest(digest, query, SearchOptions{MaxResults: maxResults}), nil
}

func (server mcpServer) answer(ctx context.Context, query string) (any, string, error) {
	if server.options.Endpoint == "" {
		return nil, "", errors.New("answer requires --endpoint for the remote SSE answer path; without --endpoint, use search for keyless local retrieval")
	}

	var text strings.Builder
	var fallback *KeywordResponse
	err := NewEndpointClient(server.options.Endpoint).StreamAnswer(ctx, query, func(event AnswerEvent) error {
		switch event.Event {
		case "token":
			var payload struct {
				Text string `json:"text"`
			}
			if err := json.Unmarshal(event.Data, &payload); err != nil {
				return err
			}
			text.WriteString(payload.Text)
		case "keyword":
			var response KeywordResponse
			if err := json.Unmarshal(event.Data, &response); err != nil {
				return err
			}
			fallback = &response
		}
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	if fallback != nil {
		payload := map[string]any{"fallback": fallback}
		return payload, formatJSON(payload), nil
	}
	answer := strings.TrimSpace(text.String())
	payload := map[string]any{"answer": answer}
	return payload, answer, nil
}

func decodeObject(raw json.RawMessage, out any) error {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		trimmed = "{}"
	}
	if !strings.HasPrefix(trimmed, "{") {
		return errors.New("expected object parameters")
	}
	decoder := json.NewDecoder(strings.NewReader(trimmed))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		return err
	}
	return nil
}

func mcpStructuredResult(structured any, text string) mcpToolResult {
	return mcpToolResult{
		Content:           []mcpContent{{Type: "text", Text: text}},
		StructuredContent: structured,
	}
}

func mcpErrorResult(err error) mcpToolResult {
	return mcpToolResult{
		Content: []mcpContent{{
			Type: "text",
			Text: err.Error(),
		}},
		IsError: true,
	}
}

func formatJSON(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprint(value)
	}
	return string(data)
}

func writeRPCResult(out io.Writer, id *json.RawMessage, result any) error {
	return writeRPCResponse(out, rpcResponse{JSONRPC: "2.0", ID: rawID(id), Result: result})
}

func writeRPCError(out io.Writer, id *json.RawMessage, code int, message string) error {
	return writeRPCResponse(out, rpcResponse{
		JSONRPC: "2.0",
		ID:      rawID(id),
		Error:   &rpcErrorPayload{Code: code, Message: message},
	})
}

func writeRPCResponse(out io.Writer, response rpcResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	if _, err := out.Write(append(data, '\n')); err != nil {
		return err
	}
	return nil
}

func rawID(id *json.RawMessage) json.RawMessage {
	if id == nil || len(*id) == 0 {
		return json.RawMessage("null")
	}
	return *id
}
