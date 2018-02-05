package main

import (
	"github.com/Zac-Garby/msg/server"
)

var s *server.Server

func main() {
	s = server.New()
	s.Start()
}
