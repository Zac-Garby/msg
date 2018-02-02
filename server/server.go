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

				out := serverMessage(fmt.Sprintf("Hello - welcome to the server, %s.\nType `/help` to view the available commands", name))

				if err := msg.sender.send(out); err != nil {
					log.Println("error when sending welcome msg:", err)
				}
			}

		case "chat":
			if !msg.sender.sentInfo {
				log.Println("a client tried to send a message, but hadn't sent client info beforehand")
				break
			}

			str, ok := msg.Data.(string)
			if !ok {
				log.Println("client", msg.sender.id, "tried to sender a non-string message")
				break
			}

			if strings.HasPrefix(str, "/script") {
				lines := strings.Split(str, "\n")

				if len(lines) <= 1 {
					break
				}

				lines = lines[1:]

				for _, line := range lines {
					handleCommand(msg.sender, line)
				}

				break
			}

			if str[0] == '/' {
				handleCommand(msg.sender, str[1:])
				break
			}

			broadcast(s, &message{
				Type: "chat",
				Data: map[string]string{
					"sender": msg.sender.name,
					"text":   str,
				},
			})
		}
	}
}

func handleCommand(sender *client, str string) {
	var (
		out   string
		split = strings.Fields(str)
	)

	if len(split) == 0 || len(str) == 0 {
		out = "expected a command after `/`"
	} else {
		name := split[0]

		cmd, ok := commands[name]
		if !ok {
			out = fmt.Sprintf("command not found: %s", name)
		} else {
			out = cmd(sender, split)
		}
	}

	if out == "" {
		return
	}

	if err := sender.send(serverMessage(out)); err != nil {
		log.Println("handleCommand: errored when sending output message")
	}
}

func serverMessage(content string) *message {
	return &message{
		Type: "server-msg",
		Data: content,
	}
}
