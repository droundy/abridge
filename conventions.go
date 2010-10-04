package bridge

import (
	"fmt"
	"math"
	"rand"
	"regexp"
)

type Color byte

const (
	explain = false
)

type Score float64
const (
	SuitLengthProblem Score = 100
	PointValueProblem Score = 30
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
	mkscore func (bidder Seat, ms []string, c ConventionCard, e *Ensemble) (score func(h Hand) (badness Score, explanation string))
	score func(bidder Seat, h Hand, ms []string, c ConventionCard, e *Ensemble) (badness Score, explanation string)
}

type ScoringRule struct {
	name string
	score func(h Hand) (badness Score, explanation string)
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
	StrongTwoResponse,
	CheapNTRebid, CheapRebid, TwoLevelRebidSuit, RebidSuit, Splinter,
	MichaelsCuebid, MichaelsCuebidMinorQuery, MichaelsCuebidMinorQueryResponse,
	Blackwood, BlackwoodResponse, Gerber, GerberResponse,
	MajorSupport, MajorInvitation,
	Jacobi, JacobiResponse, JacobiSuperAccept, JacobiRejection,
	OneNT, Stayman, StaymanResponse, StaymanTwo, StaymanTwoResponse, TwoNT,
	Gambling3NT, Gambling3NTquery, Gambling3NTforcingquery, Gambling3NTresponse, Gambling3NTforcedresponse,
	OneLevelOvercall, PreemptOvercall,
	// The following need to follow PreemptOvercall:
	TwoLevelOvercall, ThreeLevelOvercall,
	DirectOneNTOvercall, BalancingOneNTOvercall,
	PassOvercall, PassHigherOvercall,
	TakeOutDouble, NewSuitForcing, Forced, Natural, LimitPass }

func makeScoringRule(bidder Seat, bid string, cc ConventionCard, e *Ensemble) *ScoringRule {
	if sc,ok := e.scorers[bid]; ok {
		return sc
	}
	for _,c := range Convention {
		ms := c.match.FindStringSubmatch(bid)
		if ms != nil {
			if c.mkscore == nil {
				score := func(h Hand) (Score,string) {
					return c.score(bidder, h, ms, cc, e)
				}
				e.scorers[bid] = &ScoringRule{c.name, score}
				return e.scorers[bid]
			} else {
				if sc := c.mkscore(bidder, ms, cc, e); sc != nil {
					e.scorers[bid] = &ScoringRule{c.name, sc}
					return e.scorers[bid]
				}
			}
		}
	}
	return nil
}

func makeUnforcedScoringRule(bidder Seat, bid string, cc ConventionCard, e *Ensemble) *ScoringRule {
	if sc,ok := e.unforced[bid]; ok {
		return sc
	}
	for _,c := range Convention {
		ms := c.match.FindStringSubmatch(bid)
		if ms != nil && c.name != "Forced" {
			if c.mkscore == nil {
				score := func(h Hand) (Score,string) {
					return c.score(bidder, h, ms, cc, e)
				}
				e.unforced[bid] = &ScoringRule{c.name, score}
				return e.unforced[bid]
			} else {
				if sc := c.mkscore(bidder, ms, cc, e); sc != nil {
					e.unforced[bid] = &ScoringRule{c.name, sc}
					return e.unforced[bid]
				}
			}
		}
	}
	return nil
}

func RateNextBids(bid string, cc [2]ConventionCard) map[string]float64 {
	e := GetValidTables(South, bid, 200, cc)
	bv, bs := LastBid(bid)
	out := make(map[string]float64)
	out[" P"] = 0
	candouble := regexp.MustCompile(".[CDHSN]( P P)?$").MatchString(bid)
	canredouble := regexp.MustCompile(" X( P P)?$").MatchString(bid)
	if candouble {
		out[" X"] = 0
	}
	if canredouble {
		out["XX"] = 0
	}

	for bidlevel:=bv;bidlevel<8;bidlevel++ {
		for sv:=Color(Clubs); sv<=NoTrump; sv++ {
			if bidlevel > bv || (bidlevel == bv && sv > bs) {
				out[fmt.Sprintf("%d%v", bidlevel, SuitLetter[sv])] = 0
			}
		}
	}
	bidder := Seat((len(bid)/2) % 4)
	for b := range out {
		rule := makeScoringRule(bidder, bid + b, cc[bidder&1], e)
		numbad := float64(0)
		numgood := numbad
		for _,t := range e.tables {
			if badness,_,_ := simpleScore(t[bidder], rule); badness == 0 {
				numgood++
			} else {
				numbad++
			}
		}
		out[b] = numgood/(numbad + numgood)
	}
	return out
}

func RateBid(h Hand, bid string, cc [2]ConventionCard) (badness Score, explanation string, convention string) {
	e := GetValidTables(South, bid[0:len(bid)-2], 200, cc)
	bidder := Seat((len(bid)/2-1) % 4)
	rule := makeScoringRule(bidder, bid, cc[bidder&1], e)
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

func simpleScore(h Hand, rule *ScoringRule) (badness Score, explanation string, convention string) {
	if rule == nil {
		return 0, "", "no convention matches"
	}
	b,e := rule.score(h)
	return b, e, rule.name
}

func TableScore(t Table, bidders []Seat, rules []*ScoringRule) (badness Score, conventions []string) {
	conventions = make([]string, len(rules))
	for i, bidder := range bidders {
		b,_,c := simpleScore(t[bidder], rules[i])
		conventions[i] = c
		badness += b
	}
	return
}

// The []string output describes the bids made...
func GetValidTables(dealer Seat, bid string, num int, cc [2]ConventionCard) *Ensemble {
	seats, bids := subBids(dealer, bid)
	esold := makeEnsemble(num) // This is the ensemble before each bid
	for i := range esold.tables {
		esold.tables[i] = Shuffle() // Things start out random!
	}
	var conventions []string
	rules := make([]*ScoringRule, 0, len(seats))
	for bidnum := range seats {
		rules = rules[0:bidnum+1]
		bidder := (dealer + Seat(bidnum))%4
		rules[bidnum] = makeScoringRule(bidder, bids[bidnum], cc[bidder&1], esold)
		ccrot := cc
		if bidder&1 != 0 {
			ccrot = [2]ConventionCard{ cc[1], cc[0] }
		}
		if ecached,ok := lookupEnsembleFromCache(bids[bidnum], ccrot); ok {
			esold = ecached.RotateFromSouth(dealer)
			continue
		}
		es := makeEnsemble(num) // This is the ensemble after this bid
		//fmt.Println("I am working on bid of", bids[bidnum])
		for i := range es.tables {
			t := esold.tables[i % len(esold.tables)]
			// Initialize ensemble based on previous bidding
			oldbadness,cs := TableScore(t, seats[0:bidnum+1], rules)
			conventions = cs
			badness := oldbadness
			numswaps := 52*20
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
				fmt.Printf("Badness %4g -> %4g for bid %s\n", oldbadness, badness, bids[bidnum])
			}
			es.tables[i] = t
			es.Conventions = conventions
		}
		cacheEnsemble(bids[bidnum], ccrot, es.RotateToSouth(dealer))
		esold = es
	}
	return esold // return final ensemble (which is now old)
}
