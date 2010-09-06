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
	
	http.HandleFunc("/hello", helloServer)
	http.HandleFunc("/", helloServer)
	err := http.ListenAndServe("0.0.0.0:12345", nil)
	if err != nil {
		log.Exit("ListenAndServe: ", err.String())
	}
}

var bids = ""
var dealer = bridge.Seat(bridge.South)

// hello world, the web server
func helloServer(c *http.Conn, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		bid,ok := req.Form["bid"]
		if ok && len(bid) == 1 && len(bid[0]) == 2 {
			bids = bids + bid[0]
		} else if _,ok := req.Form["refresh"]; !ok {
			bids = ""
			dealer = (dealer + 1) % 4
		}
		if d, ok := req.Form["dealer"]; ok && len(d) == 1 {
			dealer = bridge.StringToSeat(d[0])
		}
		for k,v := range req.Form {
			fmt.Println(k, v)
		}
	}
	fmt.Println(req.Method, req.RawURL)
	io.WriteString(c, `
<html>
<head>

<meta http-equiv="content-type" content="text/html; charset=utf-8">
<title>Bridge bidding</title>

</head>

<body>
`)
	fmt.Fprintln(c, "<table><tr><td>")
	showbids(c, bids)
	io.WriteString(c, `</td><td>`)
	bidbox(c, bids)
	io.WriteString(c, `</td><td>`)
	analyzebids(c, bids)
	fmt.Fprintln(c, `</td></tr></table></body></html>`)
}

func analyzebids(c io.Writer, bids string) os.Error {
	fmt.Fprintln(c, "<pre>")
	//ts, ntry := bridge.ShuffleValidTables(lastbidder, bids, 100)
	ts := bridge.GetValidTables(dealer, bids, 100)
	fmt.Fprintln(c, ts)
	//fmt.Fprintf(c, "\nProbability = %.2f%%\n", 100/ntry)
	fmt.Fprintln(c, `</pre><table><tr><td></td>`)
	fmt.Fprintln(c, `<td align="center">South</td><td align="center">West</td><td align="center">North</td><td align="center">East</td>`)
	fmt.Fprintln(c, `</tr><tr><td>HCP</td>`)
	for i:=range ts[0] {
		min, hcp, max := ts.HCP(bridge.Seat(i))
		fmt.Fprintf(c, `<td align="center">%d-%.1f-%d</td>`, min, hcp, max)
	}
	fmt.Fprintln(c, `</tr><tr><td>Points</td>`)
	for i:=range ts[0] {
		min, hcp, max := ts.PointCount(bridge.Seat(i))
		fmt.Fprintf(c, `<td align="center">%d-%.1f-%d</td>`, min, hcp, max)
	}
	fmt.Fprintln(c, `</tr></table>`)
	return nil
}

func showbids(c io.Writer, bids string) os.Error {
	fmt.Fprintln(c, `<table><tr><td>South</td><td>West</td><td>North</td><td>East</td></tr><tr>`)
	for i:=bridge.Seat(0); i<dealer; i++ {
		fmt.Fprintln(c, `<td align="center">-</td>`)		
	}
	for i:=bridge.Seat(0); i<bridge.Seat(len(bids)/2); i++ {
		if (i + dealer) & 3 == 0 {
			fmt.Fprintln(c, `</tr><tr>`)
		}
		fmt.Fprintln(c, `<td align="center">`, bids[2*i:2*i+2], `</td>`)
	}
	for i:=bridge.Seat(len(bids)/2); i<50; i++ {
		if (i + dealer) & 3 == 0 {
			fmt.Fprintln(c, `</tr><tr>`)
		}
		fmt.Fprintln(c, `<td align="center"><font color="#FFFFFF">.</font></td>`)
	}
	fmt.Fprintln(c, `</tr></table>`)
	return nil
}

func bidbox(c io.Writer, bids string) os.Error {
	fmt.Fprintln(c, `<form method=post>`)
	candouble := regexp.MustCompile(".[CDHSN]( P P)?$").MatchString(bids)
	canredouble := regexp.MustCompile(" X( P P)?$").MatchString(bids)
	fmt.Fprintln(c, `<table><tr>
<td><input type="submit" name="bid" value=" P" /></td>`)
	if candouble {
		fmt.Fprintln(c, `<td><input type="submit" name="bid" value=" X" /></td>`)
	} else {
		fmt.Fprintln(c, `<td><font color="#aaaaaa">X</font></td>`)
	}
	if canredouble {
		fmt.Fprintln(c, `<td><input type="submit" name="bid" value="XX" /></td></tr>`)
	} else {
		fmt.Fprintln(c, `<td><font color="#aaaaaa">XX</font></td></tr>`)
	}
	bv, bs := bridge.LastBid(bids)
	for bidlevel:=1;bidlevel<8;bidlevel++ {
		fmt.Fprintln(c, "<tr>")
		for sv:=bridge.Color(bridge.Clubs); sv<=bridge.NoTrump; sv++ {
			if bidlevel > bv || (bidlevel == bv && sv > bs) {
				fmt.Fprintf(c, `<td><input type="submit" name="bid" value="%d%v" /></td>`,
					bidlevel, bridge.SuitLetter[sv])
			} else {
				fmt.Fprintf(c, `<td><font color="#aaaaaa">%d%v</font></td>`,
					bidlevel, bridge.SuitLetter[sv])
			}
		}
		fmt.Fprintln(c, "</tr>")
	}
	fmt.Fprintln(c, `</table><input type="submit" value="Clear" />`)
	fmt.Fprintln(c, `<input type="submit" name="refresh" value="Refresh" />`)
	if bids == "" {
		fmt.Fprint(c, `<br/>Dealer:<br/> <input type="radio" name="dealer" value="S" /> S
<input type="radio" name="dealer" value="W" /> W
<input type="radio" name="dealer" value="N" /> N
<input type="radio" name="dealer" value="E" /> E<br />
`)
	}
	fmt.Fprintln(c, `</form>`)
	return nil
}
