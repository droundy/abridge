package bridge

import (
	"regexp"
)

var MajorInvitation = BiddingRule{
	"Major support",
	regexp.MustCompile("^( P)*1([HS]) P3([HS])$"),
	func (h Hand, ms []string) Score {
		opensuit := stringToSuit(ms[2])
		mysuit := stringToSuit(ms[3])
		if mysuit != opensuit {
			return 0 // This isn't support
		}
		pts := h.PointCount()
		badness := Score(0)
		if pts < 10 {
			badness += Score(10-pts)*PointValueProblem
		} else if pts > 11 {
			badness += Score(pts-11)*PointValueProblem
		}
		mysuitlen := byte(h >> (4 + mysuit*8)) & 15
		if mysuitlen < 3 {
			badness += Score(3-mysuitlen)*SuitLengthProblem
		}
		return badness
	},
}

var MajorSupport = BiddingRule{
	"Major support",
	regexp.MustCompile("^( P)*1([HS]) P2([HS])$"),
	func (h Hand, ms []string) Score {
		opensuit := stringToSuit(ms[2])
		mysuit := stringToSuit(ms[3])
		if mysuit != opensuit {
			return 0 // This isn't support
		}
		pts := h.PointCount()
		badness := Score(0)
		if pts < 6 {
			badness += Score(6-pts)*PointValueProblem
		} else if pts > 9 {
			badness += Score(pts-9)*PointValueProblem
		}
		mysuitlen := byte(h >> (4 + mysuit*8)) & 15
		if mysuitlen < 3 {
			badness += Score(3-mysuitlen)*SuitLengthProblem
		}
		return badness
	},
}

var TwoOverOne = BiddingRule{
	"Two over one",
	regexp.MustCompile("^( P)*1([DHS]) P2([CDH])$"),
	func (h Hand, ms []string) Score {
		pts := h.PointCount()
		badness := Score(0)
		opensuit := stringToSuit(ms[2])
		mysuit := stringToSuit(ms[3])
		if mysuit == opensuit {
			return 0 // This isn't a two-over-one bid
		}
		if pts < 10 {
			badness += Score(10-pts)*PointValueProblem
		}
		mysuitlen := byte(h >> (4 + mysuit*8)) & 15
		if mysuitlen < 4 {
			badness += Score(4-mysuitlen)*SuitLengthProblem
		}
		heartlen := byte(h >> 20) & 15
		spadelen := byte(h >> 28) & 15
		if opensuit < Spades && spadelen > 3 {
			badness += Score(spadelen-3)*SuitLengthProblem
		}
		if opensuit < Hearts && heartlen > 3 && mysuitlen != Hearts {
			badness += Score(heartlen-3)*SuitLengthProblem
		}
		if opensuit == Hearts && heartlen > 2 && pts < 15 {
			b1 := Score(heartlen-2)*SuitLengthProblem
			b2 := Score(15-pts)*PointValueProblem
			badness += b1.min(b2)
		}
		if opensuit == Spades && spadelen > 2 && pts < 15 {
			b1 := Score(spadelen-2)*SuitLengthProblem
			b2 := Score(15-pts)*PointValueProblem
			badness += b1.min(b2)
		}
		return badness
	},
}

var CheapResponse = BiddingRule{
	"Cheap response to one",
	regexp.MustCompile("^( P)*1([CDHS]) P1([DHSN])$"),
	func (h Hand, ms []string) Score {
		pts := h.PointCount()
		badness := Score(0)
		if pts < 6 {
			badness += Score(6-pts)*PointValueProblem
		}
		opensuit := stringToSuit(ms[2])
		heartlen := byte(h >> 20) & 15
		spadelen := byte(h >> 28) & 15
		if opensuit == Spades && spadelen > 2 {
			// We missed an opening bid!
			badness += Score(spadelen-2)*SuitLengthProblem
		}
		if opensuit == Hearts && heartlen > 2 && !(ms[3] == "S" && pts > 9) {
			// We can only bid 1S if we really have good reason to force the
			// bid... i.e. a strongish hand.
			badness += Score(heartlen-2)*SuitLengthProblem
		}
		if ms[3] == "N" {
			if pts > 9 {
				badness += Score(pts - 9)*PointValueProblem
			}
			return badness // exit early, so we can assume mysuit is a valid suit
		}
		// Here we assume ms[3] is a real suit.
		mysuit := stringToSuit(ms[3])
		mysuitlen := byte(h >> (4 + mysuit*8)) & 15
		switch mysuit {
		case Hearts:
			if mysuitlen < 4 {
				badness += Score(4-mysuitlen)*SuitLengthProblem
			}
		case Spades:
			if mysuitlen < 4 {
				badness += Score(4-mysuitlen)*SuitLengthProblem
			}
			if heartlen > 3 && opensuit != Hearts && spadelen < 6 {
				// Skipping hearts denies 4 hearts, unless you've got 6 spades
				b1 := Score(heartlen-3)*SuitLengthProblem
				b2 := Score(7-spadelen)*SuitLengthProblem
				badness += b1.min(b2)
			}
		case Diamonds:
			if mysuitlen < 4 {
				badness += Score(4-mysuitlen)*SuitLengthProblem
			}
			if heartlen > 3 {
				badness += Score(heartlen-3)*SuitLengthProblem
			}
			if spadelen > 3 {
				badness += Score(spadelen-3)*SuitLengthProblem
			}
		}
		return badness
	},
}
