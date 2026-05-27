package code

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gitlink-org/gitlink-cli/internal/client"
	"github.com/gitlink-org/gitlink-cli/internal/i18n"
	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

func TestCodeShortcutsRouteToOpenAPIEndpoints(t *testing.T) {
	cases := []struct {
		name   string
		args   map[string]string
		method string
		path   string
		query  map[string]string
	}{
		{
			name:   "files",
			args:   map[string]string{"search": "README", "ref": "main"},
			method: "GET",
			path:   "/owner/repo/files.json",
			query:  map[string]string{"search": "README", "ref": "main"},
		},
		{
			name:   "entries",
			args:   map[string]string{"ref": "main"},
			method: "GET",
			path:   "/owner/repo/entries.json",
			query:  map[string]string{"ref": "main"},
		},
		{
			name:   "sub-entries",
			args:   map[string]string{"path": "docs/README.md", "ref": "main"},
			method: "GET",
			path:   "/owner/repo/sub_entries.json",
			query:  map[string]string{"filepath": "docs/README.md", "ref": "main"},
		},
		{
			name:   "tree",
			args:   map[string]string{"sha": "main", "recursive": "true", "page": "2", "limit": "50"},
			method: "GET",
			path:   "/v1/owner/repo/git/trees/main.json",
			query:  map[string]string{"recursive": "true", "page": "2", "limit": "50"},
		},
		{
			name:   "blob",
			args:   map[string]string{"sha": "abc123"},
			method: "GET",
			path:   "/v1/owner/repo/git/blobs/abc123.json",
		},
		{
			name:   "commits",
			args:   map[string]string{"sha": "main", "page": "1", "limit": "20"},
			method: "GET",
			path:   "/v1/owner/repo/commits.json",
			query:  map[string]string{"sha": "main", "page": "1", "limit": "20"},
		},
		{
			name:   "commit-files",
			args:   map[string]string{"sha": "abc123", "file": "README.md", "page": "1", "limit": "20"},
			method: "GET",
			path:   "/v1/owner/repo/commits/abc123/files.json",
			query:  map[string]string{"filepath": "README.md", "page": "1", "limit": "20"},
		},
		{
			name:   "commit-diff",
			args:   map[string]string{"sha": "abc123"},
			method: "GET",
			path:   "/v1/owner/repo/commits/abc123/diff.json",
		},
		{
			name:   "blame",
			args:   map[string]string{"sha": "main", "path": "README.md"},
			method: "GET",
			path:   "/v1/owner/repo/blame.json",
			query:  map[string]string{"sha": "main", "filepath": "README.md"},
		},
		{
			name:   "tags",
			args:   map[string]string{"page": "1", "limit": "20"},
			method: "GET",
			path:   "/v1/owner/repo/tags.json",
			query:  map[string]string{"page": "1", "limit": "20"},
		},
		{
			name:   "tag",
			args:   map[string]string{"name": "v1.0"},
			method: "GET",
			path:   "/v1/owner/repo/tags/v1.0.json",
		},
		{
			name:   "delete-tag",
			args:   map[string]string{"name": "v1.0"},
			method: "DELETE",
			path:   "/v1/owner/repo/tags/v1.0.json",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			server := newCodeTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tc.method || r.URL.Path != tc.path {
					t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
				}
				assertCodeQuery(t, r, tc.query)
				called = true
				writeCodeJSON(t, w, map[string]interface{}{"ok": true})
			})
			defer server.Close()

			err := runCodeShortcut(server, tc.name, tc.args)
			if err != nil {
				t.Fatalf("%s shortcut failed: %v", tc.name, err)
			}
			if !called {
				t.Fatal("expected API request")
			}
		})
	}
}

func TestCodeDeleteTagDryRunDoesNotCallAPI(t *testing.T) {
	server := newCodeTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
	})
	defer server.Close()

	err := runCodeShortcut(server, "delete-tag", map[string]string{
		"name":    "v1.0",
		"dry-run": "true",
	})
	if err != nil {
		t.Fatalf("delete-tag dry-run failed: %v", err)
	}
}

func TestCodeShortcutsValidateRequiredArgs(t *testing.T) {
	server := newCodeTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
	})
	defer server.Close()

	cases := []struct {
		name string
		args map[string]string
		want string
	}{
		{name: "sub-entries", args: map[string]string{}, want: "--path"},
		{name: "tree", args: map[string]string{}, want: "--sha"},
		{name: "blob", args: map[string]string{}, want: "--sha"},
		{name: "commit-files", args: map[string]string{}, want: "--sha"},
		{name: "commit-diff", args: map[string]string{}, want: "--sha"},
		{name: "blame", args: map[string]string{"sha": "main"}, want: "--path"},
		{name: "tag", args: map[string]string{}, want: "--name"},
		{name: "delete-tag", args: map[string]string{}, want: "--name"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := runCodeShortcut(server, tc.name, tc.args)
			if err == nil {
				t.Fatal("expected validation error")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %q, want it to mention %s", err.Error(), tc.want)
			}
		})
	}
}

func runCodeShortcut(server *httptest.Server, name string, args map[string]string) error {
	for _, shortcut := range Shortcuts() {
		if shortcut.Name != name {
			continue
		}
		ctx := &common.RuntimeContext{
			Client: &client.Client{
				HTTP:    server.Client(),
				BaseURL: server.URL,
			},
			Owner:  "owner",
			Repo:   "repo",
			Format: "json",
			Args:   args,
			Tr:     i18n.Default(),
		}
		return shortcut.Run(ctx)
	}
	return fmt.Errorf("shortcut %q not found", name)
}

func newCodeTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

func assertCodeQuery(t *testing.T, r *http.Request, want map[string]string) {
	t.Helper()
	query := r.URL.Query()
	if len(query) != len(want) {
		t.Fatalf("query = %v, want %v", query, want)
	}
	for key, value := range want {
		if got := query.Get(key); got != value {
			t.Fatalf("query %s = %q, want %q", key, got, value)
		}
	}
}

func writeCodeJSON(t *testing.T, w http.ResponseWriter, payload interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Fatalf("failed to write response: %v", err)
	}
}
