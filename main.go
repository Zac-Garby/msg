package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Zac-Garby/msg/server"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var s *server.Server

func main() {
	s = server.New()
	go s.HandleMessages()

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/ws", websocketHandler)

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static/")),
		),
	)

	r.Handle("/favicon.ico", http.FileServer(http.Dir("./static/")))

	fmt.Println("listening on localhost:3000")
	http.ListenAndServe(":3000", r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	if err := s.NewClient(conn); err != nil {
		log.Println("server err:", err)
		return
	}
}
