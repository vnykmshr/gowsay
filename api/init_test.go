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
	cfg := loadConfig()
	type fields struct {
		cfg *Config
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "production",
			fields: fields{
				cfg: &cfg,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=test", nil),
			},
		},
		{
			name: "empty",
			fields: fields{
				cfg: &cfg,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123", nil),
			},
		},
		{
			name: "default",
			fields: fields{
				cfg: &cfg,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text=all%20you%20see", nil),
			},
		},
		{
			name: "help",
			fields: fields{
				cfg: &cfg,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+ActionHelp, nil),
			},
		},
		{
			name: "surprise",
			fields: fields{
				cfg: &cfg,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+ActionSurprise, nil),
			},
		},
		{
			name: "think-default",
			fields: fields{
				cfg: &cfg,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+cow.ActionThink+"%20thinking", nil),
			},
		},
		{
			name: "think-apt",
			fields: fields{
				cfg: &cfg,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+cow.ActionThink+"%20apt%20thinking", nil),
			},
		},
		{
			name: "think-apt-greedy",
			fields: fields{
				cfg: &cfg,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+cow.ActionThink+"%20apt%20greedy%20thinking", nil),
			},
		},
		{
			name: "think-random",
			fields: fields{
				cfg: &cfg,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+cow.ActionThink+"%20random%20greedy%20thinking", nil),
			},
		},
		{
			name: "think-random-random",
			fields: fields{
				cfg: &cfg,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+cow.ActionThink+"%20random%20random%20thinking", nil),
			},
		},
		{
			name: "think-random-random-random",
			fields: fields{
				cfg: &cfg,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text="+cow.ActionThink+"%20random%20random", nil),
			},
		},
		{
			name: "say-random-random",
			fields: fields{
				cfg: &cfg,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text=random%20random%20all%20you%20see", nil),
			},
		},
		{
			name: "say-apt-young",
			fields: fields{
				cfg: &cfg,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text=apt%20young%20all%20you%20see", nil),
			},
		},
		{
			name: "say-apt",
			fields: fields{
				cfg: &cfg,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost/say?token=abc123&text=apt%20all%20you%20see", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "t1" {
				os.Setenv(FieldEnv, ValueProduction)
			} else {
				os.Setenv(FieldEnv, "")
			}
			hlm := &Module{
				cfg: tt.fields.cfg,
			}
			hlm.Gowsay(tt.args.w, tt.args.r)
		})
	}
}

func TestModule_motd(t *testing.T) {
	cfg := loadConfig()
	type fields struct {
		cfg *Config
	}
	type args struct {
		w http.ResponseWriter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "t1",
			fields: fields{
				cfg: &cfg,
			},
			args: args{
				w: httptest.NewRecorder(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hlm := &Module{
				cfg: tt.fields.cfg,
			}
			hlm.motd(tt.args.w)
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

func Test_loadConfig(t *testing.T) {
	// Clear env vars to test defaults
	os.Unsetenv("GOWSAY_TOKEN")
	os.Unsetenv("GOWSAY_COLUMNS")

	tests := []struct {
		name string
		want Config
	}{
		{
			name: "defaults",
			want: Config{
				Server: ServerConfig{
					Name: "gowsay",
				},
				App: AppConfig{
					Token:   "devel",
					Columns: 40,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
