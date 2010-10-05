package main

import (
	"fmt"
	"os"
	"github.com/droundy/bridge"
)

func main() {
	sec, nsec, _ := os.Time()

	ts := bridge.GetValidTables(0, " P P1C P1H P3H P4H P P P", 100, bridge.DefaultConventions())
	fmt.Println("Conventions are:", ts.Conventions)
	fmt.Println(ts)

	ts = bridge.GetValidTables(0, " P 1C1H P2H P3H P4H P P P", 100, bridge.DefaultConventions())
	fmt.Println("Conventions are:", ts.Conventions)
	fmt.Println(ts)

	sec2, nsec2, _ := os.Time()
	name, _ := os.Hostname()
	dt := float64(sec2-sec) + 1e-9*float64(nsec2-nsec)
	timelimits := map[string]float64{"collins": 1.7, "morland": 8, "bennet": 4}
	fmt.Println("It took ", dt,"seconds on", name)
	expected, ok := timelimits[name]
	if ok {
		if dt > expected {
			fmt.Println("Exceeded expected time by", dt-expected,"seconds.")
			os.Exit(1)
		}
		fmt.Println("Expected", expected,"seconds on", name)
	} else {
		fmt.Println("You should add statistics for", name)
		os.Exit(1)
	}
}
