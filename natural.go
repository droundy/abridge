package bridge

import (
	"regexp"
)

var lastbidregexp = regexp.MustCompile("([1234567])(.)( .)*")
func PossibleNondoubleBids(bid string) []string {
	ms := lastbidregexp.FindStringSubmatch(bid)
	if ms != nil || len(ms) < 3 {
		ms = []string{"","0","N"} // very hokey trick to avoid a special case
	}
	out := make([]string, 0, 37)
	for _,level := range []byte("1234567") {
		if level == ms[1][0] {
			mins := NoTrump+1
			switch ms[2][0] {
			case 'C': mins = Diamonds
			case 'D': mins = Hearts
			case 'H': mins = Spades
			case 'S': mins = NoTrump
			}
			for sv:=mins; sv <= NoTrump; sv++ {
				out = out[0:len(out)+1]
				out[len(out)-1] = string([]byte{level})+SuitLetter[sv]
			}
			out[len(out)-1] = string([]byte{level, ms[2][0]})
		} else if level > ms[1][0] {
			out = out[0:len(out)+1]
			out[len(out)-1] = string([]byte{level, ms[2][0]})
		}
	}
	return out
}

/*
var LimitPass = BiddingRule {
	"Limiting pass",
	regexp.MustCompile("^(..)?(..)?(..)?(.[^P]......)* P$"),
	func (ms []string) (score func(bidder Seat, h Hand, e Ensemble) (s Score, nothandled bool)) {
		possbids := PossibleNondoubleBids(ms[0])
		allrules := make([]ScoringRule,0,len(possbids)*len(Convention))
		for _,b := range possbids {
			rule := makeScoringRules(ms[0] + b)
			allrules = allrules[0:len(allrules)+len(rule)]
			for i,r := range rule {
				allrules[len(allrules)-1-i] = r
			}
		}
		return func(bidder Seat, h Hand, e Ensemble) (badness Score, nothandled bool) {
			for _,r := range allrules {
				sc,_ := simpleScore(bidder, h, r, e)
				if sc == 0 {
					badness += 137
				}
			}
			return
		}
	},
	nil,
}
*/

var Natural = BiddingRule{
	"Natural",
	regexp.MustCompile("(.)([^PX])$"), nil,
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
			if num > 4 {
				// Special case for splinter slams and slam invites:
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
							neededpts := Points(25)
							switch num {
							case 7:
								neededpts = 27
								if h & (Hand(Ace) << (8*i)) == 0 && (h >> (4+8*i)) & 15 > 0 {
									// For grand slam, we need the missing ace!
									badness += PointValueProblem
								}
							case 5: neededpts = 23 // Just invite to slam...
							}
							if theirhcprange.Min + myhcp < neededpts {
								badness += Score(neededpts - theirhcprange.Min - myhcp)*PointValueProblem;
							}
							return
						}
					}
				}
			}
			if num == 6 {
				// Special case when *I've* got a splintery hand and may have
				// been invited...
				for i:=uint(Clubs); i<=Spades; i++ {
					if i != mysuit {
						mysuit := byte(h >> (4+i*8))
						if mysuit < 2 {
							// We have a splintery situation!
							othersuits := [4]bool{true,true,true,true}
							othersuits[i] = false
							theirhcprange := e.SuitHCP(partner, othersuits)
							myhcp := Points(0)
							for j:=uint(Clubs); j<=Spades; j++ {
								if j != i {
									myhcp += HCP[byte(h>>(8*j))]
								}
							}
							neededpts := Points(25)
							if theirhcprange.Min + myhcp >= neededpts {
								// Ugly... I don't want to "add badness" in case
								// there's another way of counting that gives us
								// slam...
								return 0, false
								//badness += Score(neededpts - theirhcprange.Min - myhcp)*PointValueProblem;
							}
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
