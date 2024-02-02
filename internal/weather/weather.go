package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WeatherLocation struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type WeatherCondition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
}

type WeatherCurrent struct {
	Last_updated_epoch string           `json:"last_updated_epoch"`
	Last_updated       string           `json:"Last_updated"`
	Temp_c             float64          `json:"temp_c"`
	Temp_f             float64          `json:"temp_f"`
	Is_day             int              `json:"is_day"`
	Wind_kph           float64          `json:"wind_kph"`
	Pressure_in        float64          `json:"pressure_in"`
	Precip_in          float64          `json:"precip_in"`
	Condition          WeatherCondition `json:"condition"`
}

type ForecastDay struct {
	Date string `json:"date"`
	Day  struct {
		Avgtemp_c float64          `json:"avgtemp_c"`
		Condition WeatherCondition `json:"condition"`
	} `json:"day"`
}

type WeatherResponse struct {
	Location WeatherLocation `json:"location"`
	Current  WeatherCurrent  `json:"current"`
	Forecast struct {
		Days []ForecastDay `json:"forecastday"`
	} `json:"forecast"`
}

var WEATHER_API_KEY = "cd2a6e5bf8f849e7b5071528240102"

// Query and Number of days of weather forecast. Value ranges from 1 to 6 (Includes today's forecast)
func GetWeather(query string, days int64) (*WeatherResponse, error) {
	url := fmt.Sprintf("https://api.weatherapi.com/v1/forecast.json?q=%s&key=%s&days=%d", query, WEATHER_API_KEY, days+1)

	res, err := http.Get(url)

	resBody := WeatherResponse{}
	json.NewDecoder(res.Body).Decode(&resBody)

	if err != nil {
		return &resBody, err
	}
	res.Body.Close()

	return &resBody, err
}
