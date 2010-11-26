package main

import (
	"os"
	"io"
	"fmt"
	"http"
	"github.com/droundy/abridge"
)

// Bid the fourth hand...
func bidForMeNow(c http.ResponseWriter, req *http.Request, dat *TransitoryData) {
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

	readbidbox(req, dat)
	if _,ok := req.Form["undo"]; ok && len(dat.Bids) >= 2 {
		dat.Bids = dat.Bids[0:len(dat.Bids)-2]
		bidder := (dat.Dealer + bridge.Seat(len(dat.Bids)/2)) % 4
		if bidder == bridge.South {
			dat.Bids = dat.Bids[0:len(dat.Bids)-2]
		}
	}
	if _,ok := req.Form["refresh"]; ok {
		bridge.ClearBid(dat.Bids)
	}
	if _,ok := req.Form["clear"]; ok {
		dat.Bids = ""
		dat.Dealer = (dat.Dealer + 1) % 4
		for i := range dat.Hands {
			dat.Hands[i] = 0
		}
		dat.AmBidding = false
		defer header(c, dat, "Enter your next hand")()
		askhand(c, req, bridge.South, dat)
		return
	}

	bidder := (dat.Dealer + bridge.Seat(len(dat.Bids)/2)) % 4
	cc := [2]bridge.ConventionCard{ *getSettings(req).Cards[dat.NScard], *getSettings(req).Cards[dat.EWcard] }
	ts := bridge.GetValidTables(dat.Dealer, dat.Bids, 100, cc)
	if bidder == bridge.South {
		fmt.Println("Bids are:", dat.Bids)
		fmt.Println("Table is:")
		fmt.Println(dat.Hands)
		newbid, conv := bridge.PickBid(dat.Hands[bridge.South], bidder, dat.Bids, cc[bidder&1], ts)
		dat.Bids += newbid
		fmt.Println("Bid using", conv)
		ts = bridge.GetValidTables(dat.Dealer, dat.Bids, 100, cc)
	}
	defer header(c, dat, "Bridge bidder")()

	fmt.Fprintln(c, `<table width="100%"><tr>`)
	fmt.Fprintln(c, `<td rowspan="1">`)
	bidbox(c, req, dat)
	stats := bridge.GetValidTables(dat.Dealer, dat.Bids, 100, cc)
	fmt.Fprintln(c, `</td><td rowspan="2">`)
	fmt.Fprintln(c, stats.HTML())
	fmt.Fprintln(c, `</td><td rowspan="1">`)
	fmt.Fprintln(c, dat.Hands[0].HTML("My hand"))
	fmt.Fprintln(c, `</td><td rowspan="3">`)
	showbids(c, dat)
	fmt.Fprintln(c, `</td></tr></table>`)
	showconventions(c, dat, stats.Conventions)
	fmt.Fprintln(c, stats.ExampleHTML())
	printstatistics(c, ts)
}

// Bid my hand...
func bidforme(c http.ResponseWriter, req *http.Request) {
	dat := getTransitoryData(req)

	for k,v := range req.Form {
		fmt.Println("Form", k, "is", v)
	}
	if dat.AmBidding {
		// We already have the hands, and can short-circuit right now
		bidForMeNow(c, req, dat)
		return
	}
	t := dat.Hands
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
	dat.Hands = t
	if newh.Length() == 13 {
		dat.AmBidding = true
		bidForMeNow(c, req, dat)
		return
	}
	defer header(c, dat, "Enter your hand")()
	askhand(c, req, bridge.South, dat)
}


func askhand(c io.Writer, req *http.Request, seat bridge.Seat, dat *TransitoryData) os.Error {
	fmt.Fprintf(c, `<div>`)
	t := dat.Hands
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
		if (dat.Dealer == s) {
			fmt.Fprint(c, ` checked="checked"`)
		}
		fmt.Fprintln(c, `/> `, s.String())
	}

	fmt.Fprintln(c, `<br/><input type="submit" value="Enter" />`)
	fmt.Fprintln(c, `</div>`)
	return nil
}
