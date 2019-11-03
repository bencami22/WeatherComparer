package weathercomparer

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"weathercomparer"
)

func TestWeatherRequest(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{\"LocalObservationDateTime\":\"2019-10-31T21:51:00+01:00\",\"EpochTime\":1572555060,\"WeatherText\":\"Some clouds\",\"WeatherIcon\":36,\"HasPrecipitation\":false,\"PrecipitationType\":null,\"IsDayTime\":false,\"Temperature\":{\"Metric\":{\"Value\":16.8,\"Unit\":\"C\",\"UnitType\":17},\"Imperial\":{\"Value\":62.0,\"Unit\":\"F\",\"UnitType\":18}},\"MobileLink\":\"http://m.accuweather.com/en/it/rome/213490/current-weather/213490?lang=en-us\",\"Link\":\"http://www.accuweather.com/en/it/rome/213490/current-weather/213490?lang=en-us\"}")
	}))
	defer ts.Close()

	w := httptest.NewRecorder()
	configs := weathercomparer.WebAPIConfigs{
		BaseURL: ts.URL,
		APIKey:  "NA",
	}

	requestor := &weathercomparer.AccuWeather{Configuration: configs}
	print(requestor)
	response := weathercomparer.ProviderRequestor.WeatherRequest(requestor, "IT", "ROME")

	if response.DegreeCelsius != 1 {
		t.Fatalf("Expected Status OK but got %v", w.Body.String())
	}

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 but got %v", w.Code)
	}
}
