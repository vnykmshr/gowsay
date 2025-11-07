package cow

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"

	runewidth "github.com/mattn/go-runewidth"
	"github.com/mitchellh/go-wordwrap"
)

// Action types for cowsay
const (
	ActionSay   = "say"
	ActionThink = "think"
)

// Face represents the cow's facial expression
type Face struct {
	Eyes     string
	Tongue   string
	Thoughts string
	cowfile  string
}

// Render generates cowsay output with the specified parameters
func Render(text []string, cowName, mood, action string, columns int) string {
	inputs := wrapText(text, columns)
	width := maxWidth(inputs)
	msgs := padLines(inputs, width)

	if len(msgs) == 0 {
		msgs = append(msgs, RandomMessage())
	}

	face := newFace(cowName, mood)
	balloon := buildBalloon(face, action, msgs, width)
	cow := renderCow(face)

	return fmt.Sprintf("```\n%s%s\n```\n", balloon, cow)
}

// newFace creates a face with the specified cow and mood
func newFace(cowName, mood string) *Face {
	face := &Face{
		Eyes:    "oo",
		Tongue:  "  ",
		cowfile: cowName,
	}

	switch mood {
	case "borg":
		face.Eyes = "=="
	case "dead":
		face.Eyes = "xx"
		face.Tongue = "U "
	case "greedy":
		face.Eyes = "$$"
	case "paranoid":
		face.Eyes = "@@"
	case "stoned":
		face.Eyes = "**"
		face.Tongue = "U "
	case "tired":
		face.Eyes = "--"
	case "wired":
		face.Eyes = "OO"
	case "young":
		face.Eyes = ".."
	}

	return face
}

// renderCow renders the cow template with the given face
func renderCow(f *Face) string {
	tmpl := template.Must(template.New("cow").Parse(cows[f.cowfile]))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, f); err != nil {
		log.Println(err)
		return err.Error()
	}

	return buf.String()
}

// wrapText processes input text with word wrapping and tab expansion
func wrapText(args []string, columns int) []string {
	var msgs []string
	for _, arg := range args {
		expanded := strings.Replace(arg, "\t", "        ", -1)
		wrapped := wordwrap.WrapString(expanded, uint(columns))
		msgs = append(msgs, strings.Split(wrapped, "\n")...)
	}
	return msgs
}

// padLines pads each line to the specified width
func padLines(msgs []string, width int) []string {
	var ret []string
	for _, m := range msgs {
		s := m + strings.Repeat(" ", width-runewidth.StringWidth(m))
		ret = append(ret, s)
	}
	return ret
}

// buildBalloon constructs the speech/thought balloon
func buildBalloon(f *Face, action string, msgs []string, width int) string {
	var borders []string
	lineCount := len(msgs)

	if action == ActionThink {
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

	var lines []string

	topBorder := " " + strings.Repeat("_", width+2)
	bottomBorder := " " + strings.Repeat("-", width+2) + "\n"

	lines = append(lines, topBorder)

	if lineCount == 1 {
		lines = append(lines, fmt.Sprintf("%s %s %s", borders[0], msgs[0], borders[1]))
	} else {
		lines = append(lines, fmt.Sprintf("%s %s %s", borders[0], msgs[0], borders[1]))
		for i := 1; i < lineCount-1; i++ {
			lines = append(lines, fmt.Sprintf("%s %s %s", borders[4], msgs[i], borders[5]))
		}
		lines = append(lines, fmt.Sprintf("%s %s %s", borders[2], msgs[lineCount-1], borders[3]))
	}

	lines = append(lines, bottomBorder)
	return strings.Join(lines, "\n")
}

// maxWidth returns the maximum display width of all messages
func maxWidth(msgs []string) int {
	max := -1
	for _, m := range msgs {
		w := runewidth.StringWidth(m)
		if w > max {
			max = w
		}
	}
	return max
}
