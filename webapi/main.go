package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	//"weathercomparer"
	"github.com/bencami22/WeatherComparer/weathercomparer"

	"github.com/tkanos/gonfig"

	"github.com/gorilla/mux"
)

var configuration weathercomparer.Configuration

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {

	configuration = weathercomparer.Configuration{}
	err := gonfig.GetConf("configuration.json", &configuration)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("", get).Methods(http.MethodGet)
	api.HandleFunc("/{provider}", specificProvider).Methods(http.MethodGet)

	return http.ListenAndServe(":8080", r)
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	providers := map[string]weathercomparer.ProviderRequestor{
		"WeatherBit":  &weathercomparer.WeatherBit{Configuration: configuration.WeatherBitConfiguration},
		"OpenWeather": &weathercomparer.OpenWeather{Configuration: configuration.OpenWeatherConfiguration},
		"AccuWeather": &weathercomsparer.AccuWeather{Configuration: configuration.AccuWeatherConfiguration},
	}

	results := make(map[string]float64)

	for k, v := range providers {

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

		if err == nil {
			fmt.Println(weatherResponse)
			results[k] = weatherResponse.DegreeCelsius
		} else {
			fmt.Println(err)
		}
	}

	json, err := json.Marshal(results)
	if err != nil {
		fmt.Println(err)
	}

	n, err := w.Write(json)
	if err != nil {
		fmt.Println(n, err)
	}
}

func specificProvider(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	if val, ok := pathParams["provider"]; ok {
		w.Header().Set("Content-Type", "application/json")

		var providerRequestor weathercomparer.ProviderRequestor

		print(val)
		switch val {
		case "openweather":
			w.WriteHeader(http.StatusOK)
			providerRequestor = &weathercomparer.OpenWeather{Configuration: configuration.OpenWeatherConfiguration}
		case "accuweather":
			w.WriteHeader(http.StatusOK)
			providerRequestor = &weathercomparer.AccuWeather{Configuration: configuration.AccuWeatherConfiguration}
		case "weatherbit":
			w.WriteHeader(http.StatusOK)
			providerRequestor = &weathercomparer.WeatherBit{Configuration: configuration.WeatherBitConfiguration}
		default:
			w.WriteHeader(http.StatusNotImplemented)
			return
		}
		if providerRequestor != nil {
			weatherResponse, err := weathercomparer.ProviderRequestor.WeatherRequest(providerRequestor, "IT", "ROME")
			if err != nil {
				print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			n, err := w.Write([]byte(fmt.Sprintf(`{"degreeCelsius": "%v"}`, weatherResponse.DegreeCelsius)))
			if err != nil {
				fmt.Println(n, err)
			}

		}
	}
}
