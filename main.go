package main

import (
	"fmt"

	"github.com/vnykmshr/gowsay/src/cows"
	"github.com/vnykmshr/gowsay/src/gowsay"
	"github.com/vnykmshr/gowsay/src/help"
)

func main() {
	cow, err := cows.GetRandomCow()
	if err != nil {
		fmt.Printf("Error getting random cow: %v\n", err)
		return
	}

	fmt.Println(cow.Art)

	// cow2, err := cows.GetCow("test")
	// if err != nil {
	// 	fmt.Printf("Error getting cow: %v\n", err)
	// 	return
	// }

	// fmt.Println(cow2.Art)

	gowsay, err := gowsay.GetGowsay("help", "random", "random", 40, []string{"Hello World!"})
	fmt.Println(gowsay, err)

	fmt.Println(help.GetBanner("main"))
}
