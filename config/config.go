package config

import (
	"encoding/json"
	"io/ioutil"
)

type Discord struct {
	Token string
}

type Spotify struct {
	ID, Secret string
}

type Config struct {
	Discord Discord
	Spotify Spotify
}

func Load(file string) (*Config, error) {
	bytes, e := ioutil.ReadFile(file)

	if e != nil {
		return nil, e
	}

	conf := &Config{}

	return conf, json.Unmarshal(bytes, conf)
}
