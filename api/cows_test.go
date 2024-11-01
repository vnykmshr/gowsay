package api

import (
	"bytes"
	"html/template"
	"testing"
)

func TestGetRandomCow(t *testing.T) {
	c := getRandomCow()
	if _, ok := cows[c]; !ok {
		t.Errorf("What cow? %s", c)
	}
}

func TestCowTemplate(t *testing.T) {
	// will test all cow and mood combinations
	var buf bytes.Buffer
	for c, val := range cows {
		tmpl := template.Must(template.New("cow").Parse(val))

		for m := range moods {
			f := newFace(c, m)
			if err := tmpl.Execute(&buf, f); err != nil {
				t.Error(err)
			}
		}
	}
}

func Test_getRandomCow(t *testing.T) {
	tests := []struct {
		name    string
		present bool
	}{
		{
			name:    "t1",
			present: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, got := cows[getRandomCow()]; got != tt.present {
				t.Errorf("getRandomCow() = %v, present %v", got, tt.present)
			}
		})
	}
}
