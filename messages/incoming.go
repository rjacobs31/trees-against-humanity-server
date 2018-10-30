package messages

// IncomingMessageType is the type of a message received
// from a client.
type IncomingMessageType int

const (
	// Connect is a connection attempt.
	Connect IncomingMessageType = iota

	// Disconnect is a disconnect attempt.
	Disconnect

	// CreateGame is an attempt to create a new game.
	CreateGame

	// JoinGame is an attempt to join an existing game.
	JoinGame

	// LeaveGame is an attempt to leave a joined game.
	LeaveGame
)

// IncomingMessage is an incoming message from a client.
type IncomingMessage struct {
	Type IncomingMessageType `json:"type"`
	Data interface{}         `json:"data,omitempty"`
}
