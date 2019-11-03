package weathercomparer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//OpenWeather contains integration implementation to openweathermap.org
type OpenWeather struct {
	Configuration WebAPIConfigs
}

//WeatherRequest retrieves weather data from openweathermap.org
func (provider *OpenWeather) WeatherRequest(country string, city string) (WeatherResponse, error) {

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
	//log.Println(string(body))

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	var tempObj = result["main"].(map[string]interface{})
	var tempinfahrenheight = tempObj["temp"].(float64)
	var temp = Temperature(tempinfahrenheight)

	return WeatherResponse{DegreeCelsius: temp.toCelsius()}, nil
}
