package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/vnykmshr/gowsay/cow"
)

// NewModule create new module
func NewModule() *Module {
	cfg := loadConfig()
	return &Module{
		cfg: &cfg,
	}
}

func loadConfig() Config {
	// Load configuration from environment variables
	token := os.Getenv("GOWSAY_TOKEN")
	if token == "" {
		token = "devel"
	}

	columns := int32(40)
	if colStr := os.Getenv("GOWSAY_COLUMNS"); colStr != "" {
		if col, err := strconv.Atoi(colStr); err == nil && col > 0 {
			columns = int32(col)
		}
	}

	return Config{
		Server: ServerConfig{
			Name: "gowsay",
		},
		App: AppConfig{
			Token:   token,
			Columns: columns,
		},
	}
}

// Gowsay gowsay request handler
func (hlm *Module) Gowsay(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue(FieldToken)
	if os.Getenv(FieldEnv) == ValueProduction && token != hlm.cfg.App.Token && token != ValueDefaultToken {
		hlm.motd(w)
		return
	}

	text := r.FormValue(FieldText)
	if strings.TrimSpace(text) == "" {
		hlm.motd(w)
		return
	}

	parts := sanitize(strings.Split(text, " "))
	if len(parts) == 0 {
		hlm.motd(w)
		return
	}

	if len(parts) == 1 && (parts[0] == ActionList || parts[0] == ActionHelp) {
		response := SlackResponse{
			ResponseType: ResponseEphemeral,
			Text:         GetUsageString(),
			Attachments:  []Attachment{{Text: GetHelpString()}},
		}

		writeJSON(w, response)
		return
	}

	if len(parts) > 0 && parts[0] == ActionSurprise {
		parts = parts[1:]

		if len(parts) == 0 {
			parts = []string{cow.RandomMessage()}
		}

		response := SlackResponse{
			ResponseType: ResponseInChannel,
			Text:         cow.Render(parts, cow.RandomCow(), cow.RandomMood(), cow.ActionSay, int(hlm.cfg.App.Columns)),
		}

		writeJSON(w, response)
		return
	}

	var action = ActionDefault
	var cowName = CowDefault
	var mood = MoodDefault

	if len(parts) > 1 && parts[0] == ActionThink {
		action = cow.ActionThink
		parts = parts[1:]
	}

	if len(parts) > 1 {
		if cow.Exists(parts[0]) {
			cowName = parts[0]
			parts = parts[1:]
		}

		if cow.MoodExists(parts[0]) {
			mood = parts[0]
			parts = parts[1:]
		}

		if parts[0] == CowRandom {
			cowName = cow.RandomCow()
			parts = parts[1:]
		}

		if parts[0] == MoodRandom {
			mood = cow.RandomMood()
			parts = parts[1:]
		}
	}

	if len(parts) == 0 {
		parts = append(parts, cow.RandomMessage())
	}

	log.Printf("%s %s %s %s %s", CommandMoo, action, cowName, mood, strings.Join(parts, " "))
	response := SlackResponse{
		ResponseType: ResponseInChannel,
		Text:         cow.Render(parts, cowName, mood, action, int(hlm.cfg.App.Columns)),
	}

	writeJSON(w, response)
}

func (hlm *Module) motd(w http.ResponseWriter) {
	motd := cow.Render([]string{cow.RandomMessage()}, cow.RandomCow(), cow.RandomMood(), cow.ActionSay, int(hlm.cfg.App.Columns))
	_, err := w.Write([]byte(motd))
	if err != nil {
		log.Println(err)
	}
}

func writeJSON(w http.ResponseWriter, response interface{}) {
	w.Header().Set(FieldContentType, ValueApplicationJSON)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}

func sanitize(s []string) []string {
	var r []string
	for _, str := range s {
		if strings.TrimSpace(str) != "" {
			r = append(r, str)
		}
	}

	return r
}
