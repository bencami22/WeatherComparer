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

	//first we must get accuweathers 'LocationId'
	resp, err := http.Get("http://dataservice.accuweather.com/locations/v1/cities/" + country + "/search?apikey=" + configuration.AccuWeatherAPIKey + "&q=" + city)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))

	var result []interface{}
	json.Unmarshal([]byte(body), &result)
	var tempObj = result[0]
	log.Println(tempObj)

	return WeatherResponse{100}
}
