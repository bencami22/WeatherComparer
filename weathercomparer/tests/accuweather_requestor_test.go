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

func TestWeatherRequest_AccuWeather_Success(t *testing.T) {
	// Arrange
	testCountry := "IT"
	testCity := "Rome"
	testAPIKey := "abc123"

	mux := http.NewServeMux()
	mux.HandleFunc("/locations/v1/cities/"+
		testCountry+"/search",
		func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			res.Header().Set("Content-Type", "application/json")
			n, err := io.WriteString(res, "[{\"Version\":1,\"Key\":\"213490\",\"Type\":\"City\",\"Rank\":20,\"LocalizedName\":\"Rome\",\"EnglishName\":\"Rome\",\"PrimaryPostalCode\":\"\",\"Region\":{\"ID\":\"EUR\",\"LocalizedName\":\"Europe\",\"EnglishName\":\"Europe\"},\"Country\":{\"ID\":\"IT\",\"LocalizedName\":\"Italy\",\"EnglishName\":\"Italy\"},\"AdministrativeArea\":{\"ID\":\"62\",\"LocalizedName\":\"Lazio\",\"EnglishName\":\"Lazio\",\"Level\":1,\"LocalizedType\":\"Region\",\"EnglishType\":\"Region\",\"CountryID\":\"IT\"},\"TimeZone\":{\"Code\":\"CET\",\"Name\":\"Europe/Rome\",\"GmtOffset\":1.0,\"IsDaylightSaving\":false,\"NextOffsetChange\":\"2020-03-29T01:00:00Z\"},\"GeoPosition\":{\"Latitude\":41.892,\"Longitude\":12.511,\"Elevation\":{\"Metric\":{\"Value\":45.0,\"Unit\":\"m\",\"UnitType\":5},\"Imperial\":{\"Value\":147.0,\"Unit\":\"ft\",\"UnitType\":0}}},\"IsAlias\":false,\"SupplementalAdminAreas\":[{\"Level\":2,\"LocalizedName\":\"Roma\",\"EnglishName\":\"Roma\"},{\"Level\":3,\"LocalizedName\":\"Roma\",\"EnglishName\":\"Roma\"}],\"DataSets\":[\"Alerts\"]}]")
			if err != nil {
				t.Log(n, err)
			}
		})

	mux.HandleFunc("/currentconditions/v1/213490",
		func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			res.Header().Set("Content-Type", "application/json")
			n, err := io.WriteString(res, "[{\"LocalObservationDateTime\":\"2019-11-03T17:12:00+01:00\",\"EpochTime\":1572797520,\"WeatherText\":\"Thundershower\",\"WeatherIcon\":42,\"HasPrecipitation\":true,\"PrecipitationType\":\"Rain\",\"IsDayTime\":false,\"Temperature\":{\"Metric\":{\"Value\":21.1,\"Unit\":\"C\",\"UnitType\":17},\"Imperial\":{\"Value\":70.0,\"Unit\":\"F\",\"UnitType\":18}},\"MobileLink\":\"http://m.accuweather.com/en/it/rome/213490/current-weather/213490?lang=en-us\",\"Link\":\"http://www.accuweather.com/en/it/rome/213490/current-weather/213490?lang=en-us\"}]")
			if err != nil {
				t.Log(n, err)
			}
		})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	configs := weathercomparer.WebAPIConfigs{
		BaseURL: ts.URL,
		APIKey:  testAPIKey,
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(3)*time.Second)
	defer cancel()

	requestor := &weathercomparer.AccuWeather{Configuration: configs}

	// Act
	_response, err := weathercomparer.ProviderRequestor.WeatherRequest(ctx, requestor, testCountry, testCity)

	// Assert
	if err != nil {
		t.Errorf("Error was returned. Error is %v", err)
	}

	if _response.DegreeCelsius < 1 {
		t.Errorf("Invalid DegreeCelsius value")
	}
}
