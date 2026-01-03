package openmeteo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/adsr303/weatherline/cli"
	"github.com/adsr303/weatherline/geography"
	"github.com/adsr303/weatherline/ipapi"
)

const baseURL = "https://api.open-meteo.com/v1/forecast"

var currentWeatherParams = []string{
	"temperature_2m",
	"relative_humidity_2m",
	"apparent_temperature",
	"is_day",
	"wind_speed_10m",
	"wind_direction_10m",
	"wind_gusts_10m",
	"precipitation",
	"showers",
	"snowfall",
	"rain",
	"weather_code",
	"cloud_cover",
	"pressure_msl",
	"surface_pressure",
}

var defaultParams string = strings.Join(currentWeatherParams, ",")

type CurrentWeather struct {
	Temperature     float64 `json:"temperature_2m"`
	Humidity        float64 `json:"relative_humidity_2m"`
	FeelsLike       float64 `json:"apparent_temperature"`
	IsDay           int     `json:"is_day"`
	WindSpeed       float64 `json:"wind_speed_10m"`
	WindDirection   float64 `json:"wind_direction_10m"`
	WindGusts       float64 `json:"wind_gusts_10m"`
	Precipitation   float64 `json:"precipitation"`
	Showers         float64 `json:"showers"`
	Snowfall        float64 `json:"snowfall"`
	Rain            float64 `json:"rain"`
	WeatherCode     int     `json:"weather_code"`
	CloudCover      float64 `json:"cloud_cover"`
	Pressure        float64 `json:"pressure_msl"`
	SurfacePressure float64 `json:"surface_pressure"`
}

var compassDirections = []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}

// CompassWindDirection returns the wind direction as a compass direction string (e.g., "N", "NE", "E", etc.).
func (c *CurrentWeather) CompassWindDirection() string {
	index := int((c.WindDirection+11.25)/22.5) % 16
	return compassDirections[index]
}

// IsDaytime returns true if it is currently daytime at the location.
func (c *CurrentWeather) IsDaytime() bool {
	return c.IsDay == 1
}

// Weather represents general weather types.
type Weather int

// Weather type constants.
const (
	Clear Weather = iota
	Clouds
	Rain
	Fog
	Mist
	Haze
	Snow
	Thunderstorm
)

// WeatherType returns the general weather type based on the weather code.
func (c *CurrentWeather) WeatherType() Weather {
	switch c.WeatherCode {
	case 0, 1:
		return Clear
	case 2, 3:
		return Clouds
	case 45, 48:
		return Fog
	case 51, 53, 55, 56, 57, 61, 63, 65, 66, 67, 80, 81, 82:
		return Rain
	case 71, 73, 75, 77, 85, 86:
		return Snow
	case 95, 96, 99:
		return Thunderstorm
	default:
		return Clear
	}
}

type CurrentWeatherUnits struct {
	Temperature     string `json:"temperature_2m"`
	Humidity        string `json:"relative_humidity_2m"`
	FeelsLike       string `json:"apparent_temperature"`
	IsDay           string `json:"is_day"`
	WindSpeed       string `json:"wind_speed_10m"`
	WindDirection   string `json:"wind_direction_10m"`
	WindGusts       string `json:"wind_gusts_10m"`
	Precipitation   string `json:"precipitation"`
	Showers         string `json:"showers"`
	Snowfall        string `json:"snowfall"`
	Rain            string `json:"rain"`
	WeatherCode     string `json:"weather_code"`
	CloudCover      string `json:"cloud_cover"`
	Pressure        string `json:"pressure_msl"`
	SurfacePressure string `json:"surface_pressure"`
}

type DailyWeather struct {
	Sunrise    []string  `json:"sunrise"`
	Sunset     []string  `json:"sunset"`
	UVIndexMax []float64 `json:"uv_index_max"`
}

type WeatherResponse struct {
	Latitude       float64             `json:"latitude"`
	Longitude      float64             `json:"longitude"`
	GenerationTime float64             `json:"generationtime_ms"`
	UTCOffset      int                 `json:"utc_offset_seconds"`
	Timezone       string              `json:"timezone"`
	TimezoneAbbr   string              `json:"timezone_abbreviation"`
	Elevation      float64             `json:"elevation"`
	Current        CurrentWeather      `json:"current"`
	Units          CurrentWeatherUnits `json:"current_units"`
	Daily          DailyWeather        `json:"daily"`
}

type ErrorResponse struct {
	Reason string `json:"reason"`
}

type WeatherError struct {
	Reason string
}

func (e *WeatherError) Error() string {
	return e.Reason
}

func GetCurrentWeather(geo *ipapi.Geolocation, options *cli.Options) (WeatherResponse, error) {
	// TODO Elevation, timezone
	units := fmt.Sprintf(
		"temperature_unit=%s&wind_speed_unit=%s&precipitation_unit=%s",
		getTemperatureUnit(options, geo.CountryCode),
		getWindSpeedUnit(options, geo.CountryCode),
		getPrecipitationUnit(options, geo.CountryCode))
	requestUrl := fmt.Sprintf(
		"%s?latitude=%f&longitude=%f&timezone=%s&current=%s&daily=sunrise,sunset,uv_index_max&%s",
		baseURL, geo.Lat, geo.Lon, geo.Timezone, defaultParams, units)
	resp, err := http.Get(requestUrl)
	if err != nil {
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		var weatherErr ErrorResponse
		if err := json.Unmarshal(b, &weatherErr); err != nil {
			return WeatherResponse{}, err
		}
		return WeatherResponse{}, &WeatherError{Reason: weatherErr.Reason}
	}
	var weatherResp WeatherResponse
	if err := json.Unmarshal(b, &weatherResp); err != nil {
		return WeatherResponse{}, err
	}
	return weatherResp, nil
}

func getTemperatureUnit(options *cli.Options, countryCode string) string {
	switch options.TempUnits {
	case "local":
		if geography.UsesFahrenheit(countryCode) {
			return "fahrenheit"
		}
		return "celsius"
	default:
		return options.TempUnits
	}
}

func getWindSpeedUnit(options *cli.Options, countryCode string) string {
	switch options.Units {
	case "metric":
		return "kmh"
	case "imperial":
		return "mph"
	default:
		if geography.UsesImperial(countryCode) {
			return "mph"
		}
		return "kmh"
	}
}

func getPrecipitationUnit(options *cli.Options, countryCode string) string {
	switch options.Units {
	case "metric":
		return "mm"
	case "imperial":
		return "inch"
	default:
		if geography.UsesImperial(countryCode) {
			return "inch"
		}
		return "mm"
	}
}
