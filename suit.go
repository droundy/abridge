package bridge

import (
	"fmt"
	"os"
)

type Suit byte
type Points byte

const (
	Jack Suit = 1 << iota
	Queen
	King
	Ace
)

var (
	stringTable [224]string
	HCP [224]Points
	DistPoints [224]Points
	PointCount [224]Points
)

func init() {
	for s:=Suit(0); s<224; s++ {
		l := int(s >> 4)
		// First, let's generate the friendly output...
		str := make([]byte, int(s >> 4))
		for i := range str {
			str[i] = 'x'
		}
		cnum := 0
		if s & Ace != 0 && cnum < l {
			str[cnum] = 'A'
			cnum++
		}
		if s & King != 0 && cnum < l {
			str[cnum] = 'K'
			cnum++
		}
		if s & Queen != 0 && cnum < l {
			str[cnum] = 'Q'
			cnum++
		}
		if s & Jack != 0 && cnum < l {
			str[cnum] = 'J'
			cnum++
		}
		stringTable[s] = string(str)
		// Now let's count high card points...
		HCP[s] = Points(((Ace & s) >> 1) + ((Jack + Queen) & s) + ((King & s) >> 1) + ((King & s) >> 2))
		// Distributional points are easy...
		if l < 3 {
			DistPoints[s] = Points(3 - l)
		} else {
			DistPoints[s] = 0
		}
		// "Reasonable" point count
		PointCount[s] = HCP[s] + DistPoints[s]
		if l < 3 && Jack & s != 0 {
			PointCount[s] -= 1
		} else if l == 2 && Queen & s != 0 {
			PointCount[s] -= 1
		} else if l == 1 && Queen & s != 0 {
			PointCount[s] -= 2
		} else if l == 1 && King & s != 0 {
			PointCount[s] -= 2
		}
	}
}

// We implement fmt.Formatter interface just so we can left-align the cards.
func (s Suit) Format(f fmt.State, c int) {
	x := stringTable[s]
	f.Write([]byte(x))
	if w,ok := f.Width(); ok {
		for i:=0; i<w-len(x); i++ {
			f.Write([]byte{' '})
		}
	}
}
func (s Suit) String() string {
	return stringTable[s]
}
func ReadSuit(x string) Suit {
	if x == "~" {
		return Suit(0) // special name for void suit
	}
	out := Suit(len(x) << 4)
	for _,c := range x {
		switch c {
		case 'A','a': out = out | Ace // Be cautious in case duplicate cards show up
		case 'K','k': out = out | King
		case 'Q','q': out = out | Queen
		case 'J','j': out = out | Jack
		}
	}
	return out
}

func (s *Suit) Scan(st fmt.ScanState, x int) os.Error {
	str, e := st.Token()
	*s = ReadSuit(str)
	return e
}
