# msg

A messaging web-app written in Go, with some JavaScript too.

It's usable, and has quite a lot of useful features:

  - A login page
    - You don't need to make an account
  - Rooms
    - Rooms don't need to be created either
  - Commands
    - Executed on the server
	  lots of commands at the same time
  - Easy deployment
    - Just clone the repository and run `main.go`

Here's a screenshot:

![](screenshot.png)

There are still a number of things that need doing:

  - Admins
    - Since user's aren't really a thing, there would
	  need to be some other system
  - Back up messages to a database (probably Redis)
  - Muting other users
  - Limiting amount of messages per minute
  - Improve styles on phones
  - Host it somewhere
