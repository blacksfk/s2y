package spotify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	TOKEN_URL    = "https://accounts.spotify.com/api/token"
	PLAYLIST_URL = "https://api.spotify.com/v1/playlists/"

	// Request only the required information from the API.
	PLAYLIST_FIELDS = "tracks.items(track(name,artists(name)))"
)

// Successful authentication response.
type authResponse struct {
	Access_token, Token_type string
	Expires_in               int64
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

// Public type to be used to get a playlist.
type Client struct {
	token string
	// expires time.Time
}

// Create a new spotify client that implements the client credentials flow.
// For more information see: https://developer.spotify.com/documentation/general/guides/authorization-guide/#client-credentials-flow
func NewClient(id, secret string) (*Client, error) {
	// custom request with urlencoded body
	buf := bytes.NewBufferString("grant_type=client_credentials")
	req, e := http.NewRequest(http.MethodPost, TOKEN_URL, buf)

	if e != nil {
		return nil, e
	}

	// encode credentials in base64 (with a colon between)
	str := base64.StdEncoding.EncodeToString([]byte(id + ":" + secret))

	// add the custom, base64-encoded, auth header
	req.Header.Add("Authorization", "Basic "+str)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// send the request
	ar := &authResponse{}
	e = performReq(req, ar)

	if e != nil {
		return nil, e
	}

	return &Client{ar.Access_token}, nil
}

// Get a playlist (public or private) by the provided ID.
func (c *Client) Playlist(id string) (*Playlist, error) {
	// TODO: check validity of token

	// construct a URL with query parameters (limiting returned fields)
	url := strings.Builder{}
	url.WriteString(PLAYLIST_URL)
	url.WriteString(id)
	url.WriteString("?fields=")
	url.WriteString(PLAYLIST_FIELDS)

	req, e := http.NewRequest(http.MethodGet, url.String(), nil)

	if e != nil {
		return nil, e
	}

	// add the authorisation token
	req.Header.Add("Authorization", "Bearer "+c.token)

	// send the request
	list := &Playlist{}
	e = performReq(req, list)

	return list, e
}

// Creates a client and sends the request. v should be a pointer.
func performReq(req *http.Request, v interface{}) error {
	client := &http.Client{}
	res, e := client.Do(req)

	if e != nil {
		return e
	}

	body, e := io.ReadAll(res.Body)

	if e != nil {
		return e
	}

	if res.StatusCode >= 400 {
		// something went wrong with the request
		ae := &APIError{}
		e = json.Unmarshal(body, ae)

		if e != nil {
			// busted JSON
			return e
		}

		return ae
	}

	// request was successful
	// unmarshal into the provided pointer
	return json.Unmarshal(body, v)
}

// Spotify API PlaylistObject.
type Playlist struct {
	Tracks struct {
		Items []struct {
			Track Track
		}
	}
}

// Unwind the spotify tracks.items.track structure into something usable.
func (l Playlist) GetTracks() []Track {
	var tracks []Track

	for _, item := range l.Tracks.Items {
		tracks = append(tracks, item.Track)
	}

	return tracks
}

// Spotify API TrackObject.
type Track struct {
	Name    string
	Artists []Artist
}

// Spotify API ArtistObject.
type Artist struct {
	Name string
}
