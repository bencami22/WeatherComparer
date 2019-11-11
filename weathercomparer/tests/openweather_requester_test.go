package weathercomparer

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	//"weathercomparer"
	"github.com/bencami22/WeatherComparer/weathercomparer"
)

func TestWeatherRequest_OpenWeather_Success200(t *testing.T) {
	// Arrange
	testCountry := "IT"
	testCity := "Rome"
	testAPIKey := "abc123"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Header().Set("Content-Type", "application/json")
		n, err :=io.WriteString(res, "{\"coord\":{\"lon\":12.48,\"lat\":41.89},\"weather\":[{\"id\":802,\"main\":\"Clouds\",\"description\":\"scattered clouds\",\"icon\":\"03n\"}],\"base\":\"stations\",\"main\":{\"temp\":61.05,\"pressure\":1005,\"humidity\":87,\"temp_min\":54,\"temp_max\":66.99},\"visibility\":10000,\"wind\":{\"speed\":4.7,\"deg\":190},\"clouds\":{\"all\":40},\"dt\":1572985001,\"sys\":{\"type\":1,\"id\":6792,\"country\":\"IT\",\"sunrise\":1572932808,\"sunset\":1572969629},\"timezone\":3600,\"id\":6539761,\"name\":\"Rome\",\"cod\":200}")
		if err != nil {
			fmt.Println(n, err)
		}
	}))

	defer ts.Close()

	configs := weathercomparer.WebAPIConfigs{
		BaseURL: ts.URL,
		APIKey:  testAPIKey,
	}

	requestor := &weathercomparer.OpenWeather{Configuration: configs}

	// Act
	_response, err := weathercomparer.ProviderRequestor.WeatherRequest(requestor, testCountry, testCity)

	// Assert
	if err != nil {
		t.Errorf("Error was returned")
	}

	if _response.DegreeCelsius < 1 {
		t.Errorf("Invalid DegreeCelsius value")
	}
}

func TestWeatherRequest_OpenWeather_Failure404(t *testing.T) {
	// Arrange
	testCountry := "IT"
	testCity := "Rome"
	testAPIKey := "abc123"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		http.Error(res, "Not found", http.StatusNotFound)
	}))

	defer ts.Close()

	configs := weathercomparer.WebAPIConfigs{
		BaseURL: ts.URL,
		APIKey:  testAPIKey,
	}

	requestor := &weathercomparer.WeatherBit{Configuration: configs}

	// Act
	_response, err := weathercomparer.ProviderRequestor.WeatherRequest(requestor, testCountry, testCity)

	// Assert
	if err == nil {
		t.Errorf("Error was not returned when it should have. Got %v", err)
	}

	if _response.DegreeCelsius != 0 {
		t.Errorf("Invalid DegreeCelsius value. DegreeCelsius value is %v", _response.DegreeCelsius)
	}
}
