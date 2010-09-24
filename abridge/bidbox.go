package main

import (
	"fmt"
	"regexp"
	"os"
	"io"
	"http"
	"github.com/droundy/bridge"
)

func bidbox(c io.Writer, req *http.Request, clientname string, bidfor bridge.Seat) os.Error {
	
	fmt.Fprintf(c, `<form method="post" action="%s"><div id="bidbox">`, req.URL.Path)
	candouble := regexp.MustCompile(".[CDHSN]( P P)?$").MatchString(bids[clientname])
	canredouble := regexp.MustCompile(" X( P P)?$").MatchString(bids[clientname])
	fmt.Fprintln(c, `<table><tr>
<td><input type="submit" name="bid" value=" P" /></td>`)
	if candouble {
		fmt.Fprintln(c, `<td align="center"><input type="submit" name="bid" value=" X" /></td>`)
	} else {
		fmt.Fprintln(c, `<td align="center"><span class="disablednotrump">X</span></td>`)
	}
	if canredouble {
		fmt.Fprintln(c, `<td align="center"><input type="submit" name="bid" value="XX" /></td></tr>`)
	} else {
		fmt.Fprintln(c, `<td align="center"><span class="disablednotrump">XX</span></td></tr>`)
	}
	bv, bs := bridge.LastBid(bids[clientname])
	for bidlevel:=1;bidlevel<8;bidlevel++ {
		fmt.Fprintln(c, "<tr>")
		for sv:=bridge.Color(bridge.Clubs); sv<=bridge.NoTrump; sv++ {
			fmt.Fprint(c, `<td align="center">`)
			if bidlevel > bv || (bidlevel == bv && sv > bs) {
				fmt.Fprintf(c, `<input type="submit" name="bid" class="%s" value="%d%v" /></td>`,
					bridge.SuitName[sv], bidlevel, bridge.SuitHTML[sv])
			} else {
				fmt.Fprintf(c, `<span class="disabled%s">%d%v</span></td>`,
					bridge.SuitName[sv], bidlevel, bridge.SuitHTML[sv])
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
			fmt.Fprintf(c, `<input type="submit" disabled="disabled" value="%s" />`, v)
		}
	}
	fmt.Fprintf(c, `<input type="hidden" name="bidfor" value="%d" />`, int(bidfor))
	fmt.Fprintf(c, `<input type="hidden" name="client" value="%s" />`, clientname)
	fmt.Fprintln(c, `</div></form>`)
	return nil
}
