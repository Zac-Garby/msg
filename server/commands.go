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
name [name]     sets your username to [name]
quit            exits the server and returns to the login page`
	}

	commands["name"] = func(s *Server, c *client, args []string) string {
		if len(args) <= 1 {
			return "The command `name` expects an argument, [name]"
		}

		name := args[1]

		if msg, ok := ValidateName(name, s); !ok {
			return msg
		}

		c.Name = name
		return fmt.Sprintf("Your name has been changed to %s!", name)
	}

	commands["list"] = func(s *Server, c *client, args []string) string {
		room := c.Room

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

	commands["room"] = func(s *Server, c *client, args []string) string {
		if len(args) > 1 {
			room := args[1]
			broadcastRoom(s, c.Room, serverMessage(fmt.Sprintf("%s has left the room %s", c.Name, c.Room)))

			if msg, ok := ValidateRoom(room); !ok {
				return msg
			}

			c.Room = room

			broadcastRoom(s, c.Room, serverMessage(fmt.Sprintf("%s has joined the room %s", c.Name, c.Room)))

			return fmt.Sprintf("You are now in the room: %s", c.Room)
		} else {
			return fmt.Sprintf("You are currently in the room: %s", c.Room)
		}
	}

	commands["quit"] = func(s *Server, c *client, args []string) string {
		msg := &message{
			Type: "quit",
		}

		if err := c.send(msg); err != nil {
			return "Could not send quit message. It's possible you're not connected to the server."
		}

		return ""
	}
}
