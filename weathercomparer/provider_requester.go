package weathercomparer

type ProviderRequestor interface {
    WeatherRequest(country string) WeatherResponse
}