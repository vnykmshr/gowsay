package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/vnykmshr/gowsay/cow"
)

// NewModule creates a new API handler module with configuration from environment
func NewModule() *Module {
	token := os.Getenv("GOWSAY_TOKEN")
	if token == "" {
		token = "devel"
	}

	columns := 40
	if colStr := os.Getenv("GOWSAY_COLUMNS"); colStr != "" {
		if col, err := strconv.Atoi(colStr); err == nil && col > 0 {
			columns = col
		}
	}

	return &Module{
		token:   token,
		columns: columns,
	}
}

// Gowsay handles Slack /moo command requests
func (m *Module) Gowsay(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue(FieldToken)
	if os.Getenv(FieldEnv) == ValueProduction && token != m.token && token != ValueDefaultToken {
		m.motd(w)
		return
	}

	text := r.FormValue(FieldText)
	if strings.TrimSpace(text) == "" {
		m.motd(w)
		return
	}

	parts := sanitize(strings.Split(text, " "))
	if len(parts) == 0 {
		m.motd(w)
		return
	}

	if len(parts) == 1 && (parts[0] == ActionList || parts[0] == ActionHelp) {
		writeJSON(w, SlackResponse{
			ResponseType: ResponseEphemeral,
			Text:         GetUsageString(),
			Attachments:  []Attachment{{Text: GetHelpString()}},
		}, http.StatusOK)
		return
	}

	if len(parts) > 0 && parts[0] == ActionSurprise {
		parts = parts[1:]
		if len(parts) == 0 {
			parts = []string{cow.RandomMessage()}
		}
		output := cow.Render(parts, cow.RandomCow(), cow.RandomMood(), cow.ActionSay, m.columns)
		writeJSON(w, SlackResponse{ResponseType: ResponseInChannel, Text: fmt.Sprintf("```\n%s\n```", output)}, http.StatusOK)
		return
	}

	var action = cow.ActionSay
	var cowName = CowDefault
	var mood = MoodDefault

	if len(parts) > 1 && parts[0] == cow.ActionThink {
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

	slog.Info("slack command", "command", CommandMoo, "action", action, "cow", cowName, "mood", mood, "text", strings.Join(parts, " "))
	output := cow.Render(parts, cowName, mood, action, m.columns)
	writeJSON(w, SlackResponse{ResponseType: ResponseInChannel, Text: fmt.Sprintf("```\n%s\n```", output)}, http.StatusOK)
}

func (m *Module) motd(w http.ResponseWriter) {
	motd := cow.Render([]string{cow.RandomMessage()}, cow.RandomCow(), cow.RandomMood(), cow.ActionSay, m.columns)
	_, err := w.Write([]byte(motd))
	if err != nil {
		slog.Error("failed to write motd response", "error", err)
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
