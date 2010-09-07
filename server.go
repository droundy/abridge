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

var max_clients = 1024
var last_client = 0
var bids = make(map[string]string)
var dealer = make(map[string]bridge.Seat)

// hello world, the web server
func helloServer(c *http.Conn, req *http.Request) {
	clientname, ok := req.Header["Cookie"]
	if !ok {
		fmt.Println("Got a new client!")
		last_client = (last_client + 1) % max_clients
		clientname = fmt.Sprintf("client=%d", last_client)
		c.SetHeader("Set-Cookie", clientname)
	} else {
		fmt.Println("Welcome back,", clientname)
	}
	if req.Method == "POST" {
		req.ParseForm()
		bid,ok := req.Form["bid"]
		if ok && len(bid) == 1 && len(bid[0]) == 2 {
			bids[clientname] = bids[clientname] + bid[0]
		} else if _,ok := req.Form["refresh"]; !ok {
			bids[clientname] = ""
			dealer[clientname] = (dealer[clientname] + 1) % 4
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
	}
	fmt.Println(req.Method, req.RawURL)
	c.SetHeader("Content-Type", "text/html")
	io.WriteString(c, `
<html>
<head>

<meta http-equiv="content-type" content="text/html; charset=utf-8">
<title>Bridge bidding</title>

</head>

<body>
`)
	fmt.Fprintln(c, "<table><tr><td>")
	bidbox(c, dealer[clientname], bids[clientname])
	io.WriteString(c, `</td><td>`)
	cs,_ := analyzebids(c, dealer[clientname], bids[clientname])
	io.WriteString(c, `</td><td>`)
	showbids(c, dealer[clientname], bids[clientname])
	fmt.Fprintln(c, `</td></tr></table>`)
	showconventions(c, bids[clientname], cs)
	fmt.Fprintln(c, `</body></html>`)
}

func showconventions(c io.Writer, bids string, conventions []string) os.Error {
	fmt.Fprintln(c, `<br/>`)
	for i,cc := range conventions {
		fmt.Fprintln(c, bids[2*i:2*i+2], "=", cc, "<br/>")
	}
	return nil
}

func analyzebids(c io.Writer, dealer bridge.Seat, bids string) ([]string, os.Error) {
	fmt.Fprintln(c, "<pre>")
	//ts, ntry := bridge.ShuffleValidTables(lastbidder, bids, 100)
	ts,conventions := bridge.GetValidTables(dealer, bids, 100)
	fmt.Fprintln(c, ts)
	//fmt.Fprintf(c, "\nProbability = %.2f%%\n", 100/ntry)
	fmt.Fprintln(c, `</pre><table><tr><td></td>`)
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
	fmt.Fprintln(c, `</tr></table>`)
	return conventions, nil
}

func showbids(c io.Writer, dealer bridge.Seat, bids string) os.Error {
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

func bidbox(c io.Writer, dealer bridge.Seat, bids string) os.Error {
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
