# msg

A messaging web-app written in Go, with some JavaScript too.

It's usable, and has quite a lot of useful features:

  - A login page
    - You don't need to make an account
  - Rooms
    - Rooms don't need to be created either
  - Commands
    - Executed on the server
	- Special command `/script` allows you to execute
	  lots of commands at the same time
  - Easy deployment
    - Just clone the repository and run `main.go`

Here's a screenshot:

![](screenshot.png)

There are still a number of things that need doing:

  - Admins
    - Since user's aren't really a thing, there would
	  need to be some other system
  - Muting other users
  - Limiting amount of messages per minute
  - Improve styles on phones
  - Host it somewhere
  - Maybe rewrite backend -- split it into two parts:
    - One package to actually do stuff
	  - Handling commands
	  - Keeping track of users
    - Another package - a web server - to act as a proxy
	  between the backend and the frontend
	  - Keeps track of websockets
	  - Serves the website's pages
	  - Maybe even some routes for an API
	    - e.g. `/users?room=x` lists users in a room
    - This would make it much more extensible
