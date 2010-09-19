package bridge

import (
	"regexp"
)

var StrongTwoResponse = BiddingRule{
	"Strong two response (forcing)",
	regexp.MustCompile("^( P)*2C P(.[^P])$"),
	func (bidder Seat, ms []string, e *Ensemble) (func(Hand) (Score,string)) {
		switch ms[2] {
		case "2D":
			return func(h Hand) (badness Score, explanation string) {
				hcp := h.HCP()
				if hcp > 6 {
					badness += Score(hcp-6)*PointValueProblem
				}
				return
			}
		case "2H", "2S":
			sv := uint(Hearts)
			if ms[2] == "2S" {
				sv = Spades
			}
			return func(h Hand) (badness Score, explanation string) {
				hcp := h.HCP()
				if hcp < 7 {
					badness += Score(7-hcp)*PointValueProblem
				} else if hcp > 9 {
					badness += Score(hcp-9)*PointValueProblem
				}
				ncards := byte((h >> (4+8*sv))&15)
				if ncards < 5 {
					badness += Score(5 - ncards)*SuitLengthProblem
				}
				return
			}
		case "2N":
			return func(h Hand) (badness Score, explanation string) {
				hcp := h.HCP()
				if hcp < 7 {
					badness += Score(7-hcp)*PointValueProblem
				} else if hcp > 9 {
					badness += Score(hcp-9)*PointValueProblem
				}
				dpts := h.DistPoints()
				if dpts > 2 {
					badness += Score(dpts-1)*SuitLengthProblem
				}
				return
			}
		case "3C", "3D", "3H", "3S":
			sv := stringToSuitNumber(ms[2][1:])
			return func(h Hand) (badness Score, explanation string) {
				hcp := h.HCP()
				if hcp < 10 {
					badness += Score(10-hcp)*PointValueProblem
				}
				ncards := byte((h >> (4+8*sv))&15)
				if ncards < 5 {
					badness += Score(5 - ncards)*SuitLengthProblem
				}
				return
			}
		case "3N":
			return func(h Hand) (badness Score, explanation string) {
				hcp := h.HCP()
				if hcp < 10 {
					badness += Score(10-hcp)*PointValueProblem
				}
				dpts := h.DistPoints()
				if dpts > 2 {
					badness += Score(dpts-1)*SuitLengthProblem
				}
				return
			}
		case "4C", "4D", "4H", "4S":
			sv := stringToSuitNumber(ms[2][1:])
			return func(h Hand) (badness Score, explanation string) {
				hcp := h.HCP()
				if hcp < 13 {
					badness += Score(13-hcp)*PointValueProblem
				}
				ncards := byte((h >> (4+8*sv))&15)
				if ncards < 5 {
					badness += Score(5 - ncards)*SuitLengthProblem
				}
				return
			}
		case "4N":
			return func(h Hand) (badness Score, explanation string) {
				hcp := h.HCP()
				if hcp < 13 {
					badness += Score(13-hcp)*PointValueProblem
				}
				dpts := h.DistPoints()
				if dpts > 2 {
					badness += Score(dpts-1)*SuitLengthProblem
				}
				return
			}
		}
		return nil
	}, nil,
}

var Splinter = BiddingRule{
	"Splinter (forcing)",
	regexp.MustCompile("^( P)*1([HS]) P(3S|4[CDH])$"),
	func (bidder Seat, ms []string, e *Ensemble) (func(Hand) (Score,string)) {
		opensuit := stringToSuitNumber(ms[2])
		splintersuit := uint(Spades)
		switch ms[3] {
		case "4C": splintersuit = Clubs
		case "4D": splintersuit = Diamonds
		case "4H": splintersuit = Hearts
		}
		if opensuit == splintersuit {
			return nil // Not a splinter!
		}
		return func(h Hand) (badness Score, explanation string) {
			openlen := byte(h >> (4+opensuit*8)) & 15
			splinterlen := byte(h >> (4+splintersuit*8)) & 15
			if splinterlen > 1 {
				badness += Score(splinterlen - 1)*SuitLengthProblem
			}
			hcp_inside := HCP[Suit(h >> (splintersuit*8))]
			hcp_outside := h.HCP() - hcp_inside
			// Splinter indicates 10-12 hcp outside the singleton
			if hcp_outside < 10 {
				badness += Score(10 - hcp_outside)*PointValueProblem
			} else if hcp_outside > 12 {
				badness += Score(hcp_outside - 12)*PointValueProblem
			}
			if hcp_inside > 2 {
				// To bid a splinter, we'd better not have more than a queen in
				// the singleton.
				badness += Score(hcp_inside - 2)*PointValueProblem
			}
			if openlen < 4 {
				// A splinter bid promises four-card support.
				badness += Score(4-openlen)*SuitLengthProblem
			}
			return
		}
	}, nil,
}

var MajorInvitation = BiddingRule{
	"Major invitation",
	regexp.MustCompile("^( P)*1([HS]) P3([HS])$"),
	func (bidder Seat, ms []string, e *Ensemble) (func(Hand) (Score,string)) {
		opensuit := stringToSuitNumber(ms[2])
		mysuit := stringToSuitNumber(ms[3])
		if mysuit != opensuit {
			return nil // This isn't support
		}
		return func(h Hand) (badness Score, explanation string) {
			pts := h.PointCount()
			if pts < 10 {
				badness += Score(10-pts)*PointValueProblem
			} else if pts > 11 {
				badness += Score(pts-11)*PointValueProblem
			}
			mysuitlen := byte(h >> (4 + mysuit*8)) & 15
			if mysuitlen < 3 {
				badness += Score(3-mysuitlen)*SuitLengthProblem
			}
			return
		}
	}, nil,
}

var MajorSupport = BiddingRule{
	"Major support",
	regexp.MustCompile("^( P)*1([HS]) P2([HS])$"),
	func (bidder Seat, ms []string, e *Ensemble) (func(Hand) (Score,string)) {
		opensuit := stringToSuitNumber(ms[2])
		mysuit := stringToSuitNumber(ms[3])
		if mysuit != opensuit {
			return nil // This isn't support
		}
		return func(h Hand) (badness Score, explanation string) {
			pts := h.PointCount()
			if pts < 6 {
				badness += Score(6-pts)*PointValueProblem
			} else if pts > 9 {
				badness += Score(pts-9)*PointValueProblem
			}
			mysuitlen := byte(h >> (4 + mysuit*8)) & 15
			if mysuitlen < 3 {
				badness += Score(3-mysuitlen)*SuitLengthProblem
			}
			return
		}
	}, nil,
}

var TwoOverOne = BiddingRule{
	"Two over one (forcing)",
	regexp.MustCompile("^( P)*1([DHS]) P2([CDH])$"),
	func (bidder Seat, ms []string, e *Ensemble) (func(Hand) (Score,string)) {
		opensuit := stringToSuitNumber(ms[2])
		mysuit := stringToSuitNumber(ms[3])
		if mysuit == opensuit {
			return nil // This isn't a two-over-one bid
		}
		return func(h Hand) (badness Score, explanation string) {
			pts := h.PointCount()
			if pts < 10 {
				badness += Score(10-pts)*PointValueProblem
			}
			mysuitlen := byte(h >> (4 + mysuit*8)) & 15
			if mysuitlen < 4 {
				badness += Score(4-mysuitlen)*SuitLengthProblem
			}
			heartlen := byte(h >> 20) & 15
			spadelen := byte(h >> 28) & 15
			if opensuit < Spades && spadelen > 3 {
				badness += Score(spadelen-3)*SuitLengthProblem
			}
			if opensuit < Hearts && heartlen > 3 && mysuitlen != Hearts {
				badness += Score(heartlen-3)*SuitLengthProblem
			}
			if opensuit == Hearts && heartlen > 2 && pts < 15 {
				b1 := Score(heartlen-2)*SuitLengthProblem
				b2 := Score(15-pts)*PointValueProblem
				badness += b1.min(b2)
			}
			if opensuit == Spades && spadelen > 2 && pts < 15 {
				b1 := Score(spadelen-2)*SuitLengthProblem
				b2 := Score(15-pts)*PointValueProblem
				badness += b1.min(b2)
			}
			return
		}
	}, nil,
}

var CheapResponse = BiddingRule{
	"Cheap response to one (forcing)",
	regexp.MustCompile("^( P)*1([CDH]) P1([DHS])$"), nil,
	func (bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score, explanation string) {
		pts := h.PointCount()
		if pts < 6 {
			badness += Score(6-pts)*PointValueProblem
		}
		opensuit := stringToSuitNumber(ms[2])
		mysuit := stringToSuitNumber(ms[3])
		mysuitlen := byte(h >> (4 + mysuit*8)) & 15
		heartlen := byte(h >> 20) & 15
		spadelen := byte(h >> 28) & 15
		if opensuit == Spades && spadelen > 2 {
			// We missed an opening bid!
			badness += Score(spadelen-2)*SuitLengthProblem
		}
		if opensuit == Hearts && heartlen > 2 && !(ms[3] == "S" && pts > 9) {
			// We can only bid 1S if we really have good reason to force the
			// bid... i.e. a strongish hand.
			badness += Score(heartlen-2)*SuitLengthProblem
		}
		switch mysuit {
		case Hearts:
			if mysuitlen < 4 {
				badness += Score(4-mysuitlen)*SuitLengthProblem
			}
		case Spades:
			if mysuitlen < 4 {
				badness += Score(4-mysuitlen)*SuitLengthProblem
			}
			if heartlen > 3 && opensuit != Hearts && spadelen < 6 {
				// Skipping hearts denies 4 hearts, unless you've got 6 spades
				b1 := Score(heartlen-3)*SuitLengthProblem
				b2 := Score(7-spadelen)*SuitLengthProblem
				badness += b1.min(b2)
			}
		case Diamonds:
			if mysuitlen < 4 {
				badness += Score(4-mysuitlen)*SuitLengthProblem
			}
			if heartlen > 3 {
				badness += Score(heartlen-3)*SuitLengthProblem
			}
			if spadelen > 3 {
				badness += Score(spadelen-3)*SuitLengthProblem
			}
		}
		return
	},
}


var CheapNTResponse = BiddingRule{
	"Weak NT response",
	regexp.MustCompile("^( P)*1([CDHS]) P1N$"), nil,
	func (bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score, explanation string) {
		pts := h.PointCount()
		if pts < 6 {
			badness += Score(6-pts)*PointValueProblem
		}
		opensuit := stringToSuitNumber(ms[2])
		heartlen := byte(h >> 20) & 15
		spadelen := byte(h >> 28) & 15
		if opensuit == Spades && spadelen > 2 {
			// We missed an opening bid!
			badness += Score(spadelen-2)*SuitLengthProblem
		}
		if opensuit == Hearts && heartlen > 2 && !(ms[3] == "S" && pts > 9) {
			// We can only bid 1S if we really have good reason to force the
			// bid... i.e. a strongish hand.
			badness += Score(heartlen-2)*SuitLengthProblem
		}
		if (spadelen > 3 && opensuit < Spades) {
			badness += Score(spadelen-3)*SuitLengthProblem
		}
		if (heartlen > 3 && opensuit < Hearts) {
			badness += Score(heartlen-3)*SuitLengthProblem
		}
		if pts > 9 {
			badness += Score(pts - 9)*PointValueProblem
		}
		return
	},
}


var CheapCompetitionResponse = BiddingRule{
	"Cheap response to one over opponent (forcing)",
	regexp.MustCompile("^( P)*1([CDHS]).([^P])1([DHSN])$"), nil,
	func (bidder Seat, h Hand, ms []string, e *Ensemble) (badness Score, explanation string) {
		pts := h.PointCount()
		if pts < 8 {
			badness += Score(8-pts)*PointValueProblem
		}
		opensuit := stringToSuitNumber(ms[2])
		heartlen := byte(h >> 20) & 15
		spadelen := byte(h >> 28) & 15
		if opensuit == Spades && spadelen > 2 {
			// We missed an opening bid!
			badness += Score(spadelen-2)*SuitLengthProblem
		}
		if opensuit == Hearts && heartlen > 2 && !(ms[4] == "S" && pts > 10) {
			// We're denying a fit, unless we have 1S and a strong hand.
			badness += Score(heartlen-2)*SuitLengthProblem
		}
		if ms[4] == "N" {
			// We must have a stopper in opponent's suit!
			switch ms[3] {
			case "D":	badness += Suit(h>>8).UnStopped()
			case "H":	badness += Suit(h>>16).UnStopped()
			case "S":	badness += Suit(h>>24).UnStopped()
			}
			if (spadelen > 3 && opensuit < Spades && ms[3] != "S") {
				// We should have mentioned a 4-card spade suit, if possible...
				badness += Score(spadelen-3)*SuitLengthProblem
			}
			if (heartlen > 3 && opensuit < Hearts && ms[3] != "S" && ms[3] != "H") {
				// We should have mentioned a 4-card heart suit, if possible...
				badness += Score(heartlen-3)*SuitLengthProblem
			}
			hcp := h.HCP()
			if hcp < 9 {
				// I think we need 10 hcp to bid 1NT over an opponent overcall...
				badness += Score(9 - hcp)*PointValueProblem
			} else if hcp > 12 {
				// But if we have 13 points, we should bid higher.
				badness += Score(hcp - 12)*PointValueProblem
			}
			return // exit early, so we can assume mysuit is a valid suit
		}
		// Here we assume ms[4] is a real suit.
		mysuit := stringToSuitNumber(ms[4])
		mysuitlen := byte(h >> (4 + mysuit*8)) & 15
		switch mysuit {
		case Hearts:
			if mysuitlen < 4 {
				badness += Score(4-mysuitlen)*SuitLengthProblem
			}
		case Spades:
			if mysuitlen < 4 {
				badness += Score(4-mysuitlen)*SuitLengthProblem
			}
			if heartlen > 3 && opensuit != Hearts && spadelen < 6 {
				// Skipping hearts denies 4 hearts, unless you've got 6 spades
				b1 := Score(heartlen-3)*SuitLengthProblem
				b2 := Score(7-spadelen)*SuitLengthProblem
				badness += b1.min(b2)
			}
		case Diamonds:
			if mysuitlen < 4 {
				badness += Score(4-mysuitlen)*SuitLengthProblem
			}
			if heartlen > 3 {
				badness += Score(heartlen-3)*SuitLengthProblem
			}
			if spadelen > 3 {
				badness += Score(spadelen-3)*SuitLengthProblem
			}
		}
		return
	},
}
