package server

import (
	"fmt"
	"log"
	"net/http"
)

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	if name, err := r.Cookie("name"); err == nil && name.Value != "" {
		http.ServeFile(w, r, "static/messager.html")
	} else {
		http.ServeFile(w, r, "static/index.html")
	}
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading:", err)
		return
	}

	if err := s.newClient(conn); err != nil {
		log.Println("server err:", err)
		return
	}
}

func (s *Server) validateHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()

	nameSlice, ok := vals["name"]
	if ok && len(nameSlice) > 0 {
		name := nameSlice[0]
		reason, valid := s.backend.ValidateName(name)
		if !valid {
			fmt.Fprintf(w, reason)
			return
		}
	}

	roomSlice, ok := vals["room"]
	if ok && len(roomSlice) > 0 {
		room := roomSlice[0]
		reason, valid := s.backend.ValidateRoom(room)
		if !valid {
			fmt.Fprintf(w, reason)
			return
		}
	}

	fmt.Fprintf(w, "ok")
}
