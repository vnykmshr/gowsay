package api

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/vnykmshr/gowsay/cow"
)

// GetBanner gets banner with usage information
func GetBanner(version string) string {
	return fmt.Sprintf("%s [%s][%s]\n%s\n%s", "gowsay", version, os.Getenv(FieldEnv), GetUsageString(), GetHelpString())
}

// GetUsageString usage string
func GetUsageString() string {
	return wrapString(fmt.Sprintf("Usage: `%s [%s|%s] [cow] [mood] message`", CommandMoo, cow.ActionThink, ActionSurprise), 0)
}

// GetHelpString help string
func GetHelpString() string {
	return strings.Join([]string{
		wrapString("Cows: "+strings.Join(getCows(true), ", "), 10),
		wrapString("Moods: "+strings.Join(getMoods(true), ", "), 10),
	}, "\n")
}

func getCows(sorted bool) []string {
	var names []string
	names = append(names, "`"+CowRandom+"`")
	for _, name := range cow.List() {
		names = append(names, "`"+name+"`")
	}

	if sorted {
		sort.Strings(names)
	}

	return names
}

func getMoods(sorted bool) []string {
	var names []string
	names = append(names, "`"+MoodRandom+"`")
	for _, name := range cow.ListMoods() {
		names = append(names, "`"+name+"`")
	}

	if sorted {
		sort.Strings(names)
	}

	return names
}

// Wraps text at the specified number of columns
func wrapString(s string, limit int) string {
	if strings.TrimSpace(s) == "" {
		return ""
	}

	var result string

	strSlice := strings.Fields(s)
	if limit == 0 || limit > len(strSlice) {
		limit = len(strSlice)
	}
	for len(strSlice) >= 1 {
		// convert slice/array back to string adding \n at specified limit
		result = result + strings.Join(strSlice[:limit], " ") + "\n"

		// discard the elements that were copied over to result
		strSlice = strSlice[limit:]

		// change the limit to cater for the last few words in
		if len(strSlice) < limit {
			limit = len(strSlice)
		}
	}

	return result
}
