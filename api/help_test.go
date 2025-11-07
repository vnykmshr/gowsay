package api

import (
	"strings"
	"testing"

	"github.com/vnykmshr/gowsay/cow"
)

func TestGetUsageString(t *testing.T) {
	usage := GetUsageString()

	// Behavioral assertions: verify structure and key elements
	if !strings.Contains(usage, "Usage:") {
		t.Error("Usage string should contain 'Usage:' label")
	}
	if !strings.Contains(usage, CommandMoo) {
		t.Errorf("Usage string should contain command '%s'", CommandMoo)
	}
	if !strings.Contains(usage, cow.ActionThink) {
		t.Errorf("Usage string should contain action '%s'", cow.ActionThink)
	}
	if !strings.Contains(usage, ActionSurprise) {
		t.Errorf("Usage string should contain action '%s'", ActionSurprise)
	}
	if !strings.Contains(usage, "cow") || !strings.Contains(usage, "mood") || !strings.Contains(usage, "message") {
		t.Error("Usage string should contain parameters: cow, mood, message")
	}
}

func TestGetHelpString(t *testing.T) {
	help := GetHelpString()

	// Verify structure
	if !strings.Contains(help, "Cows:") {
		t.Error("Help string should contain 'Cows:' section")
	}
	if !strings.Contains(help, "Moods:") {
		t.Error("Help string should contain 'Moods:' section")
	}

	// Verify some known cows are present
	knownCows := []string{"default", "dragon", "apt"}
	for _, cowName := range knownCows {
		if !strings.Contains(help, "`"+cowName+"`") {
			t.Errorf("Help should contain cow '%s'", cowName)
		}
	}

	// Verify some known moods are present
	knownMoods := []string{"borg", "dead", "wired"}
	for _, mood := range knownMoods {
		if !strings.Contains(help, "`"+mood+"`") {
			t.Errorf("Help should contain mood '%s'", mood)
		}
	}

	// Verify "random" option is present
	if !strings.Contains(help, "`random`") {
		t.Error("Help should contain 'random' option")
	}

	// Verify multiline structure (wrapped text)
	if strings.Count(help, "\n") < 2 {
		t.Error("Help should be multiline with wrapped text")
	}
}

func Test_wrapString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		limit     int
		assertion func(*testing.T, string)
	}{
		{
			name:  "empty_string_returns_empty",
			input: "        ",
			limit: 3,
			assertion: func(t *testing.T, result string) {
				if result != "" {
					t.Error("Empty/whitespace input should return empty string")
				}
			},
		},
		{
			name:  "wraps_at_word_boundaries",
			input: "All you see and all you touch is all you ever be",
			limit: 3,
			assertion: func(t *testing.T, result string) {
				lines := strings.Split(strings.TrimSpace(result), "\n")
				// Should have multiple lines
				if len(lines) < 2 {
					t.Error("Should wrap text into multiple lines")
				}
				// Each line should have at most 3 words
				for _, line := range lines {
					words := strings.Fields(line)
					if len(words) > 3 {
						t.Errorf("Line has %d words, expected max 3: %q", len(words), line)
					}
				}
				// All original words should be present
				allWords := strings.Join(lines, " ")
				if !strings.Contains(allWords, "All you see") {
					t.Error("Original text should be preserved")
				}
			},
		},
		{
			name:  "limit_zero_keeps_single_line",
			input: "short text",
			limit: 0,
			assertion: func(t *testing.T, result string) {
				lines := strings.Split(strings.TrimSpace(result), "\n")
				if len(lines) != 1 {
					t.Error("Limit 0 should keep text on single line")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapString(tt.input, tt.limit)
			tt.assertion(t, result)
		})
	}
}

func Test_getCows(t *testing.T) {
	t.Run("sorted_list", func(t *testing.T) {
		cows := getCows(true)

		// Should have minimum expected count (41 cows + random)
		if len(cows) < 40 {
			t.Errorf("Expected at least 40 cows, got %d", len(cows))
		}

		// Should include random option
		hasRandom := false
		for _, c := range cows {
			if c == "`random`" {
				hasRandom = true
				break
			}
		}
		if !hasRandom {
			t.Error("Cow list should include 'random' option")
		}

		// Should be formatted with backticks
		for _, c := range cows {
			if !strings.HasPrefix(c, "`") || !strings.HasSuffix(c, "`") {
				t.Errorf("Cow name should be wrapped in backticks: %s", c)
			}
		}

		// Should be sorted when sorted=true
		for i := 1; i < len(cows); i++ {
			if cows[i-1] > cows[i] {
				t.Error("Cows should be sorted alphabetically")
				break
			}
		}
	})

	t.Run("unsorted_list", func(t *testing.T) {
		cows := getCows(false)
		if len(cows) < 40 {
			t.Errorf("Expected at least 40 cows, got %d", len(cows))
		}
	})
}

func Test_getMoods(t *testing.T) {
	t.Run("sorted_list", func(t *testing.T) {
		moods := getMoods(true)

		// Should have minimum expected count (7 moods + random)
		if len(moods) < 7 {
			t.Errorf("Expected at least 7 moods, got %d", len(moods))
		}

		// Should include random option
		hasRandom := false
		for _, m := range moods {
			if m == "`random`" {
				hasRandom = true
				break
			}
		}
		if !hasRandom {
			t.Error("Mood list should include 'random' option")
		}

		// Should be formatted with backticks
		for _, m := range moods {
			if !strings.HasPrefix(m, "`") || !strings.HasSuffix(m, "`") {
				t.Errorf("Mood name should be wrapped in backticks: %s", m)
			}
		}

		// Should be sorted when sorted=true
		for i := 1; i < len(moods); i++ {
			if moods[i-1] > moods[i] {
				t.Error("Moods should be sorted alphabetically")
				break
			}
		}
	})

	t.Run("unsorted_list", func(t *testing.T) {
		moods := getMoods(false)
		if len(moods) < 7 {
			t.Errorf("Expected at least 7 moods, got %d", len(moods))
		}
	})
}

func TestGetBanner(t *testing.T) {
	// Set up test environment
	testVersion := "v1.2.3-test"
	t.Setenv(FieldEnv, "testing")

	banner := GetBanner(testVersion)

	// Verify banner structure: should contain all key sections
	sections := []string{
		"gowsay",           // App name
		testVersion,        // Version
		"testing",          // Environment
		"Usage:",           // Usage section
		"Cows:",            // Cow list
		"Moods:",           // Mood list
	}

	for _, section := range sections {
		if !strings.Contains(banner, section) {
			t.Errorf("Banner should contain '%s'", section)
		}
	}

	// Verify multiline structure
	if strings.Count(banner, "\n") < 5 {
		t.Error("Banner should have multiple lines with usage and help info")
	}

	// Verify format: should have header line with version and env in brackets
	lines := strings.Split(banner, "\n")
	if len(lines) < 1 {
		t.Fatal("Banner should have at least one line")
	}

	firstLine := lines[0]
	if !strings.Contains(firstLine, "["+testVersion+"]") {
		t.Error("First line should contain version in brackets")
	}
	if !strings.Contains(firstLine, "[testing]") {
		t.Error("First line should contain environment in brackets")
	}
}
