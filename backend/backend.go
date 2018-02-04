package backend

import (
	"fmt"
	"regexp"

	uuid "github.com/satori/go.uuid"
)

var (
	maxNameLength = 32
	minNameLength = 2
	usernameRegex = regexp.MustCompile(`^[\p{L}\p{N}-_.]+$`)

	maxRoomLength = 64
	minRoomLength = 1
	roomNameRegex = regexp.MustCompile(`^[\p{L}\p{N}-_./<>&]+$`)

	maxMessageLength = 1024
)

// The Backend keeps track of the connected users and handles the
// messages.
type Backend struct {
	Incoming, Outgoing chan *Message
	Clients            map[uuid.UUID]*Client
}

// GetUser gets a user by name.
func (b *Backend) GetUser(name string) (*Client, bool) {
	for _, client := range b.Clients {
		if client.Name == name {
			return client, true
		}
	}

	return nil, false
}

// ValidateName checks if a name is valid. A valid name is of a
// valid length (between min and maxNameLength), matches the regex
// usernameRegex, and isn't taken.
func (b *Backend) ValidateName(name string) (reason string, ok bool) {
	if len(name) > maxNameLength {
		return fmt.Sprintf("Your name cannot be longer than %d characters", maxNameLength), false
	}

	if len(name) < minNameLength {
		return fmt.Sprintf("Your name cannot be less than %d characters", minNameLength), false
	}

	if !usernameRegex.MatchString(name) {
		return "Your username must contain only letters, numbers, hyphens, underscores, and dots", false
	}

	if _, exists := b.GetUser(name); exists {
		return fmt.Sprintf("A user already exists called %s!", name), false
	}

	return "", true
}

// ValidateRoom checks if a room name is valid. A valid room name
// is of a valid length and matches the regex roomNameRegex.
func (b *Backend) ValidateRoom(room string) (reason string, ok bool) {
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
