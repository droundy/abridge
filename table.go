package bridge

import (
	"fmt"
	"rand"
	"time"
)

type Table [4]Hand

const (
	South = iota
	West
	North
	East
)

type Seat uint
func (s Seat) String() string {
	switch s {
	case South: return "south"
	case North: return "north"
	case East:  return "east "
	case West:  return "west "
	}
	return "bugger seat"
}

func init() {
	rand.Seed(time.Seconds())
}

func (d Table) String() string {
	out := ""
	suitletter := []string{"C", "D", "H", "S"}
	out += fmt.Sprintf("         (%d)\n", d[North].HCP())
	for sv:=uint(Spades); sv>Clubs; sv-- {
		out += fmt.Sprintf("        %s: %v\n", suitletter[sv], Suit(d[North] >> (8*sv)))
	}
	out += fmt.Sprintf(" (%2d)   C: %6v(%2d)\n", d[West].HCP(), Suit(d[North]), d[East].HCP())
	for sv:=uint(Spades); sv>Clubs; sv-- {
		out += fmt.Sprintf("%s: %13v%s: %v\n", suitletter[sv], Suit(d[West] >> (8*sv)), suitletter[sv], Suit(d[East] >> (8*sv)))
	}
	out += fmt.Sprintf("C: %6v(%2d)   C: %v\n", Suit(d[West]), d[South].HCP(), Suit(d[East]))
	for sv:=uint(Spades); sv<=Spades; sv-- {
		out += fmt.Sprintf("        %s: %v\n", suitletter[sv], Suit(d[South] >> (8*sv)))
	}
	return out
}

const FullSuit = 15+(13<<4)
var Sorted = Table { FullSuit, FullSuit << 8, FullSuit << 16, FullSuit << 24 }

func (d Table) ShuffleCard(fromwhich int) Table {
	fromwhich = fromwhich % 52 // Just to be paranoid
	from := d[fromwhich&3].Nth(fromwhich>>2)
	towhich := rand.Intn(52)
	to := d[towhich&3].Nth(towhich>>2)
	d[towhich&3] += from - to
	d[fromwhich&3] += to - from
	return d
}

func Shuffle() Table {
	d := Sorted
	for i:=0; i<52; i++ {
		d = d.ShuffleCard(i)
	}
	return d
}
