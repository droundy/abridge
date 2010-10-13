package easysocket

import (
	"fmt"
	"http"
	"websocket"
	"bufio"
	"path"
	"os"
)

// We export a single function, which creates a page controlled by a
// single websocket.  It's quite primitive, and yet quite easy to use!
func Handle(url string, handler func(evts <-chan string, pages chan<- string, done <-chan os.Error)) {
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
		fmt.Fprintln(c, `
<!DOCTYPE HTML>
<html>
<head>
<script type="text/javascript">

if (! "WebSocket" in window) {
 // The browser doesn't support WebSocket
 alert("WebSocket NOT supported by your Browser!");
}

// Let us open a web socket
var ws = new WebSocket("ws://localhost:12345` + path.Join(req.URL.Path,"socket") + `");
ws.onmessage = function (evt) {
   var received_msg = evt.data;
   //alert("Message is received: " + received_msg);
   document.getElementById("everything").innerHTML=received_msg;
};
say = function(txt) {
   ws.send(txt + '\n')
}
ws.onclose = function() {
   // websocket is closed.
   alert("Connection is closed..."); 
};
</script>
</head>
<body>
<div id="header">
  <a href="javascript:say('go home')">Home</a>
  <a href="javascript:say('go analyze bids')">Analyze bids</a>
  <a href="javascript:say('go bid fourth hand')">Bid fourth hand</a>
</div>
<div id="everything">


  Everything goes here.
</div>
  <input type='submit' onclick="say('hello world')" value='Hello.'/> 
  <input type='submit' onclick="say('goodbye world')" value='Goodbye.'/> 
</body>
</html>
`)
	}
	http.HandleFunc(url, skeleton)
}
