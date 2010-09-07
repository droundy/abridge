package bridge

import (
	"regexp"
)

var Splinter = BiddingRule{
	"Splinter",
	regexp.MustCompile("^( P)*1([HS]) P(3S|4[CDH])$"),
	func (bidder Seat, h Hand, ms []string, e Ensemble) (badness Score, nothandled bool) {
		opensuit := stringToSuitNumber(ms[2])
		splintersuit := uint(Spades)
		switch ms[3] {
		case "4C": splintersuit = Clubs
		case "4D": splintersuit = Diamonds
		case "4H": splintersuit = Hearts
		}
		if opensuit == splintersuit {
			return 0, true // Not a splinter!
		}
		openlen := byte(h >> (4+opensuit*8)) & 15
		splinterlen := byte(h >> (4+splintersuit*8)) & 15
		if splinterlen > 1 {
			badness += Score(splinterlen - 1)*SuitLengthProblem
		}
		spls := byte(h >> splintersuit*8)
		hcp_inside := HCP[Suit(spls)]
		hcp_outside := h.HCP() - hcp_inside
		// Splinter indicates 10-12 hcp outside the singleton
		if hcp_outside < 10 {
			badness += Score(10 - hcp_outside)*PointValueProblem
		} else if hcp_outside > 12 {
			badness += Score(hcp_outside - 12)*PointValueProblem
		}
		if hcp_inside > 2 {
			// To bid a splinter, we'd better not have more than a queen in
			// the singleton.
			badness += Score(hcp_inside - 2)*PointValueProblem
		}
		if openlen < 4 {
			// A splinter bid promises four-card support.
			badness += Score(4-openlen)*SuitLengthProblem
		}
		return
	},
}

var MajorInvitation = BiddingRule{
	"Major support",
	regexp.MustCompile("^( P)*1([HS]) P3([HS])$"),
	func (bidder Seat, h Hand, ms []string, e Ensemble) (badness Score, nothandled bool) {
		opensuit := stringToSuitNumber(ms[2])
		mysuit := stringToSuitNumber(ms[3])
		if mysuit != opensuit {
			return 0, true // This isn't support
		}
		pts := h.PointCount()
		if pts < 10 {
			badness += Score(10-pts)*PointValueProblem
		} else if pts > 11 {
			badness += Score(pts-11)*PointValueProblem
		}
		mysuitlen := byte(h >> (4 + mysuit*8)) & 15
		if mysuitlen < 3 {
			badness += Score(3-mysuitlen)*SuitLengthProblem
		}
		return
	},
}

var MajorSupport = BiddingRule{
	"Major support",
	regexp.MustCompile("^( P)*1([HS]) P2([HS])$"),
	func (bidder Seat, h Hand, ms []string, e Ensemble) (badness Score, nothandled bool) {
		opensuit := stringToSuitNumber(ms[2])
		mysuit := stringToSuitNumber(ms[3])
		if mysuit != opensuit {
			return 0, true // This isn't support
		}
		pts := h.PointCount()
		if pts < 6 {
			badness += Score(6-pts)*PointValueProblem
		} else if pts > 9 {
			badness += Score(pts-9)*PointValueProblem
		}
		mysuitlen := byte(h >> (4 + mysuit*8)) & 15
		if mysuitlen < 3 {
			badness += Score(3-mysuitlen)*SuitLengthProblem
		}
		return
	},
}

var TwoOverOne = BiddingRule{
	"Two over one",
	regexp.MustCompile("^( P)*1([DHS]) P2([CDH])$"),
	func (bidder Seat, h Hand, ms []string, e Ensemble) (badness Score, nothandled bool) {
		pts := h.PointCount()
		opensuit := stringToSuitNumber(ms[2])
		mysuit := stringToSuitNumber(ms[3])
		if mysuit == opensuit {
			return 0, true // This isn't a two-over-one bid
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
		return
	},
}

var CheapResponse = BiddingRule{
	"Cheap response to one",
	regexp.MustCompile("^( P)*1([CDHS]) P1([DHSN])$"),
	func (bidder Seat, h Hand, ms []string, e Ensemble) (badness Score, nothandled bool) {
		pts := h.PointCount()
		if pts < 6 {
			badness += Score(6-pts)*PointValueProblem
		}
		opensuit := stringToSuitNumber(ms[2])
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
			if (spadelen > 3 && opensuit < Spades) {
				badness += Score(spadelen-3)*SuitLengthProblem
			}
			if (heartlen > 3 && opensuit < Hearts) {
				badness += Score(heartlen-3)*SuitLengthProblem
			}
			if pts > 9 {
				badness += Score(pts - 9)*PointValueProblem
			}
			return // exit early, so we can assume mysuit is a valid suit
		}
		// Here we assume ms[3] is a real suit.
		mysuit := stringToSuitNumber(ms[3])
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
		return
	},
}
