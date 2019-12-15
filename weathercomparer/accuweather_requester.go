package weathercomparer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//AccuWeather contains integration implementation to accuweather.com
type AccuWeather struct {
	Configuration WebAPIConfigs
}

//WeatherRequest retrieves weather data from accuweather.com
func (provider *AccuWeather) WeatherRequest(ctx context.Context, country string, city string) (WeatherResponse, error) {

	//first we must get accuweathers 'LocationKey'
	resp, err := http.Get(provider.Configuration.BaseURL + "/locations/v1/cities/" +
		country + "/search?apikey=" + provider.Configuration.APIKey + "&q=" + city)
	if err != nil {
		log.Fatalln(err)
		return WeatherResponse{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return WeatherResponse{}, err
	}

	if resp.StatusCode > 299 {
		return WeatherResponse{}, fmt.Errorf("Received %d from remote", resp.StatusCode)
	}

	var result []interface{}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return WeatherResponse{}, err
	}
	locationObj := result[0].(map[string]interface{})
	locationKey := locationObj["Key"].(string)

	resp, err = http.Get(provider.Configuration.BaseURL + "/currentconditions/v1/" + locationKey +
		"?apikey=" + provider.Configuration.APIKey)
	if err != nil {
		log.Fatalln(err)
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return WeatherResponse{}, err
	}

	if resp.StatusCode > 299 {
		return WeatherResponse{}, fmt.Errorf("Received %d from remote", resp.StatusCode)
	}

	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return WeatherResponse{}, err
	}

	tempObj := result[0].(map[string]interface{})
	t := tempObj["Temperature"].(map[string]interface{})
	metricTemperature := t["Metric"].(map[string]interface{})
	temperature := metricTemperature["Value"].(float64)

	return WeatherResponse{Provider: "Accuweather", DegreeCelsius: temperature}, nil
}
