package api

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"

	runewidth "github.com/mattn/go-runewidth"
	"github.com/mitchellh/go-wordwrap"
)

func newFace(cow, mood string) *face {
	f := &face{
		Eyes:    "oo",
		Tongue:  "  ",
		cowfile: cow,
	}

	switch mood {
	case "borg":
		f.Eyes = "=="
	case "dead":
		f.Eyes = "xx"
		f.Tongue = "U "
	case "greedy":
		f.Eyes = "$$"
	case "paranoid":
		f.Eyes = "@@"
	case "stoned":
		f.Eyes = "**"
		f.Tongue = "U "
	case "tired":
		f.Eyes = "--"
	case "wired":
		f.Eyes = "OO"
	case "young":
		f.Eyes = ".."
	}

	return f
}

// getGowsay get cowsay
func getGowsay(action, cow, mood string, columns int32, text []string) string {
	inputs := readInput(text, columns)
	width := maxWidth(inputs)
	msgs := setPadding(inputs, width)
	if len(msgs) == 0 {
		msgs = append(msgs, getRandomMoo())
	}

	f := newFace(cow, mood)
	balloon := constructBallon(f, action, msgs, width)
	cowsay := constructCow(f)
	return fmt.Sprintf("```\n%s%s\n```\n", balloon, cowsay)
}

func constructCow(f *face) string {
	t := template.Must(template.New("cow").Parse(cows[f.cowfile]))

	var buf bytes.Buffer
	if err := t.Execute(&buf, f); err != nil {
		log.Println(err)
		_, err := buf.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
	}

	return buf.String()
}

func readInput(args []string, columns int32) []string {
	var msgs []string
	for i := 0; i < len(args); i++ {
		expand := strings.Replace(args[i], "\t", "        ", -1)
		tmp := wordwrap.WrapString(expand, uint(columns))
		msgs = append(msgs, strings.Split(tmp, "\n")...)
	}

	return msgs
}

func setPadding(msgs []string, width int) []string {
	var ret []string
	for _, m := range msgs {
		s := m + strings.Repeat(" ", width-runewidth.StringWidth(m))
		ret = append(ret, s)
	}

	return ret
}

func constructBallon(f *face, action string, msgs []string, width int) string {
	var borders []string
	line := len(msgs)

	if action == ActionThink {
		f.Thoughts = "o"
		borders = []string{"(", ")", "(", ")", "(", ")"}
	} else {
		f.Thoughts = "\\"
		if line == 1 {
			borders = []string{"<", ">"}
		} else {
			borders = []string{"/", "\\", "\\", "/", "|", "|"}
		}
	}

	var lines []string

	topBorder := " " + strings.Repeat("_", width+2)
	bottomBoder := " " + strings.Repeat("-", width+2) + "\n"

	lines = append(lines, topBorder)
	if line == 1 {
		s := fmt.Sprintf("%s %s %s", borders[0], msgs[0], borders[1])
		lines = append(lines, s)
	} else {
		s := fmt.Sprintf(`%s %s %s`, borders[0], msgs[0], borders[1])
		lines = append(lines, s)
		i := 1
		for ; i < line-1; i++ {
			s = fmt.Sprintf(`%s %s %s`, borders[4], msgs[i], borders[5])
			lines = append(lines, s)
		}
		s = fmt.Sprintf(`%s %s %s`, borders[2], msgs[i], borders[3])
		lines = append(lines, s)
	}

	lines = append(lines, bottomBoder)
	return strings.Join(lines, "\n")
}

func maxWidth(msgs []string) int {
	max := -1
	for _, m := range msgs {
		l := runewidth.StringWidth(m)
		if l > max {
			max = l
		}
	}

	return max
}
