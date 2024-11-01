package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/bmishra/gowsay/api"
	"gopkg.in/tokopedia/grace.v1"
	"gopkg.in/tokopedia/logging.v1"
)

var version = "devel"

func main() {
	flag.Parse()

	// XXX: set logging.buildHash for app version
	logging.LogInit()

	m, err := api.NewModule()
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/", m.Gowsay)
	http.HandleFunc("/say", m.Gowsay)

	fmt.Println(api.GetBanner(version))
	log.Fatal(grace.Serve(":9000", nil))
}
