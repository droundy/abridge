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

var MichaelsCuebidMinorQueryResponse = BiddingRule{
	"Answering about Michaels minor suit",
	regexp.MustCompile("^.*1([HS])2([HS])..2N..3([CD])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		if ms[1] != ms[2] {
			return nil // not the same suit
		}
		if cc.Radio["MajorCuebid"] != "Michaels" {
			return nil // We aren't playing Michaels
		}
		ourminor := stringToSuitNumber(ms[3])
		return func(h Hand) (badness Score, explanation string) {
			length := (h >> (4+8*ourminor)) & 15
			if length < 5 {
				badness += Score(5-length)*SuitLengthProblem
			}
			return
		}
	}, nil,
}

var MichaelsCuebidMinorQuery = BiddingRule{
	"Asking about Michaels minor suit (forcing)",
	regexp.MustCompile("^.*1([HS])2([HS])..2N$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		if ms[1] != ms[2] {
			return nil // not the same suit
		}
		if cc.Radio["MajorCuebid"] != "Michaels" {
			return nil // We aren't playing Michaels
		}
		return func(h Hand) (badness Score, explanation string) {
			// This is a forcing bid, so it doesn't deny a major fit.  It
			// also doesn't promise a minor fit, as we may be thinking about
			// slam (and looking for a cross-ruff).

			// So all we promise is that we're confident that we can make
			// *some* sort of a contract above this level!
			return WorstCaseSuit(bidder, h, 2, NoTrump, cc, e)
		}
	}, nil,
}

var MichaelsCuebid = BiddingRule{
	"Michaels cue bid (forcing)",
	regexp.MustCompile("^.*1([CDHS])2([CDHS])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		if ms[1] != ms[2] {
			return nil // not the same suit
		}
		theirsuit := stringToSuitNumber(ms[2])
		minpts := Points(13)
		// Listen to the convention card:
		if theirsuit < Hearts && cc.Radio["MinorCuebid"] != "Michaels" {
			return nil
		} else if theirsuit < Hearts {
			return func(h Hand) (badness Score, explanation string) {
				pts := h.PointCount()
				if pts < minpts {
					badness += Fudge
				}
				if pts < minpts - 1 {
					badness += Score(minpts-1-pts)*PointValueProblem
				}
				length := (h >> 28) & 15
				if length < 5 {
					badness += Score(5-length)*SuitLengthProblem
				}
				length = (h >> 20) & 15
				if length < 5 {
					badness += Score(5-length)*SuitLengthProblem
				}
				return
			}
		} else if cc.Radio["MajorCuebid"] != "Michaels" {
			return nil
		}
		return func(h Hand) (badness Score, explanation string) {
			pts := h.PointCount()
			if pts < minpts {
				badness += Fudge
			}
			if pts < minpts - 1 {
				badness += Score(minpts-1-pts)*PointValueProblem
			}
			if theirsuit != Spades {
				length := (h >> 28) & 15
				if length < 5 {
					badness += Score(5-length)*SuitLengthProblem
				}
			} else {
				length := (h >> 20) & 15
				if length < 5 {
					badness += Score(5-length)*SuitLengthProblem
				}
			}
			lengthc := (h >> 4) & 15
			length := (h >> 12) & 15
			if length < lengthc {
				length = lengthc
			}
			if length < 5 {
				badness += Score(5-length)*SuitLengthProblem
			}
			return
		}
	}, nil,
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
			minpts = 12
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

var DirectOneNTOvercall = BiddingRule{
	"Direct 1NT overcall",
	regexp.MustCompile("^( P)*1([CDHS])1N$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		hcpmin := cc.Pts["DirectOvercallNTmin"]
		hcpmax := cc.Pts["DirectOvercallNTmax"]
		fivecardmajor := cc.Options["OneNT5CardMajor"]
		theirsuit := stringToSuitNumber(ms[2])
		return func(h Hand) (badness Score, explanation string) {
			hcp := h.HCP()
			dist := h.DistPoints()
			cardsintheirsuit := Suit(h >> (8*theirsuit))
			ptsintheirsuit := PointCount[cardsintheirsuit] - DistPoints[cardsintheirsuit]
			if ptsintheirsuit < 2 {
				// We don't have a stopper in their suit!
				badness += Score(2 - ptsintheirsuit)*PointValueProblem
			}
			if hcp > hcpmax {
				badness += Score(hcp-hcpmax)*PointValueProblem
			} else if hcp < hcpmin {
				badness += Score(hcpmin-hcp)*PointValueProblem
			}
			if dist > 1 {
				badness += Score(dist-1)*PointValueProblem
			}
			ls := byte(h>>28) & 15
			lh := byte(h>>20) & 15
			if !fivecardmajor {
				if ls > 4 {
					badness += Score(ls-4)*SuitLengthProblem
				}
				if lh > 4 {
					badness += Score(lh-4)*SuitLengthProblem
				}
			}
			return
		}
	}, nil,
}

var BalancingOneNTOvercall = BiddingRule{
	"Balancing 1NT overcall",
	regexp.MustCompile("^1([CDHS]) P P1N$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		hcpmin := cc.Pts["BalancingOvercallNTmin"]
		hcpmax := cc.Pts["BalancingOvercallNTmax"]
		fivecardmajor := cc.Options["OneNT5CardMajor"]
		theirsuit := stringToSuitNumber(ms[1])
		return func(h Hand) (badness Score, explanation string) {
			hcp := h.HCP()
			dist := h.DistPoints()
			cardsintheirsuit := Suit(h >> (8*theirsuit))
			ptsintheirsuit := PointCount[cardsintheirsuit] - DistPoints[cardsintheirsuit]
			if ptsintheirsuit < 1 {
				// We don't have even a marginal stopper in their suit!
				badness += Score(1 - ptsintheirsuit)*PointValueProblem
			}
			if hcp > hcpmax {
				badness += Score(hcp-hcpmax)*PointValueProblem
			} else if hcp < hcpmin {
				badness += Score(hcpmin-hcp)*PointValueProblem
			}
			if dist > 1 {
				badness += Score(dist-1)*PointValueProblem
			}
			ls := byte(h>>28) & 15
			lh := byte(h>>20) & 15
			if !fivecardmajor {
				if ls > 4 {
					badness += Score(ls-4)*SuitLengthProblem
				}
				if lh > 4 {
					badness += Score(lh-4)*SuitLengthProblem
				}
			}
			return
		}
	}, nil,
}
