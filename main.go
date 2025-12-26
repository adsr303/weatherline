package main

import (
	"fmt"

	"github.com/adsr303/weatherline/cli"
	"github.com/adsr303/weatherline/ipapi"
	"github.com/adsr303/weatherline/openmeteo"
	"github.com/alecthomas/kong"
)

func main() {
	var cli cli.CLI
	ctx := kong.Parse(&cli)
	var geo ipapi.Geolocation
	switch ctx.Command() {
	case "here":
		var err error // prevent shadowing of geo
		geo, err = ipapi.GetGeolocation()
		if err != nil {
			panic(err)
		}
	case "at <latitude> <longitude>":
		geo.Lat, geo.Lon = cli.At.Latitude, cli.At.Longitude
	}
	r, err := openmeteo.GetCurrentWeather(geo.Lat, geo.Lon)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%.1f%s\n", r.CurrentWeather.Temperature, r.Units.Temperature)
}
