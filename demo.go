package main

import (
	"fmt"
	"github.com/droundy/bridge"
)

func main() {
	fmt.Println(bridge.Suit(bridge.Ace + bridge.Queen + 5 << 4))
	var h bridge.Hand
	fmt.Scanln(&h)
	fmt.Println(h)
	fmt.Println("This hand has", h.HCP(),"high card points")
	fmt.Println("This hand has", h.DistPoints(),"distributional points")
	fmt.Println("This hand has", h.PointCount(),"realistic points")
	var s bridge.Suit
	fmt.Scanln(&s)
	fmt.Println(s)
	fmt.Println("This hand has", bridge.HCP[s],"high card points")
	fmt.Println("This hand has", bridge.DistPoints[s],"distributional points")
	fmt.Println("This hand has", bridge.PointCount[s],"realistic points")
}
