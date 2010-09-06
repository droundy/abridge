package bridge

import (
	"regexp"
)

var PassOvercall = BiddingRule{
	"One-level overcall",
	regexp.MustCompile("^( P)*1[CDH]( P..)? P$"),
	func (h Hand, ms []string) Score {
		pts := h.PointCount()
		badness := Score(0)
		if pts > 12 {
			badness += Score(pts-12)*PointValueProblem
		}
		return badness
	},
}

var OneLevelOvercall = BiddingRule{
	"One-level overcall",
	regexp.MustCompile("^( P)*1[CDH]( P..)?1([DHS])$"),
	func (h Hand, ms []string) Score {
		pts := h.PointCount()
		badness := Score(0)
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
		return badness
	},
}
