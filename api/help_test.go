package api

import (
	"strings"
	"testing"

	"github.com/vnykmshr/gowsay/cow"
)

func TestGetUsageString(t *testing.T) {
	usage := GetUsageString()

	if !strings.Contains(usage, "Usage:") {
		t.Error("Usage string should contain 'Usage:' label")
	}
	if !strings.Contains(usage, "/moo") {
		t.Error("Usage string should contain command '/moo'")
	}
	if !strings.Contains(usage, cow.ActionThink) {
		t.Errorf("Usage string should contain action '%s'", cow.ActionThink)
	}
	if !strings.Contains(usage, commandSurprise) {
		t.Errorf("Usage string should contain action '%s'", commandSurprise)
	}
	if !strings.Contains(usage, "cow") || !strings.Contains(usage, "mood") || !strings.Contains(usage, "message") {
		t.Error("Usage string should contain parameters: cow, mood, message")
	}
}

func TestGetHelpString(t *testing.T) {
	help := GetHelpString()

	if !strings.Contains(help, "Cows:") {
		t.Error("Help string should contain 'Cows:' section")
	}
	if !strings.Contains(help, "Moods:") {
		t.Error("Help string should contain 'Moods:' section")
	}

	knownCows := []string{"default", "dragon", "apt"}
	for _, cowName := range knownCows {
		if !strings.Contains(help, "`"+cowName+"`") {
			t.Errorf("Help should contain cow '%s'", cowName)
		}
	}

	knownMoods := []string{"borg", "dead", "wired"}
	for _, mood := range knownMoods {
		if !strings.Contains(help, "`"+mood+"`") {
			t.Errorf("Help should contain mood '%s'", mood)
		}
	}

	if !strings.Contains(help, "`"+commandRandom+"`") {
		t.Error("Help should contain 'random' option")
	}
}

func TestGetBanner(t *testing.T) {
	testVersion := "v1.2.3-test"
	t.Setenv(envKey, "testing")

	banner := GetBanner(testVersion)

	sections := []string{
		"gowsay",
		testVersion,
		"testing",
		"Usage:",
		"Cows:",
		"Moods:",
	}

	for _, section := range sections {
		if !strings.Contains(banner, section) {
			t.Errorf("Banner should contain '%s'", section)
		}
	}

	lines := strings.Split(banner, "\n")
	if len(lines) < 1 {
		t.Fatal("Banner should have at least one line")
	}

	firstLine := lines[0]
	if !strings.Contains(firstLine, "["+testVersion+"]") {
		t.Error("First line should contain version in brackets")
	}
	if !strings.Contains(firstLine, "[testing]") {
		t.Error("First line should contain environment in brackets")
	}
}
