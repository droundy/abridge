package main

import (
	"fmt"
	"regexp"
	"os"
	"io"
	"github.com/droundy/bridge"
)

func bidbox(c io.Writer, clientname string, bidfor bridge.Seat) os.Error {
	fmt.Fprintln(c, `<div id="bidbox"><form method="post">`)
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
