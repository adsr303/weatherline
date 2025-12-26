package cli

type Options struct {
	Units     string `help:"Units: metric or imperial" short:"u" default:"metric" enum:"metric,imperial"`
	FeelsLike bool   `help:"Show 'feels like' temperature" short:"l" default:"false"`
	UVIndex   bool   `help:"Show UV index" short:"i" default:"false"`
	Humidity  bool   `help:"Show humidity" short:"m" default:"false"`
	Wind      bool   `help:"Show wind speed and direction" short:"w" default:"false"`
	Pressure  bool   `help:"Show atmospheric pressure" short:"p" default:"false"`
	Daylight  bool   `help:"Show daylight status" short:"d" default:"false"`
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
