package api

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

func getRandomMood() string {
	keys := reflect.ValueOf(moods).MapKeys()
	return keys[rand.Intn(len(keys))].Interface().(string)
}
