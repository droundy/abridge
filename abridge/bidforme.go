package main

import (
	"os"
	"io"
	"fmt"
	"http"
	"github.com/droundy/bridge"
)

// Bid the fourth hand...
func bidForMeNow(c *http.Conn, req *http.Request, clientname string) {
	bid,ok := req.Form["bid"]
	switch {
	case !ok || len(bid) != 1:
	case bid[0][1:] == bridge.SuitHTML[bridge.Clubs]:
		bids[clientname] = bids[clientname] + bid[0][0:1] + "C"
	case bid[0][1:] == bridge.SuitHTML[bridge.Diamonds]:
		bids[clientname] = bids[clientname] + bid[0][0:1] + "D"
	case bid[0][1:] == bridge.SuitHTML[bridge.Hearts]:
		bids[clientname] = bids[clientname] + bid[0][0:1] + "H"
	case bid[0][1:] == bridge.SuitHTML[bridge.Spades]:
		bids[clientname] = bids[clientname] + bid[0][0:1] + "S"
	case bid[0][1:] == bridge.SuitHTML[bridge.NoTrump]:
		bids[clientname] = bids[clientname] + bid[0][0:1] + "N"
	case len(bid[0]) == 2:
		bids[clientname] = bids[clientname] + bid[0]
	default:
		fmt.Println("I don't recognize", bid[0])
	}

	if d, ok := req.Form["dealer"]; ok && len(d) == 1 {
		dealer[clientname] = bridge.StringToSeat(d[0])
	}
	if _,ok := req.Form["undo"]; ok && len(bids[clientname]) >= 2 {
		bids[clientname] = bids[clientname][0:len(bids[clientname])-2]
		bidder := (dealer[clientname] + bridge.Seat(len(bids[clientname])/2)) % 4
		if bidder == bridge.South {
			bids[clientname] = bids[clientname][0:len(bids[clientname])-2]
		}
	}
	if _,ok := req.Form["refresh"]; ok {
		bridge.ClearBid(bids[clientname])
	}
	if _,ok := req.Form["clear"]; ok {
		bids[clientname] = ""
		dealer[clientname] = (dealer[clientname] + 1) % 4
		var x bridge.Table
		hands[clientname] = x, false
		defer header(c, req, "Enter your next hand")()
		askhand(c, req, bridge.South, clientname)
		return
	}

	bidder := (dealer[clientname] + bridge.Seat(len(bids[clientname])/2)) % 4
	cc := getSettings(req).Card
	ts := bridge.GetValidTables(dealer[clientname], bids[clientname], 100, cc)
	if bidder == bridge.South {
		fmt.Println("Bids are:", bids[clientname])
		fmt.Println("Table is:")
		fmt.Println(hands[clientname])
		newbid, conv := bridge.PickBid(hands[clientname][bridge.South], bidder, bids[clientname], cc, ts)
		bids[clientname] += newbid
		fmt.Println("Bid using", conv)
		ts = bridge.GetValidTables(dealer[clientname], bids[clientname], 100, cc)
	}
	defer header(c, req, "Bridge bidder")()
	bidbox(c, req, clientname, bridge.South)
	stats := bridge.GetValidTables(dealer[clientname], bids[clientname], 100, cc)
	fmt.Fprintln(c, `<table><tr><td>`)
	fmt.Fprintln(c, hands[clientname][0].HTML("My hand"))
	fmt.Fprintln(c, `</td><td>`)
	fmt.Fprintln(c, stats.HTML())
	fmt.Fprintln(c, `</td></tr></table>`)
	fmt.Fprintln(c, stats.ExampleHTML())
	showbids(c, clientname)
	showconventions(c, clientname, ts.Conventions)
}

// Bid my hand...
func bidforme(c *http.Conn, req *http.Request) {
	clientname := ""
	if req.Method == "POST" {
		req.ParseForm()
		for k,v := range req.Form {
			fmt.Println("Form", k, "is", v)
		}
		xx,ok := req.Form["client"]
		if ok {
			fmt.Println("Client is", xx)
			clientname = xx[0]
		} else {
			fmt.Println("No client name.")
		}
		_,ok = req.Form["bidfor"]
		if ok {
			// We already have the hands, and can short-circuit right now
			bidForMeNow(c, req, clientname)
			return
		}
		t := hands[clientname]
		newh := bridge.Hand(0)
		for sv:=uint(0); sv < 5; sv++ {
			if hnd,ok := req.Form["southhand" + bridge.SuitLetter[sv]]; ok {
				s := bridge.Suit(t[bridge.South] >> (sv*8))
				fmt.Sscan(hnd[0], &s)
				fmt.Println("Got southhand", bridge.SuitLetter[sv], "of")
				fmt.Println(s)
				newh += bridge.Hand(s) << (sv*8)
			}
		}
		t[bridge.South] = newh
		hands[clientname] = t
		if newh.Length() == 13 {
			bidForMeNow(c, req, clientname)
			return
		}
	}
	if clientname == "" {
		last_client = (last_client + 1) % max_clients
		clientname = fmt.Sprintf("client=%d", last_client)
	}
	defer header(c, req, "Enter your hand")()
	askhand(c, req, bridge.South, clientname)
}


func askhand(c io.Writer, req *http.Request, seat bridge.Seat, clientname string) os.Error {
	fmt.Fprintf(c, `<form method="post" action="%s"><div>`, req.URL.Path)
	t := hands[clientname]
	if t[seat] != 13 {
		fmt.Fprintln(c, `<table>`)
		form := `<tr><td>%s</td><td><input type="text" name="%s" value="%s" /></td></tr>`+"\n"
		for sv:=uint(bridge.Spades); sv <= bridge.Spades; sv-- {
			fmt.Fprintf(c, form, bridge.SuitColorHTML[sv], "southhand"+bridge.SuitLetter[sv], bridge.Suit(t[seat] >> (8*sv)).String())
		}
		fmt.Fprintln(c, `</table>`)
	} else {
		fmt.Fprintln(c, "Hand already entered!")
	}

	fmt.Fprintln(c, `Dealer: `)
	for s:=bridge.Seat(0); s<4; s++ {
		fmt.Fprintf(c, `<input type="radio" name="dealer" value="%s"`, s.String())
		if (dealer[clientname] == s) {
			fmt.Fprint(c, ` checked="checked"`)
		}
		fmt.Fprintln(c, `/> `, s.String())
	}

	fmt.Fprintln(c, `<br/><input type="submit" value="Enter" />`)
	fmt.Fprintf(c, `<input type="hidden" name="client" value="%s" />`, clientname)
	fmt.Fprintln(c, `</div></form>`)
	return nil
}
