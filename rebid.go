package bridge

import (
	"regexp"
)

var CheapRebid = BiddingRule{
	"Cheap response to one",
	regexp.MustCompile("^( P)*1[CD] P1([DH]) P1([HSN])$"),
	func (bidder Seat, h Hand, ms []string, e Ensemble) (badness Score, nothandled bool) {
		pts := h.PointCount()
		if pts < 6 {
			badness += Score(6-pts)*PointValueProblem
		}
		theirsuit := stringToSuitNumber(ms[2])
		heartlen := byte(h >> 20) & 15
		spadelen := byte(h >> 28) & 15
		if theirsuit == Hearts && heartlen > 3 {
			// We have a fit!
			badness += Score(heartlen-4)*SuitLengthProblem
		}
		if ms[3] == "N" {
			if theirsuit == Hearts && spadelen > 3 {
				// Should mention spades.
				badness += Score(spadelen-4)*SuitLengthProblem
			}
			if pts > 15 {
				badness += Score(pts - 15)*PointValueProblem
			}
			return // exit early, so we can assume mysuit is a valid suit
		}
		// Here we assume ms[3] is a real suit.
		mysuit := stringToSuitNumber(ms[3])
		mysuitlen := byte(h >> (4 + mysuit*8)) & 15
		if mysuitlen < 4 {
			badness += Score(4-mysuitlen)*SuitLengthProblem
		}
		return
	},
}
