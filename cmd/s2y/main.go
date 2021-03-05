package main

import (
	"flag"
	"log"

	uf "github.com/blacksfk/microframework"
	"github.com/blacksfk/s2y/config"
	"github.com/blacksfk/s2y/http"
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

	services, e := initServices(conf)

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
	routes(s, conf, services)

	// anchors aweigh
	log.Println(s.Start())
}

func routes(s *uf.Server, conf *config.Config, services *services) {
	ctrl := http.NewController(services.sc, services.whc, conf.Discord.PublicKeyBytes())

	s.Post("/", ctrl.Index)
}

// Error logger supplied to the framework.
// TODO: log to file or some service
func logError(e error) {
	log.Println(e)
}
