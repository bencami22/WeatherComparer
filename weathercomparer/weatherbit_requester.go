package weathercomparer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tkanos/gonfig"
)

//WeatherBit contains integration implementation to weatherbit.io
type WeatherBit struct{}

//WeatherRequest retrieves weather data from weatherbit.io
func (provider *WeatherBit) WeatherRequest(country string, city string) WeatherResponse {

	configuration := Configuration{}
	err := gonfig.GetConf("configuration.json", &configuration)
	if err != nil {
		panic(err)
	}

	resp, err := http.Get("https://api.weatherbit.io/v2.0/current?city=" + city + "&country=" + country + "&key=" + configuration.WeatherBitAPIKey)
	//log.Println("https://api.weatherbit.io/v2.0/current?city=ROME&country=IT&key=" + configuration.WeatherBitAPIKey)
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
	var temp = result["data"].([]interface{})
	var bla = temp[0].(map[string]interface{})
	var finaltemp = bla["temp"].(float64)

	return WeatherResponse{finaltemp}
}
