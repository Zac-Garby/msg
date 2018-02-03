package server

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type client struct {
	id       uuid.UUID
	conn     *websocket.Conn
	sentInfo bool

	Room string `json:"room"`
	Name string `json:"name"`
}

// sends a message to the client
func (c *client) send(m *message) error {
	return c.conn.WriteJSON(m)
}
