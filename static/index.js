function login() {
	var name = document.getElementById("name").value || "anon"
	var room = document.getElementById("room").value || "/"

	document.cookie = `name=${name}`
	document.cookie = `room=${room}`
	
	location.reload()
}
