package webhook

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gitlink-org/gitlink-cli/internal/client"
	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

func TestWebhookList(t *testing.T) {
	server := newWebhookTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertRequest(t, r, "GET", "/v1/owner/repo/webhooks.json")
		writeJSON(t, w, map[string]interface{}{"total_count": 1, "webhooks": []interface{}{}})
	})
	defer server.Close()

	if err := runWebhookShortcut(t, server, "list", nil); err != nil {
		t.Fatalf("list shortcut failed: %v", err)
	}
}

func TestWebhookView(t *testing.T) {
	server := newWebhookTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertRequest(t, r, "GET", "/v1/owner/repo/webhooks/7.json")
		writeJSON(t, w, map[string]interface{}{"id": 7, "url": "https://example.com/hook"})
	})
	defer server.Close()

	if err := runWebhookShortcut(t, server, "view", map[string]string{"id": "7"}); err != nil {
		t.Fatalf("view shortcut failed: %v", err)
	}
}

func TestWebhookCreatePayload(t *testing.T) {
	var payload map[string]interface{}
	server := newWebhookTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertRequest(t, r, "POST", "/v1/owner/repo/webhooks.json")
		payload = decodeJSON(t, r)
		writeJSON(t, w, map[string]interface{}{"id": 1})
	})
	defer server.Close()

	err := runWebhookShortcut(t, server, "create", map[string]string{
		"url":           "https://example.com/hook",
		"events":        "push,issues_only,push",
		"type":          "gitea",
		"content-type":  "json",
		"http-method":   "POST",
		"secret":        "secret-token",
		"branch-filter": "master,{release*}",
		"active":        "true",
	})
	if err != nil {
		t.Fatalf("create shortcut failed: %v", err)
	}

	assertEqual(t, payload["url"], "https://example.com/hook")
	assertEqual(t, payload["type"], "gitea")
	assertEqual(t, payload["content_type"], "json")
	assertEqual(t, payload["http_method"], "POST")
	assertEqual(t, payload["secret"], "secret-token")
	assertEqual(t, payload["branch_filter"], "master,{release*}")
	assertEqual(t, payload["active"], true)
	assertStringSlice(t, payload["events"], []string{"push", "issues_only"})
}

func TestWebhookUpdatePreservesCurrentFields(t *testing.T) {
	var payload map[string]interface{}
	server := newWebhookTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/v1/owner/repo/webhooks/7.json":
			writeJSON(t, w, map[string]interface{}{
				"id":            7,
				"url":           "https://old.example.com/hook",
				"type":          "gitea",
				"content_type":  "json",
				"http_method":   "POST",
				"branch_filter": "*",
				"events":        []string{"push"},
				"active":        true,
			})
		case r.Method == "PUT" && r.URL.Path == "/v1/owner/repo/webhooks/7.json":
			payload = decodeJSON(t, r)
			writeJSON(t, w, payload)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})
	defer server.Close()

	err := runWebhookShortcut(t, server, "update", map[string]string{
		"id":     "7",
		"url":    "https://new.example.com/hook",
		"events": "push,issue_comment",
	})
	if err != nil {
		t.Fatalf("update shortcut failed: %v", err)
	}

	assertEqual(t, payload["url"], "https://new.example.com/hook")
	assertEqual(t, payload["type"], "gitea")
	assertEqual(t, payload["content_type"], "json")
	assertEqual(t, payload["http_method"], "POST")
	assertEqual(t, payload["branch_filter"], "*")
	assertEqual(t, payload["active"], true)
	assertStringSlice(t, payload["events"], []string{"push", "issue_comment"})
}

func TestWebhookDelete(t *testing.T) {
	server := newWebhookTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertRequest(t, r, "DELETE", "/v1/owner/repo/webhooks/7.json")
		writeJSON(t, w, map[string]interface{}{"status": 0, "message": "success"})
	})
	defer server.Close()

	if err := runWebhookShortcut(t, server, "delete", map[string]string{"id": "7"}); err != nil {
		t.Fatalf("delete shortcut failed: %v", err)
	}
}

func TestWebhookTestDelivery(t *testing.T) {
	server := newWebhookTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertRequest(t, r, "POST", "/v1/owner/repo/webhooks/7/tests.json")
		writeJSON(t, w, map[string]interface{}{"status": 0, "message": "success"})
	})
	defer server.Close()

	if err := runWebhookShortcut(t, server, "test", map[string]string{"id": "7"}); err != nil {
		t.Fatalf("test shortcut failed: %v", err)
	}
}

func TestWebhookTasks(t *testing.T) {
	server := newWebhookTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertRequest(t, r, "GET", "/v1/owner/repo/webhooks/7/hooktasks.json")
		writeJSON(t, w, map[string]interface{}{"total_count": 0, "hooktasks": []interface{}{}})
	})
	defer server.Close()

	if err := runWebhookShortcut(t, server, "tasks", map[string]string{"id": "7"}); err != nil {
		t.Fatalf("tasks shortcut failed: %v", err)
	}
}

func TestParseWebhookEventsRejectsInvalidEvent(t *testing.T) {
	_, err := parseWebhookEvents("push,invalid")
	if err == nil {
		t.Fatal("expected invalid event to return an error")
	}
}

func TestWebhookCreateRejectsInvalidActive(t *testing.T) {
	server := newWebhookTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("invalid active should not call API, got: %s %s", r.Method, r.URL.Path)
	})
	defer server.Close()

	err := runWebhookShortcut(t, server, "create", map[string]string{
		"url":    "https://example.com/hook",
		"events": "push",
		"active": "maybe",
	})
	if err == nil {
		t.Fatal("expected invalid active to return an error")
	}
}

func runWebhookShortcut(t *testing.T, server *httptest.Server, name string, args map[string]string) error {
	t.Helper()
	shortcut := findWebhookShortcut(t, name)
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

func findWebhookShortcut(t *testing.T, name string) *common.Shortcut {
	t.Helper()
	for _, shortcut := range Shortcuts() {
		if shortcut.Name == name {
			return shortcut
		}
	}
	t.Fatalf("shortcut %q not found", name)
	return nil
}

func newWebhookTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

func assertRequest(t *testing.T, r *http.Request, method, path string) {
	t.Helper()
	if r.Method != method || r.URL.Path != path {
		t.Fatalf("got request %s %s, want %s %s", r.Method, r.URL.Path, method, path)
	}
}

func decodeJSON(t *testing.T, r *http.Request) map[string]interface{} {
	t.Helper()
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode request body: %v", err)
	}
	return payload
}

func writeJSON(t *testing.T, w http.ResponseWriter, payload interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Fatalf("failed to write response: %v", err)
	}
}

func assertEqual(t *testing.T, got interface{}, want interface{}) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v (%T), want %v (%T)", got, got, want, want)
	}
}

func assertStringSlice(t *testing.T, got interface{}, want []string) {
	t.Helper()
	values, ok := got.([]interface{})
	if !ok {
		t.Fatalf("got %T, want []interface{}", got)
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		text, ok := value.(string)
		if !ok {
			t.Fatalf("got event %v (%T), want string", value, value)
		}
		result = append(result, text)
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("got %v, want %v", result, want)
	}
}
