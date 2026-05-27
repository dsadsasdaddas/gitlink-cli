package dataset

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

func TestDatasetShortcutsRouteToOpenAPIEndpoints(t *testing.T) {
	cases := []struct {
		name    string
		args    map[string]string
		method  string
		path    string
		query   map[string]string
		payload map[string]interface{}
	}{
		{
			name:   "view",
			args:   map[string]string{"page": "2", "limit": "10"},
			method: "GET",
			path:   "/v1/owner/repo/dataset.json",
			query:  map[string]string{"page": "2", "limit": "10"},
		},
		{
			name:   "list",
			args:   map[string]string{"ids": "1,2"},
			method: "GET",
			path:   "/v1/project_datasets.json",
			query:  map[string]string{"ids": "1,2"},
		},
		{
			name:   "create",
			args:   datasetWriteArgs(),
			method: "POST",
			path:   "/v1/owner/repo/dataset.json",
			payload: map[string]interface{}{
				"title":         "Dataset title",
				"description":   "Dataset description",
				"license_id":    "3",
				"paper_content": "Paper content",
			},
		},
		{
			name:   "update",
			args:   datasetWriteArgs(),
			method: "PUT",
			path:   "/v1/owner/repo/dataset.json",
			payload: map[string]interface{}{
				"title":         "Dataset title",
				"description":   "Dataset description",
				"license_id":    "3",
				"paper_content": "Paper content",
			},
		},
		{
			name:   "delete-attachment",
			args:   map[string]string{"uuid": "abc-123"},
			method: "DELETE",
			path:   "/attachments/abc-123.json",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			server := newDatasetTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tc.method || r.URL.Path != tc.path {
					t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
				}
				assertDatasetQuery(t, r, tc.query)
				if tc.payload != nil {
					assertDatasetPayload(t, r, tc.payload)
				}
				called = true
				writeDatasetJSON(t, w, map[string]interface{}{"status": 0, "message": "success"})
			})
			defer server.Close()

			err := runDatasetShortcut(server, tc.name, tc.args)
			if err != nil {
				t.Fatalf("%s shortcut failed: %v", tc.name, err)
			}
			if !called {
				t.Fatal("expected API request")
			}
		})
	}
}

func TestDatasetDeleteAttachmentDryRunDoesNotCallAPI(t *testing.T) {
	server := newDatasetTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
	})
	defer server.Close()

	err := runDatasetShortcut(server, "delete-attachment", map[string]string{
		"uuid":    "abc-123",
		"dry-run": "true",
	})
	if err != nil {
		t.Fatalf("delete-attachment dry-run failed: %v", err)
	}
}

func TestDatasetShortcutsValidateRequiredArgs(t *testing.T) {
	server := newDatasetTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
	})
	defer server.Close()

	cases := []struct {
		name string
		args map[string]string
		want string
	}{
		{name: "create", args: map[string]string{"description": "desc"}, want: "--title"},
		{name: "create", args: map[string]string{"title": "title"}, want: "--description"},
		{name: "update", args: map[string]string{"description": "desc"}, want: "--title"},
		{name: "delete-attachment", args: map[string]string{}, want: "--uuid"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := runDatasetShortcut(server, tc.name, tc.args)
			if err == nil {
				t.Fatal("expected validation error")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %q, want it to mention %s", err.Error(), tc.want)
			}
		})
	}
}

func datasetWriteArgs() map[string]string {
	return map[string]string{
		"title":         "Dataset title",
		"description":   "Dataset description",
		"license-id":    "3",
		"paper-content": "Paper content",
	}
}

func runDatasetShortcut(server *httptest.Server, name string, args map[string]string) error {
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

func newDatasetTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

func assertDatasetQuery(t *testing.T, r *http.Request, want map[string]string) {
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

func assertDatasetPayload(t *testing.T, r *http.Request, want map[string]interface{}) {
	t.Helper()
	var got map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
		t.Fatalf("decode request body: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("payload = %v, want %v", got, want)
	}
	for key, value := range want {
		if got[key] != value {
			t.Fatalf("payload %s = %v, want %v", key, got[key], value)
		}
	}
}

func writeDatasetJSON(t *testing.T, w http.ResponseWriter, payload interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Fatalf("failed to write response: %v", err)
	}
}
