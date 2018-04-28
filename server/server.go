package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
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

	maxMessageLength = 1024
)

// A Server is a websocket server which handles
// websocket connections from clients.
type Server struct {
	messages chan *Message
	clients  map[uuid.UUID]Client
	client   Client
}

// New creates a new server.
func New() *Server {
	s := &Server{
		messages: make(chan *Message, 1),
		clients:  make(map[uuid.UUID]Client),
	}

	s.addServerClient()

	return s
}

func (s *Server) addServerClient() {
	sc := &serverClient{
		server: s,
		id:     uuid.NewV4(),
		Room:   "/",
		Name:   "*server*",
	}

	s.clients[sc.id] = sc
	s.client = sc
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
		msg := &Message{
			sender: client,
		}

		if err := conn.ReadJSON(msg); err != nil {
			if strings.Contains(err.Error(), "use of closed network") ||
				strings.Contains(err.Error(), "unexpected EOF") ||
				strings.Contains(err.Error(), "going away") {
				break
			}

			log.Println("server err: listen:", err)
			break
		}

		s.messages <- msg
	}

	if client.SentInfo() {
		broadcast(s, serverMessage(fmt.Sprintf("%s has left the server", client.Username())))
	}

	delete(s.clients, id)

	return nil
}

func (s *Server) HandleInput(r io.Reader) {
	br := bufio.NewReader(r)

	for {
		line, err := br.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		line = strings.TrimSpace(line)
		s.handleChat(line, s.client)
	}
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

				if reason, ok := ValidateName(name, s); !ok {
					out := serverMessage(fmt.Sprintf("Your username is invalid (%s)", reason))
					if err := msg.sender.Send(out); err != nil {
						log.Println("error when sending invalid username msg:", err)
					}

					break
				}

				if reason, ok := ValidateRoom(room); !ok {
					out := serverMessage(fmt.Sprintf("Your room name is invalid (%s)", reason))
					if err := msg.sender.Send(out); err != nil {
						log.Println("error when sending invalid room msg:", err)
					}

					break
				}

				msg.sender.Rename(name)
				msg.sender.GotoRoom(room)
				msg.sender.InfoReceived()

				out := serverMessage(fmt.Sprintf(`Hello - welcome to the server, %s.
Type 'help' to view the available commands.
You can go to another room using the '/room [room-name]' command.`, name))

				if err := msg.sender.Send(out); err != nil {
					log.Println("error when sending welcome msg:", err)
				}

				broadcast(s, serverMessage(fmt.Sprintf("%s has joined the server, and is in the room: '%s'", name, room)))
			}

		case "chat":
			if !msg.sender.SentInfo() {
				log.Println("a client tried to send a message, but hadn't sent client info beforehand")
				break
			}

			str, ok := msg.Data.(string)
			if !ok {
				log.Println("client", msg.sender.ID(), "tried to send a non-string message")
				break
			}

			s.handleChat(str, msg.sender)
		}
	}
}

// handleChat handles a chat message
func (s *Server) handleChat(msg string, sender Client) {
	if len(msg) < 1 || len(msg) > maxMessageLength {
		sender.Send(serverMessage(
			fmt.Sprintf("Your message must be between 1 and %d characters", maxMessageLength),
		))
		return
	}

	if strings.HasPrefix(msg, "/script") {
		lines := strings.Split(msg, "\n")

		if len(lines) <= 1 {
			return
		}

		lines = lines[1:]

		for _, line := range lines {
			s.handleCommand(sender, line)
		}

		return
	}

	if msg[0] == '/' {
		s.handleCommand(sender, msg[1:])
		return
	}

	broadcastRoom(s, sender.RoomName(), &Message{
		Type: "chat",
		Data: map[string]interface{}{
			"sender": sender,
			"text":   msg,
		},
	})
}

// checkName checks if a client exists with a
// given name
func (s *Server) checkName(name string) bool {
	for _, c := range s.clients {
		if c.Username() == name {
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
		if c.RoomName() == room {
			names = append(names, c.Username())
		}
	}

	return names
}

func (s *Server) handleCommand(sender Client, str string) {
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

	if err := sender.Send(serverMessage(out)); err != nil {
		log.Println("handleCommand: errored when sending output message")
	}
}

func serverMessage(content string) *Message {
	return &Message{
		Type: "server-msg",
		Data: content,
	}
}
