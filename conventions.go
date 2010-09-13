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
	SuitLengthProblem Score = 1000
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
	mkscore func (bidder Seat, ms []string, e *Ensemble) (score func(h Hand) (badness Score))
	score func(bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score)
}

type ScoringRule struct {
	name string
	score func(h Hand) Score
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

var Convention = []BiddingRule{ PassOfForcing,
	StrongTwoClubs, Opening, Preempt, PassOpening,
	CheapResponse, CheapNTResponse, TwoOverOne, CheapCompetitionResponse,
	CheapNTRebid, CheapRebid, RebidSuit, Splinter,
	MajorSupport, MajorInvitation,
	OneNT, Stayman, StaymanResponse, StaymanTwo, StaymanTwoResponse, TwoNT, Gambling3NT,
	OneLevelOvercall, PreemptOvercall,
	PassOvercall, PassHigherOvercall, Forced, Natural, LimitPass }

func makeScoringRule(bidder Seat, bid string, e *Ensemble) *ScoringRule {
	if sc,ok := e.scorers[bid]; ok {
		return sc
	}
	for _,c := range Convention {
		ms := c.match.FindStringSubmatch(bid)
		if ms != nil {
			if c.mkscore == nil {
				score := func(h Hand) Score {
					return c.score(bidder, h, ms, e)
				}
				e.scorers[bid] = &ScoringRule{c.name, score}
				return e.scorers[bid]
			} else {
				if sc := c.mkscore(bidder, ms, e); sc != nil {
					e.scorers[bid] = &ScoringRule{c.name, sc}
					return e.scorers[bid]
				}
			}
		}
	}
	return nil
}

func makeUnforcedScoringRule(bidder Seat, bid string, e *Ensemble) *ScoringRule {
	if sc,ok := e.unforced[bid]; ok {
		return sc
	}
	for _,c := range Convention {
		ms := c.match.FindStringSubmatch(bid)
		if ms != nil && c.name != "Forced" {
			if c.mkscore == nil {
				score := func(h Hand) Score {
					return c.score(bidder, h, ms, e)
				}
				e.unforced[bid] = &ScoringRule{c.name, score}
				return e.unforced[bid]
			} else {
				if sc := c.mkscore(bidder, ms, e); sc != nil {
					e.unforced[bid] = &ScoringRule{c.name, sc}
					return e.unforced[bid]
				}
			}
		}
	}
	return nil
}

func RateBid(h Hand, bid string) (badness Score, convention string) {
	e := GetValidTables(South, bid[0:len(bid)-2], 1000)
	rule := makeScoringRule(Seat((len(bid)/2-1) % 4), bid, e)
	return simpleScore(h, rule)
}

func subBids(dealer Seat, bid string) (seats []Seat, bids []string) {
	seats = make([]Seat, len(bid)/2)
	bids = make([]string, len(bid)/2)
	for i := range seats {
		seats[i] = (dealer + Seat(i)) % 4
		bids[i] = bid[0:2*(i+1)]
	}
	return
}

func simpleScore(h Hand, rule *ScoringRule) (badness Score, convention string) {
	if rule == nil {
		return 0, "no convention matches"
	}
	return rule.score(h), rule.name
}

func TableScore(t Table, bidders []Seat, rules []*ScoringRule) (badness Score, conventions []string) {
	conventions = make([]string, len(rules))
	for i, bidder := range bidders {
		b,c := simpleScore(t[bidder], rules[i])
		conventions[i] = c
		badness += b
	}
	return
}

// The []string output describes the bids made...
func GetValidTables(dealer Seat, bid string, num int) *Ensemble {
	seats, bids := subBids(dealer, bid)
	esold := makeEnsemble(num) // This is the ensemble before each bid
	for i := range esold.tables {
		esold.tables[i] = Shuffle() // Things start out random!
	}
	var conventions []string
	rules := make([]*ScoringRule, 0, len(seats))
	for bidnum := range seats {
		rules = rules[0:bidnum+1]
		rules[bidnum] = makeScoringRule((dealer + Seat(bidnum))%4, bids[bidnum], esold)
		es := makeEnsemble(num) // This is the ensemble after this bid
		es.old = esold
		//fmt.Println("I am working on bid of", bids[bidnum])
		for i,t := range esold.tables {
			// Initialize ensemble based on previous bidding
			oldbadness,cs := TableScore(t, seats[0:bidnum+1], rules)
			conventions = cs
			badness := oldbadness
			numswaps := 52*10
			beta := Score(0.01)
			maxbeta := Score(1)
			betainc := Score(math.Pow(float64(maxbeta/beta), 1/float64(numswaps)))
			for i:=0; i<numswaps && (badness > 0 || i < 52); i++ {
				// Try moving a couple of cards at a time...
				t2 := t.ShuffleCard(rand.Intn(52)).ShuffleCard(rand.Intn(52))
				b2,cs := TableScore(t2, seats[0:bidnum+1], rules[0:bidnum+1])
				if b2 <= badness || rand.Float64() < math.Exp(float64(-beta*(b2 - badness))) {
					t = t2
					badness = b2
					conventions = cs
				}
				beta *= betainc // Here's our annealing schedule...
			}
			if badness > 0 {
				fmt.Printf("Badness %4g -> %4g\n", oldbadness, badness)
			}
			es.tables[i] = t
			es.Conventions = conventions
		}
		esold = es
	}
	return esold // return final ensemble (which is now old)
}
