package gowsay

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/mattn/go-runewidth"
	"github.com/mitchellh/go-wordwrap"
	"github.com/vnykmshr/gowsay/src/actions"
	"github.com/vnykmshr/gowsay/src/constants"
	"github.com/vnykmshr/gowsay/src/cows"
	"github.com/vnykmshr/gowsay/src/faces"
	"github.com/vnykmshr/gowsay/src/moods"
	"github.com/vnykmshr/gowsay/src/moos"
	"github.com/vnykmshr/gowsay/src/utils"
)

var cowTemplateCache = make(map[string]*template.Template)

func GetGowsay(action, cow, mood string, columns int32, text []string) (string, error) {
	updateRandoms(&action, &cow, &mood)
	msgs := wrapAndFindMaxWidth(text, columns)
	if len(msgs) == 0 {
		msgs = []string{moos.GetRandomMoo()}
	}
	msgs = setPadding(msgs, maxWidth(msgs))

	f, err := faces.New(cow, mood)
	if err != nil {
		return "", fmt.Errorf("error creating face: %w", err)
	}

	balloon := constructBallonWithBuilder(f, action, msgs, maxWidth(msgs))
	cowArt, err := constructCowWithCache(f)
	if err != nil {
		return "", fmt.Errorf("error constructing cow: %w", err)
	}

	return fmt.Sprintf("```\n%s%s\n```\n", balloon, *cowArt), nil
}

func updateRandoms(action, cow, mood *string) {
	if *action == constants.ActionRandom {
		*action = actions.GetRandomAction()
	}
	if *cow == constants.CowRandom {
		*cow, _ = cows.GetRandomCowName()
	}
	if *mood == constants.MoodRandom {
		*mood = moods.GetRandomMood()
	}
}

func wrapAndFindMaxWidth(text []string, columns int32) []string {
	var msgs []string
	max := 0
	for _, arg := range text {
		expand := strings.ReplaceAll(arg, "\t", "        ") // Use regular spaces
		wrapped := wordwrap.WrapString(expand, uint(columns))
		lines := strings.Split(wrapped, "\n")
		msgs = append(msgs, lines...)
		for _, line := range lines { // find max width during wrap
			if l := runewidth.StringWidth(line); l > max {
				max = l
			}
		}
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

func constructBallonWithBuilder(f *faces.Face, action string, msgs []string, width int) string {
	var b strings.Builder
	lineCount := len(msgs)
	var borders []string

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

	b.WriteString(" " + strings.Repeat("_", width+2) + "\n")
	if lineCount == 1 {
		b.WriteString(fmt.Sprintf("%s %s %s\n", borders[0], msgs[0], borders[1]))
	} else {
		b.WriteString(fmt.Sprintf("%s %s %s\n", borders[0], msgs[0], borders[1]))
		for i := 1; i < lineCount-1; i++ {
			b.WriteString(fmt.Sprintf("%s %s %s\n", borders[4], msgs[i], borders[5]))
		}
		b.WriteString(fmt.Sprintf("%s %s %s\n", borders[2], msgs[lineCount-1], borders[3]))
	}
	b.WriteString(" " + strings.Repeat("-", width+2) + "\n")

	return b.String()
}

func constructCowWithCache(f *faces.Face) (*string, error) {
	tmpl, ok := cowTemplateCache[f.Cow]
	if !ok {
		cow, err := cows.GetCow(f.Cow)
		if err != nil {
			return nil, fmt.Errorf("error getting cow: %w", err)
		}

		tmpl, err = template.New("cow").Parse(cow.Art)
		if err != nil {
			return nil, fmt.Errorf("error parsing cow template: %w", err)
		}
		cowTemplateCache[f.Cow] = tmpl
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, f); err != nil {
		return nil, fmt.Errorf("error executing template: %w", err)
	}

	return utils.GetStringPtr(buf.String()), nil
}
