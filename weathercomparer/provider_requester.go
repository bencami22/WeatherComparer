package weathercomparer

//ProviderRequestor is the interface for all Weather Provider Requester to implement
type ProviderRequestor interface {
	WeatherRequest(country string) WeatherResponse
}
