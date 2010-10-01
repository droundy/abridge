package bridge

import (
	"regexp"
)

var Jacobi = BiddingRule{
	"Jacobi transfer (forcing)",
	regexp.MustCompile("^(..)?( P..)?1N P2([DH])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (s Score, e string)) {
		if !cc.Options["Jacobi"] {
			return nil
		}
		isovercall := (len(ms[1]) > 0 && ms[1] != " P") || (len(ms[2]) > 0 && ms[2][2:] != " P")
		if isovercall && !cc.Options["NTOvercallSystemsOn"] {
			// Systems *not* on after overcall
			return nil
		}
		mysuit := stringToSuitNumber(ms[3])+1
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
	regexp.MustCompile("^(..)?( P..)?1N P2([DH]) P2([HS])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		if !cc.Options["Jacobi"] {
			return nil
		}
		isovercall := (len(ms[1]) > 0 && ms[1] != " P") || (len(ms[2]) > 0 && ms[2][2:] != " P")
		if isovercall && !cc.Options["NTOvercallSystemsOn"] {
			// Systems *not* on after overcall
			return nil
		}
		if ms[3] == "D" && ms[4] == "S" {
			return nil
		}
		return func(h Hand) (Score,string) {
			return 0, ""
		}
	}, nil,
}

var JacobiSuperAccept = BiddingRule{
	"Jacobi super accept",
	regexp.MustCompile("^(..)?( P..)?1N P2([DH]) P3([HS])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		if !cc.Options["Jacobi"] {
			return nil
		}
		isovercall := (len(ms[1]) > 0 && ms[1] != " P") || (len(ms[2]) > 0 && ms[2][2:] != " P")
		if isovercall && !cc.Options["NTOvercallSystemsOn"] {
			// Systems *not* on after overcall
			return nil
		}
		sv := stringToSuitNumber(ms[4])
		if ms[3] == "D" && sv == Spades {
			return nil
		}
		return func(h Hand) (badness Score, explanation string) {
			ns := byte(h >> (4+8*sv)) & 15
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
	regexp.MustCompile("^(..)?( P..)?1N P2([DH]) P..$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		if !cc.Options["Jacobi"] {
			return nil
		}
		isovercall := (len(ms[1]) > 0 && ms[1] != " P") || (len(ms[2]) > 0 && ms[2][2:] != " P")
		if isovercall && !cc.Options["NTOvercallSystemsOn"] {
			// Systems *not* on after overcall
			return nil
		}
		mysuit := stringToSuitNumber(ms[3])+1
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
	regexp.MustCompile("^(..)?( P..)?1N P2C$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		if !cc.Options["Stayman"] {
			return nil
		}
		isovercall := (len(ms[1]) > 0 && ms[1] != " P") || (len(ms[2]) > 0 && ms[2][2:] != " P")
		if isovercall && !cc.Options["NTOvercallSystemsOn"] {
			// Systems *not* on after overcall
			return nil
		}
		return func(h Hand) (badness Score, explanation string) {
			return // Says nothing...
		}
	}, nil,
}

var StaymanTwo = BiddingRule{
	"Stayman (forcing)",
	regexp.MustCompile("^(..)?( P..)?2N P3C$"),
	Stayman.mkscore, nil,
}

var StaymanResponse = BiddingRule{
	"Stayman response",
	regexp.MustCompile("^(..)?( P..)?1N P2C P2([DHS])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		if !cc.Options["Stayman"] {
			return nil
		}
		isovercall := (len(ms[1]) > 0 && ms[1] != " P") || (len(ms[2]) > 0 && ms[2][2:] != " P")
		if isovercall && !cc.Options["NTOvercallSystemsOn"] {
			// Systems *not* on after overcall
			return nil
		}
		switch ms[3] {
		case "D":
			return func(h Hand) (badness Score, explanation string) {
				lh := byte(h>>20) & 15
				ls := byte(h>>28) & 15
				if lh > 3 {
					badness += Score(lh - 3)*SuitLengthProblem
				}
				if ls > 3 {
					badness += Score(ls - 3)*SuitLengthProblem
				}
				return
			}
		case "H":
			return func(h Hand) (badness Score, explanation string) {
				lh := byte(h>>20) & 15
				if lh < 4 {
					badness += Score(4 - lh)*SuitLengthProblem
				}
				return
			}
		}
		return func(h Hand) (badness Score, explanation string) {
			lh := byte(h>>20) & 15
			ls := byte(h>>28) & 15
			if lh > 3 {
				badness += Score(lh - 3)*SuitLengthProblem
			}
			if ls < 4 {
				badness += Score(4 - ls)*SuitLengthProblem
			}
			return
		}
	}, nil,
}

var StaymanTwoResponse = BiddingRule{
	"Stayman",
	regexp.MustCompile("^(..)?( P..)?2N P3C P3([DHS])$"),
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
					// We must have AKQ in our suit!
					badness += Score(14 - (d&15))*BigFudge
				}
			} else {
				if (c>>4) < 7 {
					badness += Score(7 - (c>>4))*SuitLengthProblem
				}
				if (c & 15 < 14) {
					// We must have AKQ in our suit!
					badness += Score(14 - (c&15))*BigFudge
				}
			}
			return
		}
	}, nil,
}

var Gambling3NTforcingquery = BiddingRule{
	"Gambling 3NT query (forcing)",
	regexp.MustCompile("^( P)*3N P4C$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (s Score, e string)) {
		if !cc.Options["Gambling3NT"] {
			return nil
		}
		return func(h Hand) (badness Score, explanation string) {
			ntsafe := SafeContractInThisSuit(bidder, h, NoTrump, e)
			if ntsafe > 2 {
				badness += Score(ntsafe - 2)*SuitLengthProblem
			}
			nc := byte(h >> 4) & 15
			nd := byte(h >> 12) & 15
			if nc < 1 {
				badness += SuitLengthProblem
			}
			if nd < 1 {
				badness += SuitLengthProblem
			}
			partner := (bidder+2)&3
			// FIXME: The following should take into account that if I have
			// more than a jack in clubs or diamonds, my partner's long suit
			// must be the other one.
			worstscore := Score(0)
			for _,t := range e.tables {
				best := ScoreHands(h, t[partner], 5, Clubs)
				bd := ScoreHands(h, t[partner], 4, Diamonds) // Response is 4D or 5C
				if bd < best {
					best = bd
				} 
				if best > worstscore {
					worstscore = best
				}
			}
			badness += worstscore
			return
		}
	}, nil,
}


var Gambling3NTquery = BiddingRule{
	"Gambling 3NT query",
	regexp.MustCompile("^( P)*3N P([567])C$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (s Score, e string)) {
		if !cc.Options["Gambling3NT"] {
			return nil
		}
		bidlevel := int(ms[2][0]-'0')
		return func(h Hand) (badness Score, explanation string) {
			ntsafe := SafeContractInThisSuit(bidder, h, NoTrump, e)
			if ntsafe > 2 {
				badness += Score(ntsafe - 2)*SuitLengthProblem
			}
			nc := byte(h >> 4) & 15
			nd := byte(h >> 12) & 15
			if nc < 1 {
				badness += SuitLengthProblem
			}
			if nd < 1 {
				badness += SuitLengthProblem
			}
			partner := (bidder+2)&3
			// FIXME: The following should take into account that if I have
			// more than a jack in clubs or diamonds, my partner's long suit
			// must be the other one.
			worstscore := Score(0)
			for _,t := range e.tables {
				best := ScoreHands(h, t[partner], bidlevel, Clubs)
				bd := ScoreHands(h, t[partner], bidlevel, Diamonds)
				if bd < best {
					best = bd
				} 
				if best > worstscore {
					worstscore = best
				}
			}
			badness += worstscore
			return
		}
	}, nil,
}

var Gambling3NTresponse = BiddingRule{
	"Gambling 3NT response",
	regexp.MustCompile("^( P)*3N P([567])C P([567]D| P)$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (s Score, e string)) {
		if !cc.Options["Gambling3NT"] {
			return nil
		}
		if ms[3] == " P" {
			return func(h Hand) (badness Score, explanation string) {
				nc := byte(h >> 4) & 15
				nd := byte(h >> 12) & 15
				if nd > nc {
					badness += Score(nd-nc)*SuitLengthProblem
				}
				return
			}
		}
		if ms[3][0] != ms[2][0] {
			// It's a weird jump!
			return nil
		}
		return func(h Hand) (badness Score, explanation string) {
			nc := byte(h >> 4) & 15
			nd := byte(h >> 12) & 15
			if nc > nd {
				badness += Score(nc-nd)*SuitLengthProblem
			}
			return
		}
	}, nil,
}

var Gambling3NTforcedresponse = BiddingRule{
	"Gambling 3NT forced response",
	regexp.MustCompile("^( P)*3N P4C P([45]D|5C)$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (s Score, e string)) {
		if !cc.Options["Gambling3NT"] {
			return nil
		}
		hcprange := e.HCP(bidder)
		hcpinvite := (hcprange.Max - hcprange.Min)/2 + hcprange.Min
		switch ms[2] {
		case "4D":
			return func(h Hand) (badness Score, explanation string) {
				nc := byte(h >> 4) & 15
				nd := byte(h >> 12) & 15
				if nc > nd {
					badness += Score(nc-nd)*SuitLengthProblem
				}
				hcp := h.HCP()
				if hcp > hcpinvite {
					badness += Score(hcp - hcpinvite)*PointValueProblem
				}
				return
			}
		case "5D":
			return func(h Hand) (badness Score, explanation string) {
				nc := byte(h >> 4) & 15
				nd := byte(h >> 12) & 15
				if nc > nd {
					badness += Score(nc-nd)*SuitLengthProblem
				}
				hcp := h.HCP()
				if hcp < hcpinvite + 1 {
					badness += Score(hcpinvite + 1 - hcp)*PointValueProblem
				}
				return
			}
		}
		return func(h Hand) (badness Score, explanation string) {
			nc := byte(h >> 4) & 15
			nd := byte(h >> 12) & 15
			if nd > nc {
				badness += Score(nd-nc)*SuitLengthProblem
			}
			return
		}
	}, nil,
}
