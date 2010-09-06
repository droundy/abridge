package bridge

import (
	"fmt"
)

type Ensemble []Table

func (e Ensemble) String() string {
	out := ""
	Nmin, _, Nmax := e.HCP(North)
	Smin, _, Smax := e.HCP(South)
	Emin, _, Emax := e.HCP(East)
	Wmin, _, Wmax := e.HCP(West)
	tables := [4]map[Hand]int{make(map[Hand]int),make(map[Hand]int),make(map[Hand]int),make(map[Hand]int)}
	points := [4]map[Points]int{make(map[Points]int),make(map[Points]int),make(map[Points]int),make(map[Points]int)}
	hcp := [4]map[Points]int{make(map[Points]int),make(map[Points]int),make(map[Points]int),make(map[Points]int)}
	for _,t := range e {
		for i := range tables {
			tables[i][t[i]]++
			points[i][t[i].PointCount()]++
			hcp[i][t[i].HCP()]++
		}
	}
	hands := e[0]
	prob := 0
	for _,t := range e {
		probi := 0
		for i:=range tables {
			probi += tables[i][t[i]]
			probi += points[i][t[i].PointCount()]
			probi += hcp[i][t[i].HCP()]
		}
		if probi > prob {
			hands = t
			prob = probi
		}
	}
	fmt.Printf("prob is %d out of %d\n", prob-12, 12*(len(e)-1))
	out += fmt.Sprintf("        (%2d-%2d)\n", Nmin, Nmax)
	for sv:=uint(Spades); sv>Clubs; sv-- {
		out += fmt.Sprintf("        %s: %v\n", SuitLetter[sv], Suit(hands[North] >> (8*sv)))
	}
	out += fmt.Sprintf("(%2d-%2d) C: %6v(%2d-%2d)\n",Wmin,Wmax,Suit(hands[North]),Emin,Emax)
	for sv:=uint(Spades); sv>Clubs; sv-- {
		out += fmt.Sprintf("%s: %13v%s: %v\n", SuitLetter[sv], Suit(hands[West] >> (8*sv)), SuitLetter[sv], Suit(hands[East] >> (8*sv)))
	}
	out += fmt.Sprintf("C: %5v(%2d-%2d) C: %v\n",Suit(hands[West]),Smin,Smax,Suit(hands[East]))
	for sv:=uint(Spades); sv<=Spades; sv-- {
		out += fmt.Sprintf("        %s: %v\n", SuitLetter[sv], Suit(hands[South] >> (8*sv)))
	}
	return out
}

func (e Ensemble) HCP(seat Seat) (min Points, mean float64, max Points) {
	min = 100
	for _,t := range e {
		hcp := t[seat].HCP()
		if hcp < min {
			min = hcp
		}
		if hcp > max {
			max = hcp
		}
		mean += float64(hcp)
	}
	return min, mean/float64(len(e)), max
}

func (e Ensemble) PointCount(seat Seat) (min Points, mean float64, max Points) {
	min = 100
	for _,t := range e {
		hcp := t[seat].PointCount()
		if hcp < min {
			min = hcp
		}
		if hcp > max {
			max = hcp
		}
		mean += float64(hcp)
	}
	return min, mean/float64(len(e)), max
}
