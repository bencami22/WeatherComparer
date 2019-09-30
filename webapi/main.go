package main

import (
	"fmt"
	"weathercomparer"
)

func main() {
	var providers = []weathercomparer.ProviderRequestor{
		&weathercomparer.WeatherBit{},
		&weathercomparer.OpenWeather{}}
	for _, v := range providers {
		var response = weathercomparer.ProviderRequestor.WeatherRequest(v, "IT", "ROME")
		fmt.Println(response)
	}

}
