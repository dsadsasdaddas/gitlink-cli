package access

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gitlink-org/gitlink-cli/internal/client"
	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

func runAccessShortcut(t *testing.T, server *httptest.Server, name string, args map[string]string) error {
	t.Helper()
	shortcut := findAccessShortcut(t, name)
	ctx := &common.RuntimeContext{
		Client: &client.Client{HTTP: server.Client(), BaseURL: server.URL},
		Owner:  "owner",
		Repo:   "repo",
		Format: "json",
		Args:   args,
	}
	return shortcut.Run(ctx)
}

func findAccessShortcut(t *testing.T, name string) *common.Shortcut {
	t.Helper()
	for _, s := range Shortcuts() {
		if s.Name == name {
			return s
		}
	}
	t.Fatalf("shortcut %q not found", name)
	return nil
}

func writeAccessJSON(t *testing.T, w http.ResponseWriter, v interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Fatalf("write json: %v", err)
	}
}

func assertAccessRequest(t *testing.T, r *http.Request, method, path string) {
	t.Helper()
	if r.Method != method {
		t.Fatalf("expected method %s, got %s", method, r.Method)
	}
	if r.URL.Path != path {
		t.Fatalf("expected path %s, got %s", path, r.URL.Path)
	}
}

func TestJoinProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertAccessRequest(t, r, "POST", "/applied_projects.json")
		var body map[string]map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		applied := body["applied_project"]
		if applied["code"] != "invite-code" {
			t.Fatalf("expected code invite-code, got %q", applied["code"])
		}
		if applied["role"] != "reporter" {
			t.Fatalf("expected role reporter, got %q", applied["role"])
		}
		writeAccessJSON(t, w, map[string]interface{}{"id": 7, "status": "common"})
	}))
	defer server.Close()

	if err := runAccessShortcut(t, server, "join", map[string]string{"code": "invite-code", "role": "Reporter"}); err != nil {
		t.Fatalf("join failed: %v", err)
	}
}

func TestJoinProjectDefaultRole(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertAccessRequest(t, r, "POST", "/applied_projects.json")
		var body map[string]map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got := body["applied_project"]["role"]; got != "developer" {
			t.Fatalf("expected default role developer, got %q", got)
		}
		writeAccessJSON(t, w, map[string]interface{}{"id": 7})
	}))
	defer server.Close()

	if err := runAccessShortcut(t, server, "join", map[string]string{"code": "invite-code"}); err != nil {
		t.Fatalf("join failed: %v", err)
	}
}

func TestJoinProjectDryRun(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no API call expected for dry-run")
	}))
	defer server.Close()

	if err := runAccessShortcut(t, server, "join", map[string]string{"code": "invite-code", "role": "developer", "dry-run": "true"}); err != nil {
		t.Fatalf("join dry-run failed: %v", err)
	}
}

func TestJoinProjectInvalidRole(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no API call expected")
	}))
	defer server.Close()

	if err := runAccessShortcut(t, server, "join", map[string]string{"code": "invite-code", "role": "owner"}); err == nil {
		t.Fatal("expected invalid role error")
	}
}

func TestJoinProjectMissingCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no API call expected")
	}))
	defer server.Close()

	if err := runAccessShortcut(t, server, "join", map[string]string{}); err == nil {
		t.Fatal("expected missing code error")
	}
}

func TestQuitProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertAccessRequest(t, r, "POST", "/owner/repo/quit.json")
		writeAccessJSON(t, w, map[string]interface{}{"status": 0, "message": "success"})
	}))
	defer server.Close()

	if err := runAccessShortcut(t, server, "quit", map[string]string{}); err != nil {
		t.Fatalf("quit failed: %v", err)
	}
}

func TestQuitProjectDryRun(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no API call expected for dry-run")
	}))
	defer server.Close()

	if err := runAccessShortcut(t, server, "quit", map[string]string{"dry-run": "true"}); err != nil {
		t.Fatalf("quit dry-run failed: %v", err)
	}
}

func TestNormalizeRole(t *testing.T) {
	cases := map[string]string{
		"":          "developer",
		"MANAGER":   "manager",
		"developer": "developer",
		"Reporter":  "reporter",
	}
	for input, want := range cases {
		got, err := normalizeRole(input)
		if err != nil {
			t.Fatalf("normalizeRole(%q) failed: %v", input, err)
		}
		if got != want {
			t.Fatalf("normalizeRole(%q)=%q, want %q", input, got, want)
		}
	}
}
