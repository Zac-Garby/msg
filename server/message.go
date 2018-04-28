package server

import (
	uuid "github.com/satori/go.uuid"
)

type Message struct {
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
	sender Client
}

func broadcast(s *Server, m *Message) {
	for _, client := range s.clients {
		client.Send(m)
	}
}

func broadcastRoom(s *Server, room string, m *Message, blacklist ...uuid.UUID) {
outer:
	for _, client := range s.clients {
		if client.RoomName() != room {
			continue
		}

		for _, i := range blacklist {
			if i == client.ID() {
				continue outer
			}
		}

		client.Send(m)
	}
}
