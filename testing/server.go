package main

import (
	"fmt"
	"http"
	"time"
	"websocket"
	"bufio"
	//"io"
	//"os"
)

// Echo the data received on the Web Socket.
func EchoServer(ws *websocket.Conn) {
	fmt.Println("Got a connection!")
	go func() {
		r := bufio.NewReader(ws)
		for {
			x, err := r.ReadString('\n')
			//var x string
			//_,err := fmt.Fscan(ws, &x)
			if err == nil {
				fmt.Print("read: ", x)
				fmt.Fprintln(ws, `<h1>`, x[:len(x)-1], `</h1>`)
				//fmt.Fprintln(ws, x[:len(x)-1])
			} else {
				fmt.Println("error reading:", err)
				return
			}
		}
	}()
	//go io.Copy(os.Stdout, ws)
	for {
		fmt.Fprintln(ws, "I am getting bored...")
		time.Sleep(10e9)
	}
}

func home(c http.ResponseWriter, req *http.Request) {
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
var ws = new WebSocket("ws://localhost:12345/echo");

ws.onmessage = function (evt) {
   var received_msg = evt.data;
   //alert("Message is received: " + received_msg);
   document.getElementById("everything").innerHTML=received_msg;
};

ws.onclose = function() {
   // websocket is closed.
   alert("Connection is closed..."); 
};
</script>
</head>
<body>
<div id="everything">
  Everything goes here.
</div>
  <input type='submit' onclick="ws.send('hello world\n')" value='Hello.'/> 
  <input type='submit' onclick="ws.send('goodbye world\n')" value='Goodbye.'/> 
</body>
</html>
`)
}

func main() {
	http.Handle("/echo", websocket.Handler(EchoServer));
	http.HandleFunc("/", home)
	err := http.ListenAndServe(":12345", nil);
	if err != nil {
		panic("ListenAndServe: " + err.String())
	}
}
