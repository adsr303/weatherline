package main

import (
	"fmt"

	"github.com/adsr303/weatherline/ipapi"
	"github.com/adsr303/weatherline/openmeteo"
	"github.com/alecthomas/kong"
)

var CLI struct {
	Here struct {
	} `cmd:"" help:"Get weather at current location" default:"1"`
	At struct {
		Latitude  float64 `arg:"" help:"Latitude"`
		Longitude float64 `arg:"" help:"Longitude"`
	} `cmd:"" help:"Get weather at specified coordinates"`
}

func main() {
	ctx := kong.Parse(&CLI)
	fmt.Printf("%+v %+v\n", CLI, ctx.Command())
	var geo ipapi.Geolocation
	switch ctx.Command() {
	case "here":
		var err error // prevent shadowing of geo
		geo, err = ipapi.GetGeolocation()
		if err != nil {
			panic(err)
		}
	case "at <latitude> <longitude>":
		geo.Lat, geo.Lon = CLI.At.Latitude, CLI.At.Longitude
	}
	fmt.Printf("%+v\n", geo)
	r, err := openmeteo.GetCurrentWeather(geo.Lat, geo.Lon)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", r)
}
