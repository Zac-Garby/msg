package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// A Server is a websocket server which handles
// websocket connections from clients.
type Server struct {
	messages chan *message
	clients  map[uuid.UUID]*client
}

// New creates a new server.
func New() *Server {
	s := &Server{
		messages: make(chan *message, 1),
		clients:  make(map[uuid.UUID]*client),
	}

	return s
}

func (s *Server) NewClient(conn *websocket.Conn) error {
	defer conn.Close()

	id := uuid.NewV4()

	s.clients[id] = &client{
		id:       id,
		conn:     conn,
		sentInfo: false,
	}

	client := s.clients[id]

	for {
		msg := &message{
			sender: client,
		}

		if err := conn.ReadJSON(msg); err != nil {
			if strings.Contains(err.Error(), "use of closed network") ||
				strings.Contains(err.Error(), "unexpected EOF") {
				break
			}

			log.Println("server err: listen:", err)
			break
		}

		s.messages <- msg
	}

	delete(s.clients, id)

	return nil
}

func (s *Server) HandleMessages() {
	for {
		msg := <-s.messages

		switch msg.Type {
		case "client-info":
			data, ok := msg.Data.(map[string]interface{})
			if ok {
				name, ok := data["name"].(string)
				if !ok {
					break
				}

				room, ok := data["room"].(string)
				if !ok {
					break
				}

				msg.sender.name = name
				msg.sender.room = room
				msg.sender.sentInfo = true

				out := &message{
					Type: "server-msg",
					Data: fmt.Sprintf("Hello - welcome to the server, %s.\nType `/help` to view the available commands", name),
				}

				if err := msg.sender.send(out); err != nil {
					log.Println("error when sending welcome msg:", err)
				}
			}

		case "chat":
			if !msg.sender.sentInfo {
				log.Println("a client tried to send a message, but hadn't sent client info beforehand")
				break
			}

			broadcast(s, &message{
				Type: "chat",
				Data: map[string]string{
					"sender": msg.sender.name,
					"text":   msg.Data.(string),
				},
			})
		}
	}
}
