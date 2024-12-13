package help

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/vnykmshr/gowsay/src/constants"
	"github.com/vnykmshr/gowsay/src/cows"
	"github.com/vnykmshr/gowsay/src/moods"
)

var envField string

func init() {
	envField = os.Getenv(constants.FieldEnv)
}

// GetBanner gets banner with usage information
func GetBanner(version string) string {
	return fmt.Sprintf("gowsay [%s][%s]\n%s\n%s", version, envField, GetUsageString(), GetHelpString())
}

// GetUsageString returns the usage string
func GetUsageString() string {
	return fmt.Sprintf("Usage: `%s [%s|%s] [cow] [mood] message`",
		constants.CommandMoo, constants.ActionThink, constants.ActionSurprise)
}

// GetHelpString returns the help string for cows and moods
func GetHelpString() string {
	return strings.Join([]string{
		"**Cows**: " + formatList(getCows(true)),
		"**Moods**: " + formatList(getMoods(true)),
	}, "\n")
}

// getCows returns a sorted list of available cows
func getCows(sorted bool) []string {
	names := make([]string, 0, 1)
	names = append(names, "`"+constants.CowRandom+"`")

	cowNames, _ := cows.GetCowNames()
	for _, name := range cowNames {
		names = append(names, "`"+name+"`")
	}

	if sorted {
		sort.Strings(names)
	}
	return names
}

// getMoods returns a sorted list of available moods
func getMoods(sorted bool) []string {
	names := make([]string, 0, 1)
	names = append(names, "`"+constants.MoodRandom+"`")

	moodNames := moods.GetMoods()
	for _, name := range moodNames {
		names = append(names, "`"+name+"`")
	}

	if sorted {
		sort.Strings(names)
	}
	return names
}

// wrapString wraps text at the specified limit of words
func wrapString(text string, limit int) string {
	if text == "" || limit <= 0 {
		return text
	}

	words := strings.Fields(text)
	var builder strings.Builder

	for i := 0; i < len(words); i += limit {
		end := i + limit
		if end > len(words) {
			end = len(words)
		}
		builder.WriteString(strings.Join(words[i:end], " "))
		builder.WriteByte('\n')
	}

	return builder.String()
}

// formatList joins elements with a comma and returns a wrapped string
func formatList(items []string) string {
	return wrapString(strings.Join(items, ", "), 10)
}
