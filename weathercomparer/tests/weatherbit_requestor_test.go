package weathercomparer

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	//"weathercomparer"
	"github.com/bencami22/WeatherComparer/weathercomparer"
)

func TestWeatherRequest_WeatherBit_Success200(t *testing.T) {
	// Arrange
	testCountry := "IT"
	testCity := "Rome"
	testAPIKey := "abc123"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Header().Set("Content-Type", "application/json")
		n, err := io.WriteString(res, "{\"data\":[{\"rh\":83,\"pod\":\"n\",\"lon\":12.51133,\"pres\":1005.3,\"timezone\":\"Europe/Rome\",\"ob_time\":\"2019-11-04 19:50\",\"country_code\":\"IT\",\"clouds\":0,\"ts\":1572897000,\"solar_rad\":0,\"state_code\":\"07\",\"city_name\":\"Rome\",\"wind_spd\":2.24,\"last_ob_time\":\"2019-11-04T19:50:00\",\"wind_cdir_full\":\"south-southeast\",\"wind_cdir\":\"SSE\",\"slp\":1012.4,\"vis\":5,\"h_angle\":-90,\"sunset\":\"16:00\",\"dni\":0,\"dewpt\":14.3,\"snow\":0,\"uv\":0,\"precip\":0,\"wind_dir\":160,\"sunrise\":\"05:46\",\"ghi\":0,\"dhi\":0,\"aqi\":40,\"lat\":41.89193,\"weather\":{\"icon\":\"c01n\",\"code\":\"800\",\"description\":\"Clear sky\"},\"datetime\":\"2019-11-04:20\",\"temp\":17.2,\"station\":\"AT416\",\"elev_angle\":-44.59,\"app_temp\":17.3}],\"count\":1}")
		if err != nil {
			t.Log(n, err)
		}
	}))

	defer ts.Close()

	configs := weathercomparer.WebAPIConfigs{
		BaseURL: ts.URL,
		APIKey:  testAPIKey,
	}

	requestor := &weathercomparer.WeatherBit{Configuration: configs}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(3)*time.Second)
	defer cancel()

	// Act
	_response, err := weathercomparer.ProviderRequestor.WeatherRequest(requestor, ctx, testCountry, testCity)

	// Assert
	if err != nil {
		t.Errorf("Error was returned")
	}

	if _response.DegreeCelsius < 1 {
		t.Errorf("Invalid DegreeCelsius value")
	}
}

func TestWeatherRequest_WeatherBit_Failure404(t *testing.T) {
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

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(3)*time.Second)
	defer cancel()

	// Act
	_response, err := weathercomparer.ProviderRequestor.WeatherRequest(requestor, ctx, testCountry, testCity)

	// Assert
	if err == nil {
		t.Errorf("Error was not returned when it should have. Got %v", err)
	}

	if _response.DegreeCelsius != 0 {
		t.Errorf("Invalid DegreeCelsius value. DegreeCelsius value is %v", _response.DegreeCelsius)
	}
}
