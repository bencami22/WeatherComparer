package weathercomparer

type OpenWeather struct {}

func (provider *OpenWeather) WeatherRequest(country string) WeatherResponse{
	return WeatherResponse{10}
}