package bridge

import (
	"regexp"
)

var PassOvercall = BiddingRule{
	"Pass an opportunity to overcall a one-suit bid",
	regexp.MustCompile("^( P)*1[CDH]( P..)? P$"), nil,
	func (bidder Seat, h Hand, ms []string, cc ConventionCard, e *Ensemble) (badness Score, explanation string) {
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
	func (bidder Seat, h Hand, ms []string, cc ConventionCard, e *Ensemble) (badness Score, explanation string) {
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
	regexp.MustCompile("^( P)*1[CDH]( P..)?1([DHS])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		mysuit := stringToSuitNumber(ms[3])
		hcpmin := cc.Pts["Overcallmin"]
		hcpmax := cc.Pts["Overcallmax"]
		suitpromise := byte(5)
		if cc.Options["FourCardOvercalls"] {
			suitpromise = 4
		}
		minpts := Points(13)
		hcpbadness := PointValueProblem
		if cc.Options["VeryLightOvercalls"] {
			hcpbadness = Fudge
			minpts = 7
		}
		return func(h Hand) (badness Score, explanation string) {
			mysuitlen := byte(h >> (4 + mysuit*8)) & 15
			// Give me a couple of extra points for each card in suit beyond five:
			pts := h.PointCount() + 2*(Points(mysuitlen) - 5)
			hcp := h.HCP()
			if pts < minpts {
				badness += Fudge
			}
			if pts < minpts - 1 {
				badness += Score(minpts-1-pts)*PointValueProblem
			}
			if hcp < hcpmin {
				badness += Score(hcpmin-hcp)*hcpbadness
			}
			if hcp > hcpmax {
				badness += Score(hcp-hcpmax)*hcpbadness
			}
			if mysuitlen < suitpromise {
				badness += Score(suitpromise - mysuitlen)*SuitLengthProblem
			}
			return
		}
	}, nil,
}

var TwoLevelOvercall = BiddingRule{
	"Two-level overcall",
	regexp.MustCompile("^( P)*[12]([CDHSN])( P..)?2([CDHS])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		mysuit := stringToSuitNumber(ms[4])
		minpts := Points(13)
		hcpmin := cc.Pts["Overcallmin"]
		hcpmax := cc.Pts["Overcallmax"]
		hcpbadness := PointValueProblem
		if cc.Options["VeryLightOvercalls"] {
			hcpbadness = Fudge
		}
		return func(h Hand) (badness Score, explanation string) {
			mysuitlen := byte(h >> (4 + mysuit*8)) & 15
			// Give me an extra point for each card in suit beyond five:
			pts := h.PointCount() + 2*(Points(mysuitlen) - 5)
			hcp := h.HCP()
			if pts < minpts {
				badness += Fudge
			}
			if pts < minpts - 1 {
				badness += Score(minpts-1-pts)*PointValueProblem
			}
			if mysuitlen < 5 {
				badness += Score(5 - mysuitlen)*SuitLengthProblem
			}
			if hcp < hcpmin {
				badness += Score(hcpmin-hcp)*hcpbadness
			}
			if hcp > hcpmax {
				badness += Score(hcp-hcpmax)*hcpbadness
			}
			return
		}
	}, nil,
}

var ThreeLevelOvercall = BiddingRule{
	"Three-level overcall",
	regexp.MustCompile("^(.. P)?[123]([CDHSN])( P..)?3([CDHS])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		mysuit := stringToSuitNumber(ms[4])
		minpts := Points(16)
		hcpmin := cc.Pts["Overcallmin"]
		hcpmax := cc.Pts["Overcallmax"]
		hcpbadness := PointValueProblem
		if cc.Options["VeryLightOvercalls"] {
			hcpbadness = Fudge
		}
		return func(h Hand) (badness Score, explanation string) {
			mysuitlen := byte(h >> (4 + mysuit*8)) & 15
			// Give me a couple of extra points for each card in suit beyond five:
			pts := h.PointCount() + 2*(Points(mysuitlen) - 5)
			hcp := h.HCP()
			if pts < minpts {
				badness += Fudge
			}
			if pts < minpts - 1 {
				badness += Score(minpts-1-pts)*PointValueProblem
			}
			if mysuitlen < 5 {
				badness += Score(5 - mysuitlen)*SuitLengthProblem
			}
			if hcp < hcpmin {
				badness += Score(hcpmin-hcp)*hcpbadness
			}
			if hcp > hcpmax {
				badness += Score(hcp-hcpmax)*hcpbadness
			}
			return
		}
	}, nil,
}
