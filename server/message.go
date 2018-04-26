package server

import (
	uuid "github.com/satori/go.uuid"
)

type message struct {
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
	sender *client
}

func broadcast(s *Server, m *message) {
	for _, client := range s.clients {
		client.send(m)
	}
}

func broadcastRoom(s *Server, room string, m *message, blacklist ...uuid.UUID) {
outer:
	for _, client := range s.clients {
		if client.Room != room {
			continue
		}

		for _, i := range blacklist {
			if i == client.id {
				continue outer
			}
		}

		client.send(m)
	}
}
