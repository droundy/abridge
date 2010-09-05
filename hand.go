package bridge

import (
	"fmt"
	"os"
)

type Hand uint32

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
func (h Hand) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n", Suit(h>>24), Suit(h>>16), Suit(h>>8), Suit(h))
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
