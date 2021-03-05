package main

import (
	"github.com/blacksfk/s2y/api/spotify"
	"github.com/blacksfk/s2y/api/webhook"
	"github.com/blacksfk/s2y/config"
)

type services struct {
	whc webhook.Client
	sc  *spotify.Client
}

func initServices(conf *config.Config) (*services, error) {
	sc, e := spotify.NewClient(conf.Spotify.ID, conf.Spotify.Secret)

	if e != nil {
		return nil, e
	}

	return &services{webhook.Client{conf.Discord.ID}, sc}, nil
}
