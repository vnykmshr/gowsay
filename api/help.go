package api

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/vnykmshr/gowsay/cow"
)

// GetBanner returns the startup banner with usage information
func GetBanner(version string) string {
	return fmt.Sprintf("gowsay [%s][%s]\n%s\n%s", version, os.Getenv(envKey), GetUsageString(), GetHelpString())
}

// GetUsageString returns the usage string
func GetUsageString() string {
	return fmt.Sprintf("Usage: `/moo [%s|surprise] [cow] [mood] message`", cow.ActionThink)
}

// GetHelpString returns the help string with available cows and moods
func GetHelpString() string {
	cows := append([]string{"`" + commandRandom + "`"}, formatList(cow.List())...)
	moods := append([]string{"`" + commandRandom + "`"}, formatList(cow.ListMoods())...)
	sort.Strings(cows)
	sort.Strings(moods)

	return fmt.Sprintf("Cows: %s\nMoods: %s", strings.Join(cows, ", "), strings.Join(moods, ", "))
}

func formatList(items []string) []string {
	formatted := make([]string, len(items))
	for i, item := range items {
		formatted[i] = "`" + item + "`"
	}
	return formatted
}
