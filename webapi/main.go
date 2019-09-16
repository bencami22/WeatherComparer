package main

import(
	"fmt"
	"weathercomparer"
)

func main(){
	var providers = []weathercomparer.ProviderRequestor{
		&weathercomparer.WeatherUnderground{},
		&weathercomparer.OpenWeather{}	}
	for _, v := range providers{
		var response =weathercomparer.ProviderRequestor.WeatherRequest(v, "MT")
		fmt.Println(response)
	} 

}