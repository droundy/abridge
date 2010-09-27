package main

import (
	"os"
	"io"
	"fmt"
	"http"
	"github.com/droundy/bridge"
)

// Bid the fourth hand...
func bidfor(c *http.Conn, req *http.Request, dat *TransitoryData) {
	bid,ok := req.Form["bid"]
	switch {
	case !ok || len(bid) != 1:
	case bid[0][1:] == bridge.SuitHTML[bridge.Clubs]:
		dat.Bids = dat.Bids + bid[0][0:1] + "C"
	case bid[0][1:] == bridge.SuitHTML[bridge.Diamonds]:
		dat.Bids = dat.Bids + bid[0][0:1] + "D"
	case bid[0][1:] == bridge.SuitHTML[bridge.Hearts]:
		dat.Bids = dat.Bids + bid[0][0:1] + "H"
	case bid[0][1:] == bridge.SuitHTML[bridge.Spades]:
		dat.Bids = dat.Bids + bid[0][0:1] + "S"
	case bid[0][1:] == bridge.SuitHTML[bridge.NoTrump]:
		dat.Bids = dat.Bids + bid[0][0:1] + "N"
	case len(bid[0]) == 2:
		dat.Bids = dat.Bids + bid[0]
	default:
		fmt.Println("I don't recognize", bid[0])
	}

	if _,ok := req.Form["undo"]; ok && len(dat.Bids) >= 2 {
		dat.Bids = dat.Bids[0:len(dat.Bids)-2]
		bidder := (dat.Dealer + bridge.Seat(len(dat.Bids)/2)) % 4
		if bidder == dat.Bidfor {
			dat.Bids = dat.Bids[0:len(dat.Bids)-2]
		}
	}
	if _,ok := req.Form["clear"]; ok {
		dat.Bids = ""
		dat.Dealer = (dat.Dealer + 1) % 4
		for i := range dat.Hands {
			dat.Hands[i] = 0
		}
		dat.AmBidding = false
		defer header(c, dat, "Enter your next hand")()
		askhands(c, req, dat)
		return
	}

	bidder := (dat.Dealer + bridge.Seat(len(dat.Bids)/2)) % 4
	cc := *getSettings(req).Card()
	ts := bridge.GetValidTables(dat.Dealer, dat.Bids, 100, cc)
	if bidder == dat.Bidfor {
		fmt.Println("Bids are:", dat.Bids)
		fmt.Println("Table is:")
		fmt.Println(dat.Hands)
		newbid, conv := bridge.PickBid(dat.Hands[dat.Bidfor], bidder, dat.Bids, cc, ts)
		dat.Bids += newbid
		fmt.Println("Bid using", conv)
		ts = bridge.GetValidTables(dat.Dealer, dat.Bids, 100, cc)
	}
	defer header(c, dat, "Bridge bidder")()

	fmt.Fprintln(c, `<table width="100%"><tr>`)
	fmt.Fprintln(c, `<td rowspan="1">`)
	bidbox(c, req, dat) // the second argument is bogus (but allows reusing bidbox)
	stats := bridge.GetValidTables(dat.Dealer, dat.Bids, 100, *getSettings(req).Card())	
	fmt.Fprintln(c, `</td><td rowspan="2">`)
	fmt.Fprintln(c, stats.HTML())
	fmt.Fprintln(c, `</td><td rowspan="3">`)
	showbids(c, dat)
	fmt.Fprintln(c, `</td></tr></table>`)
	showconventions(c, dat, stats.Conventions)
	fmt.Fprintln(c, stats.ExampleHTML())
	printstatistics(c, ts)
}

// Bid the fourth hand...
func bidder(c *http.Conn, req *http.Request) {
	dat := getTransitoryData(req)

	for k,v := range req.Form {
		fmt.Println("Form", k, "is", v)
	}
	if d, ok := req.Form["dealer"]; ok && len(d) == 1 {
		dat.Dealer = bridge.StringToSeat(d[0])
	}
	if dat.AmBidding {
		// We already have the hands, and can short-circuit right now
		bidfor(c, req, dat)
		return
	}
	t := dat.Hands
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
	dat.Hands = t

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

	dat.Hands = t
	fmt.Println("Table is:")
	fmt.Println(t)

	if numhands > 2 {
		// We're ready to start bidding!
		dat.Bidfor = missinghand
		dat.AmBidding = true
		bidfor(c, req, dat)
		return
	}
	defer header(c, dat, "Enter your hand")()
	askhands(c, req, dat)
}


func askhands(c io.Writer, req *http.Request, dat *TransitoryData) os.Error {
	fmt.Fprintf(c, `<div><table><tr><td></td><td align="center">`)
	askonehand(c, bridge.North, dat)
	fmt.Fprintln(c, `</td></tr><tr><td align="center">`)
	askonehand(c, bridge.West, dat)
	fmt.Fprintln(c, `</td><td></td><td align="center">`)
	askonehand(c, bridge.East, dat)
	fmt.Fprintln(c, `</td></tr><tr><td></td><td align="center">`)
	askonehand(c, bridge.South, dat)
	fmt.Fprintln(c, `</td></tr></table>`)

	fmt.Fprintln(c, `Dealer: `)
	for s:=bridge.Seat(0); s<4; s++ {
		fmt.Fprintf(c, `<input type="radio" name="dealer" value="%s"`, s.String())
		if (dat.Dealer == s) {
			fmt.Fprint(c, ` checked="checked"`)
		}
		fmt.Fprintln(c, `/> `, s.String())
	}
	fmt.Fprintln(c, `<br/>`)
	fmt.Fprintln(c, `<input type="submit" value="Enter" />`)
	fmt.Fprintln(c, `</div>`)
	return nil
}

func askonehand(c io.Writer, seat bridge.Seat, dat *TransitoryData) os.Error {
	t := dat.Hands
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
