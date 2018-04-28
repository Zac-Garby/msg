package server

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// A serverClient implements Client and acts very much like a regular client, but is controlled by the
// server instead. This means that, while running a server, you will be able to write messages directly
// into the terminal.
type serverClient struct {
	server *Server
	id     uuid.UUID
	Room   string `json:"room"`
	Name   string `json:"name"`
}

// Send sends a message to the client
func (s *serverClient) Send(m *Message) error {
	switch m.Type {
	case "chat":
		data, ok := m.Data.(map[string]interface{})
		if !ok {
			fmt.Println("invalid message data:", m.Data)
			return nil
		}

		genericSender, ok := data["sender"]
		if !ok {
			fmt.Println("sender key non-existant in", data)
			return nil
		}

		sender, ok := genericSender.(Client)
		if !ok {
			fmt.Println("sender not a client in", data)
			return nil
		}

		text, ok := data["text"]
		if !ok {
			fmt.Println("text key non-existant in", data)
			return nil
		}

		fmt.Printf("[%s] %s\n", sender.Username(), text)

	case "server-msg":
		fmt.Println("~~~", m.Data)

	case "quit":
		fmt.Println("cannot quit using /quit. use ctrl+c instead")
	}

	return nil
}

// ID returns the UUID of the client
func (s *serverClient) ID() uuid.UUID { return s.id }

// SentInfo checks whether or not the client has joined properly
func (s *serverClient) SentInfo() bool { return true }

// InfoReceived sets the SentInfo flag to true
func (s *serverClient) InfoReceived() {}

// RoomName gets the name of the room which the client is in
func (s *serverClient) RoomName() string { return s.Room }

// Username gets the name of the client
func (s *serverClient) Username() string { return s.Name }

// GotoRoom causes the client to go to the specified room
func (s *serverClient) GotoRoom(room string) { s.Room = room }

// Rename changes the name of client
func (s *serverClient) Rename(name string) { s.Name = name }
