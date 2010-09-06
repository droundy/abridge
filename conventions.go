package bridge

import (
	"fmt"
	"math"
	"rand"
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

var Convention = []BiddingRule{ Opening, Preempt, PassOpening, CheapResponse, TwoOverOne,
	MajorSupport, MajorInvitation, OneNT, TwoNT, Gambling3NT }

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

func ShuffleValidTables(seat Seat, bid string, num int) (ts Ensemble, numt float64) {
	ts = make([]Table, num)
	for n := range ts {
		i := 0
		for {
			i++
			t := Shuffle()
			score := TableScore(t, seat, bid)
			if score == 0 {
				ts[n] = t
				numt += float64(i)
				break
			} else {
				//fmt.Print("Bad table with score ", score, ":\n", t)
			}
		}
	}
	numt /= float64(num)
	fmt.Println("It took on average", numt, "tries.")
	return
}

var last_ts = make([]Table, 0)
func GetValidTables(seat Seat, bid string, num int) (ts Ensemble) {
	ts = make([]Table, num)
	if len(last_ts) != num {
		last_ts = make([]Table, num)
		for n := range last_ts {
			last_ts[n] = Shuffle()
		}
	}
	for i,t := range last_ts {
		ts[i] = t
	}
	for n, t := range ts {
		e := TableScore(t, seat, bid)
		orige := e
		numswaps := 52*10
		beta := Score(0.01)
		for i:=0; i<numswaps && (e > 0 || i < 52); i++ {
			// Try moving a couple of cards at a time...
			t2 := t.ShuffleCard(rand.Intn(52)).ShuffleCard(rand.Intn(52))
			e2 := TableScore(t2, seat, bid)
			if e2 <= e || rand.Float64() < math.Exp(float64(-beta*(e2 - e))) {
				t = t2
				e = e2
			}
			beta *= 1.01 // Here's our annealing schedule...
		}
		ts[n] = t
		fmt.Printf("Energy %4g -> %4g\n", orige, e)
	}
	for i,t := range ts {
		last_ts[i] = t
	}
	return
}
