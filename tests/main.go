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

	res := aurora.Check(results, kpresults, 0)

	fmt.Println("Aurora Score is:", res["Score"])
	fmt.Println("This is based on the following parameters. -1 = No Data, 0 = Green, 1 = Yellow, 2 = Orange, 3 = Red:")
	fmt.Println("Bz status:", res["Bz"], "Speed status:", res["Speed"], "Density status:", res["Density"], "Kp status:", res["Kp"])

}
