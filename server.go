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
	err := http.ListenAndServe(":12345", nil)
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
		} else {
			bids = ""
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
	io.WriteString(c, "\n<table><tr><td><pre>\n\nbids are: `" + bids + "`\n")
	lastbidder := (dealer - 1 + bridge.Seat(len(bids)/2)) & 3
	fmt.Fprintln(c, "dealer is", dealer)
	fmt.Fprintln(c, "lastbidder is", lastbidder)
	io.WriteString(c, bridge.ShuffleValidTable(lastbidder, bids).String())
	io.WriteString(c, `</pre></td><td>`)
	bidbox(c, bids)
	fmt.Fprintln(c, `</td></tr></table></body></html>`)
}

func bidbox(c io.Writer, bids string) os.Error {
	fmt.Fprintln(c, `<form method=post>`)
	if bids == "" {
		fmt.Fprint(c, `
Dealer: <input type="radio" name="dealer" value="S" /> South
<input type="radio" name="dealer" value="W" /> West
<input type="radio" name="dealer" value="N" /> North
<input type="radio" name="dealer" value="E" /> East<br />
`)
	}
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
	fmt.Fprintln(c, `</table>
<input type="submit" value="Clear" />
</form>`)
	return nil
}
