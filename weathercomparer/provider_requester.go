package weathercomparer

import "context"

//ProviderRequestor is the interface for all Weather Provider Requester to implement
type ProviderRequestor interface {
	WeatherRequest(ctx context.Context, country string, city string) (WeatherResponse, error)
}
