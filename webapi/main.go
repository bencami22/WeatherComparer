package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	//"weathercomparer"
	"github.com/bencami22/WeatherComparer/weathercomparer"

	"github.com/tkanos/gonfig"

	"github.com/gorilla/mux"
)

const (
	amountOfConcurrentWorkers = 2
	timeoutSeconds            = 3
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

//WeatherResponseWrapper used to return both api response and error from goroutine.
type WeatherResponseWrapper struct {
	WeatherResponse weathercomparer.WeatherResponse
	Error           error
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	providers := map[string]weathercomparer.ProviderRequestor{
		"WeatherBit":  &weathercomparer.WeatherBit{Configuration: configuration.WeatherBitConfiguration},
		"OpenWeather": &weathercomparer.OpenWeather{Configuration: configuration.OpenWeatherConfiguration},
		"AccuWeather": &weathercomparer.AccuWeather{Configuration: configuration.AccuWeatherConfiguration},
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	results := make(map[string]float64)

	jobs := make(chan weathercomparer.ProviderRequestor, len(providers))
	jobResults := make(chan WeatherResponseWrapper, len(providers))
	var wg sync.WaitGroup

	for i := 0; i < amountOfConcurrentWorkers; i++ {
		go worker(ctx, i, jobs, jobResults, &wg)
	}

	for _, v := range providers {
		wg.Add(1)
		jobs <- v
	}

	wg.Wait()

	close(jobs)

	for i := 1; i < len(providers); i++ {
		result := <-jobResults
		if result.Error != nil {
			fmt.Println("Error")
			continue
		}
		results["dsss"] = result.WeatherResponse.DegreeCelsius
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

func worker(ctx context.Context, id int, providers <-chan weathercomparer.ProviderRequestor, results chan<- WeatherResponseWrapper, wg *sync.WaitGroup) {
	for p := range providers {
		weatherResponse, err := weathercomparer.ProviderRequestor.WeatherRequest(ctx, p, "IT", "ROME")
		results <- WeatherResponseWrapper{WeatherResponse: weatherResponse, Error: err}
		wg.Done()
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
