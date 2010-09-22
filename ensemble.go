package bridge

import (
	"fmt"
)

type Ensemble struct {
	tables []Table
	hcp [4]*PointRange
	pts [4]*PointRange
	ptsandsuit [4][4]*PointRange
	suits [4][4]*Range
	Conventions []string // the names of all our conventions
	scorers, unforced map[string]*ScoringRule
}

// RotateToSouth generates an entirely new Ensemble rather than
// modifying its receiver in-place.  But it doesn't do a copy, if
// dealer is South, so you do need to be careful.
func (e *Ensemble) RotateToSouth(dealer Seat) (out *Ensemble) {
	if dealer == South {
		return e
	}
	out = makeEnsemble(len(e.tables));
	for snew:=Seat(South); snew<4; snew++ {
		sold := (snew+dealer) & 3
		out.hcp[snew] = e.hcp[sold]
		out.pts[snew] = e.pts[sold]
		for sv := range out.suits[snew] {
			out.suits[snew][sv] = e.suits[sold][sv]
			out.ptsandsuit[snew][sv] = e.ptsandsuit[sold][sv]
		}
		out.Conventions = e.Conventions
		for i,t := range e.tables {
			out.tables[i][snew] = t[sold]
		}
	}
	return
}

// RotateFromSouth generates an entirely new Ensemble rather than
// modifying its receiver in-place.  But it doesn't do a copy, if
// dealer is South, so you do need to be careful.
func (e *Ensemble) RotateFromSouth(dealer Seat) (out *Ensemble) {
	return e.RotateToSouth((4 - dealer)&3)
}

func (e *Ensemble) ShowHands() string {
	out := ""
	for _,t := range e.tables {
		out += "\n" + t.String()
	}
	return out
}

func (e *Ensemble) String() string {
	out := ""
	N := e.HCP(North)
	S := e.HCP(South)
	E := e.HCP(East)
	W := e.HCP(West)
	Np := e.PointCount(North)
	Sp := e.PointCount(South)
	Ep := e.PointCount(East)
	Wp := e.PointCount(West)
	/*
	tables := [4]map[Hand]int{make(map[Hand]int),make(map[Hand]int),make(map[Hand]int),make(map[Hand]int)}
	points := [4]map[Points]int{make(map[Points]int),make(map[Points]int),make(map[Points]int),make(map[Points]int)}
	hcp := [4]map[Points]int{make(map[Points]int),make(map[Points]int),make(map[Points]int),make(map[Points]int)}
	for _,t := range e.tables {
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
*/
	out += fmt.Sprintf("           [%2d-%2d]\n", Np.Min, Np.Max)
	out += fmt.Sprintf("           (%2d-%2d)\n", N.Min, N.Max)
	for sv:=uint(Spades); sv>Diamonds; sv-- {
		out += fmt.Sprintf("          %s: %v\n", SuitLetter[sv], e.SuitLength(North, sv))
	}
	out += fmt.Sprintf(" [%2d-%2d]  D: %9v[%2d-%2d]\n",Wp.Min,Wp.Max,e.SuitLength(North, Diamonds),Ep.Min,Ep.Max)
	out += fmt.Sprintf(" (%2d-%2d)  C: %9v(%2d-%2d)\n",W.Min,W.Max,e.SuitLength(North, Clubs),E.Min,E.Max)
	for sv:=uint(Spades); sv>Diamonds; sv-- {
		out += fmt.Sprintf("%s: %18v%s: %v\n",
			SuitLetter[sv], e.SuitLength(West, sv),
			SuitLetter[sv], e.SuitLength(East, sv))
	}
	out += fmt.Sprintf("D: %8v[%2d-%2d]   D: %v\n",e.SuitLength(West, Diamonds),Sp.Min,Sp.Max,e.SuitLength(East,Diamonds))
	out += fmt.Sprintf("C: %8v(%2d-%2d)   C: %v\n",e.SuitLength(West, Clubs),S.Min,S.Max,e.SuitLength(East,Clubs))
	for sv:=uint(Spades); sv<=Spades; sv-- {
		out += fmt.Sprintf("          %s: %v\n", SuitLetter[sv], e.SuitLength(South,sv))
	}
	//out += fmt.Sprintf("Spades north: %g\n", e.SuitLength(North, Spades).Mean)
	//out += fmt.Sprintf("Spades south: %g\n", e.SuitLength(South, Spades).Mean)
	out += "\n\n" + e.tables[0].String()
	return out
}

func (e *Ensemble) HtmlSeat(seat Seat) string {
	out := `<div class="analysis">`
	out += "<table>\n"
	out += fmt.Sprintf("<tr><td>&nbsp;&nbsp;&nbsp;<em>%d-%d&nbsp;Points</em></td></tr>\n", e.PointCount(seat).Min, e.PointCount(seat).Max)
	out += fmt.Sprintf("<tr><td>&nbsp;&nbsp;&nbsp;<em>%d-%d&nbsp;HCP</em></td></tr>\n", e.HCP(seat).Min, e.HCP(seat).Max)
	for sv := uint(Spades); sv <= Spades; sv-- {
		out += `<tr><td><div class="bridgecards">` + SuitColorHTML[sv] + "&nbsp;" + e.SuitLength(seat,sv).HTML()
		out += `</div></td></tr>`
	}
	out += "</table></div>\n"
	return out
}

func (e *Ensemble) HTML() string {
	out := fmt.Sprintln(`<div class="analysis">`)
  out += fmt.Sprintf(`<table><tr>`)
  out += fmt.Sprintf(`<td>%s</td>`, e.HtmlSeat(West))
  out += fmt.Sprintf("<td><table><tr><td>%s</td></tr>\n", e.HtmlSeat(North))
  out += fmt.Sprintf("<tr><td>&nbsp;</td></tr>\n")
  out += fmt.Sprintf("<tr><td>&nbsp;</td></tr>\n")
  out += fmt.Sprintf("<tr><td>%s</td></tr></table></td>\n", e.HtmlSeat(South))
  out += fmt.Sprintf(`<td>%s</td>`, e.HtmlSeat(East))
	out += "</tr></table></div>\n"
	return out
}

func (e *Ensemble) ExampleHTML() string {
	return e.tables[0].HTML("Example hand")
}

func (e *Ensemble) Invalidate() {
	for i := range e.hcp {
		e.hcp[i] = nil
		e.pts[i] = nil
		for j := range e.suits[i] {
			e.suits[i][j] = nil
		}
	}
}

func (e *Ensemble) HCP(seat Seat) (r PointRange) {
	if e.hcp[seat] != nil {
		return *e.hcp[seat]
	}
	r.Min = 100
	for _,t := range e.tables {
		hcp := t[seat].HCP()
		if hcp < r.Min {
			r.Min = hcp
		}
		if hcp > r.Max {
			r.Max = hcp
		}
		r.Mean += float64(hcp)
	}
	r.Mean /= float64(len(e.tables))
	e.hcp[seat] = &r
	return r
}

func (e *Ensemble) PointCount(seat Seat) (r PointRange) {
	if e.pts[seat] != nil {
		return *e.pts[seat]
	}
	r.Min = 100
	for _,t := range e.tables {
		pts := t[seat].PointCount()
		if pts < r.Min {
			r.Min = pts
		}
		if pts > r.Max {
			r.Max = pts
		}
		r.Mean += float64(pts)
	}
	r.Mean /= float64(len(e.tables))
	e.pts[seat] = &r
	return r
}

type Range struct {
	Min, Max byte
	Mean float64
}
type PointRange struct {
	Min, Max Points
	Mean float64
}
func (r Range) Format(f fmt.State, c int) {
	i := byte(0)
	for ; i<r.Min; i++ {
		f.Write([]byte{'X'})
	}
	for ; float64(i)+0.5 < r.Mean; i++ {
		f.Write([]byte{'x'})
	}
	for ; i < r.Max; i++ {
		f.Write([]byte{'.'})
	}
	if w,ok := f.Width(); ok {
		for ; i < byte(w); i++ {
			f.Write([]byte{' '})
		}
	}
}
func (r Range) HTML() string {
	out := ""
	i := byte(0)
	for ; i<r.Min; i++ {
		out += "<strong>x</strong>"
	}
	for ; float64(i)+0.5 < r.Mean; i++ {
		out += "x"
	}
	for ; i < r.Max; i++ {
		out += "<em>x</em>"
	}
	return out
}

func (e *Ensemble) SuitLength(seat Seat, suit uint) (r Range) {
	suit = suit % 4
	if e.suits[seat][suit] != nil {
		return *e.suits[seat][suit]
	}
	r.Min = byte(100)
	for _,t := range e.tables {
		num := byte((t[seat] >> (4+8*suit)) & 15)
		if num < r.Min {
			r.Min = num
		}
		if num > r.Max {
			r.Max = num
		}
		r.Mean += float64(num)
	}
	r.Mean /= float64(len(e.tables))
	e.suits[seat][suit] = &r
	return
}

func (e *Ensemble) PointsAndSuit(seat Seat, suit uint) (r PointRange) {
	suit = suit % 4 // ...just to be careful
	if e.ptsandsuit[seat][suit] != nil {
		return *e.ptsandsuit[seat][suit]
	}
	r.Min = Points(100)
	for _,t := range e.tables {
		suitlen := Points((t[seat] >> (4+8*suit)) & 15)
		num := suitlen + t[seat].PointCount()
		switch suitlen {
		case 0: num -= 3
		case 1: num -= 2
		case 2: num -= 1
		}
		if num < r.Min {
			r.Min = num
		}
		if num > r.Max {
			r.Max = num
		}
		r.Mean += float64(num)
	}
	r.Mean /= float64(len(e.tables))
	e.ptsandsuit[seat][suit] = &r
	return
}

func (e *Ensemble) SuitHCP(seat Seat, suits [4]bool) (r PointRange) {
	r.Min = 100
	for _,t := range e.tables {
		hcp := Points(0)
		for sv:=uint(0); sv<4; sv++ {
			if suits[sv] {
				hcp += HCP[byte(t[seat] >> (8*sv))]
			}
		}
		if hcp < r.Min {
			r.Min = hcp
		}
		if hcp > r.Max {
			r.Max = hcp
		}
		r.Mean += float64(hcp)
	}
	r.Mean /= float64(len(e.tables))
	return
}

func makeEnsemble(num int) *Ensemble {
	var foo Ensemble
	foo.tables = make([]Table, num)
	foo.scorers = make(map[string]*ScoringRule)
	foo.unforced = make(map[string]*ScoringRule)
	return &foo
}
