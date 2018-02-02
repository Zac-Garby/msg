package server

import (
	"fmt"
	"strings"
)

// This function is called whenever a command is called
// by a client. It is passed a slice of the given
// arguments, and returns a string to send back to
// the user. If the returned string is empty, no
// message is sent to the user.
type command func(server *Server, sender *client, args []string) string

var commands = make(map[string]command)

func init() {
	commands["help"] = func(s *Server, c *client, args []string) string {
		return `Available commands:

help            prints help about the commands
list            lists the users in the current room
room [room]     with no arguments, prints the current room. otherwise, switches to [room]
quit            exits to the landing page
name [name]     sets your username to [name]
fg   [colour]   sets your username's foreground colour to [colour] (a valid CSS colour)
bg   [colour]   sets your username's background colour to [colour] (a valid CSS colour)`
	}

	commands["name"] = func(s *Server, c *client, args []string) string {
		if len(args) <= 1 {
			return "The command `name` expects an argument, [name]"
		}

		name := args[1]

		if len(name) > maxNameLength {
			return fmt.Sprintf("Your name cannot be longer than %d characters", maxNameLength)
		}

		if len(name) < minNameLength {
			return fmt.Sprintf("Your name cannot be less than %d characters", minNameLength)
		}

		if s.checkName(name) {
			return fmt.Sprintf("A user already exists called %s!", name)
		}

		c.name = name
		return fmt.Sprintf("Your name has been changed to %s!", name)
	}

	commands["list"] = func(s *Server, c *client, args []string) string {
		room := c.room

		if len(args) > 1 {
			room = args[1]
		}

		var (
			names = s.usersInRoom(room)
			out   = strings.Join(names, "\n")
		)

		if len(names) == 0 {
			return fmt.Sprintf("There are no users in %s", room)
		}

		return fmt.Sprintf("Users currently in %s (%d)\n%s", room, len(names), out)
	}
}
