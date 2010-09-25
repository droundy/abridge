package main

import (
	"fmt"
	"os"
	"github.com/droundy/bridge"
)

func TestBids(handstring string, good, bad []string) int {
	retval := 0
	var h bridge.Hand
	fmt.Sscan(handstring, &h)
	fmt.Print("\nTesting bids with hand:\n", h)
	for _,g := range good {
		sc,expl,con := bridge.RateBid(h, g, bridge.DefaultConvention)
		fmt.Println("Good bid", g,"has score", sc,"using convention", con)
		if sc != 0 {
			fmt.Println("FAIL: But", g, "is an ok bid, which was rejected because:")
			fmt.Println(expl)
			retval++
		}
	}
	for _,b := range bad {
		sc,_,con := bridge.RateBid(h, b, bridge.DefaultConvention)
		fmt.Println("Bad bid", b,"has score", sc,"using convention", con)
		if sc == 0 {
			fmt.Println("FAIL: But", b, "is a bad bid!")
			retval++
		}
	}
	return retval
}

func main() {
	exitval := 0

	var h bridge.Hand
	fmt.Scanln(&h)
	fmt.Println("Hand is:")
	fmt.Println(h)
	fmt.Println("This hand has", h.HCP(),"high card points")
	fmt.Println("This hand has", h.DistPoints(),"distributional points")
	fmt.Println("This hand has", h.PointCount(),"realistic points")
	fmt.Printf("Card %d is: %s", 3, h.Nth(3))

	fmt.Println("Watch me shuffle:")
	fmt.Print("Deal is now:\n", bridge.Shuffle())

	exitval += TestBids("Axxx xxx Kxxx xx", []string{"1C P1S", "1C P1S P3S P"}, []string{"1C P2S", "1C P1S P3S"})

	exitval += TestBids("AQxx Axx Kxxx xx", []string{"1C P1S", "1C P1S P2S P4S"}, []string{"1C P2S", "1C P1S P2S P3S", "1C P1S P2S P5S", "1C P1S P2S P2N", "1C P1S P2S P6S", "1C P1S P2S P3N", "1C P1S P2S P P"})

	exitval += TestBids("KJxx xxxxx xx xx", []string{"1C P1H", "1C P1H P1S P2S"},
		[]string{"1C P2S", "1C P1H P1S P3S",
		         "1C P1H P1S P1N", "1C P1H P1S P2N"})

	exitval += TestBids("AKQx Axx Kxxx KQ", []string{"1C P1S P4S"}, []string{"1C P2S", "1C P1S P3S", "1C P1S P2S", "1C P1S P2N", "1C P1S P6S", "1C P1S P3N"})

	exitval += TestBids("x aq qjxx qxxxxx", []string{"1C P1S P2C"}, []string{"1C P1S P1N","1C P1S P2S"})

	exitval += TestBids("xxx aq qjxx axxx", []string{"1C P1S P1N"}, []string{"1C P1S P2C","1C P1S P2S"})

	exitval += TestBids("Ax KQxx Jxx Qxxx", []string{"1C P1S P1N"}, []string{"1C P1S P2C","1C P1S P2S"})

	exitval += TestBids("Qxxx qxxxx qx qj", []string{"1N P2D", "1N P2D P2H P2N"}, []string{"1N P2H", "1N P2D P2H P3N", "1N P2D P2H P3H", "1N P2D P2H P4H"})

	ts := bridge.GetValidTables(bridge.South, "1C P1S P2S P", 100, bridge.DefaultConvention)
	fmt.Println(ts)

	/*
	fmt.Print("Valid table for a 1S opener:\n",
		bridge.ShuffleValidTable(bridge.South, " P P P1S"))
	fmt.Print("Valid table for a 1S ... 1N opener:\n",
		bridge.ShuffleValidTable(bridge.South, " P P1S P1N"))
	fmt.Print("Valid table for a 1H ... 2C opener:\n",
		bridge.ShuffleValidTable(bridge.South, " P P1H P2C"))
	fmt.Print("Valid table for a 1H ... 2H opener:\n",
		bridge.ShuffleValidTable(bridge.South, " P P P1H P2H"))
	fmt.Print("Valid table for a 1S ... 3S opener:\n",
		bridge.ShuffleValidTable(bridge.South, " P1S P3S"))
	 */

	switch exitval {
	case 0: fmt.Println("Passed tests.")
	case 1: fmt.Println("Failed with one error!")
	default: fmt.Println("Failed with", exitval, "errors!")
	}
	os.Exit(exitval)
}
