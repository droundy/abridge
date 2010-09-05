package bridge

import (
	"fmt"
	"regexp"
)

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

var Convention = []BiddingRule{ Opening, Preempt, PassOpening, CheapResponse, TwoOverOne }

func TableScore(t Table, seat int, bid string) Score {
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

func ShuffleValidTable(seat int, bid string) (t Table) {
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
