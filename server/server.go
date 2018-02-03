package server

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var (
	maxNameLength = 32
	minNameLength = 2
	usernameRegex = regexp.MustCompile(`^[\p{L}\p{N}-_.]+$`)

	maxRoomLength = 64
	minRoomLength = 1
	roomNameRegex = regexp.MustCompile(`^[\p{L}\p{N}-_./<>&]+$`)
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

				if reason, ok := validateName(name, s); !ok {
					out := serverMessage(fmt.Sprintf("Your username is invalid (%s)", reason))
					if err := msg.sender.send(out); err != nil {
						log.Println("error when sending invalid username msg:", err)
					}

					break
				}

				if reason, ok := validateRoom(room); !ok {
					out := serverMessage(fmt.Sprintf("Your room name is invalid (%s)", reason))
					if err := msg.sender.send(out); err != nil {
						log.Println("error when sending invalid room msg:", err)
					}

					break
				}

				msg.sender.Name = name
				msg.sender.Room = room
				msg.sender.sentInfo = true

				out := serverMessage(fmt.Sprintf(`Hello - welcome to the server, %s.
Type 'help' to view the available commands.
You can change your name via the '/name [name]' command.`, name))

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
					s.handleCommand(msg.sender, line)
				}

				break
			}

			if str[0] == '/' {
				s.handleCommand(msg.sender, str[1:])
				break
			}

			broadcastRoom(s, msg.sender.Room, &message{
				Type: "chat",
				Data: map[string]interface{}{
					"sender": msg.sender,
					"text":   str,
				},
			})
		}
	}
}

// checkName checks if a client exists with a
// given name
func (s *Server) checkName(name string) bool {
	for _, c := range s.clients {
		if c.Name == name {
			return true
		}
	}

	return false
}

// usersInRoom gets a list of the usernames of
// the users in a given room.
func (s *Server) usersInRoom(room string) []string {
	var names []string

	for _, c := range s.clients {
		if c.Room == room {
			names = append(names, c.Name)
		}
	}

	return names
}

func (s *Server) handleCommand(sender *client, str string) {
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
			out = cmd(s, sender, split)
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
