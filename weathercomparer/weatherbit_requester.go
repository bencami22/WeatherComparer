package weathercomparer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//WeatherBit contains integration implementation to weatherbit.io
type WeatherBit struct {
	Configuration WebAPIConfigs
}

//WeatherRequest retrieves weather data from weatherbit.io
func (provider *WeatherBit) WeatherRequest(country string, city string) WeatherResponse {
	resp, err := http.Get(provider.Configuration.BaseURL + "current?city=" + city +
		"&country=" + country + "&key=" + provider.Configuration.APIKey)

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
