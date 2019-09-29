package weathercomparer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//OpenWeather contains integration implementation to openweathermap.org
type OpenWeather struct{}

//WeatherRequest retrieves weather data from openweathermap.org
func (provider *OpenWeather) WeatherRequest(country string) WeatherResponse {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=ROME,IT&appid=")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//log.Println(string(body))

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	var temp = result["main"].(map[string]interface{})
	var finaltemp = temp["temp"].(float64)

	return WeatherResponse{finaltemp}
}
