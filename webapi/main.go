package main

import (
	"fmt"
	"os"

	"github.com/bencami22/WeatherComparer/weathercomparer"
	//"weathercomparer"

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
		tempChannel := make(chan weathercomparer.WeatherResponse, 1)
		errChannel := make(chan error, 1)
		go func() {
			defer close(tempChannel)
			defer close(errChannel)
			weatherResponse, err := weathercomparer.ProviderRequestor.WeatherRequest(v, "IT", "ROME")
			if err == nil {
				tempChannel <- weatherResponse
				return
			}
			print(err)
			errChannel <- err
		}()
		weatherResponse := <-tempChannel
		err := <-errChannel
		fmt.Println(weatherResponse)
		fmt.Println(err)

	}
	return nil
}
