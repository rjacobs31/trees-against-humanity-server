package messages

// OutgoingMessageType is the type of a message received
// from a client.
type OutgoingMessageType int

const (
	// FullGamesList will contain the full list of available games.
	FullGamesList OutgoingMessageType = iota
)

// OutgoingMessage is an outgoing message from the server.
type OutgoingMessage struct {
	Type OutgoingMessageType `json:"type"`
	Data interface{}         `json:"data,omitempty"`
}
