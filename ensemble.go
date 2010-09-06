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
	/*
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
*/
	out += fmt.Sprintf("           (%2d-%2d)\n", Nmin, Nmax)
	for sv:=uint(Spades); sv>Clubs; sv-- {
		out += fmt.Sprintf("          %s: %v\n", SuitLetter[sv], e.SuitLength(North, sv))
	}
	out += fmt.Sprintf(" (%2d-%2d)  C: %9v(%2d-%2d)\n",Wmin,Wmax,e.SuitLength(North, Clubs),Emin,Emax)
	for sv:=uint(Spades); sv>Clubs; sv-- {
		out += fmt.Sprintf("%s: %18v%s: %v\n",
			SuitLetter[sv], e.SuitLength(West, sv),
			SuitLetter[sv], e.SuitLength(East, sv))
	}
	out += fmt.Sprintf("C: %8v(%2d-%2d)   C: %v\n",e.SuitLength(West, Clubs),Smin,Smax,e.SuitLength(East,Clubs))
	for sv:=uint(Spades); sv<=Spades; sv-- {
		out += fmt.Sprintf("          %s: %v\n", SuitLetter[sv], e.SuitLength(South,sv))
	}
	//out += fmt.Sprintf("Spades north: %g\n", e.SuitLength(North, Spades).mean)
	//out += fmt.Sprintf("Spades south: %g\n", e.SuitLength(South, Spades).mean)
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

type Range struct {
	min, max byte
	mean float64
}
func (r Range) Format(f fmt.State, c int) {
	i := byte(0)
	for ; i<r.min; i++ {
		f.Write([]byte{'X'})
	}
	for ; float64(i)+0.5 < r.mean; i++ {
		f.Write([]byte{'x'})
	}
	for ; i < r.max; i++ {
		f.Write([]byte{'.'})
	}
	if w,ok := f.Width(); ok {
		for ; i < byte(w); i++ {
			f.Write([]byte{' '})
		}
	}
}

func (e Ensemble) SuitLength(seat Seat, suit uint) (r Range) {
	suit = suit % 4
	r.min = byte(100)
	for _,t := range e {
		num := byte((t[seat] >> (4+8*suit)) & 15)
		if num < r.min {
			r.min = num
		}
		if num > r.max {
			r.max = num
		}
		r.mean += float64(num)
	}
	r.mean /= float64(len(e))
	return
}
