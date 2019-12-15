package weathercomparer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//WeatherBit contains integration implementation to weatherbit.io
type WeatherBit struct {
	Configuration WebAPIConfigs
}

//WeatherRequest retrieves weather data from weatherbit.io
func (provider *WeatherBit) WeatherRequest(ctx context.Context, country string, city string) (WeatherResponse, error) {
	resp, err := http.Get(provider.Configuration.BaseURL + "/current?city=" + city +
		"&country=" + country + "&key=" + provider.Configuration.APIKey)

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

	var result map[string]interface{}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return WeatherResponse{}, err
	}
	var temp = result["data"].([]interface{})
	var bla = temp[0].(map[string]interface{})
	var finaltemp = bla["temp"].(float64)

	return WeatherResponse{Provider: "WeatherBit", DegreeCelsius: finaltemp}, nil
}
