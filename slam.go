package bridge

import (
	"regexp"
	"fmt"
)

var Blackwood = BiddingRule{
	"Blackwood (forcing)",
	regexp.MustCompile("^.*4N$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (s Score, e string)) {
		if !cc.Options["Blackwood"] {
			return nil
		}
		suit := e.FindFit(bidder)
		if suit == NoTrump {
			// We haven't agreed on a fit, so this can't be Blackwood!
			fmt.Println("There is no fit!")
			return nil
		}
		return func(h Hand) (badness Score, explanation string) {
			safe := SafeContractInThisSuit(bidder, h, suit, e)
			if safe < 5 {
				// We really ought to be safe at the five level to bid this!
				badness += Score(5 - safe)*SuitLengthProblem
			}
			return
		}
	}, nil,
}

var BlackwoodResponse = BiddingRule{
	"Blackwood response",
	regexp.MustCompile("^.*4N P5([CDHS])$"),
	func (bidder Seat, ms []string, cc ConventionCard, e *Ensemble) (score func(h Hand) (s Score, e string)) {
		if e.Conventions[len(e.Conventions)-2] != "Blackwood (forcing)" {
			// Not responding to Blackwood!
			return nil
		}
		aces := 0
		switch ms[1] {
		case "D": aces = 1
		case "H": aces = 2
		case "S": aces = 3
		}
		altaces := aces
		if aces == 0 {
			altaces = 4
		}
		return func(h Hand) (badness Score, explanation string) {
			numaces := 0
			if Suit(h) & Ace != 0 {
				numaces++
			}
			if Suit(h>>8) & Ace != 0 {
				numaces++
			}
			if Suit(h>>16) & Ace != 0 {
				numaces++
			}
			if Suit(h>>24) & Ace != 0 {
				numaces++
			}
			if numaces > aces {
				badness += Score(numaces - aces)*SuitLengthProblem
			}
			if numaces < aces {
				badness += Score(aces - numaces)*SuitLengthProblem
			}
			if numaces == altaces {
				badness = 0
			}
			return
		}
	}, nil,
}