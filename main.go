package main

import (
	"fmt"
	"time"

	"travelpass.averagehightemp/m/v1/averagehightemp"
)

func main() {
	places := [3]averagehightemp.Place{
		{PlaceName: "Salt Lake City", Woeid: "2487610"},
		{PlaceName: "Los Angeles", Woeid: "2442047"},
		{PlaceName: "Boise", Woeid: "2366355"},
		//{"nowhere", "0"},
	}
	channel := make(chan averagehightemp.Result)
	for _, place := range places {
		fmt.Printf("Starting query for %s\n", place.PlaceName)
		go averagehightemp.GetAverageHighTemp(place, channel)
	}

	for index := 0; index < len(places); index++ {
		select {
		case result := <-channel:
			fmt.Printf("%s Average Max Temp: %s\n", result.PlaceName, result.AverageHighTemp)
		case <-time.After(time.Second * 2):
			fmt.Println("timed out before receiving all results")
		}
	}
}
