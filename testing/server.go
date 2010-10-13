package main

import (
	"fmt"
	"http"
	//"time"
	"github.com/droundy/bridge/easysocket"
	"os"
)

// Echo the data received on the Web Socket.
func BridgeServer(evts <-chan string, pages chan<- string, done <-chan os.Error) {
	fmt.Println("about to send intro page...")
	pages <- `
<h1> Intro to aBridge</h1>

This is a neat thing.
`
	// The ticks gives a demo of how we could handle some sort of a
	// timeout.
	//ticks := time.NewTicker(10e9)
	for {
		select {
		case x := <- evts:
			fmt.Println("got event:", x)
			pages <- `<h1>` + x + `</h1>`
			//ticks.Stop()
			//ticks = time.NewTicker(10e9) // Start counting again!
		case err := <- done:
			fmt.Println("All done!", err)
			return
		//case _ = <- ticks.C:
		//	pages <- `I am getting bored...`
		}
	}
}

func main() {
	easysocket.Handle("/", BridgeServer);
	err := http.ListenAndServe(":12345", nil);
	if err != nil {
		panic("ListenAndServe: " + err.String())
	}
}
