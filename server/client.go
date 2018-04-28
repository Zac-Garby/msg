package server

import (
	"fmt"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// A Client is an entity which acts like a messaging client. It will usually be a user connected
// via the website, but could also be the server, in which case Send() won't actually use a
// network.
type Client interface {
	Send(m *Message) error
	ID() uuid.UUID
	SentInfo() bool
	InfoReceived()
	RoomName() string
	Username() string
	GotoRoom(room string)
	Rename(name string)
}

type client struct {
	id       uuid.UUID
	conn     *websocket.Conn
	sentInfo bool

	Room string `json:"room"`
	Name string `json:"name"`
}

// Send sends a message to the client
func (c *client) Send(m *Message) error {
	return c.conn.WriteJSON(m)
}

// ID returns the UUID of the client
func (c *client) ID() uuid.UUID { return c.id }

// SentInfo checks whether or not the client has joined properly
func (c *client) SentInfo() bool { return c.sentInfo }

// InfoReceived sets the SentInfo flag to true
func (c *client) InfoReceived() { c.sentInfo = true }

// RoomName gets the name of the room which the client is in
func (c *client) RoomName() string { return c.Room }

// Username gets the name of the client
func (c *client) Username() string { return c.Name }

// GotoRoom causes the client to go to the specified room
func (c *client) GotoRoom(room string) { c.Room = room }

// Rename changes the name of client
func (c *client) Rename(name string) { c.Name = name }

// ValidateName checks if a name is valid,
// i.e. if it is the right length, matches the
// regex, and isn't already taken.
func ValidateName(name string, s *Server) (reason string, ok bool) {
	if len(name) > maxNameLength {
		return fmt.Sprintf("Your name cannot be longer than %d characters", maxNameLength), false
	}

	if len(name) < minNameLength {
		return fmt.Sprintf("Your name cannot be less than %d characters", minNameLength), false
	}

	if !usernameRegex.MatchString(name) {
		return "Your username must contain only letters, numbers, hyphens, underscores, and dots", false
	}

	if s.checkName(name) {
		return fmt.Sprintf("A user already exists called %s!", name), false
	}

	return "", true
}

// ValidateRoom checks if a room is valid
func ValidateRoom(room string) (reason string, ok bool) {
	if len(room) > maxRoomLength {
		return fmt.Sprintf("A room name cannot be longer than %d characters", maxNameLength), false
	}

	if len(room) < minRoomLength {
		return fmt.Sprintf("A room name cannot be less than %d characters", minNameLength), false
	}

	if !roomNameRegex.MatchString(room) {
		return "A room name must only contain letters, numbers, and any of: -_./<>&", false
	}

	return "", true
}
