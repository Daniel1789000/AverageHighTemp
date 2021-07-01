package averagehightemp

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestGetAverageHighTemp(t *testing.T) {
	places := [4]Place{
		{"Salt Lake City", "2487610"},
		{"Los Angeles", "2442047"},
		{"Boise", "2366355"},
		{"nowhere", "0"}}
	channel := make(chan Result)
	for _, place := range places {
		fmt.Printf("Starting query for %s\n", place.PlaceName)
		go GetAverageHighTemp(place, channel)
	}

	for index := 0; index < len(places); index++ {
		select {
		case result := <-channel:
			fmt.Printf("%s : %s\n", result.PlaceName, result.AverageHighTemp)
			if result.PlaceName == places[len(places)-1].PlaceName {
				if !strings.Contains(result.AverageHighTemp, "no results") {
					t.Errorf("expected no results error but got %s", result.AverageHighTemp)
				}
			} else if strings.HasPrefix(result.AverageHighTemp, "error") {
				t.Errorf("expected a number but got %s", result.AverageHighTemp)
			}
		case <-time.After(time.Second * 2):
			fmt.Println("timed out before receiving all results")
		}
	}
}
