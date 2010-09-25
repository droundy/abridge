package bridge

import (
	"fmt"
	"regexp"
)

func stringToSuitNumber(s string) uint {
	switch s {
	case "S","s": return Spades
	case "H","h": return Hearts
	case "D","d": return Diamonds
	case "C","c": return Clubs
	case "N","n": return NoTrump
	}
	panic(fmt.Sprint("Bad string in stringToSuitNumber: ", s))
	return 0
}

var Opening = BiddingRule{
	"Opening",
	regexp.MustCompile("^( P)*1([CDHS])$"),
	nil,
	func (bidder Seat, h Hand, ms []string, c ConventionCard, e *Ensemble) (badness Score, explanation string) {
		pts := h.PointCount()
		if pts < 13 {
			explanation = "Opening with 12 points is a fudge!\n"
			badness += Fudge
		}
		if pts < 12 {
			explanation = "Opening with under 12 points is just wrong!\n"
			badness += Score(12-pts)*PointValueProblem
		}
		ls := byte(h >> 28)
		lh := byte(h >> 20) & 15
		ld := byte(h >> 12) & 15
		lc := byte(h >> 4) & 15
		switch stringToSuitNumber(ms[2]) {
		case Spades:
			if ls < 5 {
				badness += Score(5-ls)*SuitLengthProblem
			}
			if ls < lh {
				badness += Score(lh-ls)*SuitLengthProblem
			}
		case Hearts:
			if lh < 5 {
				badness += Score(5-lh)*SuitLengthProblem
			}
			if lh < ls {
				badness += Score(ls-lh)*SuitLengthProblem
			}
		case Diamonds:
			if lh > 4 {
				badness += Score(lh-4)*SuitLengthProblem
			}
			if ls > 4 {
				badness += Score(ls-4)*SuitLengthProblem
			}
			if lc > ld {
				badness += Score(lc-ld)*SuitLengthProblem
			}
		case Clubs:
			if lh > 4 {
				badness += Score(lh-4)*SuitLengthProblem
			}
			if ls > 4 {
				badness += Score(ls-4)*SuitLengthProblem
			}
			if ld > lc {
				badness += Score(ld-lc)*SuitLengthProblem
			}
		}
		return
	},
}

var PreemptOvercall = BiddingRule{
	"Preemptive overcall of suit bid",
	// preempts over 1NT should be stronger and/or longer...
	regexp.MustCompile("^(.. P)?([12])([CDHSN])( P P)?([23])([CDHS])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (func(Hand) (Score,string)) {
		if cc.Radio["JumpOvercall"] == "Weak" {
			theirsuit := stringToSuitNumber(ms[3])
			mysuit := stringToSuitNumber(ms[6])
			if mysuit == theirsuit {
				return nil // It's a cue bid!
			}
			mynum := ms[5][0]-'0'
			theirnum := ms[2][0]-'0'
			if mynum == theirnum {
				return nil // It's not a jump
			}
			if mynum == theirnum+1 && mysuit < theirsuit {
				return nil // It's not a jump
			}
			// Reorder matches for an ordinary preempt
			ms[2] = ms[5]
			ms[3] = ms[6]
			return Preempt.mkscore(bidder, ms, cc, e)
		}
		return nil
	}, nil, // preemptive overcalls at 3 level are like ordinary preempts.
}

var Preempt = BiddingRule{
	"Preempt",
	regexp.MustCompile("^( P)*([234])([CDHS])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (func(Hand) (Score,string)) {
		if ms[2] == "2" && ms[3] == "C" {
			return nil // it's not a weak two bid
		}
		verylightpreempts := cc.Options["VeryLightPreempts"]
		contract := int(ms[2][0] - '0')
		goal := byte(contract + 4)
		mysuit := stringToSuitNumber(ms[3])
		if goal > 6 {
			switch cc.Radio["WeakThree"] {
			case "Sound":
				return func(h Hand) (badness Score, explanation string) {
					pts := h.PointCount()
					hcp := h.HCP()
					if pts > 12 {
						badness += Score(pts-12)*PointValueProblem
					} else if hcp < 5 {
						badness += Score(5 - hcp)*PointValueProblem
					}
					cardsinsuit := Suit(h >> (8*mysuit))
					numinsuit := byte(cardsinsuit >> 4) & 15
					lev := SafeContractInThisSuit(bidder, h, mysuit, e)
					if lev < contract - 3 {
						badness += Score(contract - 3 - lev)*SuitLengthProblem
					}
					ptsinsuit := HCP[cardsinsuit]
					if ptsinsuit < 4 {
						badness += Score(4 - ptsinsuit)*PointValueProblem
					}
					if numinsuit < goal {
						badness += Score(goal-numinsuit)*SuitLengthProblem
					}
					return
				}
			case "Light":
				return func(h Hand) (badness Score, explanation string) {
					pts := h.PointCount()
					hcp := h.HCP()
					if pts > 12 {
						badness += Score(pts-12)*PointValueProblem
					} else if hcp < 5 {
						badness += Score(5 - hcp)*PointValueProblem
					}
					cardsinsuit := Suit(h >> (8*mysuit))
					numinsuit := byte(cardsinsuit >> 4) & 15
					ptsinsuit := HCP[cardsinsuit]
					if ptsinsuit < 4 {
						badness += Score(4 - ptsinsuit)*PointValueProblem
					}
					if numinsuit < goal {
						badness += Score(goal-numinsuit)*SuitLengthProblem
					}
					return
				}
			case "VeryLight":
				// I'll treat VeryLight same as weak twos.
			}
		}
		return func(h Hand) (badness Score, explanation string) {
			pts := h.PointCount()
			hcp := h.HCP()
			if pts > 12 {
				badness += Score(pts-12)*PointValueProblem
			} else if hcp < 5 {
				badness += Score(5 - hcp)*PointValueProblem
			}
			cardsinsuit := Suit(h >> (8*mysuit))
			numinsuit := byte(cardsinsuit >> 4) & 15
			if !verylightpreempts {
				ptsinsuit := HCP[cardsinsuit]
				if ptsinsuit < 3 {
					badness += Score(3 - ptsinsuit)*PointValueProblem
				}
			}
			if numinsuit < goal {
				badness += Score(goal-numinsuit)*SuitLengthProblem
			} else if goal == 6 {
				badness += Score(numinsuit-goal)*Fudge
			}
			return
		}
	}, nil,
}

var StrongTwoClubs = BiddingRule{
	"Strong two clubs (forcing)",
	regexp.MustCompile("^( P)*2C$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (func(Hand) (Score, string)) {
		if !cc.Options["StrongTwoClubs"] {
			return nil
		}
		return func(h Hand) (badness Score, explanation string) {
			pts := h.PointCount()
			if pts < 23 {
				badness += Score(23-pts)*PointValueProblem
			}
			return
		}
	}, nil,
}

var PassOpening = BiddingRule{
	"Pass opening",
	regexp.MustCompile("^( P)* P$"),
	nil,
	func (bidder Seat, h Hand, ms []string, c ConventionCard, e *Ensemble) (badness Score, explanation string) {
		pts := h.PointCount()
		hcp := h.HCP()
		if pts > 12 {
			badness += Score(pts-12)*PointValueProblem
		}
		if (byte(h >> 4) & 15) > 6 && hcp >= 5 { // should bid weak
			badness += Score((byte(h >> 4) & 15) - 6)*BigFudge
		}
		for sv:=uint(Diamonds); sv <= Spades; sv++ {
			l := byte(h >> (4 + sv*8)) & 15
			if l > 5 && hcp >= 5 { // should bid weak
				badness += Score(l - 5)*BigFudge
			}
		}
		return
	},
}
