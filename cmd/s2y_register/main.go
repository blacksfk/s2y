package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/blacksfk/s2y"
	"github.com/blacksfk/s2y/config"
)

const (
	URL = "https://discord.com/api/v8/applications/"
)

func main() {
	file := flag.String("config", "config.json", "JSON configuration file")
	guildID := flag.String("guild", "", "Guild ID to test commands in")

	// parse the flags
	flag.Parse()

	conf, e := config.Load(*file)

	if e != nil {
		log.Fatal(e)
	}

	// command definitions
	cmd := s2y.NewCommand(
		"s2y",
		"Convert a spotify playlist into a public youtube playlist.",
		s2y.CommandOption{
			Type:        s2y.CMD_TYPE_STR,
			Name:        "playlist",
			Description: "Spotify playlist. Can be a complete URL, the Spotify URI, or just the playlist ID.",
			Required:    true,
		},
	)

	// marshal the body
	body, e := json.Marshal(cmd)

	if e != nil {
		log.Fatal(e)
	}

	// construct the URL
	url := strings.Builder{}
	url.WriteString(URL)
	url.WriteString(conf.Discord.ID)

	if len(*guildID) > 0 {
		// guild ID was provided for testing purposes
		url.WriteString("/guilds/")
		url.WriteString(*guildID)
	}

	url.WriteString("/commands")

	// create the request
	req, e := http.NewRequest(http.MethodPost, url.String(), bytes.NewReader(body))

	if e != nil {
		log.Fatal(e)
	}

	// attach required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bot "+conf.Discord.Token)

	// create the client and send the request
	client := http.Client{}
	res, e := client.Do(req)

	if e != nil {
		// network failure
		log.Fatal(e)
	}

	// unmarshal the response body
	body, e = io.ReadAll(res.Body)
	defer res.Body.Close()

	if e != nil {
		log.Fatal(e)
	}

	fmt.Println(string(body))
}
