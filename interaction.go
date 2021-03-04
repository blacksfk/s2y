package s2y

// See: https://golang.org/ref/spec#Constant_declarations
// and: https://golang.org/ref/spec#Iota
// and: https://golang.org/doc/effective_go.html#constants
const (
	// Discord interaction response types.
	// See: https://discord.com/developers/docs/interactions/slash-commands#interaction-response-interactionresponsetype
	pong = iota + 1
	ack
	msg
	msgSrc
	ackSrc

	// Discord interaction request types.
	ping = pong
	cmd  = ack
)

// Represents a Discord interaction request.
type Interaction struct {
	ID    string
	Type  int64
	Token string
	Data  CommandData
}

// Whether or not the interaction is a ping request.
func (i Interaction) IsPing() bool {
	return i.Type == ping
}

// Represents the response to a Discord interaction request.
type InteractionResponse struct {
	Type int                     `json:"type"`
	Data InteractionResponseData `json:"data"`
}

func NewPongResponse() InteractionResponse {
	return InteractionResponse{Type: pong}
}

func NewAckResponse() InteractionResponse {
	return InteractionResponse{Type: ack}
}

func NewMsgResponse(content string) InteractionResponse {
	return InteractionResponse{
		Type: msg,
		Data: InteractionResponseData{
			Tts:     false,
			Content: content,
		},
	}
}

func NewMsgWithSrcResponse(i Interaction) InteractionResponse {
	return InteractionResponse{Type: msgSrc}
}

func NewAckWithSrcResponse(i Interaction) InteractionResponse {
	return InteractionResponse{Type: ackSrc}
}

type InteractionResponseData struct {
	Tts     bool   `json:"tts"`
	Content string `json:"content"`
}
