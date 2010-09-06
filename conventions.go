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
	score func(bidder Seat, h Hand, ms []string, e Ensemble) (s Score, nothandled bool)
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
	CheapRebid,
	MajorSupport, MajorInvitation, OneNT, TwoNT, Gambling3NT,
	OneLevelOvercall, PassOvercall, Natural }

func subBids(dealer Seat, bid string) (seats []Seat, bids []string) {
	seats = make([]Seat, len(bid)/2)
	bids = make([]string, len(bid)/2)
	for i := range bids {
		seats[i] = (dealer + Seat(i)) % 4
		bids[i] = bid[0:2*(i+1)]
	}
	return
}

func simpleScore(bidder Seat, h Hand, bid string, e Ensemble) (badness Score) {
	for _,c := range Convention {
		ms := c.match.FindStringSubmatch(bid)
		if ms != nil {
			b,unhandled := c.score(bidder, h, ms, e)
			badness += b
			if !unhandled {
				break
			}
			//fmt.Printf("Got badness %g from %s\n", b, c.name)
		}
	}
	return
}

func TableScore(t Table, bidders []Seat, bids []string, es []Ensemble) (badness Score) {
	for i, bidder := range bidders {
		badness += simpleScore(bidder, t[bidder], bids[i], es[i])
	}
	return
}

var last_ts = make([]Table, 0)
func GetValidTables(dealer Seat, bid string, num int) Ensemble {
	seats, bids := subBids(dealer, bid)
	es := make([]Ensemble, len(bids)+1) // This is the ensemble after each bid
	es[0] = make(Ensemble, num)
	for i := range es[0] {
		es[0][i] = Shuffle() // Things start out random!
	}
	for bidnum := range seats {
		es[bidnum+1] = make(Ensemble, num)
		for i,eold := range es[bidnum] {
			t := eold // Initialize ensemble based on previous bidding
			oldbadness := TableScore(t, seats[0:bidnum+1], bids[0:bidnum+1], es[0:bidnum+1])
			badness := oldbadness
			numswaps := 52*10
			beta := Score(0.01)
			for i:=0; i<numswaps && (badness > 0 || i < 52); i++ {
				// Try moving a couple of cards at a time...
				t2 := t.ShuffleCard(rand.Intn(52)).ShuffleCard(rand.Intn(52))
				b2 := TableScore(t2, seats[0:bidnum+1], bids[0:bidnum+1], es[0:bidnum+1])
				if b2 <= badness || rand.Float64() < math.Exp(float64(-beta*(b2 - badness))) {
					t = t2
					badness = b2
				}
				beta *= 1.01 // Here's our annealing schedule...
			}
			fmt.Printf("Badness %4g -> %4g\n", oldbadness, badness)
			es[bidnum+1][i] = t
		}
	}
	return es[len(bids)] // return final ensemble
}
