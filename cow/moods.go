package cow

import "math/rand"

// Mood represents a facial expression configuration
type Mood struct {
	Eyes   string
	Tongue string
}

var moods = map[string]Mood{
	"borg":     {Eyes: "==", Tongue: "  "},
	"dead":     {Eyes: "xx", Tongue: "U "},
	"greedy":   {Eyes: "$$", Tongue: "  "},
	"paranoid": {Eyes: "@@", Tongue: "  "},
	"stoned":   {Eyes: "**", Tongue: "U "},
	"tired":    {Eyes: "--", Tongue: "  "},
	"wired":    {Eyes: "OO", Tongue: "  "},
	"young":    {Eyes: "..", Tongue: "  "},
}

var moodNames = []string{"borg", "dead", "greedy", "paranoid", "stoned", "tired", "wired", "young"}

// GetMood returns the mood configuration for the given name
func GetMood(name string) (Mood, bool) {
	mood, ok := moods[name]
	return mood, ok
}

// RandomMood returns a random mood name
func RandomMood() string {
	return moodNames[rand.Intn(len(moodNames))]
}

// ListMoods returns a list of all available mood names
func ListMoods() []string {
	result := make([]string, len(moodNames))
	copy(result, moodNames)
	return result
}

// MoodExists checks if a mood with the given name exists
func MoodExists(name string) bool {
	_, ok := moods[name]
	return ok
}
