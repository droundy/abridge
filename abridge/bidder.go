package main

import (
	"os"
	"io"
	"fmt"
	"http"
	"github.com/droundy/bridge"
)

// Bid the fourth hand...
func bidfor(c *http.Conn, req *http.Request, clientname string, bidfor bridge.Seat) {
	bid,ok := req.Form["bid"]
	fmt.Println("All bids were", bids)
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
	fmt.Println("Bids are", bids[clientname])
	fmt.Println("allbids are", bids)

	if d, ok := req.Form["dealer"]; ok && len(d) == 1 {
		dealer[clientname] = bridge.StringToSeat(d[0])
	}
	if _,ok := req.Form["undo"]; ok && len(bids[clientname]) >= 2 {
		bids[clientname] = bids[clientname][0:len(bids[clientname])-2]
		bidder := (dealer[clientname] + bridge.Seat(len(bids[clientname])/2)) % 4
		if bidder == bidfor {
			bids[clientname] = bids[clientname][0:len(bids[clientname])-2]
		}
	}
	if _,ok := req.Form["clear"]; ok {
		bids[clientname] = ""
		dealer[clientname] = (dealer[clientname] + 1) % 4
		var x bridge.Table
		hands[clientname] = x, false
		defer header(c, req, "Enter your next hand")()
		askhand(c, clientname)
		return
	}

	bidder := (dealer[clientname] + bridge.Seat(len(bids[clientname])/2)) % 4
	ts := bridge.GetValidTables(dealer[clientname], bids[clientname], 100)
	if bidder == bidfor {
		fmt.Println("Bids are:", bids[clientname])
		fmt.Println("Table is:")
		fmt.Println(hands[clientname])
		newbid, conv := bridge.PickBid(hands[clientname][bidfor], bidder, bids[clientname], ts)
		bids[clientname] += newbid
		fmt.Println("Bid using", conv)
		ts = bridge.GetValidTables(dealer[clientname], bids[clientname], 100)
	}
	defer header(c, req, "Bridge bidder")()
	bidbox(c, clientname, bidfor)
	stats := bridge.GetValidTables(dealer[clientname], bids[clientname], 100)
	fmt.Fprintln(c, stats.HTML())
	fmt.Fprintln(c, stats.ExampleHTML())
	showbids(c, clientname)
	showconventions(c, clientname, ts.Conventions)
}

// Bid the fourth hand...
func bidder(c *http.Conn, req *http.Request) {
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
		bidforstr,ok := req.Form["bidfor"]
		if ok {
			// We already have the hands, and can short-circuit right now
			bidfor(c, req, clientname, bridge.Seat(bidforstr[0][0])-'0')
			return
		}
		t := hands[clientname]
		hnd,ok := req.Form["southhand"]
		if ok {
			fmt.Sscan(hnd[0], &t[bridge.South])
			fmt.Println("Got southhand of")
			fmt.Println(t[bridge.South])
		}
		hnd,ok = req.Form["northhand"]
		if ok {
			fmt.Sscan(hnd[0], &t[bridge.North])
			fmt.Println("Got northhand of")
			fmt.Println(t[bridge.North])
		}
		hnd,ok = req.Form["easthand"]
		if ok {
			fmt.Sscan(hnd[0], &t[bridge.East])
			fmt.Println("Got easthand of")
			fmt.Println(t[bridge.East])
		}
		hnd,ok = req.Form["westhand"]
		if ok {
			fmt.Sscan(hnd[0], &t[bridge.West])
			fmt.Println("Got westhand of")
			fmt.Println(t[bridge.West])
		}

		// Figure out if we can fill out the last hand:
		numhands := 0
		var missinghand bridge.Seat
		tothands := bridge.Hand(0)
		for x:=bridge.Seat(bridge.South);x<4;x++ {
			if t[x] != 0 {
				tothands += t[x]
				numhands++
			} else {
				missinghand = x
			}
		}
		if numhands == 3 {
			t[missinghand] = bridge.AllCards - tothands
		}

		hands[clientname] = t
		fmt.Println("Table is:")
		fmt.Println(t)

		if numhands > 2 {
			// We're ready to start bidding!
			bidfor(c, req, clientname, missinghand)
			return
		}
	}
	if clientname == "" {
		last_client = (last_client + 1) % max_clients
		clientname = fmt.Sprintf("client=%d", last_client)
	}
	defer header(c, req, "Enter your hand")()
	askhand(c, clientname)
}


func askhand(c io.Writer, clientname string) os.Error {
	fmt.Fprintln(c, `<form method="post"><fieldset><table><tr><td></td><td align="center">`)
	t := hands[clientname]
	if t[bridge.North] == 0 {
		fmt.Fprintln(c, `North: <input type="text" name="northhand" value="" />`)
	} else {
		fmt.Fprintln(c, `North known`)
	}
	fmt.Fprintln(c, `</td></tr><tr><td align="center">`)
	if t[bridge.West] == 0 {
		fmt.Fprintln(c, `West: <input type="text" name="westhand" value="" />`)
	} else {
		fmt.Fprintln(c, `West known`)
	}
	fmt.Fprintln(c, `</td><td></td><td align="center">`)
	if t[bridge.East] == 0 {
		fmt.Fprintln(c, `East: <input type="text" name="easthand" value="" />`)
	} else {
		fmt.Fprintln(c, `East known`)
	}
	fmt.Fprintln(c, `</td></tr><tr><td></td><td align="center">`)
	if t[bridge.South] == 0 {
		fmt.Fprintln(c, `South: <input type="text" name="southhand" value="" />`)
	} else {
		fmt.Fprintln(c, `South known`)
	}
	fmt.Fprintln(c, `</td></tr></table>`)
	fmt.Fprintln(c, `<input type="submit" value="Enter" />`)
	fmt.Fprintf(c, `<input type="hidden" name="client" value="%s" />`, clientname)
	fmt.Fprintln(c, `</fieldset></form>`)
	return nil
}
