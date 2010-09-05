package main

import (
	"fmt"
	"github.com/droundy/bridge"
)

func main() {
	fmt.Println(bridge.Suit(bridge.Ace + bridge.Queen + 5 << 4))
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
}
