package server

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type client struct {
	id         uuid.UUID
	conn       *websocket.Conn
	room, name string
	sentInfo   bool
}

// sends a message to the client
func (c *client) send(m *message) error {
	return c.conn.WriteJSON(m)
}
