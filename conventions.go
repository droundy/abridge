package bridge

import (
	"fmt"
	"regexp"
)

type Color byte

type Score float64
const (
	SuitLengthProblem Score = 100
	PointValueProblem Score = 100
	BigFudge Score = 3
	Fudge Score = 1
)
func (s Score) min(s2 Score) Score {
	if s < s2 {
		return s
	}
	return s2
}

type BiddingRule struct {
	name string
	match *regexp.Regexp
	score func(h Hand, ms []string) Score 
}

func LastBid(bid string) (val int, s Color) {
	for {
		if len(bid) < 2 {
			return 0, NoTrump
		}
		switch bid[len(bid)-2] {
		case '1': val = 1
		case '2': val = 2
		case '3': val = 3
		case '4': val = 4
		case '5': val = 5
		case '6': val = 6
		case '7': val = 7
		}
		switch xx := bid[len(bid)-1]; xx {
		case 'N': return val, NoTrump
		case 'X','P':
		default:
			return val, Color(stringToSuitNumber(string([]byte{xx})))
		}
		bid = bid[0:len(bid)-2]
	}
	return
}

var Convention = []BiddingRule{ Opening, Preempt, PassOpening, CheapResponse, TwoOverOne, MajorSupport, MajorInvitation }

func TableScore(t Table, seat Seat, bid string) Score {
	badness := Score(0)
	for bid != "" {
		for _,c := range Convention {
			ms := c.match.FindStringSubmatch(bid)
			if ms != nil {
				b := c.score(t[seat], ms)
				badness += b
				//fmt.Printf("Got badness %g from %s by seat %v\n", b, c.name, Seat(seat))
			}
		}
		bid = bid[0:len(bid)-2] // chop off most recent bid
		seat = (seat + 3) % 4 // subtract one off seat
	}
	return badness
}

func ShuffleValidTable(seat Seat, bid string) (t Table) {
	i := 0
	for {
		i++
		t = Shuffle()
		score := TableScore(t, seat, bid)
		if score == 0 {
			fmt.Println("It took", i, "tries.")
			return t
		} else {
			//fmt.Print("Bad table with score ", score, ":\n", t)
		}
	}
	return
}
