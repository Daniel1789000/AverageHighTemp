package averagehightemp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type WeatherData struct {
	ConsolidatedWeather []ConsolidatedWeather `json:"consolidated_weather"`
	Time                time.Time             `json:"time"`
	SunRise             time.Time             `json:"sun_rise"`
	SunSet              time.Time             `json:"sun_set"`
	TimezoneName        string                `json:"timezone_name"`
	Parent              Parent                `json:"parent"`
	Sources             []Source              `json:"sources"`
	Title               string                `json:"title"`
	LocationType        string                `json:"location_type"`
	Woeid               int64                 `json:"woeid"`
	LattLong            string                `json:"latt_long"`
	Timezone            string                `json:"timezone"`
}
type ConsolidatedWeather struct {
	ID                   int64     `json:"id"`
	WeatherStateName     string    `json:"weather_state_name"`
	WeatherStateAbbr     string    `json:"weather_state_abbr"`
	WindDirectionCompass string    `json:"wind_direction_compass"`
	Created              time.Time `json:"created"`
	ApplicableDate       string    `json:"applicable_date"`
	MinTemp              float64   `json:"min_temp"`
	MaxTemp              float64   `json:"max_temp"`
	TheTemp              float64   `json:"the_temp"`
	WindSpeed            float64   `json:"wind_speed"`
	WindDirection        float64   `json:"wind_direction"`
	AirPressure          float64   `json:"air_pressure"`
	Humidity             int64     `json:"humidity"`
	Visibility           float64   `json:"visibility"`
	Predictability       int64     `json:"predictability"`
}
type Parent struct {
	Title        string `json:"title"`
	LocationType string `json:"location_type"`
	Woeid        int64  `json:"woeid"`
	LattLong     string `json:"latt_long"`
}

type Source struct {
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	URL       string `json:"url"`
	CrawlRate int64  `json:"crawl_rate"`
}
type Place struct {
	PlaceName string
	Woeid     string
}

type Result struct {
	PlaceName       string
	AverageHighTemp string
}

func GetAverageHighTemp(place Place, channel chan<- Result) {
	result := Result{place.PlaceName,
		"error:unknown"}
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.metaweather.com/api/location/"+place.Woeid+"/", nil)
	if err != nil {
		result.AverageHighTemp = "error:" + err.Error()
	} else {
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		resp, err := httpClient.Do(req)
		if err != nil {
			result.AverageHighTemp = "error:" + err.Error()
		} else {
			defer resp.Body.Close()
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				result.AverageHighTemp = "error:" + err.Error()
			} else {
				var responseObject WeatherData
				json.Unmarshal(bodyBytes, &responseObject)
				numConsolidatedWeather := len(responseObject.ConsolidatedWeather)
				if numConsolidatedWeather > 0 {
					var acum float64 = 0
					for _, consolidatedWeather := range responseObject.ConsolidatedWeather {
						acum += float64(consolidatedWeather.MaxTemp)
					}
					aveTemp := float64(acum / float64(numConsolidatedWeather))
					result.AverageHighTemp = fmt.Sprintf("%.2f", aveTemp)
				} else {
					result.AverageHighTemp = "error:no results"
				}
			}
		}
	}
	channel <- result
}
