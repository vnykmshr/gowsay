package api

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"sort"
	"strconv"

	"github.com/vnykmshr/gowsay/cow"
	"github.com/vnykmshr/gowsay/web"
)

// MooRequest represents a request to generate cowsay
type MooRequest struct {
	Text    string `json:"text"`
	Cow     string `json:"cow,omitempty"`
	Mood    string `json:"mood,omitempty"`
	Action  string `json:"action,omitempty"`
	Columns int    `json:"columns,omitempty"`
}

// MooResponse represents the cowsay output
type MooResponse struct {
	Output string `json:"output"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// APIMoo handles /api/moo endpoint - accepts both JSON and query params
func (m *Module) APIMoo(w http.ResponseWriter, r *http.Request) {
	var req MooRequest

	// Try to parse JSON body first
	if r.Header.Get("Content-Type") == "application/json" && r.Method == "POST" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else {
		// Fall back to query parameters or form data
		req.Text = r.FormValue("text")
		req.Cow = r.FormValue("cow")
		req.Mood = r.FormValue("mood")
		req.Action = r.FormValue("action")
		if colStr := r.FormValue("columns"); colStr != "" {
			if col, err := strconv.Atoi(colStr); err == nil && col > 0 {
				req.Columns = col
			}
		}
	}

	// Set defaults
	if req.Cow == "" {
		req.Cow = "default"
	}
	if req.Action == "" {
		req.Action = cow.ActionSay
	}
	if req.Columns == 0 {
		req.Columns = int(m.cfg.App.Columns)
	}

	// Handle random
	if req.Cow == "random" {
		req.Cow = cow.RandomCow()
	}
	if req.Mood == "random" {
		req.Mood = cow.RandomMood()
	}

	// Validate
	if req.Text == "" {
		writeJSONError(w, "text parameter is required", http.StatusBadRequest)
		return
	}
	if !cow.Exists(req.Cow) {
		writeJSONError(w, fmt.Sprintf("cow '%s' not found", req.Cow), http.StatusBadRequest)
		return
	}
	if req.Mood != "" && !cow.MoodExists(req.Mood) {
		writeJSONError(w, fmt.Sprintf("mood '%s' not found", req.Mood), http.StatusBadRequest)
		return
	}
	if req.Action != cow.ActionSay && req.Action != cow.ActionThink {
		req.Action = cow.ActionSay
	}

	// Render
	output := cow.Render([]string{req.Text}, req.Cow, req.Mood, req.Action, req.Columns)

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MooResponse{Output: output})
}

// APICows handles /api/cows endpoint - lists all available cows
func (m *Module) APICows(w http.ResponseWriter, r *http.Request) {
	cows := cow.List()
	sort.Strings(cows)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{
		"cows": cows,
	})
}

// APIMoods handles /api/moods endpoint - lists all available moods
func (m *Module) APIMoods(w http.ResponseWriter, r *http.Request) {
	moods := cow.ListMoods()
	sort.Strings(moods)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{
		"moods": moods,
	})
}

// Health handles /health endpoint
func Health(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{
			Status:  "ok",
			Version: version,
		})
	}
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// ServeWeb serves the web UI
func ServeWeb() http.Handler {
	fsys, err := fs.Sub(web.Files, ".")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(fsys))
}
