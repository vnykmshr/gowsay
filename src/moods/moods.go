package moods

import (
	"math/rand"
	"time"
)

var (
	moods []string
	rnd   = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func init() {
	moods = []string{
		"borg",
		"dead",
		"greedy",
		"paranoid",
		"stoned",
		"wired",
		"young",
	}
}

func GetMoods() []string {
	return moods
}

func GetRandomMood() string {
	return moods[rnd.Intn(len(moods))]
}
