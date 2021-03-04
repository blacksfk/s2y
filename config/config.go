package config

import (
	"encoding/hex"
	"encoding/json"
	"os"
)

// Discord configuration.
type Discord struct {
	// Application ID
	ID string

	// Application public key
	PublicKey string

	// Application public key in byte form
	// Converted from PublicKey (string) during load time
	publicKeyBytes []byte

	// Application bot token (command creation authorisation)
	Token string
}

// Get the public key in byte slice form.
func (d Discord) PublicKeyBytes() []byte {
	return d.publicKeyBytes
}

// Spotify configuration.
type Spotify struct {
	ID, Secret string
}

// Wrapper struct around Discord and Spotify configuration structs.
type Config struct {
	Address string
	Discord Discord
	Spotify Spotify
}

// Load the contents of file as JSON into a Config struct.
func Load(file string) (*Config, error) {
	bytes, e := os.ReadFile(file)

	if e != nil {
		return nil, e
	}

	conf := &Config{}
	e = json.Unmarshal(bytes, conf)

	if e != nil {
		return nil, e
	}

	// convert hex string into byte slice
	conf.Discord.publicKeyBytes, e = hex.DecodeString(conf.Discord.PublicKey)

	return conf, e
}
