package mirror

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gitlink-org/gitlink-cli/internal/client"
	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

func runMirrorShortcut(t *testing.T, server *httptest.Server, name string, args map[string]string) error {
	t.Helper()
	shortcut := findMirrorShortcut(t, name)
	ctx := &common.RuntimeContext{
		Client: &client.Client{HTTP: server.Client(), BaseURL: server.URL},
		Format: "json",
		Args:   args,
	}
	return shortcut.Run(ctx)
}

func findMirrorShortcut(t *testing.T, name string) *common.Shortcut {
	t.Helper()
	for _, s := range Shortcuts() {
		if s.Name == name {
			return s
		}
	}
	t.Fatalf("shortcut %q not found", name)
	return nil
}

func writeMirrorJSON(t *testing.T, w http.ResponseWriter, v interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Fatalf("write json: %v", err)
	}
}

func assertMirrorRequest(t *testing.T, r *http.Request, method, path string) {
	t.Helper()
	if r.Method != method {
		t.Fatalf("expected method %s, got %s", method, r.Method)
	}
	if r.URL.Path != path {
		t.Fatalf("expected path %s, got %s", path, r.URL.Path)
	}
}

func TestCreateMirror(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMirrorRequest(t, r, "POST", "/projects/migrate.json")
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		assertNumber(t, payload, "user_id", 42)
		assertString(t, payload, "name", "Demo Mirror")
		assertString(t, payload, "repository_name", "demo-mirror")
		assertString(t, payload, "clone_addr", "https://example.com/demo.git")
		assertString(t, payload, "description", "mirror repo")
		assertNumber(t, payload, "project_category_id", 3)
		assertNumber(t, payload, "project_language_id", 5)
		assertBool(t, payload, "is_mirror", true)
		assertBool(t, payload, "private", true)
		assertString(t, payload, "auth_username", "bot")
		assertString(t, payload, "auth_password", "secret-token")
		writeMirrorJSON(t, w, map[string]interface{}{"id": 99, "identifier": "demo-mirror"})
	}))
	defer server.Close()

	err := runMirrorShortcut(t, server, "create", map[string]string{
		"user-id":         "42",
		"name":            "Demo Mirror",
		"repository-name": "demo-mirror",
		"clone-addr":      "https://example.com/demo.git",
		"description":     "mirror repo",
		"category-id":     "3",
		"language-id":     "5",
		"private":         "true",
		"auth-username":   "bot",
		"auth-password":   "secret-token",
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
}

func TestCreateMirrorDryRun(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no API call expected for dry-run")
	}))
	defer server.Close()

	err := runMirrorShortcut(t, server, "create", map[string]string{
		"user-id":         "42",
		"name":            "Demo Mirror",
		"repository-name": "demo-mirror",
		"clone-addr":      "https://example.com/demo.git",
		"auth-username":   "bot",
		"auth-password":   "secret-token",
		"dry-run":         "true",
	})
	if err != nil {
		t.Fatalf("create dry-run failed: %v", err)
	}
}

func TestCreateMirrorRequiresCredentialsPair(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no API call expected")
	}))
	defer server.Close()

	err := runMirrorShortcut(t, server, "create", map[string]string{
		"user-id":         "42",
		"name":            "Demo Mirror",
		"repository-name": "demo-mirror",
		"clone-addr":      "https://example.com/demo.git",
		"auth-password":   "secret-token",
	})
	if err == nil {
		t.Fatal("expected credential validation error")
	}
}

func TestCreateMirrorInvalidUserID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no API call expected")
	}))
	defer server.Close()

	err := runMirrorShortcut(t, server, "create", map[string]string{
		"user-id":         "0",
		"name":            "Demo Mirror",
		"repository-name": "demo-mirror",
		"clone-addr":      "https://example.com/demo.git",
	})
	if err == nil {
		t.Fatal("expected invalid user-id error")
	}
}

func TestSyncMirror(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMirrorRequest(t, r, "POST", "/repositories/99/sync_mirror.json")
		writeMirrorJSON(t, w, map[string]interface{}{"status": 0, "message": "success"})
	}))
	defer server.Close()

	if err := runMirrorShortcut(t, server, "sync", map[string]string{"repo-id": "99"}); err != nil {
		t.Fatalf("sync failed: %v", err)
	}
}

func TestSyncMirrorDryRun(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no API call expected for dry-run")
	}))
	defer server.Close()

	if err := runMirrorShortcut(t, server, "sync", map[string]string{"repo-id": "99", "dry-run": "true"}); err != nil {
		t.Fatalf("sync dry-run failed: %v", err)
	}
}

func TestSyncMirrorInvalidRepoID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no API call expected")
	}))
	defer server.Close()

	if err := runMirrorShortcut(t, server, "sync", map[string]string{"repo-id": "abc"}); err == nil {
		t.Fatal("expected invalid repo-id error")
	}
}

func assertString(t *testing.T, payload map[string]interface{}, key, want string) {
	t.Helper()
	if got, _ := payload[key].(string); got != want {
		t.Fatalf("expected %s=%q, got %q", key, want, got)
	}
}

func assertBool(t *testing.T, payload map[string]interface{}, key string, want bool) {
	t.Helper()
	if got, _ := payload[key].(bool); got != want {
		t.Fatalf("expected %s=%t, got %t", key, want, got)
	}
}

func assertNumber(t *testing.T, payload map[string]interface{}, key string, want int) {
	t.Helper()
	got, ok := payload[key].(float64)
	if !ok || int(got) != want {
		t.Fatalf("expected %s=%d, got %#v", key, want, payload[key])
	}
}
