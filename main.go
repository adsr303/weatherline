package main

import (
	"fmt"
	"strings"

	"github.com/adsr303/weatherline/cli"
	"github.com/adsr303/weatherline/ipapi"
	"github.com/adsr303/weatherline/openmeteo"
	"github.com/alecthomas/kong"
)

func main() {
	var cliArgs cli.CLI
	ctx := kong.Parse(&cliArgs)
	var geo ipapi.Geolocation
	switch ctx.Command() {
	case cli.HereCommand:
		var err error // prevent shadowing of geo
		geo, err = ipapi.GetGeolocation()
		if err != nil {
			panic(err)
		}
	case cli.AtCommand:
		geo.Lat, geo.Lon = cliArgs.At.Latitude, cliArgs.At.Longitude
	}
	r, err := openmeteo.GetCurrentWeather(geo.Lat, geo.Lon, &cliArgs.Options, geo.CountryCode)
	if err != nil {
		panic(err)
	}
	// TODO "in City"
	parts := []string{fmt.Sprintf("Weather: %.0f%s", r.CurrentWeather.Temperature, r.Units.Temperature)}
	if cliArgs.FeelsLike {
		parts = append(parts, fmt.Sprintf("(Feels like %.0f%s)", r.CurrentWeather.FeelsLike, r.Units.FeelsLike))
	}
	if cliArgs.Wind {
		// TODO direction as text
		parts = append(parts, fmt.Sprintf("Wind: %.0f %s %f", r.CurrentWeather.WindSpeed, r.Units.WindSpeed, r.CurrentWeather.WindDirection))
	}
	if cliArgs.Humidity {
		parts = append(parts, fmt.Sprintf("Humidity: %.0f%s", r.CurrentWeather.Humidity, r.Units.Humidity))
	}
	if cliArgs.Pressure {
		parts = append(parts, fmt.Sprintf("Pressure: %.0f %s", r.CurrentWeather.Pressure, r.Units.Pressure))
	}
	if cliArgs.UVIndex {
		parts = append(parts, fmt.Sprintf("Max UVI: %.1f", r.Daily.UVIndexMax[0]))
	}
	if cliArgs.Daylight {
		// TODO format
		parts = append(parts, fmt.Sprintf("Sunrise: %s", r.Daily.Sunrise[0]))
		parts = append(parts, fmt.Sprintf("Sunset: %s", r.Daily.Sunset[0]))
	}
	fmt.Println(strings.Join(parts, " - "))
}
