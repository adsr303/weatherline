package cli

type Options struct {
	Units     string `help:"Units: metric or imperial" short:"u" default:"local" enum:"local,metric,imperial"`
	TempUnits string `help:"Temperature units: celsius or fahrenheit" short:"t" default:"local" enum:"local,celsius,fahrenheit"`
	FeelsLike bool   `help:"Show 'feels like' temperature" short:"l" negatable:"" default:"false"`
	UVIndex   bool   `help:"Show UV index" short:"i" negatable:"" default:"true"`
	Wind      bool   `help:"Show wind speed and direction" short:"w" negatable:"" default:"true"`
	Humidity  bool   `help:"Show humidity" short:"m" negatable:"" default:"true"`
	Pressure  bool   `help:"Show atmospheric pressure" short:"p" negatable:"" default:"true"`
	Daylight  bool   `help:"Show daylight status" short:"d" negatable:"" default:"false"`
}

type CLI struct {
	Options
	Here struct {
	} `cmd:"" help:"Get weather at current location" default:"1"`
	At struct {
		Latitude  float64 `arg:"" help:"Latitude"`
		Longitude float64 `arg:"" help:"Longitude"`
	} `cmd:"" help:"Get weather at specified coordinates"`
}

const HereCommand = "here"
const AtCommand = "at <latitude> <longitude>"
