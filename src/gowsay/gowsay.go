package gowsay

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/mattn/go-runewidth"
	"github.com/mitchellh/go-wordwrap"
	"github.com/vnykmshr/gowsay/src/constants"
	"github.com/vnykmshr/gowsay/src/cows"
	"github.com/vnykmshr/gowsay/src/faces"
	"github.com/vnykmshr/gowsay/src/moods"
	"github.com/vnykmshr/gowsay/src/moos"
	"github.com/vnykmshr/gowsay/src/utils"
)

func GetGowsay(action, cow, mood string, columns int32, text []string) (string, error) {
	msgs := setPadding(readInput(text, columns), maxWidth(readInput(text, columns)))
	if len(msgs) == 0 {
		msgs = []string{moos.GetRandomMoo()}
	}

	if cow == constants.CowRandom {
		cow, _ = cows.GetRandomCowName()
	}

	if mood == constants.MoodRandom {
		mood = moods.GetRandomMood()
	}

	f, err := faces.New(cow, mood)
	if err != nil {
		return "", fmt.Errorf("error creating face: %w", err)
	}

	balloon := constructBallon(f, action, msgs, maxWidth(msgs))
	cowArt, err := constructCow(f)
	if err != nil {
		return "", fmt.Errorf("error constructing cow: %w", err)
	}

	return fmt.Sprintf("```\n%s%s\n```\n", balloon, *cowArt), nil
}

func readInput(args []string, columns int32) []string {
	var msgs []string
	for _, arg := range args {
		expand := strings.ReplaceAll(arg, "\t", "        ")
		wrapped := wordwrap.WrapString(expand, uint(columns))
		msgs = append(msgs, strings.Split(wrapped, "\n")...)
	}
	return msgs
}

func maxWidth(msgs []string) int {
	max := 0
	for _, m := range msgs {
		if l := runewidth.StringWidth(m); l > max {
			max = l
		}
	}
	return max
}

func setPadding(msgs []string, width int) []string {
	ret := make([]string, len(msgs))
	for i, m := range msgs {
		padding := width - runewidth.StringWidth(m)
		ret[i] = m + strings.Repeat(" ", padding)
	}
	return ret
}

func constructBallon(f *faces.Face, action string, msgs []string, width int) string {
	var borders []string
	lineCount := len(msgs)
	lines := make([]string, 0, lineCount+3) // Pre-allocate with sufficient capacity

	if action == constants.ActionThink {
		f.Thoughts = "o"
		borders = []string{"(", ")", "(", ")", "(", ")"}
	} else {
		f.Thoughts = "\\"
		if lineCount == 1 {
			borders = []string{"<", ">"}
		} else {
			borders = []string{"/", "\\", "\\", "/", "|", "|"}
		}
	}

	lines = append(lines, " "+strings.Repeat("_", width+2))
	if lineCount == 1 {
		lines = append(lines, fmt.Sprintf("%s %s %s", borders[0], msgs[0], borders[1]))
	} else {
		lines = append(lines, fmt.Sprintf("%s %s %s", borders[0], msgs[0], borders[1]))
		for i := 1; i < lineCount-1; i++ {
			lines = append(lines, fmt.Sprintf("%s %s %s", borders[4], msgs[i], borders[5]))
		}
		lines = append(lines, fmt.Sprintf("%s %s %s", borders[2], msgs[lineCount-1], borders[3]))
	}
	lines = append(lines, " "+strings.Repeat("-", width+2)+"\n")

	return strings.Join(lines, "\n")
}

func constructCow(f *faces.Face) (*string, error) {
	cow, err := cows.GetCow(f.Cow)
	if err != nil {
		return nil, fmt.Errorf("error construct cow: %w", err)
	}

	t := template.Must(template.New("cow").Parse(cow.Art))
	var buf bytes.Buffer
	if err := t.Execute(&buf, f); err != nil {
		return nil, fmt.Errorf("error executing template: %w", err)
	}

	return utils.GetStringPtr(buf.String()), nil
}
