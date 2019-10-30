package main

import (
	"fmt"
	"os"

	"github.com/bencami22/WeatherComparer/weathercomparer"

	"github.com/tkanos/gonfig"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {

	configuration := weathercomparer.Configuration{}
	err := gonfig.GetConf("configuration.json", &configuration)
	if err != nil {
		panic(err)
	}

	var providers = []weathercomparer.ProviderRequestor{
		&weathercomparer.WeatherBit{Configuration: configuration.WeatherBitConfiguration},
		&weathercomparer.OpenWeather{Configuration: configuration.OpenWeatherConfiguration},
		&weathercomparer.AccuWeather{Configuration: configuration.AccuWeatherConfiguration},
	}
	for _, v := range providers {
		tempChannel := make(chan weathercomparer.WeatherResponse)
		go func() {
			tempChannel <- weathercomparer.ProviderRequestor.WeatherRequest(v, "IT", "ROME")
			close(tempChannel)
		}()
		fmt.Println(<-tempChannel)
	}
	return nil
}
