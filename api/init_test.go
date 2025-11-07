package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/vnykmshr/gowsay/cow"
)

func TestNewModule(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "t1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if m := NewModule(); m == nil {
				t.Errorf("NewModule() m = %v, wantErr %v", m, tt.wantErr)
			}
		})
	}
}

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
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+ActionHelp, nil),
			},
		},
		{
			name: "surprise",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+ActionSurprise, nil),
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
				os.Setenv(FieldEnv, ValueProduction)
			} else {
				os.Setenv(FieldEnv, "")
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
	tests := []struct {
		name string
		args args
	}{
		{
			name: "t1",
			args: args{
				w: httptest.NewRecorder(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				token:   "test",
				columns: 40,
			}
			m.motd(tt.args.w)
		})
	}
}

func Test_writeJSON(t *testing.T) {
	type args struct {
		w        http.ResponseWriter
		response interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "t1",
			args: args{
				w:        httptest.NewRecorder(),
				response: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeJSON(tt.args.w, tt.args.response, http.StatusOK)
			if got, ok := tt.args.w.Header()[FieldContentType]; !ok || !reflect.DeepEqual(got, []string{ValueApplicationJSON}) {
				t.Errorf("writeJSON() Header: %s: got: %s, want: %s", FieldContentType, got, []string{ValueApplicationJSON})
			}
		})
	}
}

func Test_sanitize(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "t1",
			args: args{
				s: []string{"abc", " ", " abc"},
			},
			want: []string{"abc", " abc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitize(tt.args.s); !reflect.DeepEqual(got, tt.want) {
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
			// Set up env vars
			if tt.tokenEnv != "" {
				os.Setenv("GOWSAY_TOKEN", tt.tokenEnv)
			} else {
				os.Unsetenv("GOWSAY_TOKEN")
			}
			if tt.columnsEnv != "" {
				os.Setenv("GOWSAY_COLUMNS", tt.columnsEnv)
			} else {
				os.Unsetenv("GOWSAY_COLUMNS")
			}

			got := NewModule()
			if got.token != tt.wantToken {
				t.Errorf("NewModule().token = %v, want %v", got.token, tt.wantToken)
			}
			if got.columns != tt.wantColumns {
				t.Errorf("NewModule().columns = %v, want %v", got.columns, tt.wantColumns)
			}

			// Clean up
			os.Unsetenv("GOWSAY_TOKEN")
			os.Unsetenv("GOWSAY_COLUMNS")
		})
	}
}
