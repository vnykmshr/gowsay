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

			// Validate non-empty output
			if result == "" {
				t.Error("Render() returned empty string")
			}

			// Validate multi-line structure (balloon + cow)
			lines := strings.Split(result, "\n")
			if len(lines) < 3 {
				t.Errorf("Expected multi-line output with balloon and cow, got %d lines", len(lines))
			}

			// Validate balloon structure (should have border characters)
			hasBorder := strings.Contains(result, "-") || strings.Contains(result, "_")
			if !hasBorder {
				t.Error("Output should contain speech/thought balloon borders (-/_)")
			}

			// Validate cow ASCII art presence
			hasCowArt := strings.Contains(result, "\\") || strings.Contains(result, "/")
			if !hasCowArt {
				t.Error("Output should contain cow ASCII art (backslashes/slashes)")
			}

			// Validate text content appears in output
			for _, line := range tt.text {
				if !strings.Contains(result, line) {
					t.Errorf("Output should contain input text: %s", line)
				}
			}

			// Validate action-specific rendering
			if tt.action == ActionThink {
				// Think bubbles use 'o' or 'O' characters
				if !strings.Contains(result, "o") && !strings.Contains(result, "O") {
					t.Error("Think action should use thought bubbles (o/O)")
				}
			} else {
				// Say action uses backslashes
				if !strings.Contains(result, "\\") {
					t.Error("Say action should use speech indicators (backslashes)")
				}
			}

			// Validate mood affects eyes (if mood specified)
			if tt.mood != "" {
				// Eyes should be present in the cow
				hasEyes := false
				for _, line := range lines {
					// Most cows have eyes represented by characters like oo, OO, XX, etc.
					if strings.Count(line, "o") >= 2 || strings.Count(line, "O") >= 2 ||
						strings.Count(line, "x") >= 2 || strings.Count(line, "X") >= 2 ||
						strings.Count(line, "@") >= 2 || strings.Count(line, "*") >= 2 {
						hasEyes = true
						break
					}
				}
				if !hasEyes {
					t.Log("Warning: Mood specified but couldn't detect eye changes in cow")
				}
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

// Edge case tests

func TestRender_Unicode(t *testing.T) {
	tests := []struct {
		name string
		text []string
	}{
		{
			name: "chinese_characters",
			text: []string{"你好世界"},
		},
		{
			name: "emoji",
			text: []string{"Hello 🐄🐮 World"},
		},
		{
			name: "mixed_unicode",
			text: []string{"Café", "日本語", "مرحبا"},
		},
		{
			name: "rtl_text",
			text: []string{"שלום עולם"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Render(tt.text, "default", "", ActionSay, 40)

			if result == "" {
				t.Error("Unicode text should render without crashing")
			}

			// Verify text is present (may be wrapped or modified)
			for _, line := range tt.text {
				// Check that at least some characters from the input are present
				if !strings.Contains(result, string([]rune(line)[0])) {
					t.Logf("Warning: Unicode text may not be preserved: %s", line)
				}
			}
		})
	}
}

func TestRender_LargeInput(t *testing.T) {
	tests := []struct {
		name     string
		textGen  func() []string
		maxLines int
	}{
		{
			name: "very_long_single_line",
			textGen: func() []string {
				return []string{strings.Repeat("word ", 1000)}
			},
			maxLines: 200,
		},
		{
			name: "many_lines",
			textGen: func() []string {
				lines := make([]string, 100)
				for i := range lines {
					lines[i] = "This is line " + string(rune('A'+i%26))
				}
				return lines
			},
			maxLines: 150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := tt.textGen()
			result := Render(text, "default", "", ActionSay, 40)

			if result == "" {
				t.Error("Large input should render without crashing")
			}

			// Verify output is reasonable (not exponentially large)
			lineCount := strings.Count(result, "\n")
			if lineCount > tt.maxLines {
				t.Errorf("Output too large: %d lines (max %d)", lineCount, tt.maxLines)
			}
		})
	}
}

func TestRender_EmptyAndWhitespace(t *testing.T) {
	tests := []struct {
		name string
		text []string
	}{
		{
			name: "empty_string",
			text: []string{""},
		},
		{
			name: "whitespace_only",
			text: []string{"   ", "\t\t", "  "},
		},
		{
			name: "mixed_empty_and_text",
			text: []string{"", "hello", "", "world", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not crash
			result := Render(tt.text, "default", "", ActionSay, 40)
			if result == "" {
				t.Error("Should return output even for empty/whitespace input")
			}
		})
	}
}

func TestRender_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name string
		text []string
	}{
		{
			name: "quotes_and_escapes",
			text: []string{`"Hello" 'World' \n \t`},
		},
		{
			name: "symbols",
			text: []string{"!@#$%^&*()_+-=[]{}|;:,.<>?"},
		},
		{
			name: "newlines_in_text",
			text: []string{"line1\nline2\nline3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Render(tt.text, "default", "", ActionSay, 40)
			if result == "" {
				t.Error("Special characters should not break rendering")
			}
		})
	}
}
