package server

import (
	"fmt"

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

// validateName checks if a name is valid,
// i.e. if it is the right length, matches the
// regex, and isn't already taken.
func validateName(name string, s *Server) (reason string, ok bool) {
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

// validateRoom checks if a room is valid
func validateRoom(room string) (reason string, ok bool) {
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
