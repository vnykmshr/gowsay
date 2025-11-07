package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/vnykmshr/gowsay/api"
)

// TestServerIntegration tests the full HTTP server with all endpoints
func TestServerIntegration(t *testing.T) {
	// Set up test environment
	os.Setenv("GOWSAY_TOKEN", "test-token")
	os.Setenv("GOWSAY_COLUMNS", "40")
	os.Setenv("PORT", "9999")
	defer os.Unsetenv("GOWSAY_TOKEN")
	defer os.Unsetenv("GOWSAY_COLUMNS")
	defer os.Unsetenv("PORT")

	// Create API module and setup routes
	module := api.NewModule()
	mux := http.NewServeMux()

	// Register API routes
	mux.Handle("/api/moo", api.CORS(http.HandlerFunc(module.APIMoo)))
	mux.Handle("/api/cows", api.CORS(http.HandlerFunc(module.APICows)))
	mux.Handle("/api/moods", api.CORS(http.HandlerFunc(module.APIMoods)))
	mux.Handle("/health", api.Health("test"))
	mux.Handle("/say", http.HandlerFunc(module.Gowsay))
	mux.Handle("/", api.ServeWeb())

	// Create test server
	server := httptest.NewServer(mux)
	defer server.Close()

	client := server.Client()
	baseURL := server.URL

	t.Run("HealthCheck", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/health")
		if err != nil {
			t.Fatalf("Failed to call /health: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		if !strings.Contains(string(body), "version") {
			t.Errorf("Health response should contain version info")
		}
	})

	t.Run("CompleteAPIFlow", func(t *testing.T) {
		// Step 1: Get available cows
		resp, err := client.Get(baseURL + "/api/cows")
		if err != nil {
			t.Fatalf("Failed to get cows: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET /api/cows: expected 200, got %d", resp.StatusCode)
		}

		var cowsResult map[string][]string
		if err := json.NewDecoder(resp.Body).Decode(&cowsResult); err != nil {
			t.Fatalf("Failed to decode cows: %v", err)
		}
		resp.Body.Close()

		cows, ok := cowsResult["cows"]
		if !ok || len(cows) == 0 {
			t.Fatal("Expected non-empty cows list")
		}

		// Step 2: Get available moods
		resp, err = client.Get(baseURL + "/api/moods")
		if err != nil {
			t.Fatalf("Failed to get moods: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET /api/moods: expected 200, got %d", resp.StatusCode)
		}

		var moodsResult map[string][]string
		if err := json.NewDecoder(resp.Body).Decode(&moodsResult); err != nil {
			t.Fatalf("Failed to decode moods: %v", err)
		}
		resp.Body.Close()

		moods, ok := moodsResult["moods"]
		if !ok || len(moods) == 0 {
			t.Fatal("Expected non-empty moods list")
		}

		// Step 3: Test valid JSON POST request
		reqBody := map[string]string{
			"text":   "Hello Integration Test",
			"cow":    "default",
			"action": "say",
		}
		jsonData, _ := json.Marshal(reqBody)

		resp, err = client.Post(
			baseURL+"/api/moo",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			t.Fatalf("POST /api/moo (JSON): %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("POST /api/moo (JSON): expected 200, got %d. Body: %s", resp.StatusCode, body)
		}

		var mooResp api.MooResponse
		if err := json.NewDecoder(resp.Body).Decode(&mooResp); err != nil {
			t.Fatalf("Failed to decode moo response: %v", err)
		}
		resp.Body.Close()

		if mooResp.Output == "" || !strings.Contains(mooResp.Output, "Hello Integration Test") {
			t.Error("Moo output should contain input text")
		}

		// Step 4: Test valid GET request with query params
		resp, err = client.Get(baseURL + "/api/moo?text=Query+Test&cow=dragon&action=think")
		if err != nil {
			t.Fatalf("GET /api/moo (query): %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET /api/moo (query): expected 200, got %d", resp.StatusCode)
		}

		if err := json.NewDecoder(resp.Body).Decode(&mooResp); err != nil {
			t.Fatalf("Failed to decode query response: %v", err)
		}
		resp.Body.Close()

		if !strings.Contains(mooResp.Output, "Query Test") {
			t.Error("Query response should contain input text")
		}

		// Step 5: Test random cow selection
		reqBody = map[string]string{
			"text": "Random test",
			"cow":  "random",
		}
		jsonData, _ = json.Marshal(reqBody)

		resp, err = client.Post(
			baseURL+"/api/moo",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			t.Fatalf("POST /api/moo (random): %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("POST /api/moo (random): expected 200, got %d", resp.StatusCode)
		}
		resp.Body.Close()

		// Step 6: Test error handling - invalid cow
		reqBody = map[string]string{
			"text": "Test",
			"cow":  "nonexistent",
		}
		jsonData, _ = json.Marshal(reqBody)

		resp, err = client.Post(
			baseURL+"/api/moo",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			t.Fatalf("POST /api/moo (invalid cow): %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("POST /api/moo (invalid cow): expected 400, got %d", resp.StatusCode)
		}
		resp.Body.Close()

		// Step 7: Test error handling - missing text
		reqBody = map[string]string{
			"cow": "default",
		}
		jsonData, _ = json.Marshal(reqBody)

		resp, err = client.Post(
			baseURL+"/api/moo",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			t.Fatalf("POST /api/moo (missing text): %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("POST /api/moo (missing text): expected 400, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	})

	t.Run("WebUI", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/")
		if err != nil {
			t.Fatalf("Failed to call /: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		htmlContent := string(body)

		// Check for key elements in the web UI
		if !strings.Contains(htmlContent, "gowsay") {
			t.Error("Web UI should contain 'gowsay'")
		}
	})

	t.Run("CORS_Preflight", func(t *testing.T) {
		req, _ := http.NewRequest("OPTIONS", baseURL+"/api/moo", nil)
		req.Header.Set("Origin", "http://example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send OPTIONS request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			t.Errorf("Expected status 200 or 204 for OPTIONS, got %d", resp.StatusCode)
		}

		if resp.Header.Get("Access-Control-Allow-Origin") == "" {
			t.Error("Expected CORS headers in preflight response")
		}
	})
}

// waitForServer polls the server URL until it responds or timeout is reached
func waitForServer(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 1 * time.Second}

	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("server did not start within %v", timeout)
}

// TestServerLifecycle tests server startup and graceful shutdown
func TestServerLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping lifecycle test in short mode")
	}

	// Set up environment
	os.Setenv("PORT", "19999")
	defer os.Unsetenv("PORT")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start server in goroutine
	serverDone := make(chan error, 1)
	go func() {
		// This simulates what runServer() does
		module := api.NewModule()
		mux := http.NewServeMux()

		mux.Handle("/api/moo", api.CORS(http.HandlerFunc(module.APIMoo)))
		mux.Handle("/health", api.Health("test"))

		srv := &http.Server{
			Addr:    ":19999",
			Handler: mux,
		}

		go func() {
			<-ctx.Done()
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()
			srv.Shutdown(shutdownCtx)
		}()

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverDone <- err
		}
		close(serverDone)
	}()

	// Wait for server to start using polling
	if err := waitForServer("http://localhost:19999/health", 3*time.Second); err != nil {
		t.Fatalf("Server did not start: %v", err)
	}

	// Test that server is running
	resp, err := http.Get("http://localhost:19999/health")
	if err != nil {
		t.Fatalf("Server health check failed: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected health check to return 200, got %d", resp.StatusCode)
	}

	// Cancel context to trigger shutdown
	cancel()

	// Wait for server to shut down
	select {
	case err := <-serverDone:
		if err != nil {
			t.Errorf("Server shutdown with error: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Error("Server did not shut down within timeout")
	}
}
