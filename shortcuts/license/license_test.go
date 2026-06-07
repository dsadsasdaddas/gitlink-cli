package license

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gitlink-org/gitlink-cli/internal/client"
	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

func runShortcut(t *testing.T, server *httptest.Server, name string, args map[string]string) error {
	t.Helper()
	shortcut := findShortcut(t, name)
	ctx := &common.RuntimeContext{
		Client: &client.Client{HTTP: server.Client(), BaseURL: server.URL},
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

func findShortcut(t *testing.T, name string) *common.Shortcut {
	t.Helper()
	for _, s := range Shortcuts() {
		if s.Name == name {
			return s
		}
	}
	t.Fatalf("shortcut %q not found", name)
	return nil
}

func writeJSON(t *testing.T, w http.ResponseWriter, payload interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Fatalf("failed to write response: %v", err)
	}
}

func TestLicenseList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("got method %s, want GET", r.Method)
		}
		if r.URL.Path != "/licenses.json" {
			t.Fatalf("got path %s, want /licenses.json", r.URL.Path)
		}
		if got := r.URL.Query().Get("name"); got != "" {
			t.Fatalf("expected no name filter, got %q", got)
		}
		writeJSON(t, w, map[string]interface{}{
			"licenses": []interface{}{
				map[string]interface{}{"id": 5, "name": "AFL-1.1"},
				map[string]interface{}{"id": 6, "name": "AFL-1.2"},
				map[string]interface{}{"id": 7, "name": "AFL-2.0"},
			},
		})
	}))
	defer server.Close()

	if err := runShortcut(t, server, "list", map[string]string{}); err != nil {
		t.Fatalf("license list failed: %v", err)
	}
}

func TestLicenseListWithNameFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("got method %s, want GET", r.Method)
		}
		if r.URL.Path != "/licenses.json" {
			t.Fatalf("got path %s, want /licenses.json", r.URL.Path)
		}
		if got := r.URL.Query().Get("name"); got != "AFL" {
			t.Fatalf("expected name=AFL, got %q", got)
		}
		writeJSON(t, w, map[string]interface{}{
			"licenses": []interface{}{
				map[string]interface{}{"id": 5, "name": "AFL-1.1"},
				map[string]interface{}{"id": 6, "name": "AFL-1.2"},
			},
		})
	}))
	defer server.Close()

	if err := runShortcut(t, server, "list", map[string]string{"name": "AFL"}); err != nil {
		t.Fatalf("license list with name filter failed: %v", err)
	}
}

func TestLicenseListHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	defer server.Close()

	err := runShortcut(t, server, "list", map[string]string{})
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}

func TestLicenseListEmptyResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/licenses.json" {
			t.Fatalf("got path %s, want /licenses.json", r.URL.Path)
		}
		writeJSON(t, w, map[string]interface{}{
			"licenses": []interface{}{},
		})
	}))
	defer server.Close()

	if err := runShortcut(t, server, "list", map[string]string{"name": "NONEXISTENT"}); err != nil {
		t.Fatalf("license list with empty result failed: %v", err)
	}
}
