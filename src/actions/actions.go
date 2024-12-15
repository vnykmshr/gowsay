package actions

import (
	"math/rand"
	"time"

	"github.com/vnykmshr/gowsay/src/constants"
)

var (
	actions = []string{
		constants.ActionSay,
		constants.ActionHelp,
		constants.ActionList,
		constants.ActionThink,
		constants.ActionSurprise,
	}

	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func GetActionNames() []string {
	return actions
}

func GetRandomAction() string {
	return actions[rnd.Intn(len(actions))]
}
