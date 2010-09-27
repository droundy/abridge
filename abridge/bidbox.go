package main

import (
	"fmt"
	"regexp"
	"os"
	"io"
	"http"
	"github.com/droundy/bridge"
)

func bidbox(c io.Writer, req *http.Request, dat *TransitoryData) os.Error {
	
	fmt.Fprintf(c, `<div id="bidbox">`)
	candouble := regexp.MustCompile(".[CDHSN]( P P)?$").MatchString(dat.Bids)
	canredouble := regexp.MustCompile(" X( P P)?$").MatchString(dat.Bids)
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
	bv, bs := bridge.LastBid(dat.Bids)
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
		if dat.Dealer != bridge.Seat(s) && dat.Bids == "" {
			fmt.Fprintf(c, `<input type="submit" name="dealer" value="%s" />`, v)
		} else {
			fmt.Fprintf(c, `<input type="submit" disabled="disabled" value="%s" />`, v)
		}
	}
	fmt.Fprintln(c, `</div>`)
	return nil
}
