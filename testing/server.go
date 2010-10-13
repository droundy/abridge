package main

import (
	"fmt"
	"http"
	//"time"
	"github.com/droundy/bridge/easysocket"
	"github.com/droundy/bridge"
	"os"
)

type TransitoryData struct {
	Bids string
	Hands bridge.Table
	Dealer bridge.Seat
	Bidfor bridge.Seat
	AmBidding bool
	NScard, EWcard string
	Url string
	Write func(string)
}

func (dat *TransitoryData) Handle(evt string) {
	dat.Write(`<h1>` + evt + `</h1>`)
}

func (dat *TransitoryData) Done(err os.Error) {
	fmt.Println("All done!", err)
}

func NewClient(write func(string)) easysocket.Handler {
	dat := new(TransitoryData)
	dat.Write = write
	write(`
<h1> Intro to aBridge</h1>

This is a neat thing.
`)
	return dat
}

func main() {
	easysocket.Handle("/", NewClient);
	err := http.ListenAndServe(":12345", nil);
	if err != nil {
		panic("ListenAndServe: " + err.String())
	}
}
