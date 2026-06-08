package pipeline

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gitlink-org/gitlink-cli/internal/client"
	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

func TestPipelineList(t *testing.T) {
	server := newPipelineTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertPipelineRequest(t, r, "GET", "/pm/pipelines.json")
		assertPipelineQuery(t, r, "owner_id", "42")
		assertPipelineQuery(t, r, "page", "2")
		assertPipelineQuery(t, r, "limit", "50")
		writePipelineJSON(t, w, map[string]interface{}{"status": 0, "message": "ok"})
	})
	defer server.Close()

	err := runPipelineShortcut(t, server, "list", map[string]string{
		"owner-id": "42",
		"page":     "2",
		"limit":    "50",
	})
	if err != nil {
		t.Fatalf("list shortcut failed: %v", err)
	}
}

func TestPipelineRuns(t *testing.T) {
	server := newPipelineTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertPipelineRequest(t, r, "GET", "/v1/owner/repo/actions/runs.json")
		assertPipelineQuery(t, r, "ref", "master")
		assertPipelineQuery(t, r, "workflow", "build.yml")
		writePipelineJSON(t, w, map[string]interface{}{"runs": []interface{}{}})
	})
	defer server.Close()

	err := runPipelineShortcut(t, server, "runs", map[string]string{
		"ref":      "master",
		"workflow": "build.yml",
	})
	if err != nil {
		t.Fatalf("runs shortcut failed: %v", err)
	}
}

func TestPipelineRun(t *testing.T) {
	server := newPipelineTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertPipelineRequest(t, r, "POST", "/v1/owner/repo/actions/runs.json")
		assertPipelineQuery(t, r, "ref", "master")
		assertPipelineQuery(t, r, "workflow", "build.yml")
		writePipelineJSON(t, w, map[string]interface{}{"status": 0})
	})
	defer server.Close()

	err := runPipelineShortcut(t, server, "run", map[string]string{
		"ref":      "master",
		"workflow": "build.yml",
	})
	if err != nil {
		t.Fatalf("run shortcut failed: %v", err)
	}
}

func TestPipelineView(t *testing.T) {
	server := newPipelineTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertPipelineRequest(t, r, "GET", "/v1/owner/repo/pipelines/7.json")
		writePipelineJSON(t, w, map[string]interface{}{"id": 7})
	})
	defer server.Close()

	if err := runPipelineShortcut(t, server, "view", map[string]string{"id": "7"}); err != nil {
		t.Fatalf("view shortcut failed: %v", err)
	}
}

func TestPipelineDelete(t *testing.T) {
	server := newPipelineTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertPipelineRequest(t, r, "DELETE", "/v1/owner/repo/pipelines/7.json")
		writePipelineJSON(t, w, map[string]interface{}{"status": 0})
	})
	defer server.Close()

	if err := runPipelineShortcut(t, server, "delete", map[string]string{"id": "7"}); err != nil {
		t.Fatalf("delete shortcut failed: %v", err)
	}
}

func TestPipelineSaveYamlPayload(t *testing.T) {
	var payload map[string]interface{}
	server := newPipelineTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertPipelineRequest(t, r, "POST", "/v1/owner/repo/pipelines/save_yaml.json")
		payload = decodePipelineJSON(t, r)
		writePipelineJSON(t, w, map[string]interface{}{"status": 0})
	})
	defer server.Close()

	err := runPipelineShortcut(t, server, "save-yaml", map[string]string{
		"id":            "7",
		"pipeline-json": `{"nodes":[]}`,
	})
	if err != nil {
		t.Fatalf("save-yaml shortcut failed: %v", err)
	}

	assertPipelineEqual(t, payload["id"], float64(7))
	graph, ok := payload["pipeline_json"].(map[string]interface{})
	if !ok {
		t.Fatalf("pipeline_json = %T, want object", payload["pipeline_json"])
	}
	nodes, ok := graph["nodes"].([]interface{})
	if !ok {
		t.Fatalf("nodes = %T, want array", graph["nodes"])
	}
	if len(nodes) != 0 {
		t.Fatalf("nodes length = %d, want 0", len(nodes))
	}
}

func TestPipelineEnablePayload(t *testing.T) {
	var payload map[string]interface{}
	server := newPipelineTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertPipelineRequest(t, r, "POST", "/v1/owner/repo/actions/enable.json")
		payload = decodePipelineJSON(t, r)
		writePipelineJSON(t, w, map[string]interface{}{"status": 0})
	})
	defer server.Close()

	err := runPipelineShortcut(t, server, "enable", map[string]string{
		"id":       "7",
		"workflow": "build.yml",
	})
	if err != nil {
		t.Fatalf("enable shortcut failed: %v", err)
	}

	assertPipelineEqual(t, payload["id"], float64(7))
	assertPipelineEqual(t, payload["workflow"], "build.yml")
}

func TestPipelineDisablePayload(t *testing.T) {
	var payload map[string]interface{}
	server := newPipelineTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertPipelineRequest(t, r, "POST", "/v1/owner/repo/actions/disable.json")
		payload = decodePipelineJSON(t, r)
		writePipelineJSON(t, w, map[string]interface{}{"status": 0})
	})
	defer server.Close()

	err := runPipelineShortcut(t, server, "disable", map[string]string{
		"id":       "7",
		"workflow": "build.yml",
	})
	if err != nil {
		t.Fatalf("disable shortcut failed: %v", err)
	}

	assertPipelineEqual(t, payload["id"], float64(7))
	assertPipelineEqual(t, payload["workflow"], "build.yml")
}

func TestPipelineLogsPayload(t *testing.T) {
	var payload map[string]interface{}
	server := newPipelineTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertPipelineRequest(t, r, "POST", "/v1/owner/repo/actions/runs/99/jobs/0.json")
		payload = decodePipelineJSON(t, r)
		writePipelineJSON(t, w, map[string]interface{}{"status": 0})
	})
	defer server.Close()

	err := runPipelineShortcut(t, server, "logs", map[string]string{
		"run-id":   "99",
		"id":       "7",
		"index":    "43",
		"job":      "0",
		"cursor":   "cursor-1",
		"step":     "1",
		"expanded": "true",
	})
	if err != nil {
		t.Fatalf("logs shortcut failed: %v", err)
	}

	assertPipelineEqual(t, payload["id"], float64(7))
	assertPipelineEqual(t, payload["index"], "43")
	assertPipelineEqual(t, payload["job"], float64(0))
	assertPipelineEqual(t, payload["owner"], "owner")
	assertPipelineEqual(t, payload["repo"], "repo")
	cursors, ok := payload["log_cursors"].([]interface{})
	if !ok {
		t.Fatalf("log_cursors = %T, want array", payload["log_cursors"])
	}
	if len(cursors) != 1 {
		t.Fatalf("log_cursors length = %d, want 1", len(cursors))
	}
	cursor, ok := cursors[0].(map[string]interface{})
	if !ok {
		t.Fatalf("cursor entry = %T, want object", cursors[0])
	}
	assertPipelineEqual(t, cursor["cursor"], "cursor-1")
	assertPipelineEqual(t, cursor["expanded"], true)
	assertPipelineEqual(t, cursor["step"], float64(1))
}

func TestPipelineResults(t *testing.T) {
	server := newPipelineTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertPipelineRequest(t, r, "GET", "/v1/owner/repo/pipelines/run_results.json")
		assertPipelineQuery(t, r, "run_id", "99")
		writePipelineJSON(t, w, map[string]interface{}{"reports": []interface{}{}})
	})
	defer server.Close()

	if err := runPipelineShortcut(t, server, "results", map[string]string{"run-id": "99"}); err != nil {
		t.Fatalf("results shortcut failed: %v", err)
	}
}

func TestPipelineDryRunDoesNotCallAPI(t *testing.T) {
	dryRunCases := []struct {
		name string
		args map[string]string
	}{
		{name: "run", args: map[string]string{"dry-run": "true", "ref": "master"}},
		{name: "delete", args: map[string]string{"dry-run": "true", "id": "7"}},
		{name: "save-yaml", args: map[string]string{"dry-run": "true", "id": "7", "pipeline-json": `{"nodes":[]}`}},
		{name: "enable", args: map[string]string{"dry-run": "true", "id": "7", "workflow": "build.yml"}},
		{name: "disable", args: map[string]string{"dry-run": "true", "id": "7", "workflow": "build.yml"}},
	}

	for _, tc := range dryRunCases {
		t.Run(tc.name, func(t *testing.T) {
			server := newPipelineTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				t.Fatalf("dry-run should not call API, got %s %s", r.Method, r.URL.Path)
			})
			defer server.Close()

			if err := runPipelineShortcut(t, server, tc.name, tc.args); err != nil {
				t.Fatalf("%s dry-run failed: %v", tc.name, err)
			}
		})
	}
}

func TestPipelineRejectsInvalidID(t *testing.T) {
	server := newPipelineTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("invalid id should not call API, got %s %s", r.Method, r.URL.Path)
	})
	defer server.Close()

	err := runPipelineShortcut(t, server, "view", map[string]string{"id": "abc"})
	if err == nil {
		t.Fatal("expected invalid id to return an error")
	}
}

func runPipelineShortcut(t *testing.T, server *httptest.Server, name string, args map[string]string) error {
	t.Helper()
	shortcut := findPipelineShortcut(t, name)
	ctx := &common.RuntimeContext{
		Client: &client.Client{
			HTTP:    server.Client(),
			BaseURL: server.URL,
		},
		Owner:  "owner",
		Repo:   "repo",
		Format: "json",
		Args:   args,
	}
	if ctx.Args == nil {
		ctx.Args = map[string]string{}
	}
	return shortcut.Run(ctx)
}

func findPipelineShortcut(t *testing.T, name string) *common.Shortcut {
	t.Helper()
	for _, shortcut := range Shortcuts() {
		if shortcut.Name == name {
			return shortcut
		}
	}
	t.Fatalf("shortcut %q not found", name)
	return nil
}

func newPipelineTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

func assertPipelineRequest(t *testing.T, r *http.Request, method, path string) {
	t.Helper()
	if r.Method != method || r.URL.Path != path {
		t.Fatalf("got request %s %s, want %s %s", r.Method, r.URL.Path, method, path)
	}
}

func assertPipelineQuery(t *testing.T, r *http.Request, key, want string) {
	t.Helper()
	if got := r.URL.Query().Get(key); got != want {
		t.Fatalf("query %s = %q, want %q", key, got, want)
	}
}

func decodePipelineJSON(t *testing.T, r *http.Request) map[string]interface{} {
	t.Helper()
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode request body: %v", err)
	}
	return payload
}

func writePipelineJSON(t *testing.T, w http.ResponseWriter, payload interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Fatalf("failed to write response: %v", err)
	}
}

func assertPipelineEqual(t *testing.T, got interface{}, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v (%T), want %v (%T)", got, got, want, want)
	}
}

func TestPipelineShortcutNames(t *testing.T) {
	got := map[string]bool{}
	for _, shortcut := range Shortcuts() {
		got[shortcut.Name] = true
	}
	want := []string{"list", "runs", "run", "view", "delete", "save-yaml", "enable", "disable", "logs", "results"}
	for _, name := range want {
		if !got[name] {
			t.Fatalf("missing shortcut %q in %v", name, got)
		}
	}
	if len(got) != len(want) {
		t.Fatalf("shortcut count = %d, want %d: %v", len(got), len(want), got)
	}
}

func TestPipelineRunIDIsPathEscaped(t *testing.T) {
	server := newPipelineTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("got method %s, want POST", r.Method)
		}
		if got, want := r.URL.EscapedPath(), "/v1/owner/repo/actions/runs/run%2F99/jobs/0.json"; got != want {
			t.Fatalf("escaped path = %q, want %q", got, want)
		}
		writePipelineJSON(t, w, map[string]interface{}{"status": 0})
	})
	defer server.Close()

	err := runPipelineShortcut(t, server, "logs", map[string]string{
		"run-id":   "run/99",
		"id":       "7",
		"index":    "43",
		"expanded": "true",
	})
	if err != nil {
		t.Fatalf("logs shortcut failed: %v", err)
	}
}

func ExampleShortcuts() {
	for _, shortcut := range Shortcuts() {
		fmt.Println(shortcut.Name)
	}
	// Output:
	// list
	// runs
	// run
	// view
	// delete
	// save-yaml
	// enable
	// disable
	// logs
	// results
}
