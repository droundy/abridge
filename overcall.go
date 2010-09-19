package bridge

import (
	"regexp"
)

var PassOvercall = BiddingRule{
	"Pass an opportunity to overcall a one-suit bid",
	regexp.MustCompile("^( P)*1[CDH]( P..)? P$"), nil,
	func (bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score, explanation string) {
		pts := h.PointCount()
		if pts > 12 {
			badness += Score(pts-12)*PointValueProblem
		}
		return
	},
}

var PassHigherOvercall = BiddingRule{
	"Pass a higher overcall",
	regexp.MustCompile("^( P)*(1N|2N|[23][CDHS])( P..)? P$"), nil,
	func (bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score, explanation string) {
		pts := h.PointCount()
		longest_suit := byte(0)
		for sv := uint(0); sv<4; sv++ {
			l := byte((h >> (4+8*sv)) & 15)
			if l > longest_suit {
				longest_suit = l
			}
		}
		if pts > 14 {
			badness += Score(pts-14)*PointValueProblem
		} else if longest_suit > 6 && pts > 7 {
			badness += Score(longest_suit - 6)*SuitLengthProblem
			badness += Score(pts - 8)*PointValueProblem
		}
		return
	},
}

var OneLevelOvercall = BiddingRule{
	"One-level overcall",
	regexp.MustCompile("^( P)*1[CDH]( P..)?1([DHS])$"), nil,
	func (bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score, explanation string) {
		pts := h.PointCount()
		if pts < 13 {
			badness += Fudge
		}
		if pts < 12 {
			badness += Score(12-pts)*PointValueProblem
		}
		mysuit := stringToSuitNumber(ms[3])
		mysuitlen := byte(h >> (4 + mysuit*8)) & 15
		if mysuitlen < 5 {
			badness += Score(5 - mysuitlen)*SuitLengthProblem
		}
		return
	},
}
