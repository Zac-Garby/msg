package server

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

func broadcastRoom(s *Server, room string, m *message) {
	for _, client := range s.clients {
		if client.Room != room {
			continue
		}

		client.send(m)
	}
}
