package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPIMoo(t *testing.T) {
	m := NewModule()

	tests := []struct {
		name       string
		method     string
		body       string
		query      string
		wantStatus int
		wantError  bool
	}{
		{
			name:       "valid JSON request",
			method:     "POST",
			body:       `{"text":"test","cow":"default","action":"say"}`,
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name:       "valid query params",
			method:     "GET",
			query:      "?text=test&cow=default",
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name:       "missing text",
			method:     "POST",
			body:       `{"cow":"default"}`,
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name:       "invalid cow",
			method:     "POST",
			body:       `{"text":"test","cow":"invalid"}`,
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name:       "invalid mood",
			method:     "POST",
			body:       `{"text":"test","mood":"invalid"}`,
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name:       "random cow and mood",
			method:     "POST",
			body:       `{"text":"test","cow":"random","mood":"random"}`,
			wantStatus: http.StatusOK,
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.query != "" {
				req = httptest.NewRequest(tt.method, "/api/moo"+tt.query, nil)
			} else {
				req = httptest.NewRequest(tt.method, "/api/moo", strings.NewReader(tt.body))
				if tt.body != "" {
					req.Header.Set("Content-Type", "application/json")
				}
			}

			w := httptest.NewRecorder()
			m.APIMoo(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.wantError {
				var errResp ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
					t.Errorf("failed to decode error response: %v", err)
				}
				if errResp.Error == "" {
					t.Error("expected error message, got empty")
				}
			} else {
				var resp MooResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if resp.Output == "" {
					t.Error("expected output, got empty")
				}
			}
		})
	}
}

func TestAPICows(t *testing.T) {
	m := NewModule()
	req := httptest.NewRequest("GET", "/api/cows", nil)
	w := httptest.NewRecorder()

	m.APICows(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string][]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	cows, ok := resp["cows"]
	if !ok {
		t.Fatal("response missing 'cows' field")
	}

	if len(cows) < 40 {
		t.Errorf("expected at least 40 cows, got %d", len(cows))
	}
}

func TestAPIMoods(t *testing.T) {
	m := NewModule()
	req := httptest.NewRequest("GET", "/api/moods", nil)
	w := httptest.NewRecorder()

	m.APIMoods(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string][]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	moods, ok := resp["moods"]
	if !ok {
		t.Fatal("response missing 'moods' field")
	}

	if len(moods) < 5 {
		t.Errorf("expected at least 5 moods, got %d", len(moods))
	}
}

func TestHealth(t *testing.T) {
	handler := Health("test-version")
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp HealthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("status = %s, want ok", resp.Status)
	}

	if resp.Version != "test-version" {
		t.Errorf("version = %s, want test-version", resp.Version)
	}
}

func TestCORS(t *testing.T) {
	handler := CORS(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name   string
		method string
	}{
		{"OPTIONS request", "OPTIONS"},
		{"GET request", "GET"},
		{"POST request", "POST"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()

			handler(w, req)

			if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
				t.Errorf("CORS origin = %s, want *", got)
			}

			if tt.method == "OPTIONS" {
				if w.Code != http.StatusOK {
					t.Errorf("OPTIONS status = %d, want %d", w.Code, http.StatusOK)
				}
			}
		})
	}
}

// Edge case and concurrency tests

func TestAPIMoo_Concurrent(t *testing.T) {
	m := NewModule()

	// Test concurrent requests don't cause race conditions
	const numRequests = 50
	done := make(chan bool, numRequests)
	errors := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(id int) {
			body := `{"text":"concurrent test","cow":"default","action":"say"}`
			req := httptest.NewRequest("POST", "/api/moo", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			m.APIMoo(w, req)

			if w.Code != http.StatusOK {
				errors <- fmt.Errorf("request %d: status = %d", id, w.Code)
			}

			done <- true
		}(i)
	}

	// Wait for all requests to complete
	for i := 0; i < numRequests; i++ {
		<-done
	}

	close(errors)
	for err := range errors {
		t.Error(err)
	}
}

func TestAPIMoo_LargePayload(t *testing.T) {
	m := NewModule()

	// Test with large text input
	largeText := strings.Repeat("This is a very long line of text. ", 1000)
	body := fmt.Sprintf(`{"text":"%s","cow":"default"}`, largeText)

	req := httptest.NewRequest("POST", "/api/moo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m.APIMoo(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Large payload: status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestAPIMoo_Unicode(t *testing.T) {
	m := NewModule()

	tests := []struct {
		name string
		text string
	}{
		{"chinese", "你好世界"},
		{"emoji", "Hello 🐄🐮"},
		{"arabic", "مرحبا بك"},
		{"hebrew", "שלום"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := fmt.Sprintf(`{"text":"%s","cow":"default"}`, tt.text)
			req := httptest.NewRequest("POST", "/api/moo", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			m.APIMoo(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Unicode test %s: status = %d", tt.name, w.Code)
			}
		})
	}
}
