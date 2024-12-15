package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/vnykmshr/gowsay/src/actions"
	"github.com/vnykmshr/gowsay/src/constants"
	"github.com/vnykmshr/gowsay/src/cows"
	"github.com/vnykmshr/gowsay/src/gowsay"
	"github.com/vnykmshr/gowsay/src/moods"
	"golang.org/x/exp/slices"
)

type ValidOptions struct {
	Actions []string
	Cows    []string
	Moods   []string
}

var validOptions ValidOptions

func init() {
	validOptions.Actions = append(validOptions.Actions, constants.ActionRandom)
	validOptions.Actions = append(validOptions.Actions, actions.GetActionNames()...)

	cowNames, err := cows.GetCowNames()
	if err != nil {
		fmt.Printf("Warning: Failed to load cow names: %v\n", err)
		validOptions.Cows = []string{constants.CowDefault}
	} else {
		validOptions.Cows = append(validOptions.Cows, constants.CowRandom)
		validOptions.Cows = append(validOptions.Cows, cowNames...)
	}

	validOptions.Moods = append(validOptions.Moods, constants.MoodRandom)
	validOptions.Moods = append(validOptions.Moods, moods.GetMoods()...)
}

func UsageMessage() {
	fmt.Printf(`Usage: gowsay-cli [action] [cow] [mood] message
Arguments:
  action: Optional. One of [%s]. Default is '%s'
  cow:    Optional. Cow name (%s). Default is '%s'.
  mood:   Optional. Mood name (%s). Default is '%s'.
  message: Required. The message you want to display.
Example:
  moo think default "Hello, World!"
  moo surprise apt dead "Oh no!"
  moo think "Just thinking about life..."
`, strings.Join(validOptions.Actions, ", "), constants.ActionRandom, strings.Join(validOptions.Cows, ", "), constants.CowRandom, strings.Join(validOptions.Moods, ", "), constants.MoodRandom)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Missing required arguments.")
		UsageMessage()
		os.Exit(1)
	}

	action, cow, mood, message, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		UsageMessage()
		os.Exit(1)
	}

	displayMessage(action, cow, mood, message)
}

func parseArgs(args []string) (string, string, string, string, error) {
	action := constants.ActionRandom
	cow := constants.CowRandom
	mood := constants.MoodRandom

	var remainingArgs []string

	for i, arg := range args {
		if slices.Contains(validOptions.Actions, arg) && action == constants.ActionRandom {
			action = arg
		} else if slices.Contains(validOptions.Cows, arg) && cow == constants.CowRandom {
			cow = arg
		} else if slices.Contains(validOptions.Moods, arg) && mood == constants.MoodRandom {
			mood = arg
		} else {
			remainingArgs = args[i:]
			break // Stop parsing options
		}
	}

	if len(remainingArgs) == 0 {
		return "", "", "", "", fmt.Errorf("message is required")
	}

	message := strings.Join(remainingArgs, " ")
	return action, cow, mood, message, nil
}

func displayMessage(action, cow, mood, message string) {
	fmt.Printf("Action : %s\nCow    : %s\nMood   : %s\nMessage: %s\n\nCow Art:\n", action, cow, mood, message)
	fmt.Println(gowsay.GetGowsay(action, cow, mood, 40, []string{message}))
}
