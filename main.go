package main

import (
	"fmt"
	"math"
	"time"

	"github.com/adsr303/weatherline/cli"
	"github.com/adsr303/weatherline/geography"
	"github.com/adsr303/weatherline/ipapi"
	"github.com/adsr303/weatherline/openmeteo"
	"github.com/alecthomas/kong"
	"github.com/fatih/color"
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
	formatOutput(cliArgs, geo, r)
}

func formatOutput(cliArgs cli.CLI, geo ipapi.Geolocation, weather openmeteo.WeatherResponse) {
	// TODO Handle no city case
	color.Set(color.BgBlue, color.Bold)
	fmt.Print(" ")
	printEntry("Weather in "+geo.City, fmt.Sprintf("%d%s", toWholeDegrees(weather.Current.Temperature), weather.Units.Temperature), false)
	if cliArgs.FeelsLike {
		printEntry("Feels like", fmt.Sprintf("%d%s", toWholeDegrees(weather.Current.FeelsLike), weather.Units.FeelsLike), true)
	}
	if cliArgs.Wind {
		printEntry("Wind", fmt.Sprintf("%.f %s %s", weather.Current.WindSpeed, weather.Units.WindSpeed, weather.Current.CompassWindDirection()), true)
	}
	if cliArgs.Humidity {
		printEntry("Humidity", fmt.Sprintf("%.f%s", weather.Current.Humidity, weather.Units.Humidity), true)
	}
	if cliArgs.Pressure {
		printEntry("Pressure", fmt.Sprintf("%.f %s", weather.Current.Pressure, weather.Units.Pressure), true)
	}
	if cliArgs.UVIndex {
		printEntry("Max UVI", fmt.Sprintf("%.1f", weather.Daily.UVIndexMax[0]), true)
	}
	if cliArgs.Daylight {
		printEntry("Sunrise", toLocalHour(weather.Daily.Sunrise[0], geo.CountryCode), true)
		printEntry("Sunset", toLocalHour(weather.Daily.Sunset[0], geo.CountryCode), true)
	}
	fmt.Print(" ")
	color.Unset()
	fmt.Println()
}

func printEntry(label, value string, dash bool) {
	if dash {
		color.Set(color.FgBlue)
		fmt.Print(" - ")
	}
	color.Set(color.FgHiCyan)
	fmt.Print(label)
	color.Set(color.FgMagenta)
	fmt.Print(": ")
	color.Set(color.FgYellow)
	fmt.Print(value)
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
