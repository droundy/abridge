package bridge

import (
	"strings"
	"regexp"
	"fmt"
)

func PickBid(h Hand, bidder Seat, oldbid string, cc ConventionCard, e *Ensemble) (bid string, convention string) {
	bid = " P"
	pbids := PossibleNondoubleBids(oldbid)
	for _,b := range pbids {
		r := makeScoringRule(bidder, oldbid + b, cc, e)
		if sc,_ := r.score(h); sc == 0 {
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
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (s Score,e string)) {
		if !strings.HasSuffix(e.Conventions[len(e.Conventions)-2], "(forcing)") {
			return nil
		}
		return func(h Hand) (Score,string) {
			return 1000, ""
		}
	},
	nil,
}

var LimitPass = BiddingRule {
	"Limiting pass",
	regexp.MustCompile("^(..)?(..)?(..)?(.[^P]......)* P$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (s Score, e string)) {
		if len(ms[0]) > 8 && ms[0][len(ms[0])-5] == 'P' {
			// We have already bid, and our partner passed, so we can always
			// pass now.
			return nil
		}
		possbids := PossibleNondoubleBids(ms[0])
		allrules := make([]*ScoringRule,0,len(possbids))
		allbids := make([]string,0,len(possbids))
		for _,b := range possbids {
			if rule := makeScoringRule(bidder, ms[0] + b, cc, e); rule != nil {
				allrules = allrules[0:len(allrules)+1]
				allrules[len(allrules)-1] = rule
				allbids = allbids[0:len(allbids)+1]
				allbids[len(allbids)-1] = b
			}
		}
		return func(h Hand) (badness Score, explanation string) {
			pts := h.PointCount()
			if pts < 6 {
				// We can always safely pass with less than 6 points!
				return
			}
			for i,r := range allrules {
				sc,_ := r.score(h)
				if sc == 0 {
					// We can't be sure of the problem, but probably we need
					// fewer points, so let's bias the search that direction!
					badness += 137 + Score(pts-5)*PointValueProblem
					explanation += allbids[i] + " is a better bid\n"
					return
				}
			}
			return
		}
	},
	nil,
}

var suitlevels = map[int]Points{ 1:13, 2:19, 3:23, 4:26, 5:29, 6:33, 7:37, 8:60 }
var ntlevels = map[int]Points{ 1:15, 2:20, 3:26, 4:33, 5:33, 6:33, 7:37, 8:60 }

func ScoreHands(h1, h2 Hand, bidlevel int, suitvalue uint) (badness Score) {
	hcp1 := h1.HCP()
	hcp2 := h2.HCP()
	hcp := hcp1 + hcp2
	if suitvalue == NoTrump {
		if hcp < ntlevels[bidlevel] {
			// One point short is a fudge...
			badness += Score(ntlevels[bidlevel] - hcp)*Fudge
		}
		if hcp < ntlevels[bidlevel]-1 {
			// More than one point short is pretty bad
			badness += Score(ntlevels[bidlevel] - 1 - hcp)*PointValueProblem
		}
		return
	}
	n1 := byte(h1 >> (4+8*suitvalue)) & 15
	n2 := byte(h2 >> (4+8*suitvalue)) & 15
	n := n1 + n2
	pts1 := h1.PointCount()
	pts2 := h2.PointCount()
	if n < 8 {
		badness += Score(8 - n)*SuitLengthProblem
	}
	// Grant an extra point for each trump beyond eight...
	adjusted_pts := Points(n1 + n2) + pts1 + pts2 - 8
	switch n1 {
	case 0: adjusted_pts -= 4 // a trump void counts negative!
	case 1: adjusted_pts -= 2
	case 2: adjusted_pts -= 1
	}
	switch n2 {
	case 0: adjusted_pts -= 4 // a trump void counts negative!
	case 1: adjusted_pts -= 2
	case 2: adjusted_pts -= 1
	}
	if adjusted_pts < suitlevels[bidlevel] {
		// One point short is a fudge...
		badness += Score(suitlevels[bidlevel] - adjusted_pts)*Fudge
	}
	if adjusted_pts < suitlevels[bidlevel] - 1 {
		// More than one point short is pretty bad
		badness += Score(suitlevels[bidlevel] - 1 - adjusted_pts)*PointValueProblem
	}
	return
}

func WorstCaseScenario(bidder Seat, h Hand, bidlevel int, suitvalue uint, cc ConventionCard, e *Ensemble) (badness Score, explanation string) {
	partner := (bidder+2)&3
	worsthand := Hand(0)
	for _,t := range e.tables {
		h2bad := Score(1000000)
		for sv := uint(Clubs); sv <= NoTrump; sv++ {
			thislevel := bidlevel
			if sv <= suitvalue {
				thislevel++
			}
			if b := ScoreHands(h,t[partner],thislevel,sv); b < h2bad {
				h2bad = b
			}
		}
		if badness < h2bad {
			// This is the worst hand we've seen yet!
			badness = h2bad
			worsthand = t[partner]
		}
	}
	return badness, worsthand.String()
}

func WorstCaseSuit(bidder Seat, h Hand, bidlevel int, suitvalue uint, cc ConventionCard, e *Ensemble) (badness Score, explanation string) {
	partner := (bidder+2)&3
	worsthand := Hand(0)
	for _,t := range e.tables {
		h2bad := Score(1000000)
		for sv := uint(Clubs); sv < NoTrump; sv++ {
			thislevel := bidlevel
			if sv <= suitvalue {
				thislevel++
			}
			if b := ScoreHands(h,t[partner],thislevel,sv); b < h2bad {
				h2bad = b
			}
		}
		if badness < h2bad {
			// This is the worst hand we've seen yet!
			badness = h2bad
			worsthand = t[partner]
		}
	}
	return badness, worsthand.String()
}

func SafeContractInThisSuit(bidder Seat, h Hand, suitvalue uint, e *Ensemble) (bidlevel int) {
	partner := (bidder+2)&3
	bidlevel = 7
	for _,t := range e.tables {
		b := Score(100)
		for b > 0 && bidlevel > 0 {
			b = ScoreHands(h, t[partner], bidlevel, suitvalue)
			if b > 0 {
				bidlevel--
			}
		}
		if bidlevel <= 0 {
			return
		}
	}
	return
}

var TakeOutDouble = BiddingRule {
	"Takeout double (forcing)",
	regexp.MustCompile("^.*([123])([CDHSN])( P P)? X$"),	
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		suitvalue := stringToSuitNumber(ms[2])
		bidlevel := int(ms[1][0] - '0') // the level they are bid to
		return func(h Hand) (badness Score, explanation string) {
			badness, explanation = WorstCaseSuit(bidder, h, bidlevel, suitvalue, cc, e)
			if badness > 0 {
				// We are willing to fudge by one point, since WorstCaseSuit
				// is pretty pessimistic.  Better might be to use a median or
				// something.
				if badness < PointValueProblem {
					badness = 0
				} else {
					badness -= PointValueProblem
				}
			}
			if bidlevel == 1 && suitvalue < Spades {
				lsp := byte(h >> 28) & 15
				if lsp > 4 {
					badness += Score(lsp-4)*SuitLengthProblem
				} else if lsp < 4 {
					badness += Score(4-lsp)*SuitLengthProblem
				}
			}
			if bidlevel == 1 && suitvalue < Hearts {
				lh := byte(h >> 20) & 15
				if lh > 4 {
					badness += Score(lh-4)*SuitLengthProblem
				} else if lh < 4 {
					badness += Score(4-lh)*SuitLengthProblem
				}
			}
			return 
		}
	}, nil,
}

var NewSuitForcing = BiddingRule {
	"New suit (forcing)",
	regexp.MustCompile("^.+(.)([^PX])( P)*(.)([CDHS])$"),	
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		suitvalue := stringToSuitNumber(ms[5])
		bidlevel := int(ms[4][0] - '0') // the level we are bid to
		oldlevel := int(ms[1][0] - '0') // the previous level bid
		if bidlevel > oldlevel + 1 {
			// We've obviously jumped!
			return nil
		}
		if (ms[2] != "N" && stringToSuitNumber(ms[2]) < suitvalue) && bidlevel > oldlevel {
			// We've jumped, which doesn't count as "new suit forcing".
			return nil
		}
		allpassed := true
		for i:=1; 4*i < len(ms[0]); i++ {
			this := ms[0][len(ms[0])-1-4*i]
			if this == ms[5][0] {
				// This suit is not new!
				return nil
			}
			if this != 'P' {
				allpassed = false
			}
		}
		if allpassed {
			// We haven't done anything but pass yet (so this bid is not
			// forcing)
			return nil
		}
		// At this point we know it's a "new suit forcing" bid.
		return func(h Hand) (badness Score, explanation string) {
			badness, explanation = WorstCaseScenario(bidder, h, bidlevel, suitvalue, cc, e)
			suitlen := byte(h >> (4+8*suitvalue)) & 15
			if suitlen < 4 {
				badness += Score(4-suitlen)*SuitLengthProblem
			}
			return
		}
	}, nil,
}

var Natural = BiddingRule{
	"Natural",
	regexp.MustCompile("(.)([CDHSN])$"),	
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		partner := (bidder+2)&3
		// gamelevel is the bid needed for game.
		mysuit := stringToSuitNumber(ms[2])
		gamelevel := 5 - int(mysuit/2) // fun formula!  :)
		// pointlevels are the points that are needed for various bids
		pointlevels := suitlevels
		num := int(ms[1][0] - '0') // the level we are bid to

		gotsplinter := false
		var theirsplinterrange PointRange
		var splintersuit uint = 0
		var splinterpts = Points(25)
		switch ms[2] {
		case "N":
			// redefine pointlevel appropriately
			pointlevels = ntlevels
		case "S","H","D","C":
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

		score = func(h Hand) (badness Score, explanation string) {
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
			switch ms[2] {
			case "N":
				if minS > 7 {
					if explain {
						explanation = "There is a spade fit!\n"
					}
					badness += Score(minS - 7)*SuitLengthProblem
				}
				if minH > 7 {
					if explain {
						explanation += "There is a hearts fit!\n"
					}
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
					if explain {
						explanation += "We don't have a fit!\n"
					}
					badness += Score(8 - mysuitlen)*SuitLengthProblem
				}
				// We get an extra point per trump beyond 8:
				partnersuitpts := e.PointsAndSuit(partner, mysuit)
				minpts = pts + (partnersuitpts.Min + Points(myownsuitlen)) - 8
				maxpts = pts + (partnersuitpts.Max + Points(myownsuitlen)) - 8
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
									return 0, ""
									//badness += Score(neededpts - theirhcprange.Min - myhcp)*PointValueProblem;
								}
							}
						}
					}
				}
			}	
			if minpts < pointlevels[num] {
				// we need to guarantee pointlevels[num] to bid at this level
				if explain {
					explanation += fmt.Sprintln("We don't have the", pointlevels[num], "needed!")
					explanation += fmt.Sprintln("We only have", minpts, "together, and I have", )
				}
				badness += Score(pointlevels[num]-minpts)*PointValueProblem
			} else if minpts >= pointlevels[num+1] && (num < gamelevel || num == 5 || num == 6) {
				// if we can guarantee pointlevels[num+1], then we should bid at
				// *that* level (if it's at-or-below game, or a slam bid)
				if explain {
					explanation += fmt.Sprintln("We have more than the", pointlevels[num+1], "needed for the next bid up!")
					explanation += fmt.Sprintln("We have at least", minpts)
				}
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
				if explain {
					explanation += "We have no chance at game!\n"
				}
				badness += Score(pointlevels[gamelevel] - maxpts)*PointValueProblem
			}
			if num > gamelevel && num < 6 && minpts < pointlevels[6]-2 {
				// We don't want to invite slam (i.e. bid over game) unless
				// there's a good reason to suspect that we've got the points
				// for slam.
				if explain {
					explanation += "No point bidding beyond game if we know we can't make slam!\n"
				}
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
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (Score,string)) {
		// FIXME: Forced ought perhaps to use WorstCaseScenario (or
		// something like it) to see which undesirable bids might possibly
		// lead to better contracts and favor those (e.g. to favor a suit
		// bid over a NT bid, even if we don't know we have a fit).
		if len(e.Conventions) < 2 {
			return nil
		}
		if !strings.HasSuffix(e.Conventions[len(e.Conventions)-2], "(forcing)") {
			return nil
		}
		//fmt.Println("Working on Forced...", ms[0])
		scorer := makeUnforcedScoringRule(bidder, ms[0], cc, e)
		allbutme := ms[0][0:len(ms[0])-2]
		pbids := PossibleNondoubleBids(allbutme)
		hirules := make([]*ScoringRule,0,len(pbids))
		lowrules := make([]*ScoringRule,0,5)
		lownames := make([]string,0,5)
		hinames := make([]string,0,len(pbids))
		for i,b := range pbids {
			if b != ms[2] {
				//fmt.Println("Making rule for", allbutme + b)
				if rule := makeUnforcedScoringRule(bidder, allbutme + b, cc, e); rule != nil {
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
				return func(h Hand) (badness Score, explanation string) {
					badness,explanation = scorer.score(h)
					if badness > 0 {
						//fmt.Print("I am running Forced on:\n", h, "for a bid of ", ms[2],"\n")
						//fmt.Println("Native badness is", badness)
						// Need to figure out if this is our only option
						bestlow := Score(1000000)
						for _,r := range lowrules {
							b,_ := r.score(h)
							//fmt.Println("I see that", lownames[rnum], "gives", b)
							if b < bestlow {
								bestlow = b
								if b == 0 {
									//fmt.Println("Okay, no point fudging!")
									return
								}
							}
						}
						for _,r := range hirules {
							if b,_ := r.score(h); b == 0 {
								//fmt.Println("A bid of", hinames[rnum],"would be great.")
								return
							}
						}
						if badness < bestlow {
							//fmt.Println("I'm going with", ms[2],"which seems my best option.")
							explanation += "...but it seems to be the best option.\n"
							return 0, explanation
						}
						//fmt.Println("I'm returning a modified value of", badness - bestlow)
						return badness - bestlow, explanation
					}
					return
				}
			}
		}
		return scorer.score
	}, nil,
}
