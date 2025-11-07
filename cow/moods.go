package cow

import (
	"math/rand"
	"reflect"
)

var moods map[string]int

func init() {
	moods = map[string]int{
		"borg":     1,
		"dead":     1,
		"greedy":   1,
		"paranoid": 1,
		"stoned":   1,
		"wired":    1,
		"young":    1,
	}
}

// RandomMood returns a random mood name
func RandomMood() string {
	keys := reflect.ValueOf(moods).MapKeys()
	return keys[rand.Intn(len(keys))].Interface().(string)
}

// ListMoods returns a list of all available mood names
func ListMoods() []string {
	names := make([]string, 0, len(moods))
	for name := range moods {
		names = append(names, name)
	}
	return names
}

// MoodExists checks if a mood with the given name exists
func MoodExists(name string) bool {
	_, ok := moods[name]
	return ok
}
