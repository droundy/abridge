package main

import (
	"regexp"
	"sync"
	"github.com/droundy/bridge"
)

type Rule struct {
	Pattern *regexp.Regexp
	Code func([]string) string
}

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

var EventHandlers = make(map[string][]Rule)
var once sync.Once
var isgo = regexp.MustCompile(`^go (.+)$`)

func (dat *ClientData) Page(evt string) string {
	once.Do(func () {
		EventHandlers["Home"] = []Rule {
			Rule {
			Pattern: regexp.MustCompile(`.*`),
			Code: Home,
			},
		}
	})
	rs,ok := EventHandlers[dat.MyPage]
	if !ok {
		return "Error: bad dat.MyPage " + dat.MyPage
	}
	for _,r := range rs {
		if ms := r.Pattern.FindStringSubmatch(evt); ms != nil {
			return r.Code(ms)
		}
	}
	return "Unknown event type: " + evt
}

func (dat *ClientData) Handle(evt string) {
	if ms := isgo.FindStringSubmatch(evt); ms != nil {
		dat.MyPage = ms[1]
	}
	out := ""
	out += dat.Header()
	out += dat.Page(evt)
	out += `<h3>` + evt + `</h3>`
	dat.Write(out)
}
