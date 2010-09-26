package bridge

import (
	"regexp"
)

var TwoLevelRebidSuit = BiddingRule{
	"Rebid in my suit after cheap unlimited response",
	regexp.MustCompile("^( P)*1([CDHS]) P1([DHS]) P2([CDHS])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (func(Hand) (Score,string)) {
		if ms[2] != ms[4] {
			return nil // This isn't a rebid of my suit
		}
		mysuit := stringToSuitNumber(ms[2])
		theirsuit := stringToSuitNumber(ms[3])
		return func(h Hand) (badness Score, explanation string) {
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
		}
	}, nil,
}

var RebidSuit = BiddingRule{
	"Rebid my suit",
	regexp.MustCompile("^.*([12])([CDHS])......([123])([CDHS])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (func(Hand) (Score,string)) {
		if ms[2] != ms[4] {
			return nil // This isn't a rebid of my suit
		}
		if ms[3][0] != ms[1][0] + 1 {
			return nil // This is a jump rebid, which is something else entirely!
		}
		for i:=len(ms[0])-5; i>0; i -= 8 {
			pardbid := ms[0][i:i+1]
			if pardbid == ms[2] {
				return nil // Partner has bid my suit.
			}
		}
		num := int(ms[3][0]) - '0'
		mysuit := stringToSuitNumber(ms[2])
		gamelevel := 5 - int(mysuit/2)
		pard := (bidder+2)&3
		mymin := e.SuitLength(bidder, mysuit).Min
		pardmin := e.SuitLength(pard, mysuit).Min
		if mymin + 1 + pardmin >= 8 {
			// We either already know we have a fit, or *I* already know we
			// have a fit, and either way, the bid is an ordinary "Natural"
			// bid.
			return nil
		}
		pardpts := e.PointCount(pard)
		return func(h Hand) (badness Score, explanation string) {
			mysuitlen := byte(h >> (4+mysuit*8)) & 15
			pts := h.PointCount()
			if pts + pardpts.Min < suitlevels[num] {
				// We'd better have the points to play at this level!
				badness += Score(suitlevels[num]-pts-pardpts.Min)*PointValueProblem
			}
			if pardmin + mysuitlen >= 8 {
				if pts + pardpts.Min >= suitlevels[num+1] && num < gamelevel && pts + pardpts.Max >= suitlevels[gamelevel] {
					// I should have bid higher, since we might have reached game! (and it was safe)
					
				}
			}
			if pts > 15 {
				badness += Score(pts-15)*PointValueProblem
			}
			if mysuitlen <  mymin + 1 {
				badness += Score(1 + mymin - mysuitlen)*SuitLengthProblem
			}
			return
		}
	}, nil,
}

var CheapNTRebid = BiddingRule{
	"Cheap no-trump rebid",
	regexp.MustCompile("^( P)*1([CDH]) P1([DHS]) P1N$"), nil,
	func (bidder Seat, h Hand, ms []string, cc ConventionCard, e *Ensemble) (badness Score, explanation string) {
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
		return
	},
}

var CheapRebid = BiddingRule{
	"Cheap rebid (forcing)",
	regexp.MustCompile("^( P)*1[CD] P1([DH]) P1([HS])$"), nil,
	func (bidder Seat, h Hand, ms []string, cc ConventionCard, e *Ensemble) (badness Score, explanation string) {
		theirsuit := stringToSuitNumber(ms[2])
		mysuit := stringToSuitNumber(ms[3])
		mysuitlen := byte(h >> (4 + mysuit*8)) & 15
		theirsuitlen := byte(h >> (4+theirsuit*8)) & 15
		heartlen := byte(h >> 20) & 15
		if theirsuit == Hearts && theirsuitlen > 3 {
			// We have a major fit!
			badness += Score(theirsuitlen-3)*SuitLengthProblem
		}
		if theirsuit < Hearts && mysuit == Spades && heartlen > 3 {
			// Should mention hearts.
			badness += Score(heartlen-3)*SuitLengthProblem
		}
		if mysuitlen < 4 {
			// I'd better have four of my suit
			badness += Score(4-mysuitlen)*SuitLengthProblem
		}
		return
	},
}
