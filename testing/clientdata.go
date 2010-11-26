package main

import (
	"fmt"
	"regexp"
	"sync"
	"strings"
	"github.com/droundy/abridge"
)

type Rule struct {
	Pattern *regexp.Regexp
	Code func(dat *ClientData, matches []string) string
}

type ClientData struct {
	MyPage string

	Cookie Settings
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

var bidexpression = regexp.MustCompile(`([^=]*)(=.*)?`)
var settingexpression = regexp.MustCompile(`(set style |set card |check |uncheck |rename to |select |setpts )?([^=]*)(=.*)?`)

func (dat *ClientData) Page(evt string) string {
	once.Do(func () {
		EventHandlers["Home"] = []Rule {
			Rule {
			Pattern: regexp.MustCompile(`.*`),
			Code: Home,
			},
		}
		EventHandlers["Bid fourth hand"] = []Rule {
			Rule {
			Pattern: regexp.MustCompile(`.*`),
			Code: Home,
			},
		}
		EventHandlers["Analyze bids"] = []Rule {
			Rule {
			Pattern: bidexpression,
			Code: AnalyzeBids,
			},
		}
		EventHandlers["Bid for me"] = []Rule {
			Rule {
			Pattern: regexp.MustCompile(`.*`),
			Code: Home,
			},
		}
		EventHandlers["Settings"] = []Rule {
			Rule {
			Pattern: settingexpression,
			Code: SettingsPage,
			},
		}
	})
	rs,ok := EventHandlers[dat.MyPage]
	if !ok {
		return "Error: bad dat.MyPage " + dat.MyPage
	}
	for _,r := range rs {
		if ms := r.Pattern.FindStringSubmatch(evt); ms != nil {
			return r.Code(dat, ms)
		}
	}
	return "Unknown event type: " + evt
}

func (dat *ClientData) WriteCookie() {
	dat.Write("write-cookie" + dat.Cookie.Write())
}

func (dat *ClientData) Handle(evt string) {
	fmt.Println("Got event:", evt)
	if evt == "First time" {
		dat.Write("read-cookie")
		return
	} else if strings.HasPrefix(evt, "cookie is ") {
		fmt.Println("got cookie:", evt)
		dat.Cookie = readCookie(evt[len("cookie is "):])
		dat.WriteCookie()
		return
	}
	if ms := isgo.FindStringSubmatch(evt); ms != nil {
		dat.MyPage = ms[1]
	}
	out := ""
	out += dat.Header()
	out += dat.Page(evt)
	out += `<h3>` + evt + `</h3>`
	dat.Write(out)
}
