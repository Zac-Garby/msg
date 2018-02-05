package backend

import (
	"fmt"
	"strings"
)

func (s *Backend) handleCommand(sender *Client, str string) {
	var (
		out  string
		args = strings.Fields(str)
	)

	if len(args) == 0 || len(str) == 0 {
		out = "expected a command after '/'"
	} else {
		var (
			name    = args[0]
			cmd, ok = commands[name]
		)

		if ok {
			out = cmd(b, sender, args)
		} else {
			out = fmt.Sprintf("command not found: %s", name)
		}
	}

	if len(out) > 0 {
		ServerMessage(out).Send(b, sender)
	}
}

// A command function is called whenever a command is called
// by a client.
type command func(b *Backend, sender *Client, args []string) string

var commands = make(map[string]command)

func init() {
	commands["help"] = func(b *Backend, c *Client, args []string) string {
		return `Available commands:

help          prints this message
list [room]   lists the users in either the current room, or [room] (if specified)
room [room]   with no arguments, prints the current room. otherwise, puts you in [room]
name [name]   with no arguments, prints your username. otherwise, sets your username to [name]
quit          logs you out of the server and exits to the login page
`
	}

	commands["list"] = func(b *Backend, c *Client, args []string) string {
		room := c.Room

		if len(args) > 1 {
			room = args[1]
		}

		var (
			names = s.UsersInRoom(room)
			out   = strings.Join(names, "\n")
		)

		if len(names) == 0 {
			return fmt.Sprintf("There are no users in %s", room)
		}

		return fmt.Sprintf("There are %d users in %s:\n%s", len(names), room, out)
	}

	commands["room"] = func(b *Backend, c *Client, args []string) string {
		if len(args) == 1 {
			return fmt.Sprintf("You are currently in '%s", c.Room)
		}

		room := args[1]
		ServerMessage("%s has left the room '%s'", c.Name, c.Room).Broadcast(b, c.Room)

		if msg, ok := b.ValidateRoom(room); !ok {
			return msg
		}

		c.Room = room
		ServerMessage("%s has joined the room '%s'", c.Name, c.Room).Broadcast(b, c.Room)

		return ""
	}

	commands["name"] = func(b *Backend, c *Client, args []string) string {
		if len(args) == 1 {
			return fmt.Sprintf("You are called '%s'", c.Name)
		}

		var (
			old  = c.Name
			name = args[1]
		)

		if msg, ok := b.ValidateName(name); !ok {
			return msg
		}

		c.Name = name
		ServerMessage("%s has changed their name to %s", old, c.Name).Broadcast(b, c.Room)
	}

	commands["quit"] = func(b *Backend, c *Client, args []string) string {
		&Message{
			Type: MsgQuit,
		}.Send(b, c)

		return ""
	}
}
