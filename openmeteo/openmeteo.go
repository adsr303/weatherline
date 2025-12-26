package openmeteo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const BaseURL = "https://api.open-meteo.com/v1/forecast"

var CurrentWeatherParams = []string{
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

var DefaultParams string = strings.Join(CurrentWeatherParams, ",")

type CurrentWeather struct {
	Temperature         float64 `json:"temperature_2m"`
	RelativeHumidity    float64 `json:"relative_humidity_2m"`
	ApparentTemperature float64 `json:"apparent_temperature"`
	IsDay               int     `json:"is_day"`
	WindSpeed           float64 `json:"wind_speed_10m"`
	WindDirection       float64 `json:"wind_direction_10m"`
	WindGusts           float64 `json:"wind_gusts_10m"`
	Precipitation       float64 `json:"precipitation"`
	Showers             float64 `json:"showers"`
	Snowfall            float64 `json:"snowfall"`
	Rain                float64 `json:"rain"`
	WeatherCode         int     `json:"weather_code"`
	CloudCover          float64 `json:"cloud_cover"`
	PressureMSL         float64 `json:"pressure_msl"`
	SurfacePressure     float64 `json:"surface_pressure"`
}

type CurrentWeatherUnits struct {
	Temperature         string `json:"temperature_2m"`
	RelativeHumidity    string `json:"relative_humidity_2m"`
	ApparentTemperature string `json:"apparent_temperature"`
	IsDay               string `json:"is_day"`
	WindSpeed           string `json:"wind_speed_10m"`
	WindDirection       string `json:"wind_direction_10m"`
	WindGusts           string `json:"wind_gusts_10m"`
	Precipitation       string `json:"precipitation"`
	Showers             string `json:"showers"`
	Snowfall            string `json:"snowfall"`
	Rain                string `json:"rain"`
	WeatherCode         string `json:"weather_code"`
	CloudCover          string `json:"cloud_cover"`
	PressureMSL         string `json:"pressure_msl"`
	SurfacePressure     string `json:"surface_pressure"`
}

type WeatherResponse struct {
	Latitude       float64             `json:"latitude"`
	Longitude      float64             `json:"longitude"`
	GenerationTime float64             `json:"generationtime_ms"`
	UTCOffset      int                 `json:"utc_offset_seconds"`
	Timezone       string              `json:"timezone"`
	TimezoneAbbr   string              `json:"timezone_abbreviation"`
	Elevation      float64             `json:"elevation"`
	CurrentWeather CurrentWeather      `json:"current_weather"`
	Units          CurrentWeatherUnits `json:"hourly_units"`
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

func GetCurrentWeather(latitude, longitude float64) (WeatherResponse, error) {
	requestUrl := fmt.Sprintf("%s?latitude=%.2f&longitude=%.2f&current=%s",
		BaseURL, latitude, longitude, DefaultParams)
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
