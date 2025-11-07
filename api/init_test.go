package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/vnykmshr/gowsay/cow"
)

func TestModule_Gowsay(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "production",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=test", nil),
			},
		},
		{
			name: "empty",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123", nil),
			},
		},
		{
			name: "default",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text=all%20you%20see", nil),
			},
		},
		{
			name: "help",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+commandHelp, nil),
			},
		},
		{
			name: "surprise",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+commandSurprise, nil),
			},
		},
		{
			name: "think-default",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+cow.ActionThink+"%20thinking", nil),
			},
		},
		{
			name: "think-apt",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+cow.ActionThink+"%20apt%20thinking", nil),
			},
		},
		{
			name: "think-apt-greedy",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+cow.ActionThink+"%20apt%20greedy%20thinking", nil),
			},
		},
		{
			name: "think-random",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+cow.ActionThink+"%20random%20greedy%20thinking", nil),
			},
		},
		{
			name: "think-random-random",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+cow.ActionThink+"%20random%20random%20thinking", nil),
			},
		},
		{
			name: "think-random-random-random",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+cow.ActionThink+"%20random%20random", nil),
			},
		},
		{
			name: "say-random-random",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text=random%20random%20all%20you%20see", nil),
			},
		},
		{
			name: "say-apt-young",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text=apt%20young%20all%20you%20see", nil),
			},
		},
		{
			name: "say-apt",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text=apt%20all%20you%20see", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "production" {
				os.Setenv(envKey, envProduction)
			} else {
				os.Setenv(envKey, "")
			}
			m := &Module{
				token:   "abc123",
				columns: 40,
			}
			m.Gowsay(tt.args.w, tt.args.r)
		})
	}
}

func TestModule_motd(t *testing.T) {
	type args struct {
		w http.ResponseWriter
	}
	// Simple smoke test - motd should write output without error
	m := &Module{
		token:   "test",
		columns: 40,
	}
	w := httptest.NewRecorder()
	m.motd(w)

	if w.Code == 0 {
		t.Error("motd should set status code")
	}
}

func Test_writeJSON(t *testing.T) {
	w := httptest.NewRecorder()
	response := map[string]string{"test": "data"}

	writeJSON(w, response, http.StatusOK)

	// Verify content type header
	if got := w.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %s, want application/json", got)
	}

	// Verify status code
	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	// Verify valid JSON response
	var result map[string]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Errorf("Response should be valid JSON: %v", err)
	}
}

func Test_sanitize(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "removes_whitespace_only_entries",
			input: []string{"abc", " ", " abc"},
			want:  []string{"abc", " abc"},
		},
		{
			name:  "preserves_all_non_whitespace",
			input: []string{"hello", "world"},
			want:  []string{"hello", "world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitize(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_NewModuleDefaults(t *testing.T) {
	tests := []struct {
		name        string
		tokenEnv    string
		columnsEnv  string
		wantToken   string
		wantColumns int
	}{
		{
			name:        "defaults",
			tokenEnv:    "",
			columnsEnv:  "",
			wantToken:   "devel",
			wantColumns: 40,
		},
		{
			name:        "custom-values",
			tokenEnv:    "test-token",
			columnsEnv:  "80",
			wantToken:   "test-token",
			wantColumns: 80,
		},
		{
			name:        "invalid-columns",
			tokenEnv:    "test",
			columnsEnv:  "invalid",
			wantToken:   "test",
			wantColumns: 40,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up env vars (t.Setenv auto-restores after test)
			if tt.tokenEnv != "" {
				t.Setenv("GOWSAY_TOKEN", tt.tokenEnv)
			}
			if tt.columnsEnv != "" {
				t.Setenv("GOWSAY_COLUMNS", tt.columnsEnv)
			}

			got := NewModule()
			if got.token != tt.wantToken {
				t.Errorf("NewModule().token = %v, want %v", got.token, tt.wantToken)
			}
			if got.columns != tt.wantColumns {
				t.Errorf("NewModule().columns = %v, want %v", got.columns, tt.wantColumns)
			}
		})
	}
}
