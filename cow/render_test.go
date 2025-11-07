package cow

import (
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name    string
		text    []string
		cow     string
		mood    string
		action  string
		columns int
	}{
		{
			name:    "simple message",
			text:    []string{"hello world"},
			cow:     "default",
			mood:    "",
			action:  ActionSay,
			columns: 40,
		},
		{
			name:    "think action",
			text:    []string{"thinking"},
			cow:     "default",
			mood:    "dead",
			action:  ActionThink,
			columns: 40,
		},
		{
			name:    "multi-line",
			text:    []string{"line 1", "line 2", "line 3"},
			cow:     "dragon",
			mood:    "wired",
			action:  ActionSay,
			columns: 40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Render(tt.text, tt.cow, tt.mood, tt.action, tt.columns)
			if result == "" {
				t.Error("Render() returned empty string")
			}
			// Should contain the cow/balloon characters, not markdown
			if !strings.Contains(result, "\\") && !strings.Contains(result, "/") && !strings.Contains(result, "-") {
				t.Error("Render() should return ASCII art with balloon characters")
			}
		})
	}
}

func TestRandomCow(t *testing.T) {
	cow := RandomCow()
	if cow == "" {
		t.Error("RandomCow() returned empty string")
	}
	if !Exists(cow) {
		t.Errorf("RandomCow() returned non-existent cow: %s", cow)
	}
}

func TestList(t *testing.T) {
	cows := List()
	if len(cows) == 0 {
		t.Error("List() returned empty list")
	}
	if len(cows) < 40 {
		t.Errorf("List() returned too few cows: got %d, want at least 40", len(cows))
	}
}

func TestExists(t *testing.T) {
	tests := []struct {
		name string
		cow  string
		want bool
	}{
		{"default exists", "default", true},
		{"dragon exists", "dragon", true},
		{"nonexistent", "nonexistent", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exists(tt.cow); got != tt.want {
				t.Errorf("Exists(%s) = %v, want %v", tt.cow, got, tt.want)
			}
		})
	}
}

func TestRandomMood(t *testing.T) {
	mood := RandomMood()
	if mood == "" {
		t.Error("RandomMood() returned empty string")
	}
	if !MoodExists(mood) {
		t.Errorf("RandomMood() returned non-existent mood: %s", mood)
	}
}

func TestListMoods(t *testing.T) {
	moods := ListMoods()
	if len(moods) == 0 {
		t.Error("ListMoods() returned empty list")
	}
	if len(moods) < 5 {
		t.Errorf("ListMoods() returned too few moods: got %d, want at least 5", len(moods))
	}
}

func TestMoodExists(t *testing.T) {
	tests := []struct {
		name string
		mood string
		want bool
	}{
		{"borg exists", "borg", true},
		{"dead exists", "dead", true},
		{"nonexistent", "nonexistent", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MoodExists(tt.mood); got != tt.want {
				t.Errorf("MoodExists(%s) = %v, want %v", tt.mood, got, tt.want)
			}
		})
	}
}

func TestRandomMessage(t *testing.T) {
	msg := RandomMessage()
	if msg == "" {
		t.Error("RandomMessage() returned empty string")
	}
}

func TestMaxWidth(t *testing.T) {
	tests := []struct {
		name string
		msgs []string
		want int
	}{
		{"empty", []string{}, -1},
		{"single", []string{"hello"}, 5},
		{"multiple", []string{"hi", "hello", "hey"}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxWidth(tt.msgs); got != tt.want {
				t.Errorf("maxWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}
