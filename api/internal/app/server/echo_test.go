package server_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api/internal/app/server"
	"api/internal/app/tasks"
	"api/internal/app/tasks/domain"
	"api/internal/app/tasks/repository/memory"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func newServer() *echo.Echo {
	log := logrus.New()
	log.SetOutput(io.Discard)
	e := server.New(log)
	tasks.Init(e.Group("/api"), memory.New())
	return e
}
func request(e *echo.Echo, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}
func TestTaskLifecycle(t *testing.T) {
	e := newServer()
	if r := request(e, "GET", "/hc", ""); r.Code != 200 || r.Body.String() != "OK" {
		t.Fatalf("health: %d %s", r.Code, r.Body)
	}
	if r := request(e, "GET", "/api/tasks", ""); r.Code != 200 || strings.TrimSpace(r.Body.String()) != "[]" {
		t.Fatalf("empty list: %d %s", r.Code, r.Body)
	}
	created := request(e, "POST", "/api/tasks", `{"title":"  example  "}`)
	if created.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", created.Code, created.Body)
	}
	var task domain.Task
	if err := json.Unmarshal(created.Body.Bytes(), &task); err != nil {
		t.Fatal(err)
	}
	if task.Title != "example" || task.ID == "" || task.CreatedAt.IsZero() {
		t.Fatalf("unexpected task: %+v", task)
	}
	if created.Header().Get("Location") != "/api/tasks/"+task.ID {
		t.Fatal("missing location")
	}
	got := request(e, "GET", "/api/tasks/"+task.ID, "")
	if got.Code != 200 || got.Body.String() != created.Body.String() {
		t.Fatalf("get: %d %s", got.Code, got.Body)
	}
	list := request(e, "GET", "/api/tasks", "")
	var all []domain.Task
	if err := json.Unmarshal(list.Body.Bytes(), &all); err != nil || list.Code != 200 || len(all) != 1 || all[0].ID != task.ID {
		t.Fatalf("list: %d %s %v", list.Code, list.Body, err)
	}
}
func TestHTTPFailures(t *testing.T) {
	e := newServer()
	for _, tt := range []struct {
		name, method, path, body string
		status                   int
	}{
		{"malformed", "POST", "/api/tasks", `{"title":`, 400},
		{"unknown field", "POST", "/api/tasks", `{"title":"ok","id":"injected"}`, 400},
		{"extra object", "POST", "/api/tasks", `{"title":"ok"} {}`, 400},
		{"invalid title", "POST", "/api/tasks", `{"title":" "}`, 422},
		{"missing task", "GET", "/api/tasks/missing", "", 404},
		{"unknown route", "GET", "/missing", "", 404},
		{"wrong method", "DELETE", "/api/tasks", "", 405},
		{"oversized", "POST", "/api/tasks", strings.Repeat("x", (1<<20)+1), 413},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := request(e, tt.method, tt.path, tt.body)
			if r.Code != tt.status {
				t.Fatalf("want %d got %d: %s", tt.status, r.Code, r.Body)
			}
			var body map[string]string
			if err := json.Unmarshal(r.Body.Bytes(), &body); err != nil || body["error"] == "" {
				t.Fatalf("bad error response %s", r.Body)
			}
		})
	}
}
func TestInternalErrorIsHidden(t *testing.T) {
	e := newServer()
	e.GET("/failure", func(c echo.Context) error { return echo.NewHTTPError(500, "database password") })
	r := request(e, "GET", "/failure", "")
	if r.Code != 500 || strings.Contains(r.Body.String(), "password") {
		t.Fatalf("leaked error: %s", r.Body)
	}
}
