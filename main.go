package main

import (
	"fmt"

	"github.com/adsr303/weatherline/openmeteo"
)

func main() {
	r, err := openmeteo.GetCurrentWeather(52.52, 13.41)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", r)
}
