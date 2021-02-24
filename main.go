package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/blacksfk/s2y/config"
)

func main() {
	// load configuration JSON file
	var file string

	flag.StringVar(&file, "config", "config.json", "JSON configuration file")
	flag.Parse()

	conf, e := config.Load(file)

	if e != nil {
		log.Fatal(e)
	}

	fmt.Printf("%+v\n", conf)
}
