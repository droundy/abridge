package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
  "http"
	"log"
	"github.com/droundy/bridge"
)

func main() {
	fmt.Println("This is only a test...")
	
	http.HandleFunc("/bidder", bidder)
	http.HandleFunc("/about", about)
	http.HandleFunc("/favicon.ico", faviconServer)
	http.HandleFunc("/style.css", styleServer)
	http.HandleFunc("/", analyzer)
	err := http.ListenAndServe("0.0.0.0:12345", nil)
	if err != nil {
		log.Exit("ListenAndServe: ", err.String())
	}
}

var max_clients = 1024
var last_client = 0
var bids = make(map[string]string)
var hands = make(map[string]bridge.Table)
var dealer = make(map[string]bridge.Seat)

// Analyze a bid sequence...
func analyzer(c *http.Conn, req *http.Request) {
	clientname := ""
	if req.Method == "POST" {
		req.ParseForm()
		xx,ok := req.Form["client"]
		if ok {
			clientname = xx[0]
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
			if _,ok := req.Form["undo"]; ok && len(bids[clientname]) >= 2 {
				bids[clientname] = bids[clientname][0:len(bids[clientname])-2]
			}
			if _,ok := req.Form["clear"]; ok {
				bids[clientname] = ""
				dealer[clientname] = (dealer[clientname] + 1) % 4
			}
			if _,ok := req.Form["refresh"]; ok {
				bridge.ClearBid(bids[clientname])
			}
			if d, ok := req.Form["dealer"]; ok && len(d) == 1 {
				dealer[clientname] = bridge.StringToSeat(d[0])
			}
			for k,v := range req.Form {
				fmt.Println(k, v)
			}
			for k,v := range req.Header {
				fmt.Println("Header: ", k, v)
			}
		} else {
			fmt.Println("No client name.")
		}
	}
	if clientname == "" {
		last_client = (last_client + 1) % max_clients
		clientname = fmt.Sprintf("client=%d", last_client)
	}
	fmt.Println(req.Method, req.RawURL)
	defer header(c, req, "Bridge bidding")()

	bidbox(c, clientname, 0) // the second argument is bogus (but allows reusing bidbox)
	ts := bridge.GetValidTables(dealer[clientname], bids[clientname], 100)
	fmt.Fprintln(c, ts.HTML())
	printstatistics(c, ts)
	showbids(c, clientname)
	showconventions(c, clientname, ts.Conventions)
	fmt.Fprintln(c, ts.ExampleHTML())
}

func showconventions(c io.Writer, clientname string, conventions []string) os.Error {
	if len(conventions) == 0 {
		return nil
	}
	fmt.Fprintln(c, `<div id="conventions"><h3>Conventions</h3>`)
	for i,cc := range conventions {
		fmt.Fprintln(c, htmlbid(bids[clientname][2*i:2*i+2]), "=", cc, "<br/>")
	}
	fmt.Fprintln(c, `</div>`)
	return nil
}

func printstatistics(c io.Writer, ts *bridge.Ensemble) os.Error {
	fmt.Fprintln(c, `<div id="statistics">`)
	fmt.Fprintln(c, `<table><tr><td></td>`)
	fmt.Fprintln(c, `<td align="center">South</td><td align="center">West</td><td align="center">North</td><td align="center">East</td>`)
	fmt.Fprintln(c, `</tr><tr><td>HCP</td>`)
	for i:=0; i<4; i++ {
		hcp := ts.HCP(bridge.Seat(i))
		fmt.Fprintf(c, `<td align="center">%d-%.1f-%d</td>`, hcp.Min, hcp.Mean, hcp.Max)
	}
	fmt.Fprintln(c, `</tr><tr><td>Points</td>`)
	for i:=0; i<4; i++ {
		pts := ts.PointCount(bridge.Seat(i))
		fmt.Fprintf(c, `<td align="center">%d-%.1f-%d</td>`, pts.Min, pts.Mean, pts.Max)
	}
	fmt.Fprintln(c, `</tr></table></div>`)
	return nil
}

func htmlbid(bid string) string {
	if len(bid) != 2 {
		panic("bad bidlength in htmlbid") // this means a bug!
	}
	out := bid[0:1]
	switch bid[1] {
	case 'S': out += bridge.SuitHTML[bridge.Spades]
	case 'C': out += bridge.SuitHTML[bridge.Clubs]
	case 'N': out += bridge.SuitHTML[bridge.NoTrump]
	case 'D': out = `<font color="#ff0000">` + out + bridge.SuitHTML[bridge.Diamonds]+`</font>`
	case 'H': out = `<font color="#ff0000">` + out + bridge.SuitHTML[bridge.Hearts]+`</font>`
	default: out = bid
	}
	return out
}

func showbids(c io.Writer, clientname string) os.Error {
	fmt.Fprintln(c, `<div id="bidtable"><table><tr><td>South</td><td>West</td><td>North</td><td>East</td></tr><tr>`)
	for i:=bridge.Seat(0); i<dealer[clientname]; i++ {
		fmt.Fprintln(c, `<td align="center">-</td>`)		
	}
	for i:=bridge.Seat(0); i<bridge.Seat(len(bids[clientname])/2); i++ {
		if (i + dealer[clientname]) & 3 == 0 {
			fmt.Fprintln(c, `</tr><tr>`)
		}
		fmt.Fprintln(c, `<td align="center">`, htmlbid(bids[clientname][2*i:2*i+2]), `</td>`)
	}
	for i:=bridge.Seat(len(bids[clientname])/2); i<5; i++ {
		if (i + dealer[clientname]) & 3 == 0 {
			fmt.Fprintln(c, `</tr><tr>`)
		}
		fmt.Fprintln(c, `<td align="center"><font color="#FFFFFF">.</font></td>`)
	}
	fmt.Fprintln(c, `</tr></table></div>`)
	return nil
}

func bidbox(c io.Writer, clientname string, bidfor bridge.Seat) os.Error {
	fmt.Fprintln(c, `<div id="bidbox"><form method=post>`)
	candouble := regexp.MustCompile(".[CDHSN]( P P)?$").MatchString(bids[clientname])
	canredouble := regexp.MustCompile(" X( P P)?$").MatchString(bids[clientname])
	fmt.Fprintln(c, `<table><tr>
<td><input type="submit" name="bid" value=" P" /></td>`)
	if candouble {
		fmt.Fprintln(c, `<td align="center"><input type="submit" name="bid" value=" X" /></td>`)
	} else {
		fmt.Fprintln(c, `<td align="center"><font color="#aaaaaa">X</font></td>`)
	}
	if canredouble {
		fmt.Fprintln(c, `<td align="center"><input type="submit" name="bid" value="XX" /></td></tr>`)
	} else {
		fmt.Fprintln(c, `<td align="center"><font color="#aaaaaa">XX</font></td></tr>`)
	}
	bv, bs := bridge.LastBid(bids[clientname])
	for bidlevel:=1;bidlevel<8;bidlevel++ {
		fmt.Fprintln(c, "<tr>")
		for sv:=bridge.Color(bridge.Clubs); sv<=bridge.NoTrump; sv++ {
			fmt.Fprint(c, `<td align="center">`)
			if bidlevel > bv || (bidlevel == bv && sv > bs) {
				if sv > bridge.Clubs && sv < bridge.Spades {
					fmt.Fprintf(c, `<input type="submit" name="bid" style="color:red" value="%d%v" /></td>`,
						bidlevel, bridge.SuitHTML[sv])
				} else {
					fmt.Fprintf(c, `<input type="submit" name="bid" value="%d%v" /></td>`,
						bidlevel, bridge.SuitHTML[sv])
				}
			} else {
				if sv > bridge.Clubs && sv < bridge.Spades {
					fmt.Fprintf(c, `<font color="#ffaaaa">%d%v</font></td>`,
						bidlevel, bridge.SuitHTML[sv])
				} else {
					fmt.Fprintf(c, `<font color="#aaaaaa">%d%v</font></td>`,
						bidlevel, bridge.SuitHTML[sv])
				}
			}
		}
		fmt.Fprintln(c, "</tr>")
	}
	fmt.Fprintln(c, `</table><input type="submit" name="clear" value="Next hand" />`)
	fmt.Fprintln(c, `<input type="submit" name="undo" value="Undo" />`)
	fmt.Fprintln(c, `<input type="submit" name="refresh" value="Refresh" />`)
	fmt.Fprintln(c, `<br/>`)
	fmt.Fprintln(c, `<br/>`)
	var seats = []string{"S", "W", "N", "E"}
	for s,v := range seats {
		if dealer[clientname] != bridge.Seat(s) && bids[clientname] == "" {
			fmt.Fprintf(c, `<input type="submit" name="dealer" value="%s" />`, v)
		} else {
			fmt.Fprintf(c, `<input type="submit" disabled="1" value="%s" />`, v)
		}
	}
	fmt.Fprintf(c, `<input type="hidden" name="bidfor" value="%d" />`, int(bidfor))
	fmt.Fprintf(c, `<input type="hidden" name="client" value="%s" />`, clientname)
	fmt.Fprintln(c, `</form></div>`)
	return nil
}

func faviconServer(c *http.Conn, req *http.Request) {
	// The following is a literal storing my favicon...
	fmt.Fprint(c, "\x00\x00\x01\x00\x01\x00\x10\x10\x10\x00\x00\x00\x00\x00h\x03\x00\x00\x16\x00\x00\x00(\x00\x00\x00\x10\x00\x00\x00 \x00\x00\x00\x01\x00\x18\x00\x00\x00\x00\x00\x00\x03\x00\x00\x12\v\x00\x00\x12\v\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xdc\xdc\u0703\x81\x83\x83\x81\x83\x94\x81\x94\x83\x81\x83\x83\x81\x83\x83\x81\x83\x94\x81\x94\x83\x81\x83\x83\x81\x83\x83\x81\x83\x94\x81\x94\x83\x81\x83\x83\x81\x83\x83\x81\x83\xdc\xdc\u0703\x81\x83\xff\xff\xff\xff\xff\xff\xff\xff\xff\xc5\xc2\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff101\xff\xff\xff\xff\xff\xff\xff\xff\xff\x83\x81\x83\x83}\x83\xff\xff\xff\xff\xff\xff\xff\xff\xff\x00\x00\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xc5\xc2\xc5101\xff\xff\xff101\xff\xff\xff101\xc5\xc2\u0143\x81\x83\x8b\x85\x8b\xff\xff\xff\xff\xff\xff\x00\x00\xff\x00\x00\xff\x00\x00\xff\xff\xff\xff\xff\xff\xff\xcd\xc6\xcd101101101101101\xcd\xc6\u0343\x81\x83{}{\xff\xff\xff\x00\x00\xff\x00\x00\xff\x00\x00\xff\x00\x00\xff\x00\x00\xff\xff\xff\xff\xff\xff\xff101101101101101\xff\xff\xff\x94\x81\x94\x8b\x85\x8b\xff\xff\xff\x00\x00\xff\x00\x00\xff\x00\x00\xff\x00\x00\xff\x00\x00\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff101101101\xff\xff\xff\xff\xff\xff\x83\x81\x83\x83}\x83\xff\xff\xff\x00\x00\xff\x00\x00\xff\xff\xff\xff\x00\x00\xff\x00\x00\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff101\xff\xff\xff\xff\xff\xff\xff\xff\xff\x83\x81\x83\x8b\x85\x8b\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\x83\x81\x83{}{\xff\xff\xff\xff\xff\xff\xcd\xc6\xcd101\xc5\xc6\xc5\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\x00\x00\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\x94\x81\x94\x8b\x85\x8b\xff\xff\xff\xff\xff\xff\xff\xff\xff101\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xc5\xc2\xff\x00\x00\xff\xc5\xc2\xff\xff\xff\xff\xff\xff\xff\x83\x81\x83\x83}\x83\xff\xff\xff101\xc5\xc6\xc5101\xcd\xc6\xcd101\xff\xff\xff\xff\xff\xff\xc5\xc2\xff\x00\x00\xff\x00\x00\xff\x00\x00\xff\xc5\xc2\xff\xff\xff\xff\x83\x81\x83\x8b\x85\x8b\xa4\xa1\xa4101101101101101\xa4\xa1\xa4\xff\xff\xff\x00\x00\xff\x00\x00\xff\x00\x00\xff\x00\x00\xff\x00\x00\xff\xff\xff\xff\x83\x81\x83{}{\xff\xff\xff101\xcd\xc6\xcd101\xc5\xc6\xc5101\xff\xff\xff\xff\xff\xff\xc5\xc2\xff\x00\x00\xff\x00\x00\xff\x00\x00\xff\xc5\xc2\xff\xff\xff\xff\x94\x81\x94\x8b\x85\x8b\xff\xff\xff\xff\xff\xff101101101\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xc5\xc2\xff\x00\x00\xff\xc5\xc2\xff\xff\xff\xff\xff\xff\xff\x83\x81\x83\x83}\x83\xff\xff\xff\xff\xff\xff\xff\xff\xff\xa4\xa1\xa4\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\x00\x00\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xa4\x99\xa4\xdc\xdc\u0703\x81\x83\x8b\x85\x8b\x83\x81\x83\x8b\x85\x8b\x83\x81\x83\x8b\x85\x8b\x83\x81\x83\x8b\x85\x8b\x83\x81\x83\x8b\x85\x8b\x83\x81\x83\x8b\x85\x8b\x83\x81\x83\xb4\xae\xb4\xdc\xdc\u0700\x01\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x00\x00\xff\xff\x80\x01\xff\xff")
}

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
	fmt.Fprintln(c, `<form method=post><table><tr><td></td><td align="center">`)
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
	fmt.Fprintln(c, `</form>`)
	return nil
}
