// Package ipapi implements querying the IP geolocation service at ip-api.com.
package ipapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Geolocation struct {
	Status     string // "success" or "failure"
	Message    string // Non-empty if Status=="failure"
	Country    string
	RegionName string
	City       string
	Lat        float64
	Lon        float64
	Timezone   string
	Offset     int16 // Timezone UTC offset in seconds
}

func GetGeolocation() (Geolocation, error) {
	var geo Geolocation
	res, err := http.Get("http://ip-api.com/json/?fields=status,message,country,regionName,city,lat,lon,timezone,offset,isp,org,as,query")
	if err != nil {
		return geo, fmt.Errorf("retrieving geolocation: %w", err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close() //nolint:errcheck
	if res.StatusCode >= 300 {
		return geo, fmt.Errorf("retrieving geolocation: %d (%s)", res.StatusCode, res.Status)
	}
	if err != nil {
		return geo, fmt.Errorf("retrieving geolocation: %w", err)
	}
	err = json.Unmarshal(body, &geo)
	if err != nil {
		return geo, fmt.Errorf("parsing geolocation: %w", err)
	}
	return geo, nil
}
