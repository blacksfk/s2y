package main

import (
	"flag"
	"log"

	uf "github.com/blacksfk/microframework"
	"github.com/blacksfk/s2y/config"
	"github.com/blacksfk/s2y/spotify"
	"github.com/husobee/vestigo"
)

func main() {
	ptr := flag.String("config", "config.json", "JSON configuration file")

	// parse the flags
	flag.Parse()

	conf, e := config.Load(*ptr)

	if e != nil {
		log.Fatal(e)
	}

	// create a spotify client
	client, e := spotify.NewClient(conf.Spotify.ID, conf.Spotify.Secret)

	if e != nil {
		log.Fatal(e)
	}

	// cors settings
	cors := &vestigo.CorsAccessControl{
		AllowOrigin:  []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"content-type"},
	}

	// server configuration
	sConfig := &uf.Config{
		Address:      conf.Address,
		ErrorLogger:  logError,
		Cors:         cors,
		AccessLogger: uf.LogStdout,
	}

	// create server
	s := uf.NewServer(sConfig)

	// add routes
	routes(s, conf, client)

	// anchors aweigh
	log.Println(s.Start())
}

func routes(s *uf.Server, conf *config.Config, client *spotify.Client) {
	// TODO: add routes
}

// Error logger supplied to the framework.
// TODO: log to file or some service
func logError(e error) {
	log.Println(e)
}
