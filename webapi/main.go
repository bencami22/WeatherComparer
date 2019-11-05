package main

import (
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

	configuration := weathercomparer.Configuration{}
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

	w.Write([]byte(`{"message": "post called"}`))
}

func specificProvider(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	if val, ok := pathParams["provider"]; ok {
		w.Header().Set("Content-Type", "application/json")

		var providerRequestor weathercomparer.ProviderRequestor = nil
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

			w.Write([]byte(fmt.Sprintf(`{"degreeCelsius": "%v"}`, weatherResponse.DegreeCelsius)))

		}
	}
}
