package bridge

import (
	"fmt"
	"os"
)

type Hand uint32
type NumberOfCards byte

const (
	Clubs = iota
	Diamonds
	Hearts
	Spades
	NoTrump
)

func (h Hand) HCP() Points {
	return HCP[Suit(h)] + HCP[Suit(h>>8)] + HCP[Suit(h >>16)] + HCP[Suit(h>>24)]
}
func (h Hand) DistPoints() Points {
	return DistPoints[255 & h] + DistPoints[255 & (h>>8)] + DistPoints[255 & (h >>16)] + DistPoints[255 & (h>>24)]
}
func (h Hand) PointCount() Points {
	return PointCount[255 & h] + PointCount[255 & (h>>8)] + PointCount[255 & (h >>16)] + PointCount[255 & (h>>24)]
}

func (h Hand) HTML(title string) string { // FIXME!
	out := `<div class="bridgehand">`
	if title != "" {
		out += fmt.Sprintf("<strong>%s</strong>\n", title)
	}
	out += "<table>\n"
	out += fmt.Sprintf("<tr><td>   <em>%d Points</em></td></tr>\n", h.PointCount())
	out += fmt.Sprintf("<tr><td>   <em>%d HCP</em></td></tr>\n", h.HCP())
	for sv := uint(Spades); sv <= Spades; sv-- {
		s := Suit(h >> (8*sv))
		out += `<tr><td><div class="bridgecards">` + SuitColorHTML[sv] + " " + s.String()
		out += `</div></td></tr>`
	}
	out += "</table></div>\n"
	return out
}

func (h Hand) String() string {
	if h.Length() == 0 {
		return "Empty"
	} else if h.Length() == 1 {
		for sv := uint(Clubs); sv <= Spades; sv++ {
			s := Suit(h >> (8*sv))
			if s != 0 {
				suitletter := []string{"C", "D", "H", "S"}
				return fmt.Sprintf("%s: %v\n", suitletter[sv], s)
			}
		}
	}
	return fmt.Sprintf("S: %v\nH: %v\nD: %v\nC: %v\n", Suit(h>>24), Suit(h>>16), Suit(h>>8), Suit(h))
}
func (h *Hand) Scan(st fmt.ScanState, x int) os.Error {
	str, e := st.Token()
	if e != nil {
		return e
	}
	l := len(str)
	hSpades := ReadSuit(str)
	str, e = st.Token()
	if e != nil {
		return e
	}
	hHearts := ReadSuit(str[l:])
	l = len(str)
	str, e = st.Token()
	if e != nil {
		return e
	}
	hDiamonds := ReadSuit(str[l:])
	l = len(str)
	str, e = st.Token()
	if e != nil {
		return e
	}
	hClubs := ReadSuit(str[l:])
	*h = Hand(hClubs) + (Hand(hDiamonds) << 8) + (Hand(hHearts) << 16) + (Hand(hSpades) << 24)
	return e
}

func (h Hand) Length() NumberOfCards {
	return NumberOfCards(((h >> 4) & 15) + ((h >> 12)&15) + ((h >> 20)&15) + ((h >> 28)&15))
}

func (h Hand) Nth(i int) Hand {
	for sv := uint(Clubs); sv <= Spades; sv++ {
		//fmt.Print("I am looking for the ", i," card\n", h)
		numinsuit := int((h>>4) & 15)
		if i < numinsuit {
			s := Suit(h)
			for c := Ace; c > 0; c = c >> 1 {
				if s & c != 0 {
					if i == 0 {
						return Hand(c + 16) << (sv*8)
					}
					i-- // count down from this card
				}
			}
			return 16 << (sv*8) // a number card
		}
		i -= numinsuit
		h = h >> 8 // shift suits down
	}
	return 0 // Not enough cards to pick from!
}
