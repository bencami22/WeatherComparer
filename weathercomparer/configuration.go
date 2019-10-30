package weathercomparer

//Configuration is to load configuration values from file
type Configuration struct {
	WeatherBitConfiguration  WebAPIConfigs
	OpenWeatherConfiguration WebAPIConfigs
	AccuWeatherConfiguration WebAPIConfigs
}

//WebAPIConfigs contains configs used to communicate with third parties over HTTP
type WebAPIConfigs struct {
	BaseURL string
	APIKey  string
}
