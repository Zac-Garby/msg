package server

type message struct {
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
	sender *client
}

func broadcast(s *Server, m *message) error {
	for _, client := range s.clients {
		if err := client.send(m); err != nil {
			return err
		}
	}

	return nil
}

func broadcastRoom(s *Server, room string, m *message) error {
	for _, client := range s.clients {
		if client.Room != room {
			continue
		}

		if err := client.send(m); err != nil {
			return err
		}
	}

	return nil
}
