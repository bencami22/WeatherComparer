package weathercomparer

//WeatherResponse represent the common data response from ProviderRequestors
type WeatherResponse struct {
	Provider      string
	DegreeCelsius float64
}
