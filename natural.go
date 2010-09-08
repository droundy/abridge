package bridge

import (
	"regexp"
)

var Natural = BiddingRule{
	"Natural",
	regexp.MustCompile("(.)([^PX])$"),
	func (bidder Seat, h Hand, ms []string, e Ensemble) (badness Score, nothandled bool) {
		pts := h.PointCount()
		hcp := h.HCP()
		partner := (bidder+2)&3
		hcprange := e.HCP(partner)
		ptsrange := e.PointCount(partner)
		minpts := pts + ptsrange.Min
		maxpts := pts + ptsrange.Max
		rspades := e.SuitLength(partner, Spades)
		rhearts := e.SuitLength(partner, Hearts)
		minS := rspades.Min + byte((h >> 28)&15)
		minH := rhearts.Min + byte((h >> 20)&15)
		// gamelevel is the bid needed for game.
		gamelevel := 4
		// pointlevels are the points that are needed for various bids
		pointlevels := map[int]Points{ 2:19, 3:23, 4:26, 5:29, 6:33, 7:37, 8:60 }
		num := int(ms[1][0] - '0') // the level we are bid to
		switch ms[2] {
		case "N":
			if minS > 7 {
				badness += Score(minS - 7)*SuitLengthProblem
			}
			if minH > 7 {
				badness += Score(minH - 7)*SuitLengthProblem
			}
			// redefine gamelevel and pointlevels appropriately
			gamelevel = 3
			pointlevels = map[int]Points{ 2:20, 3:26, 4:33, 5:33, 6:33, 7:37, 8:60 }
			// in notrump, hcp are what is relevant
			minpts = hcp + hcprange.Min
			maxpts = hcp + hcprange.Max
		case "S","H","D","C":
			mysuit := stringToSuitNumber(ms[2])
			myownsuitlen := byte((h >> (4+8*mysuit))&15)
			rsuit := e.SuitLength(partner, mysuit)
			mysuitlen := myownsuitlen + rsuit.Min
			if mysuitlen < 8 {
				// We always want a guaranteed fit.
				badness += Score(8 - mysuitlen)*SuitLengthProblem
			}
			if mysuitlen > 8 {
				for i:=2; i<6; i++ {
					pointlevels[i] -= Points(mysuitlen-8) // we need a one fewer point per extra trump?
				}
			}
			if mysuit < Hearts {
				gamelevel = 5
			}
			if num == 6 || num == 7 {
				// Special case for splinter slams:
				for i:=uint(Clubs); i<=Spades; i++ {
					if i != mysuit {
						rsuit := e.SuitLength(partner, i)
						if rsuit.Max == 1 {
							// We have a splinter situation!
							othersuits := [4]bool{true,true,true,true}
							othersuits[i] = false
							theirhcprange := e.SuitHCP(partner, othersuits)
							myhcp := Points(0)
							for j:=uint(Clubs); j<=Spades; j++ {
								if j != i {
									myhcp += HCP[byte(h>>(8*j))]
								}
							}
							if theirhcprange.Min + myhcp < 25 {
								badness += Score(25 - theirhcprange.Min - myhcp)*PointValueProblem;
							}
							if num == 7 {
								if h & (Hand(Ace) << (8*i)) == 0 && (h >> (4+8*i)) & 15 > 0 {
									// For grand slam, we need the missing ace!
									badness += PointValueProblem
								}
								if theirhcprange.Min + myhcp < 27 {
									badness += Score(27 - theirhcprange.Min - myhcp)*PointValueProblem
								}
							}
							return
						}
					}
				}
			}
		}
		if minpts < pointlevels[num] {
			// we need to guarantee pointlevels[num] to bid at this level
			badness += Score(pointlevels[num]-minpts)*PointValueProblem
		} else if minpts >= pointlevels[num+1] && (num < gamelevel || num == 5 || num == 6) {
			// if we can guarantee pointlevels[num+1], then we should bid at
			// *that* level (if it's at-or-below game, or a slam bid)
			badness += Score(minpts - pointlevels[num+1])*PointValueProblem
		}
		/*
		if num < gamelevel && (maxpts+minpts+1)/2 < pointlevels[num+1] {
			// We never want to bid a natural bid unless there is 50% of
			// partner's range that could lead to the next bid up.
			// This assumes we are bidding for game, and leaves out
			// competitive bidding...
			badness += Score(pointlevels[num+1] - (maxpts+minpts+1)/2)*PointValueProblem
		}
		 */
		if num < gamelevel && maxpts < pointlevels[gamelevel] {
			// We never want to bid a natural bid unless there is at least
			// some chance of game.  This assumes we are bidding for game,
			// and leaves out competitive bidding...
			badness += Score(pointlevels[gamelevel] - maxpts)*PointValueProblem
		}
		return
	},
}
