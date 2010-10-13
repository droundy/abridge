package main

import (
	"github.com/droundy/bridge"
)

type ClientData struct {
	MyPage string

	Bids string
	Hands bridge.Table
	Dealer bridge.Seat
	Bidfor bridge.Seat
	AmBidding bool
	NScard, EWcard string

	Write func(string)
}

func (dat *ClientData) Home() string {
	return `
<div id="header">
  <a href="javascript:say('go home')">Home</a>
  <a href="javascript:say('go analyze bids')">Analyze bids</a>
  <a href="javascript:say('go bid fourth hand')">Bid fourth hand</a>
</div>

<h1> Intro to aBridge</h1>

This is a neat thing.
<br/>

  <input type='submit' onclick="say('hello world')" value='Hello.'/> 
  <input type='submit' onclick="say('goodbye world')" value='Goodbye.'/> 
`
}

func (dat *ClientData) Handle(evt string) {
	out := ""
	out += dat.Home()
	out += `<h3>` + evt + `</h3>`
	dat.Write(out)
}
