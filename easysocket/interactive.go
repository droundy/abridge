package easysocket

import (
	"fmt"
	"http"
	"websocket"
	"bufio"
	"path"
	"os"
)

type Handler interface {
	Handle(evt string)
	Done(os.Error)
}

// We export a single function, which creates a page controlled by a
// single websocket.  It's quite primitive, and yet quite easy to use!
func Handle(url string, handler func(write func(string)) Handler) {
	myh := func(ws *websocket.Conn) {
		h := handler(func (p string) { fmt.Fprintln(ws, p) })
		r := bufio.NewReader(ws)
		for {
			x, err := r.ReadString('\n')
			if err == nil {
				h.Handle(x[:len(x)-1])
			} else {
				h.Done(os.NewError("Error from r.ReadString: " + err.String()))
				return
			}
		}
	}
	http.Handle(path.Join(url, "socket"), websocket.Handler(myh))

	skeleton := func(c http.ResponseWriter, req *http.Request) {
		c.SetHeader("Content-Type", "text/html")
		fmt.Fprintln(c, skeletonpage(req))
	}
	http.HandleFunc(url, skeleton)

}

// We export a single function, which creates a page controlled by a
// single websocket.  It's quite primitive, and yet quite easy to use!
func HandleChans(url string, handler func(evts <-chan string, pages chan<- string, done <-chan os.Error)) {
	myh := func(ws *websocket.Conn) {
		evts := make(chan string)
		pages := make(chan string)
		done := make(chan os.Error)
		go handler(evts, pages, done)
		go func() {
			r := bufio.NewReader(ws)
			for {
				x, err := r.ReadString('\n')
				if err == nil {
					evts <- x[:len(x)-1]
				} else {
					done <- os.NewError("Error from r.ReadString: " + err.String())
					return
				}
			}
		}()
		for {
			x := <- pages
			_,err := fmt.Fprintln(ws, x)
			if err != nil {
				done <- os.NewError("Error in fmt.Fprintln: " + err.String())
				return
			}
		}
	}
	http.Handle(path.Join(url, "socket"), websocket.Handler(myh))

	skeleton := func(c http.ResponseWriter, req *http.Request) {
		c.SetHeader("Content-Type", "text/html")
		fmt.Fprintln(c, skeletonpage(req))
	}
	http.HandleFunc(url, skeleton)
}

func skeletonpage(req *http.Request) string {
	return `<!DOCTYPE HTML>
<html>
<head>
<script type="text/javascript">

// Define helper cookie functions:
function createCookie(name,value,days) {
	if (days) {
		var date = new Date();
		date.setTime(date.getTime()+(days*24*60*60*1000));
		var expires = "; expires="+date.toGMTString();
	}
	else var expires = "";
	document.cookie = name+"="+value+expires+"; path=/";
}
function readCookie(name) {
	var nameEQ = name + "=";
	var ca = document.cookie.split(';');
	for(var i=0;i < ca.length;i++) {
		var c = ca[i];
		while (c.charAt(0)==' ') c = c.substring(1,c.length);
		if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
	}
	return null;
}
function eraseCookie(name) {
	createCookie(name,"",-1);
}

// Set up the websocket
if (! "WebSocket" in window) {
 // The browser doesn't support WebSocket
 alert("WebSocket NOT supported by your Browser!");
}

// Let us open a web socket
var ws = new WebSocket("ws://localhost:12345` + path.Join(req.URL.Path,"socket") + `");
function say(txt) {
   ws.send(txt + '\n')
};
ws.onclose = function() {
   // websocket is closed.
   alert("Connection is closed..."); 
};

ws.onmessage = function (evt) {
   if (evt.data.replace(/^\s+|\s+$/g,"") == 'read-cookie') {
       var cookie = readCookie('WebSocketCookie');
       if (cookie != null) {
         say('cookie is ' + readCookie('WebSocketCookie'));
       } else {
         say('cookie is unknown')
       }
       return
   }
   if (evt.data.substr(0,12) == 'write-cookie') {
      createCookie('WebSocketCookie', evt.data.substr(12), 365);
      say('got cookie');
      return
   }
   var everything = document.getElementById("everything")
   if (everything == null) {
     return
   }
   var received_msg = evt.data;
   //alert("Message is received: " + received_msg);
   everything.innerHTML=received_msg;
};

</script>
</head>
<body>
<div id="everything">

  Everything goes here.
</div>
</body>
</html>
`
}
