package main

import (
	"fmt"
	"github.com/NubeIO/rxlib/libs/pglib"
	"log"
	"math/rand"
	"time"
)

func main() {
	connString := "postgresql://postgres:postgres@localhost/postgres"

	db, err := pglib.New(connString)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := createTemperatureSamples()
	fmt.Println(data)
	err = db.InsertBulk("temperature_data", data)
	if err != nil {
		log.Fatalf("Error inserting data: %v", err)
	}

	// Query the data from the table
	query := "SELECT * FROM temperature_data limit 2"
	results := db.Select(query)

	db.Print(results)

}

func createTemperatureSamples() []map[string]interface{} {
	var samples []map[string]interface{}
	startTime := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.AddDate(0, 0, 14) // 14 days of data
	currentTime := startTime
	for currentTime.Before(endTime) {
		// Generate samples for every 5 to 15-minute interval
		interval := time.Duration(rand.Intn(11)+5) * time.Minute
		currentTime = currentTime.Add(interval)

		// Check if current time is between 8am and 5pm on Monday to Friday
		var sensor3Value int
		if currentTime.Weekday() >= time.Monday && currentTime.Weekday() <= time.Friday && currentTime.Hour() >= 8 && currentTime.Hour() < 17 {
			sensor3Value = 1
		} else {
			sensor3Value = 0
		}

		for i := 0; i < 3; i++ {
			sensorName := fmt.Sprintf("sensor_%d", i+1)
			sample := make(map[string]interface{})
			sample["time"] = currentTime.Format(time.RFC3339) // Format time as string
			if i == 0 {
				sample["value"] = rand.Float64()*10 + 20 // Random temperature between 20 and 30
				sample["tags"] = []string{"temp", "ahu", "number", "supply"}
			} else if i == 1 {
				sample["value"] = rand.Float64()*10 + 20 // Random temperature between 20 and 30
				sample["tags"] = []string{"number", "temp", "ahu", "zone"}
			} else {
				sample["value"] = sensor3Value
				sample["tags"] = []string{"status", "enable", "ahu", "bool"}
			}
			sample["sensor_name"] = sensorName
			samples = append(samples, sample)
		}
	}
	return samples
}

/*

CREATE TABLE temperature_data (
    time        TIMESTAMPTZ       NOT NULL,
    temperature DOUBLE PRECISION NOT NULL,
    tags        TEXT[]            NOT NULL,
    sensor_name TEXT
);
SELECT create_hypertable('temperature_data', 'time');

*/
