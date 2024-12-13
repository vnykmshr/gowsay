package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/vnykmshr/gowsay/src/constants"
	"github.com/vnykmshr/gowsay/src/cows"
	"github.com/vnykmshr/gowsay/src/gowsay"
	"github.com/vnykmshr/gowsay/src/moods"
)

const (
	defaultAction = constants.ActionDefault
	defaultCow    = constants.CowDefault
	defaultMood   = constants.MoodDefault
)

var (
	validActions map[string]bool
	validCows    map[string]bool
	validMoods   map[string]bool
)

func init() {
	// Initialize the maps
	validActions = map[string]bool{
		constants.ActionSay:      true,
		constants.ActionThink:    true,
		constants.ActionList:     true,
		constants.ActionHelp:     true,
		constants.ActionSurprise: true,
	}

	// Load cow names
	cowNames, err := cows.GetCowNames()
	if err != nil {
		fmt.Printf("Warning: Failed to load cow names: %v\n", err)
		validCows = map[string]bool{
			constants.CowDefault: true,
		}
	} else {
		validCows = make(map[string]bool)
		for _, cow := range cowNames {
			validCows[cow] = true
		}
	}

	// Load moods
	moodNames := moods.GetMoods()
	validMoods = make(map[string]bool)
	for _, mood := range moodNames {
		validMoods[mood] = true
	}
}

// UsageMessage displays how to use the command
func UsageMessage() {
	fmt.Printf(`Usage: moo [action] [cow] [mood] message
Arguments:
  action:  Optional. One of [%s]. Default is '%s'
  cow:     Optional. Cow name (%s). Default is '%s'.
  mood:    Optional. Mood name (%s). Default is '%s'.
  message: Required. The message you want to display.
Example:
  moo think default "Hello, World!"
  moo surprise apt dead "Oh no!"
  moo think "Just thinking about life..."
`, getKeysAsCSV(validActions), defaultAction, getKeysAsCSV(validCows), defaultCow, getKeysAsCSV(validMoods), defaultMood)
}

// Main function to handle the command-line arguments
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Missing required arguments.")
		UsageMessage()
		os.Exit(1)
	}

	// Parse the arguments
	action, cow, mood, message, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		UsageMessage()
		os.Exit(1)
	}

	// Display the result
	displayMessage(action, cow, mood, message)
}

// parseArgs parses the command-line arguments and returns action, cow, mood, and message
func parseArgs(args []string) (string, string, string, string, error) {
	// Defaults
	action := defaultAction
	cow := defaultCow
	mood := defaultMood

	index := 0

	// Check if first argument is an action
	if index < len(args) && validActions[args[index]] {
		action = args[index]
		index++
	}

	// Check if next argument is a cow
	if index < len(args) && validCows[args[index]] {
		cow = args[index]
		index++
	}

	// Check if next argument is a mood
	if index < len(args) && validMoods[args[index]] {
		mood = args[index]
		index++
	}

	// Ensure a message is present
	if index >= len(args) {
		return "", "", "", "", fmt.Errorf("message is required")
	}

	// Collect the message from remaining arguments
	message := strings.Join(args[index:], " ")
	return action, cow, mood, message, nil
}

// displayMessage formats and displays the cow's message
func displayMessage(action, cow, mood, message string) {
	fmt.Printf("\nAction : %s\nCow    : %s\nMood   : %s\nMessage: %s\n\nCow Art:\n", action, cow, mood, message)
	fmt.Println(gowsay.GetGowsay(action, cow, mood, 40, []string{message}))
}

// getKeysAsCSV returns the keys of a map as a comma-separated string
func getKeysAsCSV(m map[string]bool) string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return strings.Join(keys, ", ")
}
