package backend

import uuid "github.com/satori/go.uuid"

// A Client represents a client connected to the
// server. Whether he has sent client info or not
// is specified by the connected flag, and if false,
// Name and Room probably aren't set.
type Client struct {
	id        uuid.UUID
	connected bool

	Name string `json:"name"`
	Room string `json:"room"`
}

// Send sends a message to the client, by adding
// the message to the outgoing queue.
func (c *Client) Send(b *Backend, m *Message) {
	m.Client = c
	b.Outgoing <- m
}
