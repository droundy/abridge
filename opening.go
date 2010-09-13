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
	}
	panic(fmt.Sprint("Bad string in stringToSuitNumber: ", s))
	return 0
}

var Opening = BiddingRule{
	"Opening",
	regexp.MustCompile("^( P)*1([CDHS])$"),
	nil,
	func (bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score) {
		pts := h.PointCount()
		if pts < 13 {
			badness += Fudge
		}
		if pts < 12 {
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
	regexp.MustCompile("^( P)*1[CDHS](3)([CDHS])$"),
	nil,
	Preempt.score, // preemptive overcalls at 3 level are like ordinary preempts.
}

var Preempt = BiddingRule{
	"Preempt",
	regexp.MustCompile("^( P)*([23])([CDHS])$"),
	func (bidder Seat, ms []string, e *Ensemble) (func(Hand) (Score)) {
		if ms[2] == "2" && ms[3] == "C" {
			return nil // it's not a weak two bid
		}
		goal := byte(6)
		if ms[2] == "3" {
			goal = 7
		}
		return func(h Hand) (badness Score) {
			pts := h.PointCount()
			hcp := h.HCP()
			if pts > 12 {
				badness += Score(pts-12)*PointValueProblem
			} else if hcp < 5 {
				badness += Score(5 - hcp)*PointValueProblem
			}
			numinsuit := goal
			switch stringToSuitNumber(ms[3]) {
			case Clubs: numinsuit = byte(h >> 4) & 15
			case Spades: numinsuit = byte(h >> 28)
			case Hearts: numinsuit = byte(h >> 20) & 15
			case Diamonds: numinsuit = byte(h >> 12) & 15
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
	"Strong two clubs",
	regexp.MustCompile("^( P)*2C$"),
	func (bidder Seat, ms []string, e *Ensemble) (func(Hand) (Score)) {
		return func(h Hand) (badness Score) {
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
	func (bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score) {
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
