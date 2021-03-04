package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/blacksfk/s2y"
)

const (
	URL = "https://discord.com/api/v8/webhooks/"
)

type EditMsg struct {
	Content string `json:"content"`
}

// Something went wrong with the request or the server.
type APIError struct {
	Message string
	Status  int
}

// APIError implements error.
func (ae *APIError) Error() string {
	return fmt.Sprintf("%d: %s", ae.Status, ae.Message)
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

	// perform the request with the default client
	res, e := http.DefaultClient.Do(req)

	if e != nil {
		return e
	}

	// TODO: refactor with spotify package request handling
	if res.StatusCode >= 400 {
		// something went wrong with the request or server
		ae := &APIError{}

		defer res.Body.Close()
		body, e = io.ReadAll(res.Body)

		if e != nil {
			return e
		}

		e = json.Unmarshal(body, ae)

		if e != nil {
			return e
		}

		return ae
	}

	// if nothing went wrong then the edit request was successful
	return nil
}
