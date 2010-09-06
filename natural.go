package bridge

import (
	"regexp"
)

var Natural = BiddingRule{
	"Natural",
	regexp.MustCompile("(.)(.)$"),
	func (bidder Seat, h Hand, ms []string, e Ensemble) (badness Score, nothandled bool) {
		pts := h.PointCount()
		hcp := h.HCP()
		partner := (bidder+2)&3
		minPhcp,_,_ := e.HCP(partner)
		minhcp := hcp + minPhcp
		minPpoints,_,_ := e.PointCount(partner)
		minpts := pts + minPpoints
		rspades := e.SuitLength(partner, Spades)
		rhearts := e.SuitLength(partner, Hearts)
		minS := rspades.min + byte((h >> 28)&15)
		minH := rhearts.min + byte((h >> 20)&15)
		switch ms[2] {
		case "N":
			if minS > 7 {
				badness += Score(minS - 7)*SuitLengthProblem
			}
			if minH > 7 {
				badness += Score(minH - 7)*SuitLengthProblem
			}
			pointlevels := map[int]Points{ 2:22, 3:26, 4:33, 5:33, 6:33, 7:37, 8:60 }
			num := int(ms[1][0] - '0')
			if minhcp < pointlevels[num] {
				badness += Score(pointlevels[num]-minhcp)*PointValueProblem
			} else if minhcp >= pointlevels[num+1] {
				badness += Score(minhcp - pointlevels[num+1])*PointValueProblem
			}
		case "S","H","D","C":
			mysuit := stringToSuitNumber(ms[2])
			myownsuitlen := byte((h >> (4+8*mysuit))&15)
			rsuit := e.SuitLength(partner, mysuit)
			mysuitlen := myownsuitlen + rsuit.min
			if mysuitlen < 8 {
				// We always want a guaranteed fit.
				badness += Score(8 - mysuitlen)*SuitLengthProblem
			}
			pointlevels := map[int]Points{ 2:19, 3:23, 4:26, 5:29, 6:33, 7:37, 8:60 }
			num := int(ms[1][0] - '0')
			if minpts < pointlevels[num] {
				badness += Score(pointlevels[num]-minpts)*PointValueProblem
			} else if minpts >= pointlevels[num+1] {
				badness += Score(minpts - pointlevels[num+1])*PointValueProblem
			}
		}
		if pts < 13 {
			badness += Fudge
		}
		if pts < 12 {
			badness += Score(12-pts)*PointValueProblem
		}
		return
	},
}
