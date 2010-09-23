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
		askhands(c, clientname)
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
		for seat:=bridge.Seat(0); seat<4; seat++ {
			if t[seat].Length() != 13 {
				newh := bridge.Hand(0)
				for sv:=uint(0); sv < 5; sv++ {
					if hnd,ok := req.Form[seat.String()+" hand "+bridge.SuitLetter[sv]]; ok {
						s := bridge.Suit(t[seat] >> (sv*8))
						fmt.Sscan(hnd[0], &s)
						newh += bridge.Hand(s) << (sv*8)
					}
				}
				t[seat] = newh
			}
		}
		hands[clientname] = t

		// Figure out if we can fill out the last hand:
		numhands := 0
		var missinghand bridge.Seat
		tothands := bridge.Hand(0)
		for x:=bridge.Seat(bridge.South);x<4;x++ {
			if t[x].Length() == 13 {
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
	askhands(c, clientname)
}


func askhands(c io.Writer, clientname string) os.Error {
	fmt.Fprintln(c, `<form method="post"><fieldset><table><tr><td></td><td align="center">`)
	askonehand(c, bridge.North, clientname)
	fmt.Fprintln(c, `</td></tr><tr><td align="center">`)
	askonehand(c, bridge.West, clientname)
	fmt.Fprintln(c, `</td><td></td><td align="center">`)
	askonehand(c, bridge.East, clientname)
	fmt.Fprintln(c, `</td></tr><tr><td></td><td align="center">`)
	askonehand(c, bridge.South, clientname)
	fmt.Fprintln(c, `</td></tr></table>`)
	fmt.Fprintln(c, `<input type="submit" value="Enter" />`)
	fmt.Fprintf(c, `<input type="hidden" name="client" value="%s" />`, clientname)
	fmt.Fprintln(c, `</fieldset></form>`)
	return nil
}

func askonehand(c io.Writer, seat bridge.Seat, clientname string) os.Error {
	t := hands[clientname]
	if t[seat].Length() != 13 {
		fmt.Fprintln(c, `<table>`)
		form := `<tr><td>%s</td><td><input type="text" name="%s" value="%s" /></td></tr>`+"\n"
		for sv:=uint(bridge.Spades); sv <= bridge.Spades; sv-- {
			fmt.Fprintf(c, form, bridge.SuitColorHTML[sv], seat.String()+" hand "+bridge.SuitLetter[sv], bridge.Suit(t[seat] >> (8*sv)).String())
		}
		fmt.Fprintln(c, `</table>`)
	} else {
		fmt.Fprintln(c, seat.String(), "known")
	}
	return nil
}
