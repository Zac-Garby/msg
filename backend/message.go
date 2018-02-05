package backend

import (
	"fmt"
	"log"
)

// The types of various messages
const (
	MsgClientInfo = "client-info"
	MsgChat       = "chat"
	MsgServer     = "server-msg"
	MsgQuit       = "quit"
)

// A Message represents a message, incoming or outgoing. If
// it's incoming, client is the client who sent the message.
// If it's outgoing, client is the client to send it to.
type Message struct {
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
	Client *Client     `json:"-"`
}

func ServerMessage(format string, args ...interface{}) *Message {
	return &Message{
		Type: MsgServer,
		Data: fmt.Sprintf(format, args...),
	}
}

// Broadcast broadcasts a message to all rooms in rooms. If
// rooms is empty, the message is broadcasted to everyone.
func (m *Message) Broadcast(b *Backend, rooms ...string) {
	for _, client := range b.Clients {
		ok := false
		for _, room := range rooms {
			if room == client.Room {
				ok = true
				break
			}
		}

		if ok || len(rooms) == 0 {
			client.Send(b, m)
		}
	}
}

// Send sends a message to a particular client.
func (m *Message) Send(b *Backend, c *Client) {
	c.Send(b, m)
}

// HandleIncoming handles incoming messages one by one, and does
// something for each one.
func (b *Backend) HandleIncoming() {
	for {
		msg := <-b.Incoming

		switch msg.Type {
		case MsgClientInfo:
			if data, ok := msg.Data.(map[string]interface{}); ok {
				name, ok := data["name"].(string)
				if !ok {
					break
				}

				room, ok := data["room"].(string)
				if !ok {
					break
				}

				go b.handleClientInfo(msg.Client, name, room)
			}

		case MsgChat:
			if !msg.Client.Connected {
				log.Println("a client tried to send a message, but hadn't sent client info beforehand")
				break
			}

			if str, ok := msg.Data.(string); ok {
				go b.handleChatMessage(msg.Client, str)
			}
		}
	}
}

func (b *Backend) handleClientInfo(sender *Client, name, room string) {
	if reason, ok := b.ValidateName(name); !ok {
		ServerMessage("Your username is invalid; %s", reason).Send(b, sender)
		return
	}

	if reason, ok := b.ValidateRoom(room); !ok {
		ServerMessage("Your room name is invalid; %s", reason).Send(b, sender)
		return
	}

	sender.Name = name
	sender.Room = room
	sender.Connected = true

	ServerMessage(`Hello, %s - welcome to the server!
Type 'help' to view the available commands.
You can go to another room using the '/room [room-name]' command.`, name).Send(b, sender)

	ServerMessage("%s has joined the server, in room '%s'", name, room).Broadcast(b)
}

func (b *Backend) handleChatMessage(sender *Client, str string) {
	if len(str) < 1 || len(str) > maxMessageLength {
		ServerMessage("Your message must be between 1 and %d characters long", maxMessageLength).Send(b, sender)
		return
	}

	if str[0] == '/' {
		b.handleCommand(sender, str[1:])
		return
	}

	msg := &Message{
		Type: MsgChat,
		Data: map[string]interface{}{
			"sender": sender,
			"text":   str,
		},
	}

	msg.Broadcast(b, sender.Room)
}
