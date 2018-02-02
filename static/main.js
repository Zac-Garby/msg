var sock

// connects to the websocket server hosted at
// ws://{host}.
function connect() {	
	sock = new WebSocket(`ws://${window.location.host}/ws`)
	
	sock.onopen = function(evt) {
		sock.send(JSON.stringify({
			type: "client-info",
			data: {
				name: "unnamed",
				room: "/",
			},
		}))
	}

	sock.onmessage = function(evt) {
		onMessage(evt.data)
	}
}

function onMessage(data) {
	var msg = JSON.parse(data)

	switch (msg.type) {
	case "chat":
		putMessage(msg.data.sender, "username", "", msg.data.text)
		break

	case "server-msg":
		putMessage("server", "server-username", "server-msg", msg.data)
		break
	}
}

function handleKey(evt) {
	if (evt.keyCode === 13 && evt.shiftKey) {
		evt.preventDefault()

		sock.send(JSON.stringify({
			type: "chat",
			data: evt.target.value,
		}))

		evt.target.value = ""
	}
}

function putMessage(sender, nameClass, liClass, content) {
	content = escapeHTML(content)
	
	var name = document.createElement("span")
	name.className = nameClass
	name.innerHTML = sender + ":"
	
	var text = document.createElement("pre")
	text.className = "text"
	text.innerHTML = content
	
	var li = document.createElement("li")
	li.className = liClass
	li.appendChild(name)
	li.appendChild(text)

	var log = document.getElementById("chat-log")
	log.appendChild(li)
	log.scrollTop = log.scrollHeight
}

function escapeHTML(html) {
	return html
         .replace(/&/g, "&amp;")
         .replace(/</g, "&lt;")
         .replace(/>/g, "&gt;")
         .replace(/"/g, "&quot;")
         .replace(/'/g, "&#039;");
}
