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
	http.HandleFunc("/settings", settings)
	http.HandleFunc("/favicon.ico", faviconServer)
	http.HandleFunc("/style.css", styleServer)
	http.HandleFunc("/", analyzer)
	err := http.ListenAndServe("0.0.0.0:12345", nil)
	if err != nil {
		log.Exit("ListenAndServe: ", err.String())
	}
}

// Analyze a bid sequence...
func analyzer(c http.ResponseWriter, req *http.Request) {
	dat := getTransitoryData(req)

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
	}
	if _,ok := req.Form["clear"]; ok {
		dat.Bids = ""
		dat.Dealer = (dat.Dealer + 1) % 4
	}
	if _,ok := req.Form["refresh"]; ok {
		bridge.ClearBid(dat.Bids)
	}
	if d, ok := req.Form["dealer"]; ok && len(d) == 1 {
		dat.Dealer = bridge.StringToSeat(d[0])
	}
	//for k,v := range req.Form {
	//	fmt.Println(k, v)
	//}
	//for k,v := range req.Header {
	//	fmt.Println("Header: ", k, v)
	//}

	fmt.Println(req.Method, req.RawURL)
	defer header(c, dat, "Bridge bidding")()

	fmt.Fprintln(c, `<table width="100%"><tr>`)
	fmt.Fprintln(c, `<td rowspan="1">`)
	bidbox(c, req, dat)
	ts := bridge.GetValidTables(dat.Dealer, dat.Bids, 100, *getSettings(req).Card())	
	fmt.Fprintln(c, `</td><td rowspan="2">`)
	fmt.Fprintln(c, ts.HTML())
	fmt.Fprintln(c, `</td><td rowspan="3">`)
	showbids(c, dat)
	fmt.Fprintln(c, `</td></tr></table>`)
	showconventions(c, dat, ts.Conventions)
	fmt.Fprintln(c, ts.ExampleHTML())
	printstatistics(c, ts)
}

func showconventions(c io.Writer, dat *TransitoryData, conventions []string) os.Error {
	if len(conventions) == 0 {
		return nil
	}
	fmt.Fprintln(c, `<div id="conventions"><h3>Conventions</h3>`)
	for i,cc := range conventions {
		fmt.Fprintln(c, htmlbid(dat.Bids[2*i:2*i+2]), "=", cc, "<br/>")
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
	for sv,sch := range bridge.SuitLetter {
		if bid[1] == sch[0] {
			return fmt.Sprintf(`<span class="%s">%c%s</span>`,
				bridge.SuitName[sv], bid[0], bridge.SuitHTML[sv])
		}
	}
	return bid;
}

func showbids(c io.Writer, dat *TransitoryData) os.Error {
	fmt.Fprintln(c, `<div id="bidtable"><table><tr><td>South</td><td>West</td><td>North</td><td>East</td></tr><tr>`)
	for i:=bridge.Seat(0); i<dat.Dealer; i++ {
		fmt.Fprintln(c, `<td align="center">-</td>`)		
	}
	for i:=bridge.Seat(0); i<bridge.Seat(len(dat.Bids)/2); i++ {
		if (i + dat.Dealer) & 3 == 0 {
			fmt.Fprintln(c, `</tr><tr>`)
		}
		fmt.Fprintln(c, `<td align="center">`, htmlbid(dat.Bids[2*i:2*i+2]), `</td>`)
	}
	fmt.Fprintln(c, `</tr></table></div>`)
	return nil
}
