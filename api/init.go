package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/tokopedia/logging.v1"
)

// NewModule create new module
func NewModule() (*Module, error) {
	var cfg Config
	ok := logging.ReadModuleConfig(&cfg, "/etc/gowsay", "gowsay") || logging.ReadModuleConfig(&cfg, "files/etc/gowsay", "gowsay")
	if !ok {
		log.Println("failed to read config, loading defaults")
		cfg = getDefaultConfig()
	}

	return &Module{
		cfg: &cfg,
	}, nil
}

func getDefaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Name: "gowsay",
		},
		App: AppConfig{
			Token:   "devel",
			Columns: 40,
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
			Attachments:  []Attachment{Attachment{Text: GetHelpString()}},
		}

		writeJSON(w, response)
		return
	}

	if len(parts) > 0 && parts[0] == ActionSurprise {
		parts = parts[1:]

		if len(parts) == 0 {
			parts = []string{getRandomMoo()}
		}

		response := SlackResponse{
			ResponseType: ResponseInChannel,
			Text:         getGowsay(ActionSay, getRandomCow(), getRandomMood(), hlm.cfg.App.Columns, parts),
		}

		writeJSON(w, response)
		return
	}

	var action = ActionDefault
	var cow = CowDefault
	var mood = MoodDefault

	if len(parts) > 1 && parts[0] == ActionThink {
		action = parts[0]
		parts = parts[1:]
	}

	if len(parts) > 1 {
		if _, ok := cows[parts[0]]; ok {
			cow = parts[0]
			parts = parts[1:]
		}

		if _, ok := moods[parts[0]]; ok {
			mood = parts[0]
			parts = parts[1:]
		}

		if parts[0] == CowRandom {
			cow = getRandomCow()
			parts = parts[1:]
		}

		if parts[0] == MoodRandom {
			mood = getRandomMood()
			parts = parts[1:]
		}
	}

	if len(parts) == 0 {
		parts = append(parts, getRandomMoo())
	}

	log.Printf("%s %s %s %s %s", CommandMoo, action, cow, mood, strings.Join(parts, " "))
	response := SlackResponse{
		ResponseType: ResponseInChannel,
		Text:         getGowsay(action, cow, mood, hlm.cfg.App.Columns, parts),
	}

	writeJSON(w, response)
}

func (hlm *Module) motd(w http.ResponseWriter) {
	motd := getGowsay(ActionSay, getRandomCow(), getRandomMood(), hlm.cfg.App.Columns, []string{getRandomMoo()})
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
