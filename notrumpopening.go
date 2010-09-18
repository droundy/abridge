package bridge

import (
	"regexp"
)

var Jacobi = BiddingRule{
	"Jacobi transfer (forcing)",
	regexp.MustCompile("^( P)*1N P2([DH])$"), nil,
	func (bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score) {
		mysuit := stringToSuitNumber(ms[2])+1
		mysuitlen := byte(h>>(4+mysuit*8)) & 15
		if mysuitlen < 5 {
			badness += Score(5 - mysuitlen)*SuitLengthProblem
		}
		return
	},
}

var JacobiResponse = BiddingRule{
	"Jacobi response",
	regexp.MustCompile("^( P)*1N P2([DH]) P2([HS])$"),
	func (bidder Seat, ms []string, e *Ensemble) (score func(h Hand) Score) {
		if ms[2] == "D" && ms[3] == "S" {
			return nil
		}
		return func(h Hand) Score {
			return 0
		}
	}, nil,
}

var JacobiRejection = BiddingRule{
	"Jacobi rejection (bad bid)",
	regexp.MustCompile("^( P)*1N P2([DH]) P..$"), nil,
	func (bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score) {
		mysuit := stringToSuitNumber(ms[2])+1
		mysuitlen := byte(h>>(4+mysuit*8)) & 15
		if mysuitlen > 2 {
			badness += Score(mysuitlen-2)*SuitLengthProblem
		}
		return SuitLengthProblem + badness
	},
}

var Stayman = BiddingRule{
	"Stayman (forcing)",
	regexp.MustCompile("^( P)*1N P2C$"), nil,
	func (bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score) {
		return // Says nothing...
	},
}

var StaymanTwo = BiddingRule{
	"Stayman (forcing)",
	regexp.MustCompile("^( P)*2N P3C$"), nil,
	Stayman.score,
}

var StaymanResponse = BiddingRule{
	"Stayman",
	regexp.MustCompile("^( P)*1N P2C P2([DHS])$"), nil,
	func (bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score) {
		lh := byte(h>>20) & 15
		ls := byte(h>>28) & 15
		switch ms[2] {
		case "D":
			if lh > 3 {
				badness += Score(lh - 3)*SuitLengthProblem
			}
			if ls > 3 {
				badness += Score(ls - 3)*SuitLengthProblem
			}
		case "H":
			if lh < 4 {
				badness += Score(4 - lh)*SuitLengthProblem
			}
		case "S":
			if lh > 3 {
				badness += Score(lh - 3)*SuitLengthProblem
			}
			if ls < 4 {
				badness += Score(4 - ls)*SuitLengthProblem
			}
		}
		return
	},
}

var StaymanTwoResponse = BiddingRule{
	"Stayman",
	regexp.MustCompile("^( P)*2N P3C P3([DHS])$"), nil,
	StaymanResponse.score,
}

var OneNT = BiddingRule{
	"1NT opening",
	regexp.MustCompile("^( P)*1N$"), nil,
	func (bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score) {
		hcp := h.HCP()
		dist := h.DistPoints()
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
		return
	},
}

var TwoNT = BiddingRule{
	"2NT opening",
	regexp.MustCompile("^( P)*2N$"), nil,
	func (bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score) {
		hcp := h.HCP()
		dist := h.DistPoints()
		if hcp > 21 {
			badness += Score(hcp-21)*PointValueProblem
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
		return
	},
}

var Gambling3NT = BiddingRule{
	"Gambling 3NT",
	regexp.MustCompile("^( P)*3N$"), nil,
	func (bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score) {
		hcp := h.HCP()
		d := Suit(h >> 8)
		c := Suit(h)
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
		return
	},
}
