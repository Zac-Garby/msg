package server

import "log"

func (s *Server) handleOutgoing() {
	for {
		var (
			msg      = <-s.backend.Outgoing
			id       = msg.Client.ID
			conn, ok = s.connections[id]
		)

		if !ok {
			log.Println("attempted to send a message to a client, but no connection was found")
			continue
		}

		if err := conn.WriteJSON(msg); err != nil {
			log.Println("got error when sending message:", err)
			continue
		}
	}
}
