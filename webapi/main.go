package main

import (
	"fmt"
	"os"

	"github.com/bencami22/weathercomparer/weathercomparer"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	var providers = []weathercomparer.ProviderRequestor{
		&weathercomparer.WeatherBit{},
		&weathercomparer.OpenWeather{},
		&weathercomparer.AccuWeather{}}
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
