package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/Zac-Garby/msg/backend"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// The Server serves the website, and creates websocket connections
// to clients. When a message is received on a websocket, it's relayed
// to the backend and, in certain cases, a message is sent back.
type Server struct {
	backend     *backend.Backend
	connections map[uuid.UUID]*websocket.Conn
	upgrader    websocket.Upgrader
}

// New constructs a new Server instance.
func New() *Server {
	return &Server{
		backend:     backend.New(),
		connections: make(map[uuid.UUID]*websocket.Conn),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// Start starts serving the website and listening for websocket connections.
func (s *Server) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", s.indexHandler)
	r.HandleFunc("/ws", s.websocketHandler)
	r.HandleFunc("/validate", s.validateHandler)
	r.Handle("/favicon.ico", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))

	go s.handleOutgoing()
	go s.backend.HandleIncoming()

	log.Println("listening on localhost:3000")
	http.ListenAndServe(":3000", r)
}

func (s *Server) newClient(conn *websocket.Conn) error {
	defer conn.Close()

	id := uuid.NewV4()
	s.backend.Clients[id] = &backend.Client{
		ID:        id,
		Connected: false, // set to true after client info receieved
	}
	s.connections[id] = conn

	client := s.backend.Clients[id]

	for {
		msg, err := s.read(conn)

		if err != nil {
			// Under certain conditions, don't log an error
			if strings.Contains(err.Error(), "use of closed network") ||
				strings.Contains(err.Error(), "unexpected EOF") {
				break
			}

			log.Println("server err when listening:", err)
			break
		}

		msg.Client = client

		s.backend.Incoming <- msg
	}

	if client.Connected {
		serverMessage("%s has left the server", client.Name).Broadcast(s.backend)
	}

	delete(s.backend.Clients, id)

	return nil
}

func (s *Server) read(conn *websocket.Conn) (*backend.Message, error) {
	msg := &backend.Message{}
	return msg, conn.ReadJSON(msg)
}

func serverMessage(format string, args ...interface{}) *backend.Message {
	return backend.ServerMessage(format, args...)
}
