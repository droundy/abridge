package main

import (
	"fmt"
	"io"
	"os"
  "http"
	"log"
	"github.com/droundy/bridge"
)

func main() {
	fmt.Println("This is only a test...")
	
	http.HandleFunc("/bidder", bidder)
	http.HandleFunc("/bidforme", bidforme)
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
		fmt.Fprintf(c, `<td align="center">%d‑%.1f‑%d</td>`, hcp.Min, hcp.Mean, hcp.Max)
	}
	fmt.Fprintln(c, `</tr><tr><td>Points</td>`)
	for i:=0; i<4; i++ {
		pts := ts.PointCount(bridge.Seat(i))
		fmt.Fprintf(c, `<td align="center">%d‑%.1f‑%d</td>`, pts.Min, pts.Mean, pts.Max)
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
	case 'S': out += bridge.SuitColorHTML[bridge.Spades]
	case 'C': out += bridge.SuitColorHTML[bridge.Clubs]
	case 'N': out += bridge.SuitColorHTML[bridge.NoTrump]
	case 'D': out += bridge.SuitColorHTML[bridge.Diamonds]
	case 'H': out += bridge.SuitColorHTML[bridge.Hearts]
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
		fmt.Fprintln(c, `<td align="center"><span style="color:white">.</span></td>`)
	}
	fmt.Fprintln(c, `</tr></table></div>`)
	return nil
}
