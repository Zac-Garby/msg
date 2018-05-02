function login() {
	var name = document.getElementById("name").value || "anon"
	var room = document.getElementById("room").value || "/"

	validate(name, room, function(valid, reason) {
		if (valid) {
			document.cookie = `name=${name}`
			document.cookie = `room=${room}`
	
			location.reload(true)
		} else {
			var err = document.getElementById("error")
			err.innerHTML = reason
		}
	})	
}

function validate(name, room, callback) {
	var url = `http://${document.location.host}/validate?name=${name};room=${room}`
	var req = new XMLHttpRequest()

	req.onreadystatechange = function() {
		if (this.readyState === 4) {
			callback(this.responseText == "ok", this.responseText)
		}
	}

	req.open("GET", url, true)
	req.send()
}

function getCookie(name) {
	var value = "; " + document.cookie
	var parts = value.split("; " + name + "=")
	if (parts.length == 2) return parts.pop().split(";").shift()
}

var room = getCookie("room")
if (room) {
	document.getElementById("room").value = room
}
