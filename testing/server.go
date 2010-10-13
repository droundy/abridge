package main

import (
	"http"
	"os"
	"fmt"
	"github.com/droundy/bridge/easysocket"
)

func (dat *ClientData) Done(err os.Error) {
	fmt.Println("All done!", err)
}

func NewClient(write func(string)) easysocket.Handler {
	dat := new(ClientData)
	dat.Write = write
	dat.MyPage = "Home"
	dat.Handle("")
	return dat
}

func main() {
	easysocket.Handle("/", NewClient);
	err := http.ListenAndServe(":12345", nil);
	if err != nil {
		panic("ListenAndServe: " + err.String())
	}
}
