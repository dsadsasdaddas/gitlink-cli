package star

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

func TestStarListUsesOpenAPIEndpoint(t *testing.T) {
	server := newStarTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertStarRequest(t, r, "GET", "/users/alice/is_pinned_projects.json")
		writeStarJSON(t, w, map[string]interface{}{
			"total_count": 1,
			"projects": []interface{}{
				map[string]interface{}{"id": 9, "project_id": 17, "identifier": "demo"},
			},
		})
	})
	defer server.Close()

	err := runStarShortcut(t, server, "list", map[string]string{"login": "alice"})
	if err != nil {
		t.Fatalf("star list failed: %v", err)
	}
}

func TestStarSetPostsPinnedProjectIDs(t *testing.T) {
	var payload map[string]interface{}
	server := newStarTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertStarRequest(t, r, "POST", "/users/alice/is_pinned_projects/pin.json")
		payload = decodeStarJSON(t, r)
		writeStarJSON(t, w, map[string]interface{}{"status": 0, "message": "success"})
	})
	defer server.Close()

	err := runStarShortcut(t, server, "set", map[string]string{
		"login":       "alice",
		"project-ids": "17, 42,17",
	})
	if err != nil {
		t.Fatalf("star set failed: %v", err)
	}
	assertNumberSlice(t, payload["is_pinned_project_ids"], []float64{17, 42})
}

func TestStarSetDryRunDoesNotCallAPI(t *testing.T) {
	server := newStarTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("dry-run should not call API, got: %s %s", r.Method, r.URL.Path)
	})
	defer server.Close()

	err := runStarShortcut(t, server, "set", map[string]string{
		"login":       "alice",
		"project-ids": "17,42",
		"dry-run":     "true",
	})
	if err != nil {
		t.Fatalf("star set dry-run failed: %v", err)
	}
}

func TestStarReorderPutsPositionPayload(t *testing.T) {
	var payload map[string]interface{}
	server := newStarTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertStarRequest(t, r, "PUT", "/users/alice/is_pinned_projects/9.json")
		payload = decodeStarJSON(t, r)
		writeStarJSON(t, w, map[string]interface{}{"status": 0, "message": "success"})
	})
	defer server.Close()

	err := runStarShortcut(t, server, "reorder", map[string]string{
		"login":     "alice",
		"pinned-id": "9",
		"position":  "10",
	})
	if err != nil {
		t.Fatalf("star reorder failed: %v", err)
	}
	pinned, ok := payload["pinned_project"].(map[string]interface{})
	if !ok {
		t.Fatalf("pinned_project = %T, want object", payload["pinned_project"])
	}
	if pinned["position"] != float64(10) {
		t.Fatalf("position = %v, want 10", pinned["position"])
	}
}

func TestStarReorderDryRunDoesNotCallAPI(t *testing.T) {
	server := newStarTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("dry-run should not call API, got: %s %s", r.Method, r.URL.Path)
	})
	defer server.Close()

	err := runStarShortcut(t, server, "reorder", map[string]string{
		"login":     "alice",
		"pinned-id": "9",
		"position":  "0",
		"dry-run":   "true",
	})
	if err != nil {
		t.Fatalf("star reorder dry-run failed: %v", err)
	}
}

func TestStarShortcutsValidateRequiredArgs(t *testing.T) {
	server := newStarTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected API call: %s %s", r.Method, r.URL.Path)
	})
	defer server.Close()

	cases := []struct {
		name string
		args map[string]string
		want string
	}{
		{name: "list", args: map[string]string{}, want: "--login"},
		{name: "set", args: map[string]string{"login": "alice"}, want: "--project-ids"},
		{name: "reorder", args: map[string]string{"login": "alice", "pinned-id": "9"}, want: "--position"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := runStarShortcut(t, server, tc.name, tc.args)
			if err == nil {
				t.Fatal("expected validation error")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %q, want it to mention %s", err.Error(), tc.want)
			}
		})
	}
}

func TestStarRejectsInvalidIDs(t *testing.T) {
	server := newStarTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected API call: %s %s", r.Method, r.URL.Path)
	})
	defer server.Close()

	cases := []struct {
		name string
		args map[string]string
		want string
	}{
		{name: "set", args: map[string]string{"login": "alice", "project-ids": "0"}, want: "positive integer"},
		{name: "set", args: map[string]string{"login": "alice", "project-ids": ",,"}, want: "at least one"},
		{name: "reorder", args: map[string]string{"login": "alice", "pinned-id": "9", "position": "-1"}, want: "non-negative"},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s/%s", tc.name, tc.want), func(t *testing.T) {
			err := runStarShortcut(t, server, tc.name, tc.args)
			if err == nil {
				t.Fatal("expected validation error")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %q, want it to mention %s", err.Error(), tc.want)
			}
		})
	}
}

func runStarShortcut(t *testing.T, server *httptest.Server, name string, args map[string]string) error {
	t.Helper()
	if args == nil {
		args = map[string]string{}
	}
	for _, shortcut := range Shortcuts() {
		if shortcut.Name != name {
			continue
		}
		ctx := &common.RuntimeContext{
			Client: &client.Client{
				HTTP:    server.Client(),
				BaseURL: server.URL,
			},
			Format: "json",
			Args:   args,
			Tr:     i18n.Default(),
		}
		return shortcut.Run(ctx)
	}
	return fmt.Errorf("shortcut %q not found", name)
}

func newStarTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

func assertStarRequest(t *testing.T, r *http.Request, method, path string) {
	t.Helper()
	if r.Method != method || r.URL.Path != path {
		t.Fatalf("got request %s %s, want %s %s", r.Method, r.URL.Path, method, path)
	}
}

func writeStarJSON(t *testing.T, w http.ResponseWriter, payload interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Fatalf("failed to write response: %v", err)
	}
}

func decodeStarJSON(t *testing.T, r *http.Request) map[string]interface{} {
	t.Helper()
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode request body: %v", err)
	}
	return payload
}

func assertNumberSlice(t *testing.T, got interface{}, want []float64) {
	t.Helper()
	values, ok := got.([]interface{})
	if !ok {
		t.Fatalf("got %T, want []interface{}", got)
	}
	if len(values) != len(want) {
		t.Fatalf("got %v, want %v", values, want)
	}
	for i, value := range values {
		if value != want[i] {
			t.Fatalf("got %v, want %v", values, want)
		}
	}
}
