package bridge

import (
	"fmt"
	"os"
)

type Hand [4]Suit

const (
	Clubs = iota
	Diamonds
	Hearts
	Spades
	NoTrump
)

func (h Hand) HCP() Points {
	return HCP[h[0]] + HCP[h[1]] + HCP[h[2]] + HCP[h[3]]
}
func (h Hand) DistPoints() Points {
	return DistPoints[h[0]] + DistPoints[h[1]] + DistPoints[h[2]] + DistPoints[h[3]]
}
func (h Hand) PointCount() Points {
	return PointCount[h[0]] + PointCount[h[1]] + PointCount[h[2]] + PointCount[h[3]]
}
func (h Hand) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n", h[Spades], h[Hearts], h[Diamonds], h[Clubs])
}
func (h *Hand) Scan(st fmt.ScanState, x int) os.Error {
	str, e := st.Token()
	if e != nil {
		return e
	}
	l := len(str)
	h[Spades] = ReadSuit(str)
	str, e = st.Token()
	if e != nil {
		return e
	}
	h[Hearts] = ReadSuit(str[l:])
	l = len(str)
	str, e = st.Token()
	if e != nil {
		return e
	}
	h[Diamonds] = ReadSuit(str[l:])
	l = len(str)
	str, e = st.Token()
	if e != nil {
		return e
	}
	h[Clubs] = ReadSuit(str[l:])
	return e
}
