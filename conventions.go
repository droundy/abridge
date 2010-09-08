package bridge

import (
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

type ScoringRule struct {
	name string
	ms []string
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
	CheapRebid, Splinter,
	MajorSupport, MajorInvitation,
	OneNT, Stayman, StaymanResponse, StaymanTwo, StaymanTwoResponse, TwoNT, Gambling3NT,
	OneLevelOvercall, PreemptOvercall, PassOvercall, PassHigherOvercall, Natural }

func makeScoringRules(bid string) (out []ScoringRule) {
	out = make([]ScoringRule, 0, len(Convention))
	for _,c := range Convention {
		ms := c.match.FindStringSubmatch(bid)
		if ms != nil {
			out = out[0:len(out)+1]
			out[len(out)-1] = ScoringRule{c.name, ms, c.score}
		}
	}
	return
}

func subBids(dealer Seat, bid string) (seats []Seat, rules [][]ScoringRule) {
	seats = make([]Seat, len(bid)/2)
	rules = make([][]ScoringRule, len(bid)/2)
	for i := range rules {
		seats[i] = (dealer + Seat(i)) % 4
		rules[i] = makeScoringRules(bid[0:2*(i+1)])
	}
	return
}

func simpleScore(bidder Seat, h Hand, rule []ScoringRule, e Ensemble) (badness Score, convention string) {
	for _,r := range rule {
		b,unhandled := r.score(bidder, h, r.ms, e)
		badness += b
		if !unhandled {
			convention = r.name
			break
		}
		//fmt.Printf("Got badness %g from %s\n", b, c.name)
	}
	return
}

func TableScore(t Table, bidders []Seat, rules [][]ScoringRule, es []Ensemble) (badness Score, conventions []string) {
	conventions = make([]string, len(rules))
	for i, bidder := range bidders {
		b,c := simpleScore(bidder, t[bidder], rules[i], es[i])
		conventions[i] = c
		badness += b
	}
	return
}

// The []string output describes the bids made...
func GetValidTables(dealer Seat, bid string, num int) (*Ensemble, []string) {
	seats, rules := subBids(dealer, bid)
	es := make([]Ensemble, len(rules)+1) // This is the ensemble after each bid
	es[0] = makeEnsemble(num)
	for i := range es[0].tables {
		es[0].tables[i] = Shuffle() // Things start out random!
	}
	var conventions []string
	for bidnum := range seats {
		es[bidnum+1] = makeEnsemble(num)
		for i,eold := range es[bidnum].tables {
			t := eold // Initialize ensemble based on previous bidding
			oldbadness,cs := TableScore(t, seats[0:bidnum+1], rules[0:bidnum+1], es[0:bidnum+1])
			conventions = cs
			badness := oldbadness
			numswaps := 52*10
			beta := Score(0.01)
			maxbeta := Score(1)
			betainc := Score(math.Pow(float64(maxbeta/beta), 1/float64(numswaps)))
			for i:=0; i<numswaps && (badness > 0 || i < 52); i++ {
				// Try moving a couple of cards at a time...
				t2 := t.ShuffleCard(rand.Intn(52)).ShuffleCard(rand.Intn(52))
				b2,cs := TableScore(t2, seats[0:bidnum+1], rules[0:bidnum+1], es[0:bidnum+1])
				if b2 <= badness || rand.Float64() < math.Exp(float64(-beta*(b2 - badness))) {
					t = t2
					badness = b2
					conventions = cs
				}
				beta *= betainc // Here's our annealing schedule...
			}
			//fmt.Printf("Badness %4g -> %4g\n", oldbadness, badness)
			es[bidnum+1].tables[i] = t
		}
	}
	out := es[len(rules)] // return final ensemble
	return &out, conventions
}
