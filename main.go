package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/adsr303/weatherline/cli"
	"github.com/adsr303/weatherline/geography"
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
	// TODO Handle no city case
	parts := []string{fmt.Sprintf("Weather in %s: %d%s", geo.City, toWholeDegrees(r.CurrentWeather.Temperature), r.Units.Temperature)}
	if cliArgs.FeelsLike {
		parts = append(parts, fmt.Sprintf("Feels like %d%s", toWholeDegrees(r.CurrentWeather.FeelsLike), r.Units.FeelsLike))
	}
	if cliArgs.Wind {
		parts = append(parts, fmt.Sprintf("Wind: %.f %s %s", r.CurrentWeather.WindSpeed, r.Units.WindSpeed, r.CurrentWeather.CompassWindDirection()))
	}
	if cliArgs.Humidity {
		parts = append(parts, fmt.Sprintf("Humidity: %.f%s", r.CurrentWeather.Humidity, r.Units.Humidity))
	}
	if cliArgs.Pressure {
		parts = append(parts, fmt.Sprintf("Pressure: %.f %s", r.CurrentWeather.Pressure, r.Units.Pressure))
	}
	if cliArgs.UVIndex {
		parts = append(parts, fmt.Sprintf("Max UVI: %.1f", r.Daily.UVIndexMax[0]))
	}
	if cliArgs.Daylight {
		parts = append(parts, fmt.Sprintf("Sunrise: %s", toLocalHour(r.Daily.Sunrise[0], geo.CountryCode)))
		parts = append(parts, fmt.Sprintf("Sunset: %s", toLocalHour(r.Daily.Sunset[0], geo.CountryCode)))
	}
	fmt.Println(strings.Join(parts, " - "))
}

// toWholeDegrees converts a float temperature to an integer by rounding to the nearest whole degree.
// This prevents displaying negative zero (e.g., -0.4 -> 0).
func toWholeDegrees(temp float64) int {
	return int(math.Round(temp))
}

// toLocalHour converts a date-time string in the format "2006-01-02T15:04" to a time string
// in either 24-hour format or AM/PM format based on the country's conventions.
func toLocalHour(dateTime, countryCode string) string {
	t, err := time.Parse("2006-01-02T15:04", dateTime)
	if err != nil {
		return dateTime
	}
	if geography.UsesAMPM(countryCode) {
		return t.Format(time.Kitchen)
	}
	return t.Format("15:04")
}
