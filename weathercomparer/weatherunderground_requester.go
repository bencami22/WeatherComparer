package weathercomparer

type WeatherUnderground struct {}

func (provider *WeatherUnderground) WeatherRequest(country string) WeatherResponse{
	return WeatherResponse{20}
}