package ask

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestServeMCPFetchDocsMaterializesTree(t *testing.T) {
	path := writeMCPTestGraph(t)
	cache := filepath.Join(t.TempDir(), "cache")
	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"fetch_docs","arguments":{}}}`,
	}, "\n") + "\n"

	var out bytes.Buffer
	if err := ServeMCP(context.Background(), MCPOptions{DigestPath: path, CacheDir: cache}, strings.NewReader(input), &out); err != nil {
		t.Fatal(err)
	}

	responses := decodeMCPTestResponses(t, out.String())
	if len(responses) != 3 {
		t.Fatalf("expected 4 responses, got %d: %s", len(responses), out.String())
	}

	var initialize struct {
		ProtocolVersion string `json:"protocolVersion"`
		Instructions    string `json:"instructions"`
		ServerInfo      struct {
			Name string `json:"name"`
		} `json:"serverInfo"`
	}
	decodeMCPResult(t, responses[0], &initialize)
	if initialize.ProtocolVersion == "" || initialize.ServerInfo.Name != "hev-ask" || !strings.Contains(initialize.Instructions, "fetch_docs") {
		t.Fatalf("unexpected initialize result: %#v", initialize)
	}

	var tools struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	}
	decodeMCPResult(t, responses[1], &tools)
	if len(tools.Tools) != 1 || tools.Tools[0].Name != "fetch_docs" {
		t.Fatalf("expected only fetch_docs, got %#v", tools.Tools)
	}

	var fetched struct {
		IsError           bool            `json:"isError"`
		StructuredContent fetchDocsResult `json:"structuredContent"`
	}
	decodeMCPResult(t, responses[2], &fetched)
	if fetched.IsError || fetched.StructuredContent.Path == "" || fetched.StructuredContent.Sections != 2 || !strings.Contains(fetched.StructuredContent.Tree, "CLI > Flags") {
		t.Fatalf("unexpected fetch_docs result: %#v", fetched)
	}
	if _, err := LoadDigest(fetched.StructuredContent.Path); err != nil {
		t.Fatalf("cache path was not a readable digest tree: %v", err)
	}
}

func TestServeMCPUnknownOldToolReturnsToolError(t *testing.T) {
	path := writeMCPTestGraph(t)
	input := `{"jsonrpc":"2.0","id":"answer","method":"tools/call","params":{"name":"answer","arguments":{"query":"How does it work?"}}}` + "\n"

	var out bytes.Buffer
	if err := ServeMCP(context.Background(), MCPOptions{DigestPath: path}, strings.NewReader(input), &out); err != nil {
		t.Fatal(err)
	}

	responses := decodeMCPTestResponses(t, out.String())
	if len(responses) != 1 {
		t.Fatalf("expected 1 response, got %d: %s", len(responses), out.String())
	}
	var result struct {
		IsError bool         `json:"isError"`
		Content []mcpContent `json:"content"`
	}
	decodeMCPResult(t, responses[0], &result)
	if !result.IsError || len(result.Content) != 1 || !strings.Contains(result.Content[0].Text, `unknown tool "answer"`) {
		t.Fatalf("unexpected old-tool error result: %#v", result)
	}
}

func TestServeMCPFetchDocsUsesRemoteArchiveAndHashCache(t *testing.T) {
	digest := testDigest()
	digest.ContentHash = "remote-hash"
	archiveGets := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ask/archive" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("x-hev-ask-content-hash", digest.ContentHash)
		if r.Method == http.MethodHead {
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		archiveGets++
		w.Header().Set("content-type", "application/gzip")
		if err := WriteDigestArchive(w, digest); err != nil {
			t.Fatalf("write archive: %v", err)
		}
	}))
	defer server.Close()

	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"fetch_docs","arguments":{}}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"fetch_docs","arguments":{}}}`,
	}, "\n") + "\n"

	var out bytes.Buffer
	if err := ServeMCP(context.Background(), MCPOptions{Endpoint: server.URL + "/api/ask", CacheDir: filepath.Join(t.TempDir(), "cache")}, strings.NewReader(input), &out); err != nil {
		t.Fatal(err)
	}
	responses := decodeMCPTestResponses(t, out.String())
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d: %s", len(responses), out.String())
	}
	var first, second struct {
		IsError           bool            `json:"isError"`
		StructuredContent fetchDocsResult `json:"structuredContent"`
	}
	decodeMCPResult(t, responses[0], &first)
	decodeMCPResult(t, responses[1], &second)
	if first.IsError || second.IsError {
		t.Fatalf("unexpected MCP error: first=%#v second=%#v", first, second)
	}
	if first.StructuredContent.UpToDate || !second.StructuredContent.UpToDate {
		t.Fatalf("unexpected cache states: first=%#v second=%#v", first.StructuredContent, second.StructuredContent)
	}
	if archiveGets != 1 {
		t.Fatalf("expected one archive GET, got %d", archiveGets)
	}
	if _, err := LoadDigest(first.StructuredContent.Path); err != nil {
		t.Fatalf("remote archive did not materialize a readable tree: %v", err)
	}
}

func writeMCPTestGraph(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".hev-ask")
	if err := WriteDigest(path, testDigest()); err != nil {
		t.Fatal(err)
	}
	return path
}

type mcpTestRPCResponse struct {
	Result json.RawMessage  `json:"result"`
	Error  *rpcErrorPayload `json:"error"`
}

func decodeMCPTestResponses(t *testing.T, output string) []mcpTestRPCResponse {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	responses := make([]mcpTestRPCResponse, 0, len(lines))
	for _, line := range lines {
		var response mcpTestRPCResponse
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			t.Fatalf("decode response %q: %v", line, err)
		}
		if response.Error != nil {
			t.Fatalf("unexpected rpc error in %q: %#v", line, response.Error)
		}
		responses = append(responses, response)
	}
	return responses
}

func decodeMCPResult(t *testing.T, response mcpTestRPCResponse, out any) {
	t.Helper()
	if err := json.Unmarshal(response.Result, out); err != nil {
		t.Fatalf("decode result %s: %v", response.Result, err)
	}
}
