package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()
	handler(newStore()).ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("health status: %d", response.Code)
	}
}

func TestCreateAndList(t *testing.T) {
	app := handler(newStore())
	body := `{"service":"api","environment":"staging","version":"v1.2.0"}`
	req := httptest.NewRequest(http.MethodPost, "/api/releases", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	app.ServeHTTP(response, req)
	if response.Code != http.StatusCreated {
		t.Fatalf("create: %d: %s", response.Code, response.Body.String())
	}
	var item release
	if err := json.Unmarshal(response.Body.Bytes(), &item); err != nil {
		t.Fatal(err)
	}
	if item.ID != 1 || item.Status != "queued" {
		t.Fatalf("unexpected item: %+v", item)
	}
	list := httptest.NewRecorder()
	app.ServeHTTP(list, httptest.NewRequest(http.MethodGet, "/api/releases", nil))
	var payload struct {
		Releases []release `json:"releases"`
	}
	if err := json.Unmarshal(list.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.Releases) != 1 {
		t.Fatalf("want one release, got %d", len(payload.Releases))
	}
}

func TestRejectInvalidRequests(t *testing.T) {
	cases := []struct {
		body        string
		contentType string
		expected    int
	}{
		{`{"service":"api","environment":"invalid","version":"v1"}`, "application/json", http.StatusBadRequest},
		{`{"service":"","environment":"staging","version":"v1"}`, "application/json", http.StatusBadRequest},
		{`{"service":"api","environment":"staging","version":"v1","extra":true}`, "application/json", http.StatusBadRequest},
		{`{} {}`, "application/json", http.StatusBadRequest},
		{`{"service":"api","environment":"staging","version":"v1"}`, "text/plain", http.StatusUnsupportedMediaType},
	}
	app := handler(newStore())
	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodPost, "/api/releases", strings.NewReader(tc.body))
		req.Header.Set("Content-Type", tc.contentType)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)
		if rec.Code != tc.expected {
			t.Errorf("body %q: got %d want %d", tc.body, rec.Code, tc.expected)
		}
	}
}

func TestStoreConcurrentAdds(t *testing.T) {
	s := newStore()
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); s.add(createRelease{Service: "api", Environment: "staging", Version: "v1"}) }()
	}
	wg.Wait()
	if got := len(s.list()); got != 16 {
		t.Fatalf("got %d releases, want 16", got)
	}
}
