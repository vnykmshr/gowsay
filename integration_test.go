package main

import (
	"bytes"
	"context"
	"encoding/json"
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

	t.Run("GetCows", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/api/cows")
		if err != nil {
			t.Fatalf("Failed to call /api/cows: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string][]string
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode cows response: %v", err)
		}

		cows, ok := result["cows"]
		if !ok {
			t.Fatal("Response missing 'cows' field")
		}

		if len(cows) == 0 {
			t.Error("Expected non-empty list of cows")
		}

		// Check for some known cows
		hasCow := false
		for _, cow := range cows {
			if cow == "default" || cow == "dragon" {
				hasCow = true
				break
			}
		}
		if !hasCow {
			t.Error("Expected to find 'default' or 'dragon' in cows list")
		}
	})

	t.Run("GetMoods", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/api/moods")
		if err != nil {
			t.Fatalf("Failed to call /api/moods: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string][]string
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode moods response: %v", err)
		}

		moods, ok := result["moods"]
		if !ok {
			t.Fatal("Response missing 'moods' field")
		}

		if len(moods) == 0 {
			t.Error("Expected non-empty list of moods")
		}
	})

	t.Run("APIMoo_JSON", func(t *testing.T) {
		reqBody := map[string]string{
			"text":   "Hello Integration Test",
			"cow":    "default",
			"action": "say",
		}
		jsonData, _ := json.Marshal(reqBody)

		resp, err := client.Post(
			baseURL+"/api/moo",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			t.Fatalf("Failed to call /api/moo: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d. Body: %s", resp.StatusCode, body)
		}

		var result api.MooResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Output == "" {
			t.Error("Expected non-empty output")
		}

		if !strings.Contains(result.Output, "Hello Integration Test") {
			t.Error("Output should contain the input text")
		}
	})

	t.Run("APIMoo_QueryParams", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/api/moo?text=Query+Test&cow=dragon&action=think")
		if err != nil {
			t.Fatalf("Failed to call /api/moo with query params: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result api.MooResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !strings.Contains(result.Output, "Query Test") {
			t.Error("Output should contain the query text")
		}
	})

	t.Run("APIMoo_RandomCow", func(t *testing.T) {
		reqBody := map[string]string{
			"text": "Random test",
			"cow":  "random",
		}
		jsonData, _ := json.Marshal(reqBody)

		resp, err := client.Post(
			baseURL+"/api/moo",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			t.Fatalf("Failed to call /api/moo with random: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("APIMoo_InvalidCow", func(t *testing.T) {
		reqBody := map[string]string{
			"text": "Test",
			"cow":  "nonexistent",
		}
		jsonData, _ := json.Marshal(reqBody)

		resp, err := client.Post(
			baseURL+"/api/moo",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			t.Fatalf("Failed to call /api/moo: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid cow, got %d", resp.StatusCode)
		}
	})

	t.Run("APIMoo_MissingText", func(t *testing.T) {
		reqBody := map[string]string{
			"cow": "default",
		}
		jsonData, _ := json.Marshal(reqBody)

		resp, err := client.Post(
			baseURL+"/api/moo",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			t.Fatalf("Failed to call /api/moo: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 for missing text, got %d", resp.StatusCode)
		}
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

	// Wait for server to start
	time.Sleep(500 * time.Millisecond)

	// Test that server is running
	resp, err := http.Get("http://localhost:19999/health")
	if err != nil {
		t.Fatalf("Server did not start properly: %v", err)
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
