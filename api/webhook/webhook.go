package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/blacksfk/s2y"
	"github.com/blacksfk/s2y/api"
)

const (
	URL = "https://discord.com/api/v8/webhooks/"
)

type EditMsg struct {
	Content string `json:"content"`
}

type Client struct {
	AppID string
}

// Edit the already sent message.
func (c Client) EditOriginal(i s2y.Interaction, msg EditMsg) error {
	// build the URL string
	url := strings.Builder{}

	url.WriteString(URL)
	url.WriteString(c.AppID)
	url.WriteString("/")
	url.WriteString(i.Token)
	url.WriteString("/messages/@original")

	// marshal the body
	body, e := json.Marshal(msg)

	if e != nil {
		return e
	}

	// create the request
	req, e := http.NewRequest(http.MethodPatch, url.String(), bytes.NewReader(body))

	if e != nil {
		return e
	}

	// set custom headers
	req.Header.Add("Content-Type", "application/json")

	// send the request
	_, e = api.PerformReq(req)

	return e
}
