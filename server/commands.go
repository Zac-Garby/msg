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
type command func(server *Server, sender Client, args []string) string

var commands = make(map[string]command)

func init() {
	commands["help"] = func(s *Server, c Client, args []string) string {
		return `Available commands:

help            prints help about the commands
list            lists the users in the current room
room [room]     with no arguments, prints the current room. otherwise, switches to [room]
name [name]     sets your username to [name]
quit            exits the server and returns to the login page`
	}

	commands["name"] = func(s *Server, c Client, args []string) string {
		if len(args) <= 1 {
			return "The command `name` expects an argument, [name]"
		}

		name := args[1]

		if msg, ok := ValidateName(name, s); !ok {
			return msg
		}

		broadcastRoom(s, c.RoomName(), serverMessage(fmt.Sprintf("%s is now called %s.", c.Username(), name)), c.ID())

		c.Rename(name)

		// tell the client to update their name cookie
		c.Send(&Message{
			Type: "cookie",
			Data: fmt.Sprintf("name=%s;", c.Username()),
		})

		return fmt.Sprintf("Your name has been changed to %s!", name)
	}

	commands["list"] = func(s *Server, c Client, args []string) string {
		room := c.RoomName()

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

	commands["room"] = func(s *Server, c Client, args []string) string {
		if len(args) > 1 {
			room := args[1]
			broadcastRoom(s, c.RoomName(), serverMessage(fmt.Sprintf("%s has left the room %s", c.Username(), c.RoomName())), c.ID())

			if msg, ok := ValidateRoom(room); !ok {
				return msg
			}

			c.GotoRoom(room)

			broadcastRoom(s, c.RoomName(), serverMessage(fmt.Sprintf("%s has joined the room %s", c.Username(), c.RoomName())), c.ID())

			return fmt.Sprintf("You are now in the room: %s", c.RoomName())
		} else {
			return fmt.Sprintf("You are currently in the room: %s", c.RoomName())
		}
	}

	commands["quit"] = func(s *Server, c Client, args []string) string {
		msg := &Message{
			Type: "quit",
		}

		if err := c.Send(msg); err != nil {
			return "Could not send quit message. It's possible you're not connected to the server."
		}

		return ""
	}
}
