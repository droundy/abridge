package bridge

import (
	"strings"
	"regexp"
)

func PickBid(h Hand, bidder Seat, oldbid string, e *Ensemble) (bid string, convention string) {
	bid = " P"
	pbids := PossibleNondoubleBids(oldbid)
	for _,b := range pbids {
		r := makeScoringRule(bidder, oldbid + b, e)
		if r.score(h) == 0 {
			return b, r.name
		}
	}
	return
}

var lastbidregexp = regexp.MustCompile("([1234567])(.)( .)*$")
func PossibleNondoubleBids(bid string) []string {
	ms := lastbidregexp.FindStringSubmatch(bid)
	if ms == nil || len(ms) < 3 {
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
		} else if level > ms[1][0] {
			for sv:=Clubs; sv <= NoTrump; sv++ {
				out = out[0:len(out)+1]
				out[len(out)-1] = string([]byte{level})+SuitLetter[sv]
			}
		}
	}
	return out
}

var PassOfForcing = BiddingRule {
	"Pass of forcing bid",
	regexp.MustCompile("^(.+) P P$"),
	func (bidder Seat, ms []string, e *Ensemble) (score func(h Hand) (s Score)) {
		if !strings.HasSuffix(e.Conventions[len(e.Conventions)-2], "(forcing)") {
			return nil
		}
		return func(h Hand) Score {
			return 1000
		}
	},
	nil,
}

var LimitPass = BiddingRule {
	"Limiting pass",
	regexp.MustCompile("^(..)?(..)?(..)?(.[^P]......)* P$"),
	func (bidder Seat, ms []string, e *Ensemble) (score func(h Hand) (s Score)) {
		possbids := PossibleNondoubleBids(ms[0])
		allrules := make([]*ScoringRule,0,len(possbids))
		for _,b := range possbids {
			if rule := makeScoringRule(bidder, ms[0] + b, e); rule != nil {
				allrules = allrules[0:len(allrules)+1]
				allrules[len(allrules)-1] = rule
			}
		}
		return func(h Hand) (badness Score) {
			for _,r := range allrules {
				sc := r.score(h)
				if sc == 0 {
					badness += 137
				}
			}
			return
		}
	},
	nil,
}

var Natural = BiddingRule{
	"Natural",
	regexp.MustCompile("(.)([^PX])$"),	
	func (bidder Seat, ms []string, e *Ensemble) (score func(h Hand) Score) {
		partner := (bidder+2)&3
		// gamelevel is the bid needed for game.
		gamelevel := 4
		// pointlevels are the points that are needed for various bids
		pointlevels := map[int]Points{ 2:19, 3:23, 4:26, 5:29, 6:33, 7:37, 8:60 }
		num := int(ms[1][0] - '0') // the level we are bid to

		gotsplinter := false
		var theirsplinterrange PointRange
		var splintersuit uint = 0
		var splinterpts = Points(25)
		var mysuit uint
		switch ms[2] {
		case "N":
			// redefine gamelevel and pointlevels appropriately
			gamelevel = 3
			pointlevels = map[int]Points{ 2:20, 3:26, 4:33, 5:33, 6:33, 7:37, 8:60 }
		case "S","H","D","C":
			mysuit = stringToSuitNumber(ms[2])
			if mysuit < Hearts {
				gamelevel = 5
			}
			if num > 4 {
				// Special case for splinter slams and slam invites:
				for i:=uint(Clubs); i<=Spades; i++ {
					if i != mysuit {
						rsuit := e.SuitLength(partner, i)
						if rsuit.Max == 1 {
							gotsplinter = true
							splintersuit = i
							// We have a splinter situation!
							othersuits := [4]bool{true,true,true,true}
							othersuits[i] = false
							theirsplinterrange = e.SuitHCP(partner, othersuits)
							switch num {
							case 7:	splinterpts = 27
							case 5: splinterpts = 23 // Just invite to slam...
							}
						}
					}
				}
			}
		}

		score = func(h Hand) (badness Score) {
			partner := (bidder+2)&3
			hcprange := e.HCP(partner)
			ptsrange := e.PointCount(partner)
			rspades := e.SuitLength(partner, Spades)
			rhearts := e.SuitLength(partner, Hearts)
			pts := h.PointCount()
			hcp := h.HCP()
			minpts := pts + ptsrange.Min
			maxpts := pts + ptsrange.Max
			minS := rspades.Min + byte((h >> 28)&15)
			minH := rhearts.Min + byte((h >> 20)&15)
			extrapts := Points(0)
			switch ms[2] {
			case "N":
				if minS > 7 {
					badness += Score(minS - 7)*SuitLengthProblem
				}
				if minH > 7 {
					badness += Score(minH - 7)*SuitLengthProblem
				}
				// in notrump, hcp are what is relevant
				minpts = hcp + hcprange.Min
				maxpts = hcp + hcprange.Max
			case "S","H","D","C":
				myownsuitlen := byte((h >> (4+8*mysuit))&15)
				partnerlen := e.SuitLength(partner, mysuit)
				mysuitlen := myownsuitlen + partnerlen.Min
				if mysuitlen < 8 {
					// We always want a guaranteed fit.
					badness += Score(8 - mysuitlen)*SuitLengthProblem
				}
				if mysuitlen > 8 {
					for i:=2; i<6; i++ {
						extrapts += Points(mysuitlen-8) // we need a one fewer point per extra trump?
					}
				}
				if num > 4 && gotsplinter {
					// Special case for splinter slams and slam invites:
					if splintersuit != mysuit {
						// We have a splinter situation!
						myhcp := Points(0)
						for j:=uint(Clubs); j<=Spades; j++ {
							if j != splintersuit {
								myhcp += HCP[byte(h>>(8*j))]
							}
						}
						if num == 7 {
							if h & (Hand(Ace) << (8*splintersuit)) == 0 && (h >> (4+8*splintersuit)) & 15 > 0 {
								// For grand slam, we need the missing ace!
								badness += PointValueProblem
							}
						}
						if theirsplinterrange.Min + myhcp < splinterpts {
							badness += Score(splinterpts - theirsplinterrange.Min - myhcp)*PointValueProblem;
						}
						return
					}
				}
				if num == 6 {
					// Special case when *I've* got a splintery hand and may
					// have been invited...
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
									return 0
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
				badness += Score(1 + minpts - pointlevels[num+1])*PointValueProblem
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
			if num > gamelevel && num < 6 && minpts < pointlevels[6]-2 {
				// We don't want to invite slam (i.e. bid over game) unless
				// there's a good reason to suspect that we've got the points
				// for slam.
				badness += Score(pointlevels[6] - 2 - minpts)*Fudge
			}
			return
		}
		return
	}, nil,
}

var Forced = BiddingRule{
	"Forced",
	regexp.MustCompile("(.+) P(.[^PX])$"),	
	func (bidder Seat, ms []string, e *Ensemble) (score func(h Hand) Score) {
		if len(e.Conventions) < 2 {
			return nil
		}
		if !strings.HasSuffix(e.Conventions[len(e.Conventions)-2], "(forcing)") {
			return nil
		}
		//fmt.Println("Working on Forced...", ms[0])
		scorer := makeUnforcedScoringRule(bidder, ms[0], e)
		allbutme := ms[0][0:len(ms[0])-2]
		pbids := PossibleNondoubleBids(allbutme)
		hirules := make([]*ScoringRule,0,len(pbids))
		lowrules := make([]*ScoringRule,0,5)
		lownames := make([]string,0,5)
		hinames := make([]string,0,len(pbids))
		for i,b := range pbids {
			if b != ms[2] {
				//fmt.Println("Making rule for", allbutme + b)
				if rule := makeUnforcedScoringRule(bidder, allbutme + b, e); rule != nil {
					if i < 5 {
						lowrules = lowrules[0:len(lowrules)+1]
						lowrules[len(lowrules)-1] = rule
						lownames = lownames[0:len(lownames)+1]
						lownames[len(lownames)-1] = rule.name+ " "+ b
					} else {
						hirules = hirules[0:len(hirules)+1]
						hirules[len(hirules)-1] = rule
						hinames = hinames[0:len(hinames)+1]
						hinames[len(hinames)-1] = rule.name+ " "+ b
					}
				}
			}
		}
		for i,b := range pbids {
			if i > 5 {
				break
			}
			if ms[2] == b {
				// This may be our best option
				return func(h Hand) (badness Score) {
					badness = scorer.score(h)
					if badness > 0 {
						//fmt.Print("I am running Forced on:\n", h, "for a bid of ", ms[2],"\n")
						//fmt.Println("Native badness is", badness)
						// Need to figure out if this is our only option
						bestlow := Score(1000000)
						for _,r := range lowrules {
							b := r.score(h)
							//fmt.Println("I see that", lownames[rnum], "gives", b)
							if b < bestlow {
								bestlow = b
								if b == 0 {
									//fmt.Println("Okay, no point fudging!")
									return badness
								}
							}
						}
						for _,r := range hirules {
							if b := r.score(h); b == 0 {
								//fmt.Println("A bid of", hinames[rnum],"would be great.")
								return badness
							}
						}
						if badness < bestlow {
							//fmt.Println("I'm going with", ms[2],"which seems my best option.")
							return 0
						}
						//fmt.Println("I'm returning a modified value of", badness - bestlow)
						return badness - bestlow
					}
					return
				}
			}
		}
		return scorer.score
	}, nil,
}
