package aurora

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"

	"github.com/dutchcoders/goftp"
)

const (
	NoData = -999.9

	BzGreen  = 20
	BzYellow = 0
	BzOrange = -10
	BzRed    = -15

	SpeedGreen  = 200
	SpeedYellow = 350
	SpeedOrange = 500
	SpeedRed    = 700

	DensityGreen  = 0
	DensityYellow = 4
	DensityOrange = 10
	DensityRed    = 14

	KpGreen  = 0
	KpYellow = 3
	KpOrange = 4
	KpRed    = 5

	KpGreenWeight  = -5
	KpYellowWeight = 2
	KpOrangeWeight = 5
	KpRedWeight    = 10

	GreenWeight  = -10
	YellowWeight = 10
	OrangeWeight = 15
	RedWeight    = 25

	NoDataWeight = -50
)

var MissingData = -999.9

// Get will connect to the Space Weather site and download the latest 1 minute reports from the ACE spacecraft. It will then parse all the data into a multi-dimensional map
// and then return the results. The first key (an integer) is the line number in the document (or t minus. 0 is the latest, 1 is 1 minute before that etc.)
func Get() map[int]map[string]float64 {
	// We make our map now, because otherwise we'll get nil panics
	var results map[int]map[string]float64 = make(map[int]map[string]float64)

	// getFile is a func that connects via FTP, grabs the data, passes it to extractData, then returns the data back here
	aceMag := getFile("/pub/lists/ace/ace_mag_1m.txt")

	// Now we loop through the results. No range, as we need integers
	for i := 0; i <= len(aceMag)-1; i++ {
		// We "made" our map[int] above, now we do the same here, for the map of our map
		results[i] = make(map[string]float64)
		// Grab the data from getFile, then pull out the data. We also need to convert to float64
		results[i]["YR"], _ = strconv.ParseFloat(aceMag[i][0], 64)
		results[i]["DA"], _ = strconv.ParseFloat(aceMag[i][1], 64)
		results[i]["MO"], _ = strconv.ParseFloat(aceMag[i][2], 64)
		results[i]["Time"], _ = strconv.ParseFloat(aceMag[i][3], 64)
		results[i]["JulianDay"], _ = strconv.ParseFloat(aceMag[i][4], 64)
		results[i]["Seconds"], _ = strconv.ParseFloat(aceMag[i][5], 64)
		results[i]["S"], _ = strconv.ParseFloat(aceMag[i][6], 64)
		results[i]["Bx"], _ = strconv.ParseFloat(aceMag[i][7], 64)
		results[i]["By"], _ = strconv.ParseFloat(aceMag[i][8], 64)
		results[i]["Bz"], _ = strconv.ParseFloat(aceMag[i][9], 64)
		results[i]["Bt"], _ = strconv.ParseFloat(aceMag[i][10], 64)
		results[i]["Lat"], _ = strconv.ParseFloat(aceMag[i][11], 64)
		results[i]["Long"], _ = strconv.ParseFloat(aceMag[i][12], 64)

	}

	// Retrieve the latest ACE data from this folder

	aceSwepam := getFile("/pub/lists/ace/ace_swepam_1m.txt")
	for i := 0; i <= len(aceSwepam)-1; i++ {

		results[i]["Density"], _ = strconv.ParseFloat(aceSwepam[0][7], 64)
		results[i]["Speed"], _ = strconv.ParseFloat(aceSwepam[0][8], 64)
		results[i]["IonTemperature"], _ = strconv.ParseFloat(aceSwepam[i][9], 64)

	}

	return results
}

// GetKp lives in it's own separate func, because data is given in 15 minute intervals, not 1 minute intervals.
// We could migrate this data into the other results, but that'd make finding the latest Kp data trickier.
func GetKp() map[int]map[string]float64 {
	// We make our map now, because otherwise we'll get nil panics
	var results map[int]map[string]float64 = make(map[int]map[string]float64)

	wingKp := getFile("/pub/lists/wingkp/wingkp_list.txt")

	for i := 0; i <= len(wingKp)-1; i++ {
		results[i] = make(map[string]float64)
		results[i]["Kp1Hour"], _ = strconv.ParseFloat(wingKp[i][9], 64)
		results[i]["Kp4Hour"], _ = strconv.ParseFloat(wingKp[i][15], 64)
		results[i]["Kp"], _ = strconv.ParseFloat(wingKp[i][17], 64)

	}
	return results
}

func getFile(file string) map[int][]string {
	ftp, _ := goftp.Connect("ftp.swpc.noaa.gov:21")
	// Log in as an anonymous user
	ftp.Login("anonymous", "")
	var data map[int][]string
	_, _ = ftp.Retr(file, func(r io.Reader) error {
		data = extractData(r)
		return nil
	})
	ftp.Close()
	return data
}

func extractData(r io.Reader) map[int][]string {
	var res map[int][]string = make(map[int][]string)
	// Make a new buffer
	buf := new(bytes.Buffer)
	// And read the data from the file
	buf.ReadFrom(r)
	// This regex searches for all lines that start with a number
	reg := regexp.MustCompile("(?m)^[0-9].+$")
	// Put all our matches into a variable
	data := reg.FindAllString(buf.String(), -1)
	sort.Sort(sort.Reverse(sort.StringSlice(data)))
	// This regex searches for consecutive spaces, as our data is space-separated (no pun intended)
	reg = regexp.MustCompile("\\s+")
	for i := 0; i <= len(data)-1; i++ {

		// Split our returned data into a string array
		info := reg.Split(data[i], -1)
		res[i] = info
	}
	return res
}

// Check compares alert and warning thresholds defined above, and returns a value between 0 and 2, based on what was matched.
// 0 = Density and Speed are less than their associated Warn values, Bz is greater than the associated warn levels
// 1 = Density and Speed are greater than the warn values, but less than the alert values, Bz is less than the warn level, but greater than the alert level
// 2 = Density and Speed are greater than the alert values, Bz is less than the alert value
func Check(results map[int]map[string]float64, kpResults map[int]map[string]float64, i int) (score, bz, speed, density, kp int) {

	if results[i]["Bz"] <= NoData {
		fmt.Println("No data for Bz")
		score += NoDataWeight
	}
	if results[i]["Speed"] <= NoData {
		fmt.Println("No data for Speed")
		score += NoDataWeight
	}
	if results[i]["Density"] <= NoData {
		fmt.Println("No data for Density")
		score += NoDataWeight
	}
	if results[i]["Kp"] <= NoData {
		fmt.Println("No data for Kp")
		score += NoDataWeight
	}

	if results[i]["Bz"] <= BzGreen && results[i]["Bz"] >= BzYellow {
		bz = 0
		fmt.Println("Bz Green")
		score += GreenWeight
	} else if results[i]["Bz"] <= BzYellow && results[i]["Bz"] >= BzOrange {
		bz = 1
		fmt.Println("Bz Yellow")
		score += YellowWeight
	} else if results[i]["Bz"] <= BzOrange && results[i]["Bz"] >= BzRed {
		bz = 2
		fmt.Println("Bz Orange")
		score += OrangeWeight
	} else if results[i]["Bz"] <= BzRed {
		bz = 3
		fmt.Println("Bz Red")
		score += RedWeight
	} else {
		fmt.Println("No Bz data")
		bz = -1
	}

	if results[i]["Speed"] >= SpeedGreen && results[i]["Speed"] <= SpeedYellow {
		speed = 0
		fmt.Println("Speed Green")
		score += GreenWeight
	} else if results[i]["Speed"] >= SpeedYellow && results[i]["Speed"] <= SpeedOrange {
		speed = 1
		fmt.Println("Speed Yellow")
		score += YellowWeight
	} else if results[i]["Speed"] >= SpeedOrange && results[i]["Speed"] <= SpeedRed {
		speed = 2
		fmt.Println("Speed Orange")
		score += OrangeWeight
	} else if results[i]["Speed"] >= SpeedRed {
		speed = 3
		fmt.Println("Speed Red")
		score += RedWeight
	} else {
		fmt.Println("No Speed data")
		speed = -1
	}

	if results[i]["Density"] >= DensityGreen && results[i]["Density"] <= DensityYellow {
		density = 0
		fmt.Println("Density Green")
		score += GreenWeight
	} else if results[i]["Density"] >= DensityYellow && results[i]["Density"] <= DensityOrange {
		density = 1
		fmt.Println("Density Yellow")
		score += YellowWeight
	} else if results[i]["Density"] >= DensityOrange && results[i]["Density"] <= DensityRed {
		density = 2
		fmt.Println("Density Orange")
		score += OrangeWeight
	} else if results[i]["Density"] >= DensityRed {
		density = 3
		fmt.Println("Density Red")
		score += RedWeight
	} else {
		fmt.Println("No Density data")
		density = -1
	}

	if kpResults[i]["Kp"] >= KpGreen && kpResults[i]["Kp"] <= KpYellow {
		kp = 0
		fmt.Println("Kp Green")
		score += KpGreenWeight
	} else if kpResults[i]["Kp"] >= KpYellow && kpResults[i]["Kp"] <= KpOrange {
		kp = 1
		fmt.Println("Kp Yellow")
		score += KpYellowWeight
	} else if kpResults[i]["Kp"] >= KpOrange && kpResults[i]["Kp"] <= KpRed {
		kp = 2
		fmt.Println("Kp Orange")
		score += KpOrangeWeight
	} else if kpResults[i]["Kp"] >= KpRed {
		kp = 3
		fmt.Println("Kp Red")
		score += KpRedWeight
	} else {
		fmt.Println("No Kp data")
		kp = -1
	}

	return score, bz, speed, density, kp
}
