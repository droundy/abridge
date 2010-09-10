package bridge

import (
	"regexp"
)

var RebidSuit = BiddingRule{
	"Rebid in my suit after cheap unlimited response",
	regexp.MustCompile("^( P)*1([CDHS]) P1([DHS]) P2([CDHS])$"), nil,
	func (bidder Seat, h Hand, ms []string, e Ensemble) (badness Score, nothandled bool) {
		if ms[2] != ms[4] {
			return 0, true // This isn't a rebid of my suit
		}
		mysuit := stringToSuitNumber(ms[2])
		theirsuit := stringToSuitNumber(ms[3])
		mysuitlen := byte(h >> (4+mysuit*8)) & 15
		theirsuitlen := byte(h >> (4+theirsuit*8)) & 15

		pts := h.PointCount()
		if pts > 15 {
			badness += Score(pts-15)*PointValueProblem
		}
		if mysuitlen < 6 {
			badness += Score(6 - mysuitlen)*SuitLengthProblem
		}
		if theirsuitlen > 3 && theirsuit >= Hearts {
			// If we have support for their major, say so!
			badness += Score(theirsuitlen-3)*SuitLengthProblem
		}
		return
	},
}

var CheapRebid = BiddingRule{
	"Cheap rebid",
	regexp.MustCompile("^( P)*1([CDH]) P1([DHS]) P1([HSN])$"), nil,
	func (bidder Seat, h Hand, ms []string, e Ensemble) (badness Score, nothandled bool) {
		pts := h.PointCount()
		opensuit := stringToSuitNumber(ms[2])
		theirsuit := stringToSuitNumber(ms[3])
		theirsuitlen := byte(h >> (4+theirsuit*8)) & 15
		opensuitlen := byte(h >> (4+opensuit*8)) & 15
		heartlen := byte(h >> 20) & 15
		spadelen := byte(h >> 28) & 15
		if theirsuit >= Hearts && theirsuitlen > 3 {
			// We have a major fit!
			badness += Score(theirsuitlen-3)*SuitLengthProblem
		}
		if ms[4] == "N" {
			if theirsuit < Spades && spadelen > 3 {
				// Should mention spades.
				badness += Score(spadelen-3)*SuitLengthProblem
			}
			if theirsuit < Hearts && heartlen > 3 {
				// Should mention hearts.
				badness += Score(heartlen-3)*SuitLengthProblem
			}
			if pts > 15 {
				badness += Score(pts - 15)*PointValueProblem
			}
			if opensuitlen > 5 {
				badness += Score(opensuitlen-5)*SuitLengthProblem
			}
			return // exit early, so we can assume mysuit is a valid suit
		}
		if theirsuit < Hearts && ms[4] == "S" && heartlen > 3 {
			// Should mention hearts.
			badness += Score(heartlen-3)*SuitLengthProblem
		}
		// Here we assume ms[4] is a real suit.
		mysuit := stringToSuitNumber(ms[4])
		mysuitlen := byte(h >> (4 + mysuit*8)) & 15
		if mysuitlen < 4 {
			badness += Score(4-mysuitlen)*SuitLengthProblem
		}
		return
	},
}
