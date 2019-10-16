package weathercomparer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tkanos/gonfig"
)

//AccuWeather contains integration implementation to accuweather.com
type AccuWeather struct{}

//WeatherRequest retrieves weather data from accuweather.com
func (provider *AccuWeather) WeatherRequest(country string, city string) WeatherResponse {

	configuration := Configuration{}
	err := gonfig.GetConf("configuration.json", &configuration)
	if err != nil {
		panic(err)
	}

	//first we must get accuweathers 'LocationKey'
	resp, err := http.Get("http://dataservice.accuweather.com/locations/v1/cities/" + country + "/search?apikey=" + configuration.AccuWeatherAPIKey + "&q=" + city)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//log.Println(string(body))

	var result []interface{}
	json.Unmarshal([]byte(body), &result)
	locationObj := result[0].(map[string]interface{})
	locationKey := locationObj["Key"].(string)

	resp, err = http.Get("http://dataservice.accuweather.com/currentconditions/v1/" + locationKey + "?apikey=" + configuration.AccuWeatherAPIKey)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//log.Println(string(body))

	json.Unmarshal([]byte(body), &result)
	tempObj := result[0].(map[string]interface{})
	t := tempObj["Temperature"].(map[string]interface{})
	metricTemperature := t["Metric"].(map[string]interface{})
	temperature := metricTemperature["Value"].(float64)

	return WeatherResponse{temperature}
}
