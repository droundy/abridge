package main

import (
	"http"
	"os"
	"fmt"
	"github.com/droundy/abridge/easysocket"
)

func (dat *ClientData) Done(err os.Error) {
	fmt.Println("All done!", err)
}

func NewClient(write func(string)) easysocket.Handler {
	dat := new(ClientData)
	dat.Write = write
	dat.MyPage = "Home"
	dat.Handle("First time")
	return dat
}

func main() {
	easysocket.Handle("/", NewClient)
	http.HandleFunc("/favicon.ico", faviconServer)
	http.HandleFunc("/style.css", styleServer)
	http.HandleFunc("/style-fourcolor.css", styleServer)
	http.HandleFunc("/speech/", wavServer)
	err := http.ListenAndServe(":12345", nil);
	if err != nil {
		panic("ListenAndServe: " + err.String())
	}
}
