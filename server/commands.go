package server

// This function is called whenever a command is called
// by a client. It is passed a slice of the given
// arguments, and returns a string to send back to
// the user. If the returned string is empty, no
// message is sent to the user.
type command func(sender *client, args []string) string

var commands = make(map[string]command)

func init() {
	commands["help"] = func(_ *client, args []string) string {
		return `Available commands:

help            prints help about the commands
list            lists the users in the current room
room [room]     with no arguments, prints the current room. otherwise, switches to [room]
quit            exits to the landing page
name [name]     sets your username to [name]
fg   [colour]   sets your username's foreground colour to [colour] (a valid CSS colour)
bg   [colour]   sets your username's background colour to [colour] (a valid CSS colour)`
	}
}
