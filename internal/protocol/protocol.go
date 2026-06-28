package protocol

import (
	"encoding/json"
	"time"
)

// tells the receiver how to interpret Envelope.Data
type Kind string

const (
	// chat message sent to a channel
	KindChat Kind = "chat"

	// sent by the client right after connecting to announce itself
	KindHello Kind = "hello"

	// server-originated notice (joins, leaves, errors)
	KindSystem Kind = "system"
)

// outer frame for every message in either direction
type Envelope struct {
	Type Kind            `json:"type"`
	Data json.RawMessage `json:"data"`
}

// wraps a typed payload into an Envelope ready to be written to the wire
func Encode(kind Kind, payload any) (Envelope, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return Envelope{}, err
	}
	return Envelope{Type: kind, Data: raw}, nil
}

// payload for KindChat frames
type ChatMessage struct {
	// assigned by the server so every client refers to a message the same way
	// (needed later for edits, deletes, replies). empty when sent by a client
	ID string `json:"id,omitempty"`

	// id of the channel this message belongs to (see Channel in channel.go).
	ChannelID string `json:"channelId"`

	// username of the sender
	Username string `json:"username"`

	// the message text
	Body string `json:"body"`

	// set by the server on receipt
	SentAt time.Time `json:"sentAt"`
}

// payload for KindHello frames: the client introducing itself
type Hello struct {
	Username string `json:"username"`
}

// payload for KindSystem frames: server notices
type System struct {
	Text string `json:"text"`
}
