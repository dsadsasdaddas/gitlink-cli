package todo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gitlink-org/gitlink-cli/internal/client"
	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

func runTodoShortcut(t *testing.T, server *httptest.Server, name string, args map[string]string) error {
	t.Helper()
	shortcut := findTodoShortcut(t, name)
	ctx := &common.RuntimeContext{
		Client: &client.Client{HTTP: server.Client(), BaseURL: server.URL},
		Format: "json",
		Args:   args,
	}
	return shortcut.Run(ctx)
}

func findTodoShortcut(t *testing.T, name string) *common.Shortcut {
	t.Helper()
	for _, s := range Shortcuts() {
		if s.Name == name {
			return s
		}
	}
	t.Fatalf("shortcut %q not found", name)
	return nil
}

func writeTodoJSON(t *testing.T, w http.ResponseWriter, v interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Fatalf("write json: %v", err)
	}
}

func assertTodoRequest(t *testing.T, r *http.Request, method, path string) {
	t.Helper()
	if r.Method != method {
		t.Fatalf("expected method %s, got %s", method, r.Method)
	}
	if r.URL.Path != path {
		t.Fatalf("expected path %s, got %s", path, r.URL.Path)
	}
}

func TestTransferList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertTodoRequest(t, r, "GET", "/users/alice/applied_transfer_projects.json")
		if r.URL.Query().Get("page") != "2" {
			t.Fatalf("expected page=2, got %q", r.URL.Query().Get("page"))
		}
		if r.URL.Query().Get("per_page") != "50" {
			t.Fatalf("expected per_page=50, got %q", r.URL.Query().Get("per_page"))
		}
		writeTodoJSON(t, w, map[string]interface{}{"total_count": 0, "applied_transfer_projects": []interface{}{}})
	}))
	defer server.Close()

	if err := runTodoShortcut(t, server, "transfer-list", map[string]string{"login": "alice", "page": "2", "per-page": "50"}); err != nil {
		t.Fatalf("transfer-list failed: %v", err)
	}
}

func TestTransferAccept(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertTodoRequest(t, r, "POST", "/users/alice/applied_transfer_projects/7/accept.json")
		writeTodoJSON(t, w, map[string]interface{}{"id": 7, "status": "accepted"})
	}))
	defer server.Close()

	if err := runTodoShortcut(t, server, "transfer-accept", map[string]string{"login": "alice", "id": "7"}); err != nil {
		t.Fatalf("transfer-accept failed: %v", err)
	}
}

func TestTransferRefuseDryRun(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no API call expected for dry-run")
	}))
	defer server.Close()

	if err := runTodoShortcut(t, server, "transfer-refuse", map[string]string{"login": "alice", "id": "7", "dry-run": "true"}); err != nil {
		t.Fatalf("transfer-refuse dry-run failed: %v", err)
	}
}

func TestJoinList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertTodoRequest(t, r, "GET", "/users/bob/applied_projects.json")
		if r.URL.Query().Get("page") != "1" {
			t.Fatalf("expected page=1, got %q", r.URL.Query().Get("page"))
		}
		if r.URL.Query().Get("per_page") != "20" {
			t.Fatalf("expected per_page=20, got %q", r.URL.Query().Get("per_page"))
		}
		writeTodoJSON(t, w, map[string]interface{}{"total_count": 0, "applied_projects": []interface{}{}})
	}))
	defer server.Close()

	if err := runTodoShortcut(t, server, "join-list", map[string]string{"login": "bob", "page": "1", "per-page": "20"}); err != nil {
		t.Fatalf("join-list failed: %v", err)
	}
}

func TestJoinAccept(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertTodoRequest(t, r, "POST", "/users/bob/applied_projects/11/accept.json")
		writeTodoJSON(t, w, map[string]interface{}{"id": 11, "status": "accepted"})
	}))
	defer server.Close()

	if err := runTodoShortcut(t, server, "join-accept", map[string]string{"login": "bob", "id": "11"}); err != nil {
		t.Fatalf("join-accept failed: %v", err)
	}
}

func TestJoinRefuse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertTodoRequest(t, r, "POST", "/users/bob/applied_projects/12/refuse.json")
		writeTodoJSON(t, w, map[string]interface{}{"id": 12, "status": "refuse"})
	}))
	defer server.Close()

	if err := runTodoShortcut(t, server, "join-refuse", map[string]string{"login": "bob", "id": "12"}); err != nil {
		t.Fatalf("join-refuse failed: %v", err)
	}
}

func TestActionMissingID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no API call expected")
	}))
	defer server.Close()

	if err := runTodoShortcut(t, server, "join-accept", map[string]string{"login": "bob"}); err == nil {
		t.Fatal("expected missing id error")
	}
}

func TestActionInvalidID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no API call expected")
	}))
	defer server.Close()

	if err := runTodoShortcut(t, server, "join-accept", map[string]string{"login": "bob", "id": "0"}); err == nil {
		t.Fatal("expected invalid id error")
	}
}
