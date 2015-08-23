package main

import (
	"fmt"

	"github.com/Grayda/go-aurora"
)

func main() {
	fmt.Println("Querying Spaceweather site for data..")
	results := aurora.Get()
	fmt.Println("Latest results are:")
	for k, v := range results[0] {
		fmt.Println(k, "-", v)
	}
	fmt.Println()
	fmt.Println("Grabbing Kp data..")
	kpresults := aurora.GetKp()

	fmt.Println("Current Kp index is:")
	for k, v := range kpresults[0] {
		fmt.Println(k, "-", v)
	}
	fmt.Println()
	fmt.Println("Checking these values against their thresholds..")

	switch aurora.Check(results, 0) {
	case 2:
		fmt.Println("Gauges are red! Grab your camera!")
	case 1:
		fmt.Println("Gauges are in the orange. Prepare your camera, and watch the skies")
	case 0:
		fmt.Println("No significant activity on the gauges")
	case -1:
		fmt.Println("One or more of Density, Speed or Bz could not be determined")
	}

}
