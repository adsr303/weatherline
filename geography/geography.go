// Package geography provides utilities related to geographical data.
package geography

var fahrenheitCountries = map[string]bool{
	"BS": true, // Bahamas
	"BZ": true, // Belize
	"KY": true, // Cayman Islands
	"PW": true, // Palau
	"US": true, // United States
	"FM": true, // Micronesia
	"MH": true, // Marshall Islands
	"LR": true, // Liberia
}

// UsesFahrenheit returns true if the given country code uses Fahrenheit for temperature measurement.
func UsesFahrenheit(countryCode string) bool {
	_, exists := fahrenheitCountries[countryCode]
	return exists
}

var imperialCountries = map[string]bool{
	"LR": true, // Liberia
	"MM": true, // Myanmar
	"US": true, // United States
}

// UsesImperial returns true if the given country code uses the imperial system for measurements.
func UsesImperial(countryCode string) bool {
	_, exists := imperialCountries[countryCode]
	return exists
}

var ampmCountries = map[string]bool{
	"US": true, // United States
	"PH": true, // Philippines
	"CA": true, // Canada (partially)
	"GB": true, // United Kingdom (partially)
}

// UsesAMPM returns true if the given country code commonly uses the AM/PM time format.
func UsesAMPM(countryCode string) bool {
	_, exists := ampmCountries[countryCode]
	return exists
}
