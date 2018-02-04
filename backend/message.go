package backend

// A Message represents a message, incoming or outgoing. If
// it's incoming, client is the client who sent the message.
// If it's outgoing, client is the client to send it to.
type Message struct {
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
	Client *Client     `json:"-"`
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
			client.Send(m)
		}
	}
}
