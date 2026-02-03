package cli

type GlobalOptions struct {
	Units     string `help:"Units: metric or imperial" short:"u" default:"local" enum:"local,metric,imperial"`
	TempUnits string `help:"Temperature units: celsius or fahrenheit" short:"t" default:"local" enum:"local,celsius,fahrenheit"`
	Symbols   bool   `help:"Use symbols in output" short:"s" negatable:"" default:"false"`
}

type CLI struct {
	GlobalOptions
	Now struct {
		FeelsLike bool `help:"Show 'feels like' temperature" short:"l" negatable:"" default:"true"`
		UVIndex   bool `help:"Show UV index" short:"i" negatable:"" default:"false"`
		Wind      bool `help:"Show wind speed and direction" short:"w" negatable:"" default:"true"`
		Humidity  bool `help:"Show humidity" short:"m" negatable:"" default:"true"`
		Pressure  bool `help:"Show atmospheric pressure" short:"p" negatable:"" default:"true"`
		Daylight  bool `help:"Show daylight status" short:"d" negatable:"" default:"false"`
	} `cmd:"" help:"Get current weather at current location" default:"1"`
	Forecast struct {
		Days int `help:"Number of days to forecast (1-7)" short:"n" default:"5" enum:"1,2,3,4,5,6,7"`
	} `cmd:"" help:"Get weather forecast at current location"`
}

const NowCommand = "now"
const ForecastCommand = "forecast"
