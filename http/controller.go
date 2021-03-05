package http

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	uf "github.com/blacksfk/microframework"
	"github.com/blacksfk/s2y"
	"github.com/blacksfk/s2y/api/spotify"
	"github.com/blacksfk/s2y/api/webhook"
)

const (
	// Match spotify ID strings
	ID_MATCH = `([a-zA-Z0-9]+)`

	// Timestamp header name
	HEADER_TS = "X-Signature-Timestamp"

	// Signature header name
	HEADER_SIG = "X-Signature-Ed25519"
)

var (
	reSlice []*regexp.Regexp = []*regexp.Regexp{
		// Match complete spotify URLs
		regexp.MustCompile(`playlist/` + ID_MATCH),

		// Match spotify URIs
		regexp.MustCompile(`spotify:playlist:` + ID_MATCH),

		// Match an ID string
		regexp.MustCompile(ID_MATCH),
	}
)

type Controller struct {
	sc        *spotify.Client
	whc       webhook.Client
	publicKey ed25519.PublicKey
}

// Create a new controller providing: a spotify client, a webhook client, and the discord
// application's public key.
func NewController(sc *spotify.Client, whc webhook.Client, publicKey []byte) Controller {
	return Controller{sc, whc, ed25519.PublicKey(publicKey)}
}

// Handle incoming interactions. Each interaction is verified with the provided public key
// via the ED25519 package.
func (c Controller) Index(w http.ResponseWriter, r *http.Request) error {
	// we need the raw body for signature verification
	body, e := uf.ReadBody(r, "application/json")

	if e != nil {
		return e
	}

	// get the timestamp and signature from the headers as byte slices
	ts := []byte(r.Header.Get(HEADER_TS))
	sig, e := hex.DecodeString(r.Header.Get(HEADER_SIG))

	if e != nil {
		return e
	}

	// verify the request so we don't get banned by discord
	// append the raw body to the timestamp
	if !ed25519.Verify(c.publicKey, append(ts, body...), sig) {
		// verification failed - send 401
		return uf.Unauthorized("Invalid signature")
	}

	// unmarshal the body into an interaction
	interaction := s2y.Interaction{}

	if e = json.Unmarshal(body, &interaction); e != nil {
		return e
	}

	if interaction.IsPing() {
		// ping received, return a pong
		return uf.SendJSON(w, s2y.NewPongResponse())
	}

	// extract the ID from the value string
	var id string
	v := interaction.Data.Options[0].Value

	for _, re := range reSlice {
		if re.MatchString(v) {
			// [0] = entire string (v)
			// [1] = first capture group (id)
			id = re.FindStringSubmatch(v)[1]
			break
		}
	}

	if len(id) == 0 {
		// no match found
		str := fmt.Sprintf("Invalid input: %s (%d characters)", v, len(v))

		return uf.SendJSON(w, s2y.NewMsgResponse(str))
	}

	// match found - reply to maintain validity of the token
	e = uf.SendJSON(w, s2y.NewMsgResponse("Searching for playlist ID: "+id))

	if e != nil {
		return e
	}

	// get the playlist
	playlist, e := c.sc.Playlist(id)

	if e != nil {
		// something went wrong with the request
		str := fmt.Sprintf("Error: %s", e)

		return uf.SendJSON(w, s2y.NewMsgResponse(str))
	}

	// TODO: search for the tracks on youtube and create playlist with
	// each found track
	tracks := playlist.GetTracks()
	length := len(tracks)
	list := strings.Builder{}

	for i := 0; i < length; i++ {
		list.WriteString(tracks[i].Name)
		list.WriteString("-")
		list.WriteString(tracks[i].Artists[0].Name)

		if i < length-1 {
			// not the last track, so add a comma
			list.WriteString(",")
		}
	}

	return c.whc.EditOriginal(interaction, webhook.EditMsg{list.String()})
}
