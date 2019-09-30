package weathercomparer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tkanos/gonfig"
)

//OpenWeather contains integration implementation to openweathermap.org
type OpenWeather struct{}

//WeatherRequest retrieves weather data from openweathermap.org
func (provider *OpenWeather) WeatherRequest(country string, city string) WeatherResponse {

	configuration := Configuration{}
	err := gonfig.GetConf("configuration.json", &configuration)
	if err != nil {
		panic(err)
	}

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city + "," + country + "&units=imperial&appid=" + configuration.OpenWeatherAPIKey)
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
	var tempObj = result["main"].(map[string]interface{})
	var tempinfahrenheight = tempObj["temp"].(float64)
	var temp = Temperature(tempinfahrenheight)
	return WeatherResponse{temp.toCelsius()}
}
