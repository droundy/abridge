package bridge

import (
	"regexp"
)

var Jacobi = BiddingRule{
	"Jacobi transfer (forcing)",
	regexp.MustCompile("^( P)*1N P2([DH])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (s Score, e string)) {
		if !cc.Options["Jacobi"] {
			return nil
		}
		mysuit := stringToSuitNumber(ms[2])+1
		return func(h Hand) (badness Score, explanation string) {
			mysuitlen := byte(h>>(4+mysuit*8)) & 15
			if mysuitlen < 5 {
				badness += Score(5 - mysuitlen)*SuitLengthProblem
			}
			return
		}
	}, nil,
}

var JacobiResponse = BiddingRule{
	"Jacobi response",
	regexp.MustCompile("^( P)*1N P2([DH]) P2([HS])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		if !cc.Options["Jacobi"] {
			return nil
		}
		if ms[2] == "D" && ms[3] == "S" {
			return nil
		}
		return func(h Hand) (Score,string) {
			return 0, ""
		}
	}, nil,
}

var JacobiSuperAccept = BiddingRule{
	"Jacobi super accept",
	regexp.MustCompile("^( P)*1N P2([DH]) P3([HS])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		if !cc.Options["Jacobi"] {
			return nil
		}
		sv := stringToSuitNumber(ms[3])
		if ms[2] == "D" && sv == Spades {
			return nil
		}
		return func(h Hand) (badness Score, explanation string) {
			ns := byte(h << (4+8*sv)) & 15
			if ns < 4 {
				badness += Score(4-ns)*SuitLengthProblem
			}
			pts := h.PointCount()
			if pts < 17 {
				badness += Score(17-pts)*PointValueProblem
			}
			return
		}
	}, nil,
}

var JacobiRejection = BiddingRule{
	"Jacobi rejection (bad bid)",
	regexp.MustCompile("^( P)*1N P2([DH]) P..$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		if !cc.Options["Jacobi"] {
			return nil
		}
		mysuit := stringToSuitNumber(ms[2])+1
		return func(h Hand) (badness Score, explanation string) {
			mysuitlen := byte(h>>(4+mysuit*8)) & 15
			if mysuitlen > 2 {
				badness += Score(mysuitlen-2)*SuitLengthProblem
			}
			return SuitLengthProblem + badness, ""
		}
	}, nil,
}

var Stayman = BiddingRule{
	"Stayman (forcing)",
	regexp.MustCompile("^( P)*1N P2C$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		if !cc.Options["Stayman"] {
			return nil
		}
		return func(h Hand) (badness Score, explanation string) {
			return // Says nothing...
		}
	}, nil,
}

var StaymanTwo = BiddingRule{
	"Stayman (forcing)",
	regexp.MustCompile("^( P)*2N P3C$"),
	Stayman.mkscore, nil,
}

var StaymanResponse = BiddingRule{
	"Stayman response",
	regexp.MustCompile("^( P)*1N P2C P2([DHS])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		if !cc.Options["Stayman"] {
			return nil
		}
		return func(h Hand) (badness Score, explanation string) {
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
		}
	}, nil,
}

var StaymanTwoResponse = BiddingRule{
	"Stayman",
	regexp.MustCompile("^( P)*2N P3C P3([DHS])$"),
	StaymanResponse.mkscore, nil,
}

var OneNT = BiddingRule{
	"1NT opening",
	regexp.MustCompile("^( P)*1N$"), nil,
	func (bidder Seat, h Hand, ms []string, cc ConventionCard, e *Ensemble) (badness Score, explanation string) {
		hcp := h.HCP()
		dist := h.DistPoints()
		if hcp > cc.Pts["OneNTmax"] {
			badness += Score(hcp-cc.Pts["OneNTmax"])*PointValueProblem
		} else if hcp < cc.Pts["OneNTmin"] {
			badness += Score(cc.Pts["OneNTmin"]-hcp)*PointValueProblem
		}
		if dist > 1 {
			badness += Score(dist-1)*PointValueProblem
		}
		ls := byte(h>>28) & 15
		lh := byte(h>>20) & 15
		if !cc.Options["OneNT5CardMajor"] {
			if ls > 4 {
				badness += Score(ls-4)*SuitLengthProblem
			}
			if lh > 4 {
				badness += Score(lh-4)*SuitLengthProblem
			}
		}
		return
	},
}

var TwoNT = BiddingRule{
	"2NT opening",
	regexp.MustCompile("^( P)*2N$"), nil,
	func (bidder Seat, h Hand, ms []string, cc ConventionCard, e *Ensemble) (badness Score, explanation string) {
		hcp := h.HCP()
		dist := h.DistPoints()
		if hcp > cc.Pts["TwoNTmax"] {
			badness += Score(hcp-cc.Pts["TwoNTmax"])*PointValueProblem
		} else if hcp < cc.Pts["TwoNTmin"] {
			badness += Score(cc.Pts["TwoNTmin"]-hcp)*PointValueProblem
		}
		if dist > 1 {
			badness += Score(dist-1)*PointValueProblem
		}
		ls := byte(h>>28) & 15
		lh := byte(h>>20) & 15
		// I assume here that 2NT is same length as 1NT
		if !cc.Options["OneNT5CardMajor"] {
			if ls > 4 {
				badness += Score(ls-4)*SuitLengthProblem
			}
			if lh > 4 {
				badness += Score(lh-4)*SuitLengthProblem
			}
		}
		return
	},
}

var Gambling3NT = BiddingRule{
	"Gambling 3NT",
	regexp.MustCompile("^( P)*3N$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (s Score, e string)) {
		if !cc.Options["Gambling3NT"] {
			return nil
		}
		return func(h Hand) (badness Score, explanation string) {
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
		}
	}, nil,
}
