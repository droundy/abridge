package bridge

import (
	"regexp"
)

var OneNT = BiddingRule{
	"1NT opening",
	regexp.MustCompile("^( P)*1N$"),
	func (h Hand, ms []string) Score {
		hcp := h.HCP()
		dist := h.DistPoints()
		badness := Score(0)
		if hcp > 17 {
			badness += Score(hcp-17)*PointValueProblem
		} else if hcp < 15 {
			badness += Score(15-hcp)*PointValueProblem
		}
		if dist > 1 {
			badness += Score(dist-1)*PointValueProblem
		}
		ls := byte(h>>28) & 15
		lh := byte(h>>20) & 15
		if ls > 4 {
			badness += Score(ls-4)*SuitLengthProblem
		}
		if lh > 4 {
			badness += Score(lh-4)*SuitLengthProblem
		}
		return badness
	},
}

var TwoNT = BiddingRule{
	"2NT opening",
	regexp.MustCompile("^( P)*2N$"),
	func (h Hand, ms []string) Score {
		hcp := h.HCP()
		dist := h.DistPoints()
		badness := Score(0)
		if hcp > 22 {
			badness += Score(hcp-22)*PointValueProblem
		} else if hcp < 20 {
			badness += Score(20-hcp)*PointValueProblem
		}
		if dist > 1 {
			badness += Score(dist-1)*PointValueProblem
		}
		ls := byte(h>>28) & 15
		lh := byte(h>>20) & 15
		if ls > 4 {
			badness += Score(ls-4)*SuitLengthProblem
		}
		if lh > 4 {
			badness += Score(lh-4)*SuitLengthProblem
		}
		return badness
	},
}

var Gambling3NT = BiddingRule{
	"Gambling 3NT",
	regexp.MustCompile("^( P)*3N$"),
	func (h Hand, ms []string) Score {
		hcp := h.HCP()
		d := Suit(h >> 8)
		c := Suit(h)
		badness := Score(0)
		if hcp > 17 {
			// Too strong even for strong gambling!
			badness += Score(hcp-17)*PointValueProblem
		}
		if d > c {
			if (d>>4) < 7 {
				badness += Score(7 - (d>>4))*SuitLengthProblem
			}
			if (d & 15 < 14) {
				badness += Score(14 - (d&15))*BigFudge
			}
		} else {
			if (c>>4) < 7 {
				badness += Score(7 - (c>>4))*SuitLengthProblem
			}
			if (c & 15 < 14) {
				badness += Score(14 - (c&15))*BigFudge
			}
		}
		return badness
	},
}
