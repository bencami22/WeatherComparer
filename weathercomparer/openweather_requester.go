package weathercomparer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//OpenWeather contains integration implementation to openweathermap.org
type OpenWeather struct {
	Configuration WebAPIConfigs
}

//WeatherRequest retrieves weather data from openweathermap.org
func (provider *OpenWeather) WeatherRequest(ctx context.Context, country string, city string) (WeatherResponse, error) {

	resp, err := http.Get(provider.Configuration.BaseURL + "/weather?q=" + city + "," +
		country + "&units=imperial&appid=" + provider.Configuration.APIKey)
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

	//log.Println(string(body))

	var result map[string]interface{}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return WeatherResponse{}, err
	}
	var tempObj = result["main"].(map[string]interface{})
	var tempinfahrenheight = tempObj["temp"].(float64)
	var temp = Temperature(tempinfahrenheight)

	return WeatherResponse{Provider: "OpenWeather", DegreeCelsius: temp.toCelsius()}, nil
}
