package main

import (
	"github.com/GeertJohan/go.rice/embedded"
	"time"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "favicon.ico",
		FileModTime: time.Unix(1524601990, 0),
		Content:     string("\x00\x00\x01\x00\x01\x00\x10\x10\x02\x00\x01\x00\x01\x00\xb0\x00\x00\x00\x16\x00\x00\x00(\x00\x00\x00\x10\x00\x00\x00 \x00\x00\x00\x01\x00\x01\x00\x00\x00\x00\x00\x80\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xff\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00"),
	}
	file3 := &embedded.EmbeddedFile{
		Filename:    "index.css",
		FileModTime: time.Unix(1524840417, 0),
		Content:     string("* {\n\tbox-sizing: border-box;\n\tmargin: 0;\n\toutline: none;\n}\n\nhtml, body {\n\tfont-family: monospace;\n}\n\nul#chat-log {\n\tborder-bottom: 1px solid black;\n\tpadding-left: 0;\n\tdisplay: block;\n\twidth: 100vw;\n\theight: calc(100vh - 70px);\n\tresize: none;\n\toverflow-y: auto;\n}\n\nul#chat-log li {\n\tpadding: 10px;\n\tborder-bottom: 1px solid #afafaf;\n}\n\nul#chat-log li pre {\n\toverflow-wrap: normal;\n\tword-wrap: break-word;\n\thyphens: auto;\n}\n\nul#chat-log li .username {\n\tfont-weight: bold;\n\tmargin-right: 10px;\n}\n\nul#chat-log li.server-msg {\n\tbackground-color: #ddddff;\n\tcolor: #000055;\n}\n\ntextarea#input {\n\tdisplay: block;\n\tpadding: 10px;\n\theight: 70px;\n\twidth: 100vw;\n\tborder: none;\n\tresize: none;\n\tpadding-right: 70px;\n\tbackground-color: #aaaaff3f;\n}\n\ntextarea#input::placeholder {\n\tcolor: #0000aa;\n}\n\nbutton#send-msg {\n\tposition: absolute;\n\tbottom: 10px;\n\tright: 10px;\n\theight: 50px;\n\twidth: 50px;\n\tborder-radius: 5px;\n\tborder: 1px solid #0000aa;\n\tbackground-color: #0000ff2c;\n\tcolor: #0000cc;\n\tcursor: pointer;\n}\n\nbutton#send-msg:active {\n\tbackground-color: #0000ff5c;\n}\n\ndiv.center {\n\tposition: absolute;\n\twidth: 400px;\n\ttop: 0;\n\tbottom: 0;\n\tleft: 0;\n\tright: 0;\n\tmargin: auto;\n}\n\ndiv.center h1 {\n\ttext-align: center;\n\tmargin-top: 70px;\n\tmargin-bottom: 20px;\n}\n\ndiv.center p {\n\tposition: relative;\n\tmargin-bottom: 20px;\n}\n\ndiv.center p#error {\n\tcolor: red;\n}\n\ndiv.center input.login {\n\tposition: absolute;\n\tright: 0;\n\twidth: 300px;\n}\n\nbutton#login-btn {\n\tposition: absolute;\n\tleft: 0;\n\tright: 0;\n\tmargin: auto;\n\twidth: 10em;\n}\n"),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    "index.html",
		FileModTime: time.Unix(1525286758, 0),
		Content:     string("<!DOCTYPE html>\n<html lang=\"en\">\n  <head>\n\t<meta charset=\"utf-8\">\n\t<meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n\t<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n\n\t<link rel=\"stylesheet\" href=\"/static/index.css\" />\t\n\t<title>msg - home</title>\n  </head>\n  <body>\n\t<div class=\"center\">\n\t  <h1>log in</h1>\n\n\t  <p>\n\t\tTo log in, you don't need to make an account. Just\n\t\tenter a name and what room you want to go to. You can\n\t\tchange both of these things later via commands.\n\t  </p>\n\n\t  <p id=\"error\"></p>\n\t  \n\t  <p>\n\t\tUsername:\n\t\t<input class=\"login\" id=\"name\" placeholder=\"anon\"></input>\n\t  </p>\n\t  <p>\n\t\tRoom:\n\t\t<input class=\"login\" id=\"room\" placeholder=\"/\"></input>\n\t  </p>\n\t  <p>\n\t\t<button id=\"login-btn\" onclick=\"login()\">Log in</button>\n\t  </p>\n\t</div>\n  </body>\n  <script src=\"/static/index.js\"></script>\n</html>\n"),
	}
	file5 := &embedded.EmbeddedFile{
		Filename:    "index.js",
		FileModTime: time.Unix(1525287099, 0),
		Content:     string("function login() {\n\tvar name = document.getElementById(\"name\").value || \"anon\"\n\tvar room = document.getElementById(\"room\").value || \"/\"\n\n\tvalidate(name, room, function(valid, reason) {\n\t\tif (valid) {\n\t\t\tdocument.cookie = `name=${name}`\n\t\t\tdocument.cookie = `room=${room}`\n\t\n\t\t\tlocation.reload(true)\n\t\t} else {\n\t\t\tvar err = document.getElementById(\"error\")\n\t\t\terr.innerHTML = reason\n\t\t}\n\t})\t\n}\n\nfunction validate(name, room, callback) {\n\tvar url = `http://${document.location.host}/validate?name=${name};room=${room}`\n\tvar req = new XMLHttpRequest()\n\n\treq.onreadystatechange = function() {\n\t\tif (this.readyState === 4) {\n\t\t\tcallback(this.responseText == \"ok\", this.responseText)\n\t\t}\n\t}\n\n\treq.open(\"GET\", url, true)\n\treq.send()\n}\n\nfunction getCookie(name) {\n\tvar value = \"; \" + document.cookie\n\tvar parts = value.split(\"; \" + name + \"=\")\n\tif (parts.length == 2) return parts.pop().split(\";\").shift()\n}\n\nvar room = getCookie(\"room\")\nif (room) {\n\tdocument.getElementById(\"room\").value = room\n}\n"),
	}
	file6 := &embedded.EmbeddedFile{
		Filename:    "main.js",
		FileModTime: time.Unix(1524761805, 0),
		Content:     string("var sock\n\n// connects to the websocket server hosted at\n// ws://{host}.\nfunction connect() {\t\n\tsock = new WebSocket(`ws://${window.location.host}/ws`)\n\t\n\tsock.onopen = function(evt) {\n\t\tsock.send(JSON.stringify({\n\t\t\ttype: \"client-info\",\n\t\t\tdata: {\n\t\t\t\tname: getCookie(\"name\") || \"anon\",\n\t\t\t\troom: getCookie(\"room\") || \"/\",\n\t\t\t},\n\t\t}))\n\t}\n\n\tsock.onmessage = function(evt) {\n\t\tonMessage(evt.data)\n\t}\n}\n\nfunction onMessage(data) {\n\tvar msg = JSON.parse(data)\n\n\tswitch (msg.type) {\n\tcase \"chat\":\n\t\tputMessage(msg.data.sender.name, \"username\", \"\", msg.data.text)\n\t\tbreak\n\n\tcase \"server-msg\":\n\t\tputMessage(\"server\", \"\", \"server-msg\", msg.data)\n\t\tbreak\n\t\n\tcase \"cookie\":\n\t\tdocument.cookie = msg.data\n\t\tbreak\n\n\tcase \"quit\":\n\t\tdocument.cookie = \"name=;expires=Thu, 01 Jan 1970 00:00:01 GMT;\"\n\t\tdocument.cookie = \"room=;expires=Thu, 01 Jan 1970 00:00:01 GMT;\"\n\t\tlocation.reload(true)\n\t\tbreak\n\t}\n}\n\nfunction handleKey(evt) {\n\tif (evt.keyCode === 13 && evt.shiftKey) {\n\t\tevt.preventDefault()\n\n\t\tsendMessage()\n\t}\n}\n\nfunction sendMessage() {\n\tvar elem = document.getElementById(\"input\")\n\n\tsock.send(JSON.stringify({\n\t\ttype: \"chat\",\n\t\tdata: elem.value,\n\t}))\n\t\n\telem.value = \"\"\n}\n\nfunction putMessage(sender, nameClass, liClass, content) {\n\tcontent = escapeHTML(content)\n\t\n\tvar name = document.createElement(\"span\")\n\tname.className = nameClass\n\tname.innerHTML = sender + \":\"\n\t\n\tvar text = document.createElement(\"pre\")\n\ttext.className = \"text\"\n\ttext.innerHTML = content\n\t\n\tvar li = document.createElement(\"li\")\n\tli.className = liClass\n\tif (liClass !== \"server-msg\") li.appendChild(name)\n\tli.appendChild(text)\n\n\tvar log = document.getElementById(\"chat-log\")\n\tlog.appendChild(li)\n\tlog.scrollTop = log.scrollHeight\n}\n\nfunction escapeHTML(html) {\n\treturn html\n         .replace(/&/g, \"&amp;\")\n         .replace(/</g, \"&lt;\")\n         .replace(/>/g, \"&gt;\")\n         .replace(/\"/g, \"&quot;\")\n         .replace(/'/g, \"&#039;\");\n}\n\nfunction getCookie(name) {\n\tvar value = \"; \" + document.cookie\n\tvar parts = value.split(\"; \" + name + \"=\")\n\tif (parts.length == 2) return parts.pop().split(\";\").shift()\n}\n"),
	}
	file7 := &embedded.EmbeddedFile{
		Filename:    "messager.html",
		FileModTime: time.Unix(1525284981, 0),
		Content:     string("<!DOCTYPE html>\n<html lang=\"en\">\n  <head>\n\t<meta charset=\"utf-8\">\n\t<meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n\t<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n\t\n\t<link rel=\"stylesheet\" href=\"/static/index.css\" />\n\t<script src=\"/static/main.js\"></script>\n\t\n\t<title>msg</title>\n  </head>\n  <body onload=\"connect()\">\n\t<ul id=\"chat-log\"></ul>\n\t<textarea id=\"input\" onkeypress=\"handleKey(event)\" placeholder=\"Type a message here. Press shift+return to send it\"></textarea>\n\t<button id=\"send-msg\" onclick=\"sendMessage()\">Send</button>\n  </body>\n</html>\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1524840417, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "favicon.ico"
			file3, // "index.css"
			file4, // "index.html"
			file5, // "index.js"
			file6, // "main.js"
			file7, // "messager.html"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`./static/`, &embedded.EmbeddedBox{
		Name: `./static/`,
		Time: time.Unix(1524840417, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"favicon.ico":   file2,
			"index.css":     file3,
			"index.html":    file4,
			"index.js":      file5,
			"main.js":       file6,
			"messager.html": file7,
		},
	})
}
