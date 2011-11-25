/*

	slookup - Lookup weather stations

	Prints area weather stations (from NOAA via wunderground.com)
	to the console.

	Written and maintained by Stephen Ramsay

	Last Modified: Wed Aug 03 14:02:45 CDT 2011

	Copyright © 2011 by Stephen Ramsay

	slookup is free software; you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation; either version 3, or (at your option) any
	later version.
	
	slookup is distributed in the hope that it will be useful, but
	WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
	or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License
	for more details.
	
	You should have received a copy of the GNU General Public License
	along with slookup; see the file COPYING.  If not see
	<http://www.gnu.org/licenses/>

*/
package main

import (
	"fmt"
	"http"
	"os"
	"strings"
	"xml"
	"github.com/jteeuwen/go-pkg-optarg"
)

const URLstem = "http://api.wunderground.com/auto/wui/geo/GeoLookupXML/index.xml?query="

const VERS = "1.2.1"

type Station struct {
	City		string
	Icao		string
}

type Airport struct {
	Station []Station "nearby_weather_stations>airport>station"
}

func main() {

	optarg.Add("s", "location", "Location.  May be indicated using city, state, CITY,STATE, country, (US or Canadian) zipcode, 3- or 4-letter airport code, or LAT,LONG", "Lincoln, NE")
	optarg.Add("h", "help", "Print this message", false)
	optarg.Add("V", "version", "Print version number", false)

	var location = "Lincoln,NE"
	var help, version bool
	var URL string

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "s":
			location = opt.String()
		case "h":
			help = opt.Bool()
		case "V":
			version = opt.Bool()
		}
	}

	if help {
		optarg.Usage()
		os.Exit(0)
	}

	if version {
		fmt.Println("conditions " + VERS)
		fmt.Println("Copyright (C) 2011 by Stephen Ramsay")
		fmt.Println("Data courtesy of Weather Underground, Inc.")
		fmt.Println("is subject to Weather Underground Data Feed")
		fmt.Println("Terms of Service.  The program itself is free")
		fmt.Println("software, and you are welcome to redistribute")
		fmt.Println("it under certain conditions.  See LICENSE for")
		fmt.Println("details.")
		os.Exit(0)
	}

	// Temporarily trim whitespace locations with spaces
	// (e.g. "New York, NY" -> "NewYork,NY")
	var location_components = strings.Fields(location)
	var location_id = ""
	for i := 0; i < len(location_components); i++ {
		location_id = location_id + location_components[i]
	}

	URL = URLstem + location_id

	res, err := http.Get(URL)

	if err == nil {
		var airport Airport
		xmlErr := xml.Unmarshal(res.Body, &airport)
		checkError(xmlErr)
		printWeather(&airport)
		res.Body.Close()
	}
}

func printWeather(airport *Airport) {
	for i := 0; i < len(airport.Station); i++ {
		fmt.Println(airport.Station[i].City + ": " + airport.Station[i].Icao)
	}
}

func checkError(err os.Error) {
	if err != nil {
		fmt.Println(os.Stderr, "Fatal error ", err.String())
		os.Exit(1)
	}
}
