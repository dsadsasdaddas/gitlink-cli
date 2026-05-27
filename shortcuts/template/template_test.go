package template

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

func TestTemplateList(t *testing.T) {
	server := newTemplateTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertTemplateRequest(t, r, "GET", "/v1/owner/repo/project_templates.json")
		writeTemplateJSON(t, w, map[string]interface{}{"total_count": 0, "project_templates": []interface{}{}})
	})
	defer server.Close()

	if err := runTemplateShortcut(t, server, "list", nil); err != nil {
		t.Fatalf("list shortcut failed: %v", err)
	}
}

func TestTemplateView(t *testing.T) {
	server := newTemplateTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertTemplateRequest(t, r, "GET", "/v1/owner/repo/project_templates/7.json")
		writeTemplateJSON(t, w, map[string]interface{}{"project_template": templateFixture()})
	})
	defer server.Close()

	if err := runTemplateShortcut(t, server, "view", map[string]string{"id": "7"}); err != nil {
		t.Fatalf("view shortcut failed: %v", err)
	}
}

func TestTemplateCreatePayload(t *testing.T) {
	var payload map[string]interface{}
	server := newTemplateTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertTemplateRequest(t, r, "POST", "/v1/owner/repo/project_templates.json")
		payload = decodeTemplateJSON(t, r)
		writeTemplateJSON(t, w, map[string]interface{}{"status": 0})
	})
	defer server.Close()

	err := runTemplateShortcut(t, server, "create", map[string]string{
		"type":    "ProjectTemplates::Issue",
		"name":    "Bug report",
		"content": "Bug template",
	})
	if err != nil {
		t.Fatalf("create shortcut failed: %v", err)
	}

	assertTemplateEqual(t, payload["type"], "ProjectTemplates::Issue")
	assertTemplateEqual(t, payload["name"], "Bug report")
	assertTemplateEqual(t, payload["content"], "Bug template")
}

func TestTemplateUpdatePreservesExistingFields(t *testing.T) {
	var payload map[string]interface{}
	server := newTemplateTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/v1/owner/repo/project_templates/7.json":
			writeTemplateJSON(t, w, map[string]interface{}{"project_template": templateFixture()})
		case r.Method == "PUT" && r.URL.Path == "/v1/owner/repo/project_templates/7.json":
			payload = decodeTemplateJSON(t, r)
			writeTemplateJSON(t, w, map[string]interface{}{"status": 0})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})
	defer server.Close()

	err := runTemplateShortcut(t, server, "update", map[string]string{
		"id":      "7",
		"content": "Updated content",
	})
	if err != nil {
		t.Fatalf("update shortcut failed: %v", err)
	}

	assertTemplateEqual(t, payload["type"], "ProjectTemplates::Issue")
	assertTemplateEqual(t, payload["name"], "Bug report")
	assertTemplateEqual(t, payload["content"], "Updated content")
}

func TestTemplateDelete(t *testing.T) {
	server := newTemplateTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertTemplateRequest(t, r, "DELETE", "/v1/owner/repo/project_templates/7.json")
		writeTemplateJSON(t, w, map[string]interface{}{"status": 0})
	})
	defer server.Close()

	if err := runTemplateShortcut(t, server, "delete", map[string]string{"id": "7"}); err != nil {
		t.Fatalf("delete shortcut failed: %v", err)
	}
}

func TestTemplateDryRunDoesNotWrite(t *testing.T) {
	tests := []struct {
		name string
		args map[string]string
	}{
		{name: "create", args: map[string]string{"name": "Bug report", "content": "Bug template", "dry-run": "true"}},
		{name: "delete", args: map[string]string{"id": "7", "dry-run": "true"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := newTemplateTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				t.Fatalf("dry-run should not call API, got %s %s", r.Method, r.URL.Path)
			})
			defer server.Close()

			if err := runTemplateShortcut(t, server, tc.name, tc.args); err != nil {
				t.Fatalf("%s dry-run failed: %v", tc.name, err)
			}
		})
	}
}

func TestTemplateUpdateDryRunFetchesCurrentButDoesNotWrite(t *testing.T) {
	server := newTemplateTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			t.Fatalf("dry-run should not update template, got %s %s", r.Method, r.URL.Path)
		}
		assertTemplateRequest(t, r, "GET", "/v1/owner/repo/project_templates/7.json")
		writeTemplateJSON(t, w, map[string]interface{}{"project_template": templateFixture()})
	})
	defer server.Close()

	err := runTemplateShortcut(t, server, "update", map[string]string{
		"id":      "7",
		"name":    "Updated name",
		"dry-run": "true",
	})
	if err != nil {
		t.Fatalf("update dry-run failed: %v", err)
	}
}

func TestTemplateCreateRejectsMissingContent(t *testing.T) {
	server := newTemplateTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("missing content should not call API, got %s %s", r.Method, r.URL.Path)
	})
	defer server.Close()

	err := runTemplateShortcut(t, server, "create", map[string]string{"name": "Bug report"})
	if err == nil {
		t.Fatal("expected missing content to return an error")
	}
}

func TestTemplateUpdateRejectsNoFields(t *testing.T) {
	server := newTemplateTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("empty update should not call API, got %s %s", r.Method, r.URL.Path)
	})
	defer server.Close()

	err := runTemplateShortcut(t, server, "update", map[string]string{"id": "7"})
	if err == nil {
		t.Fatal("expected empty update to return an error")
	}
}

func TestTemplateShortcutNames(t *testing.T) {
	got := map[string]bool{}
	for _, shortcut := range Shortcuts() {
		got[shortcut.Name] = true
	}
	want := []string{"list", "view", "create", "update", "delete"}
	for _, name := range want {
		if !got[name] {
			t.Fatalf("missing shortcut %q in %v", name, got)
		}
	}
	if len(got) != len(want) {
		t.Fatalf("shortcut count = %d, want %d: %v", len(got), len(want), got)
	}
}

func runTemplateShortcut(t *testing.T, server *httptest.Server, name string, args map[string]string) error {
	t.Helper()
	shortcut := findTemplateShortcut(t, name)
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

func findTemplateShortcut(t *testing.T, name string) *common.Shortcut {
	t.Helper()
	for _, shortcut := range Shortcuts() {
		if shortcut.Name == name {
			return shortcut
		}
	}
	t.Fatalf("shortcut %q not found", name)
	return nil
}

func newTemplateTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

func templateFixture() map[string]interface{} {
	return map[string]interface{}{
		"id":      7,
		"type":    "ProjectTemplates::Issue",
		"name":    "Bug report",
		"content": "Old content",
	}
}

func assertTemplateRequest(t *testing.T, r *http.Request, method, path string) {
	t.Helper()
	if r.Method != method || r.URL.Path != path {
		t.Fatalf("got request %s %s, want %s %s", r.Method, r.URL.Path, method, path)
	}
}

func decodeTemplateJSON(t *testing.T, r *http.Request) map[string]interface{} {
	t.Helper()
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode request body: %v", err)
	}
	return payload
}

func writeTemplateJSON(t *testing.T, w http.ResponseWriter, payload interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Fatalf("failed to write response: %v", err)
	}
}

func assertTemplateEqual(t *testing.T, got interface{}, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v (%T), want %v (%T)", got, got, want, want)
	}
}

func ExampleShortcuts() {
	for _, shortcut := range Shortcuts() {
		fmt.Println(shortcut.Name)
	}
	// Output:
	// list
	// view
	// create
	// update
	// delete
}
