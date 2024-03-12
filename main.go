package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}
	return c, nil
}
func hello(w http.ResponseWriter, r *http.Request) {
	html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Hello Page</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f2f2f2;
				}
				.container {
					text-align: center;
					margin-top: 100px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Hello from Go!</h1>
			</div>
		</body>
		</html>
	`
	w.Write([]byte(html))
}

func query(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]
	data, err := queryCity(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Weather Page</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f2f2f2;
				}
				.container {
					text-align: center;
					margin-top: 50px;
				}
				.weather-info {
					margin-top: 20px;
					padding: 10px;
					border: 1px solid #ccc;
					border-radius: 5px;
					background-color: #fff;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Weather Information for {{.Name}}</h1>
				<div class="weather-info">
					<p>City: {{.Name}}</p>
					<p>Temperature: {{.Main.Kelvin}} K</p>
				</div>
			</div>
		</body>
		</html>
	`

	tmpl, err := template.New("weather").Parse(html)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func queryCity(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", query)
	http.ListenAndServe(":8081", nil)
}