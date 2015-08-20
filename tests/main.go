package main

import (
	"fmt"

	"github.com/Grayda/go-aurora"
)

func main() {
	fmt.Println("Querying Spaceweather site for data..")
	results := aurora.Get()
	for k, v := range results {
		fmt.Println(k, "\t\t\t\t", v)
	}

	fmt.Println()
	switch aurora.Check(results) {
	case 2:
		fmt.Println("Gauges are red! Grab your camera!")
	case 1:
		fmt.Println("Gauges are in the orange. Prepare your camera, and watch the skies")
	case 0:
		fmt.Println("No significant activity on the gauges")
	case -1:
		fmt.Println("No usable data returned from ACE")
	}

}
