package org

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gitlink-org/gitlink-cli/internal/client"
	"github.com/gitlink-org/gitlink-cli/internal/i18n"
	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

func TestOrgList(t *testing.T) {
	server := newOrgTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertOrgRequest(t, r, "GET", "/organizations.json")
		if r.URL.Query().Get("page") != "1" || r.URL.Query().Get("limit") != "20" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		writeOrgJSON(t, w, []interface{}{
			map[string]interface{}{"login": "org1"},
			map[string]interface{}{"login": "org2"},
		})
	})
	defer server.Close()

	err := runOrgShortcut(t, server, "list", map[string]string{"page": "1", "limit": "20"})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
}

func TestOrgInfo(t *testing.T) {
	server := newOrgTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertOrgRequest(t, r, "GET", "/organizations/myorg.json")
		writeOrgJSON(t, w, map[string]interface{}{"login": "myorg", "name": "My Org"})
	})
	defer server.Close()

	err := runOrgShortcut(t, server, "info", map[string]string{"id": "myorg"})
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}
}

func TestOrgMembers(t *testing.T) {
	server := newOrgTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertOrgRequest(t, r, "GET", "/organizations/myorg/organization_users.json")
		if r.URL.Query().Get("page") != "1" || r.URL.Query().Get("limit") != "20" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		writeOrgJSON(t, w, []interface{}{
			map[string]interface{}{"login": "user1"},
		})
	})
	defer server.Close()

	err := runOrgShortcut(t, server, "members", map[string]string{"id": "myorg", "page": "1", "limit": "20"})
	if err != nil {
		t.Fatalf("members failed: %v", err)
	}
}

func TestOrgCreate(t *testing.T) {
	server := newOrgTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertOrgRequest(t, r, "POST", "/organizations.json")
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if payload["name"] != "neworg" || payload["description"] != "A new org" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		writeOrgJSON(t, w, map[string]interface{}{"login": "neworg"})
	})
	defer server.Close()

	err := runOrgShortcut(t, server, "create", map[string]string{"name": "neworg", "description": "A new org"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
}

func TestOrgCreateNoDescription(t *testing.T) {
	server := newOrgTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertOrgRequest(t, r, "POST", "/organizations.json")
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if payload["name"] != "neworg" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		if _, ok := payload["description"]; ok {
			t.Fatalf("description should be omitted when unset: %#v", payload)
		}
		writeOrgJSON(t, w, map[string]interface{}{"login": "neworg"})
	})
	defer server.Close()

	err := runOrgShortcut(t, server, "create", map[string]string{"name": "neworg"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
}

func TestOrgTeamProjectsAddAll(t *testing.T) {
	server := newOrgTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertOrgRequest(t, r, "POST", "/organizations/gitlink/teams/7/team_projects/create_all.json")
		writeOrgJSON(t, w, map[string]interface{}{"status": 0, "message": "success"})
	})
	defer server.Close()

	err := runOrgShortcut(t, server, "team-projects-add-all", map[string]string{
		"organization": "gitlink",
		"team-id":      "7",
	})
	if err != nil {
		t.Fatalf("team-projects-add-all failed: %v", err)
	}
}

func TestOrgTeamProjectsRemoveAll(t *testing.T) {
	server := newOrgTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertOrgRequest(t, r, "DELETE", "/organizations/gitlink/teams/7/team_projects/destroy_all.json")
		writeOrgJSON(t, w, map[string]interface{}{"status": 0, "message": "success"})
	})
	defer server.Close()

	err := runOrgShortcut(t, server, "team-projects-remove-all", map[string]string{
		"organization": "gitlink",
		"team-id":      "7",
	})
	if err != nil {
		t.Fatalf("team-projects-remove-all failed: %v", err)
	}
}

func TestOrgTeamProjectsAddAllDryRunDoesNotCallAPI(t *testing.T) {
	server := newOrgTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("dry-run should not call API, got: %s %s", r.Method, r.URL.Path)
	})
	defer server.Close()

	err := runOrgShortcut(t, server, "team-projects-add-all", map[string]string{
		"organization": "gitlink",
		"team-id":      "7",
		"dry-run":      "true",
	})
	if err != nil {
		t.Fatalf("team-projects-add-all dry-run failed: %v", err)
	}
}

func TestOrgTeamProjectsRemoveAllDryRunDoesNotCallAPI(t *testing.T) {
	server := newOrgTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("dry-run should not call API, got: %s %s", r.Method, r.URL.Path)
	})
	defer server.Close()

	err := runOrgShortcut(t, server, "team-projects-remove-all", map[string]string{
		"organization": "gitlink",
		"team-id":      "7",
		"dry-run":      "true",
	})
	if err != nil {
		t.Fatalf("team-projects-remove-all dry-run failed: %v", err)
	}
}

func TestOrgListHTTPError(t *testing.T) {
	server := newOrgTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeOrgText(w, http.StatusInternalServerError, "server error")
	})
	defer server.Close()

	err := runOrgShortcut(t, server, "list", map[string]string{"page": "1", "limit": "20"})
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}

func TestOrgInfoHTTPError(t *testing.T) {
	server := newOrgTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeOrgText(w, http.StatusInternalServerError, "server error")
	})
	defer server.Close()

	err := runOrgShortcut(t, server, "info", map[string]string{"id": "myorg"})
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}

func TestOrgMembersHTTPError(t *testing.T) {
	server := newOrgTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeOrgText(w, http.StatusInternalServerError, "server error")
	})
	defer server.Close()

	err := runOrgShortcut(t, server, "members", map[string]string{"id": "myorg", "page": "1", "limit": "20"})
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}

func TestOrgCreateHTTPError(t *testing.T) {
	server := newOrgTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeOrgText(w, http.StatusInternalServerError, "server error")
	})
	defer server.Close()

	err := runOrgShortcut(t, server, "create", map[string]string{"name": "neworg"})
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}

func runOrgShortcut(t *testing.T, server *httptest.Server, name string, args map[string]string) error {
	t.Helper()
	shortcut := findOrgShortcut(t, name)
	if args == nil {
		args = map[string]string{}
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

func findOrgShortcut(t *testing.T, name string) *common.Shortcut {
	t.Helper()
	for _, shortcut := range Shortcuts() {
		if shortcut.Name == name {
			return shortcut
		}
	}
	t.Fatalf("shortcut %q not found", name)
	return nil
}

func newOrgTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

func assertOrgRequest(t *testing.T, r *http.Request, method, path string) {
	t.Helper()
	if r.Method != method || r.URL.Path != path {
		t.Fatalf("got request %s %s, want %s %s", r.Method, r.URL.Path, method, path)
	}
}

func writeOrgJSON(t *testing.T, w http.ResponseWriter, payload interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Fatalf("failed to write response: %v", err)
	}
}

func writeOrgText(w http.ResponseWriter, code int, body string) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(body))
}
