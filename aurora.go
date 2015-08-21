package aurora

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/dutchcoders/goftp"
)

// If retrieved values are greater than (or less than, in the case of Bz) these, Check() will return 1, indicating a Aurora warning
var BzWarn float64 = 0
var SpeedWarn float64 = 350
var DensityWarn float64 = 4

// If retrieved values are greater than (or less than) these values, Check() will return 2, which is a "grab your camera!" alert :)
var BzAlert = -10.0
var SpeedAlert = 500.0
var DensityAlert = 10.0
var MissingData = -999.9

// Get will connect to the Space Weather site and download the latest 1 minute reports from the ACE spacecraft. It will then parse the single MOST RECENT data
// put it into a map of map[string]float64, then return the results.
func Get() []map[string]float64 {
	var results []map[string]float64
	// Create a new FTP object and connect to the SWPC site

	// Retrieve the latest ACE data from this folder

	aceMag := getFile("/pub/lists/ace/ace_mag_1m.txt")
	results[0]["Long"], _ = strconv.ParseFloat(aceMag[0][len(aceMag)-2], 64)
	results[0]["Lat"], _ = strconv.ParseFloat(aceMag[0][len(aceMag)-3], 64)
	results[0]["Bt"], _ = strconv.ParseFloat(aceMag[0][len(aceMag)-4], 64)
	results[0]["Bz"], _ = strconv.ParseFloat(aceMag[0][len(aceMag)-5], 64)
	results[0]["By"], _ = strconv.ParseFloat(aceMag[0][len(aceMag)-6], 64)
	results[0]["Bx"], _ = strconv.ParseFloat(aceMag[0][len(aceMag)-7], 64)
	results[0]["S"], _ = strconv.ParseFloat(aceMag[0][len(aceMag)-8], 64)
	results[0]["Seconds"], _ = strconv.ParseFloat(aceMag[0][len(aceMag)-9], 64)
	results[0]["JulianDay"], _ = strconv.ParseFloat(aceMag[0][len(aceMag)-10], 64)
	results[0]["Time"], _ = strconv.ParseFloat(aceMag[0][len(aceMag)-11], 64)
	results[0]["DA"], _ = strconv.ParseFloat(aceMag[0][len(aceMag)-12], 64)
	results[0]["MO"], _ = strconv.ParseFloat(aceMag[0][len(aceMag)-13], 64)
	results[0]["YR"], _ = strconv.ParseFloat(aceMag[0][len(aceMag)-14], 64)

	// Retrieve the latest ACE data from this folder

	aceSwepam := getFile("/pub/lists/ace/ace_swepam_1m.txt")
	results[0]["IonTemperature"], _ = strconv.ParseFloat(aceSwepam[0][len(aceSwepam)-2], 64)
	results[0]["Speed"], _ = strconv.ParseFloat(aceSwepam[0][len(aceSwepam)-3], 64)
	results[0]["Density"], _ = strconv.ParseFloat(aceSwepam[0][len(aceSwepam)-4], 64)

	wingKp := getFile("/pub/lists/wingkp/wingkp_list.txt")
	results[0]["Kp"], _ = strconv.ParseFloat(wingKp[0][len(wingKp)-2], 64)
	results[0]["Kp1hour"], _ = strconv.ParseFloat(wingKp[0][len(wingKp)-4], 64)
	results[0]["Kp4hours"], _ = strconv.ParseFloat(wingKp[0][len(wingKp)-10], 64)

	return results
}

func getFile(file string) [][]string {
	ftp, _ := goftp.Connect("ftp.swpc.noaa.gov:21")
	// Log in as an anonymous user
	ftp.Login("anonymous", "")
	var data [][]string
	_, _ = ftp.Retr(file, func(r io.Reader) error {
		data = extractData(r)
		return nil
	})
	return data
}

func extractData(r io.Reader) [][]string {

	var res [][]string
	// Make a new buffer
	buf := new(bytes.Buffer)
	// And read the data from the file
	buf.ReadFrom(r)
	// This regex searches for all lines that start with a number
	reg := regexp.MustCompile("[0-9]")
	// Split our returned data into a string array
	data := reg.Split(buf.String(), -1)

	for i := len(data) - 1; i >= 0; i-- {
		fmt.Println("--=-=-=-=", data[i])
		fmt.Println(strings.HasPrefix(data[i], ";") && strings.HasPrefix(data[i], ";"))

		if strings.HasPrefix(data[i], "#") == false && len(data[i]) > 0 && strings.HasPrefix(data[i], ";") == false {
			spew.Dump(data[i])
			// This regex searches for consecutive spaces, as our data is space-separated (no pun intended)
			reg = regexp.MustCompile("\\s+")

			// Split our returned data into a string array
			res[i] = reg.Split(data[i], -1)
		}
	}
	return res
}

// Check compares alert and warning thresholds defined above, and returns a value between 0 and 2, based on what was matched.
// 0 = Density and Speed are less than their associated Warn values, Bz is greater than the associated warn levels
// 1 = Density and Speed are greater than the warn values, but less than the alert values, Bz is less than the warn level, but greater than the alert level
// 2 = Density and Speed are greater than the alert values, Bz is less than the alert value
func Check(results []map[string]float64, i int) int {

	switch {
	case (results[i]["Bz"] < BzAlert && results[i]["Bz"] != MissingData) && (results[i]["Density"] > DensityAlert && results[i]["Density"] != MissingData) && (results[i]["Speed"] > SpeedAlert && results[i]["Speed"] != MissingData):
		return 2 // Grab your camera! It's aurora photographing time!
	case (results[i]["Bz"] < BzWarn && results[i]["Bz"] != MissingData) && (results[i]["Density"] > DensityWarn && results[i]["Density"] != MissingData) && (results[i]["Speed"] > SpeedWarn && results[i]["Speed"] != MissingData):
		return 1 // Prepare your camera, as the gauges are looking good..
	case results[i]["Bz"] <= MissingData && results[i]["Density"] <= MissingData && results[i]["Speed"] <= MissingData:
		return -1 // Sensor is offline
	default:
		return 0 // Business as usual
	}
}
