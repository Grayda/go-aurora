package aurora

import (
	"bytes"
	"io"
	"regexp"
	"strconv"

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
func Get() map[string]float64 {
	var results map[string]float64 = make(map[string]float64)
	// Create a new FTP object and connect to the SWPC site

	// Retrieve the latest ACE data from this folder

	aceMag := getFile("/pub/lists/ace/ace_mag_1m.txt")
	results["Long"], _ = strconv.ParseFloat(aceMag[len(aceMag)-2], 64)
	results["Lat"], _ = strconv.ParseFloat(aceMag[len(aceMag)-3], 64)
	results["Bt"], _ = strconv.ParseFloat(aceMag[len(aceMag)-4], 64)
	results["Bz"], _ = strconv.ParseFloat(aceMag[len(aceMag)-5], 64)
	results["By"], _ = strconv.ParseFloat(aceMag[len(aceMag)-6], 64)
	results["Bx"], _ = strconv.ParseFloat(aceMag[len(aceMag)-7], 64)
	results["S"], _ = strconv.ParseFloat(aceMag[len(aceMag)-8], 64)
	results["Seconds"], _ = strconv.ParseFloat(aceMag[len(aceMag)-9], 64)
	results["JulianDay"], _ = strconv.ParseFloat(aceMag[len(aceMag)-10], 64)
	results["Time"], _ = strconv.ParseFloat(aceMag[len(aceMag)-11], 64)
	results["DA"], _ = strconv.ParseFloat(aceMag[len(aceMag)-12], 64)
	results["MO"], _ = strconv.ParseFloat(aceMag[len(aceMag)-13], 64)
	results["YR"], _ = strconv.ParseFloat(aceMag[len(aceMag)-14], 64)

	// Retrieve the latest ACE data from this folder

	aceSwepam := getFile("/pub/lists/ace/ace_swepam_1m.txt")
	results["IonTemperature"], _ = strconv.ParseFloat(aceSwepam[len(aceSwepam)-2], 64)
	results["Speed"], _ = strconv.ParseFloat(aceSwepam[len(aceSwepam)-3], 64)
	results["Density"], _ = strconv.ParseFloat(aceSwepam[len(aceSwepam)-4], 64)

	wingKp := getFile("/pub/lists/wingkp/wingkp_list.txt")
	results["Kp"], _ = strconv.ParseFloat(wingKp[len(wingKp)-2], 64)
	results["Kp1hour"], _ = strconv.ParseFloat(wingKp[len(wingKp)-4], 64)
	results["Kp4hours"], _ = strconv.ParseFloat(wingKp[len(wingKp)-10], 64)

	return results
}

func getFile(file string) []string {
	ftp, _ := goftp.Connect("ftp.swpc.noaa.gov:21")
	// Log in as an anonymous user
	ftp.Login("anonymous", "")
	var data []string
	_, _ = ftp.Retr(file, func(r io.Reader) error {
		data = extractData(r)
		return nil
	})
	return data
}

func extractData(r io.Reader) []string {
	//
	// Make a new buffer
	buf := new(bytes.Buffer)
	// And read the data from the file
	buf.ReadFrom(r)
	// This regex searches for consecutive spaces, as our data is space-separated (no pun intended)
	reg := regexp.MustCompile("\\s+")
	// Split our returned data into a string array
	return reg.Split(buf.String(), -1)
}

// Check compares alert and warning thresholds defined above, and returns a value between 0 and 2, based on what was matched.
// 0 = Density and Speed are less than their associated Warn values, Bz is greater than the associated warn levels
// 1 = Density and Speed are greater than the warn values, but less than the alert values, Bz is less than the warn level, but greater than the alert level
// 2 = Density and Speed are greater than the alert values, Bz is less than the alert value
func Check(results map[string]float64) int {

	switch {
	case (results["Bz"] < BzAlert && results["Bz"] != MissingData) && (results["Density"] > DensityAlert && results["Density"] != MissingData) && (results["Speed"] > SpeedAlert && results["Speed"] != MissingData):
		return 2 // Grab your camera! It's aurora photographing time!
	case (results["Bz"] < BzWarn && results["Bz"] != MissingData) && (results["Density"] > DensityWarn && results["Density"] != MissingData) && (results["Speed"] > SpeedWarn && results["Speed"] != MissingData):
		return 1 // Prepare your camera, as the gauges are looking good..
	case results["Bz"] <= MissingData && results["Density"] <= MissingData && results["Speed"] <= MissingData:
		return -1 // Sensor is offline
	default:
		return 0 // Business as usual
	}
}
