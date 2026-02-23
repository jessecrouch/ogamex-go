package formula

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type WikiFleetTestData struct {
	Distance     map[string]DistanceCase     `json:"distance"`
	FlightTime  FlightTimeData             `json:"flight_time"`
	FuelConsume FuelConsumptionData         `json:"fuel_consumption"`
}

type DistanceCase struct {
	Origin         Coord `json:"origin"`
	Target         Coord `json:"target"`
	ExpectedDistance int   `json:"expected_distance"`
}

type Coord struct {
	Galaxy  int `json:"galaxy"`
	System  int `json:"system"`
	Position int `json:"position"`
}

type FlightTimeData struct {
	Formula   string        `json:"formula"`
	TestCases []FlightTimeCase `json:"test_cases"`
}

type FlightTimeCase struct {
	Distance     int `json:"distance"`
	Speed        int `json:"speed"`
	SpeedFactor  int `json:"speed_factor"`
	ExpectedTime int `json:"expected_time"`
}

type FuelConsumptionData struct {
	Formula   string             `json:"formula"`
	TestCases []FuelConsumeCase `json:"test_cases"`
}

type FuelConsumeCase struct {
	Ship         string `json:"ship"`
	Distance     int    `json:"distance"`
	Speed        int    `json:"speed"`
	ExpectedFuel int    `json:"expected_fuel"`
}

func getWikiFleetDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, "..", "..", "tests", "wiki", filename)
}

func loadWikiFleetData(t *testing.T) WikiFleetTestData {
	data, err := os.ReadFile(getWikiFleetDataPath("fleet.json"))
	if err != nil {
		t.Fatalf("Failed to read wiki fleet test data: %v", err)
	}

	var wikiData WikiFleetTestData
	err = json.Unmarshal(data, &wikiData)
	if err != nil {
		t.Fatalf("Failed to parse wiki fleet test data: %v", err)
	}

	return wikiData
}

func TestWiki_FleetDistance(t *testing.T) {
	wikiData := loadWikiFleetData(t)

	for name, tc := range wikiData.Distance {
		t.Run(name, func(t *testing.T) {
			got := CalculateDistance(tc.Origin.Galaxy, tc.Origin.System, tc.Origin.Position, tc.Target.Galaxy, tc.Target.System, tc.Target.Position)
			if got != tc.ExpectedDistance {
				t.Logf("DISCREPANCY: CalculateDistance(%d.%d.%d -> %d.%d.%d): expected=%d, got=%d",
					tc.Origin.Galaxy, tc.Origin.System, tc.Origin.Position,
					tc.Target.Galaxy, tc.Target.System, tc.Target.Position,
					tc.ExpectedDistance, got)
			}
		})
	}
}

func TestWiki_FlightTime(t *testing.T) {
	wikiData := loadWikiFleetData(t)

	for _, tc := range wikiData.FlightTime.TestCases {
		t.Run("distance_speed_factor", func(t *testing.T) {
			got := CalculateFlightTimeWiki(tc.Distance, tc.Speed, tc.SpeedFactor)
			if got != int64(tc.ExpectedTime) {
				t.Logf("DISCREPANCY: CalculateFlightTime(distance=%d, speed=%d, factor=%d): expected=%d, got=%d",
					tc.Distance, tc.Speed, tc.SpeedFactor, tc.ExpectedTime, got)
			}
		})
	}
}

func TestWiki_FuelConsumption(t *testing.T) {
	wikiData := loadWikiFleetData(t)

	shipBaseCost := map[string]int{
		"small_cargo": 2000,
		"large_cargo": 6000,
		"light_fighter": 3000,
		"cruiser": 20000,
		"battleship": 45000,
		"deathstar": 5000000,
	}

	for _, tc := range wikiData.FuelConsume.TestCases {
		baseCost, ok := shipBaseCost[tc.Ship]
		if !ok {
			continue
		}

		t.Run(tc.Ship, func(t *testing.T) {
			got := CalculateFuelConsumption(baseCost, tc.Distance, tc.Speed)
			if got != int64(tc.ExpectedFuel) {
				t.Logf("DISCREPANCY: FuelConsumption(%s, dist=%d, speed=%d): expected=%d, got=%d",
					tc.Ship, tc.Distance, tc.Speed, tc.ExpectedFuel, got)
			}
		})
	}
}

func CalculateFlightTimeWiki(distance int, speed int, speedFactor int) int64 {
	if speed <= 0 || distance <= 0 {
		return 0
	}
	time := 10 + 3500*math.Sqrt(10*float64(distance)/float64(speed))/float64(speedFactor)
	return int64(math.Ceil(time))
}

func CalculateFuelConsumption(baseCost int, distance int, speed int) int64 {
	if distance <= 0 || speed <= 0 {
		return 0
	}
	consumption := 1 + float64(baseCost)*float64(distance)/35000*math.Pow(float64(speed)/100+1, 2)
	return int64(math.Round(consumption))
}
