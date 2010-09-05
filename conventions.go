package bridge

import (
	"regexp"
)

type Score float64
const (
	SuitLengthProblem Score = 100
	PointValueProblem Score = 100
	Fudge Score = 1
)

type BiddingRule struct {
	name string
	match *regexp.Regexp
	score func(t Table, seat int, ms []string) Score 
}

var Convention = []BiddingRule{ Opening }

var Opening = BiddingRule{
	"Opening",
	regexp.MustCompile("^( P)*1([CDHS])$"),
	func (t Table, seat int, ms []string) Score {
		pts := t[seat].PointCount()
		badness := Score(0)
		if pts < 13 {
			badness += Fudge
		}
		if pts < 12 {
			badness += Score(12-pts)*PointValueProblem
		}
		ls := byte(t[seat] >> 28)
		lh := byte(t[seat] >> 20) & 15
		ld := byte(t[seat] >> 12) & 15
		lc := byte(t[seat] >> 4) & 15
		switch ms[2] {
		case "S":
			if ls < 5 {
				badness += Score(5-ls)*SuitLengthProblem
			}
			if ls < lh {
				badness += Score(lh-ls)*SuitLengthProblem
			}
		case "H":
			if lh < 5 {
				badness += Score(5-lh)*SuitLengthProblem
			}
			if lh < ls {
				badness += Score(ls-lh)*SuitLengthProblem
			}
		case "D":
			if lh > 4 {
				badness += Score(lh-4)*SuitLengthProblem
			}
			if ls > 4 {
				badness += Score(ls-4)*SuitLengthProblem
			}
			if lc > ld {
				badness += Score(lc-ld)*SuitLengthProblem
			}
		case "C":
			if lh > 4 {
				badness += Score(lh-4)*SuitLengthProblem
			}
			if ls > 4 {
				badness += Score(ls-4)*SuitLengthProblem
			}
			if ld > lc {
				badness += Score(ld-lc)*SuitLengthProblem
			}
		}
		return badness
	},
}

func TableScore(t Table, seat int, bid string) Score {
	for _,c := range Convention {
		ms := c.match.FindStringSubmatch(bid)
		if ms != nil {
			return c.score(t, seat, ms)
		}
	}
	return 0
}

func ShuffleValidTable(seat int, bid string) (t Table) {
	for {
		t = Shuffle()
		if TableScore(t, seat, bid) == 0 {
			return t
		}
	}
	return
}
