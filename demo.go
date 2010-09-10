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
		sc,con := bridge.RateBid(h, g)
		fmt.Println("Good bid", g,"has score", sc,"using convention", con)
		if sc != 0 {
			fmt.Println("FAIL: But it's an ok bid!")
			retval++
		}
	}
	for _,b := range bad {
		sc,con := bridge.RateBid(h, b)
		fmt.Println("Bad bid", b,"has score", sc,"using convention", con)
		if sc == 0 {
			fmt.Println("FAIL: But it's a bad bid!")
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

	exitval += TestBids("Axxx xxx Kxxx xx", []string{"1C P1S", "1C P1S P2S P"}, []string{"1C P2S", "1C P1S P2S3S"})

	exitval += TestBids("AKQx AQx Kxxx KQ", []string{"1C P1S", "1C P1S P2S4S"}, []string{"1C P2S", "1C P1S P2S3S", "1C P1S P2S P P"})

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
	os.Exit(exitval)
}
