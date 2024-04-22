package history

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers"
	"math/rand"
	"testing"
	"time"
)

func TestFilterByDate(t *testing.T) {
	now := time.Now()

	fmt.Println("Current day:", int(now.Weekday()))

	allHistories := GenerateTwoWeeksData()
	dfOps := New(allHistories)
	startTime := "2024-04-22:08:00"
	endTime := "2024-04-22:17:00"
	filteredData := dfOps.FilterByDateDays(startTime, endTime, 2)
	fmt.Println(filteredData.Count())

}

func TestFilterByTimeDays(t *testing.T) {
	now := time.Now()

	fmt.Println("Current day:", int(now.Weekday()))

	allHistories := GenerateTwoWeeksData()

	fmt.Println(allHistories)

	dfOps := New(allHistories)
	startTime := "08:00"
	endTime := "19:00"
	filteredData := dfOps.FilterByTimeDays(startTime, endTime, 1, 2, 3)

	filteredDf := filteredData.ToDF()
	fmt.Println(filteredDf)

}

func TestConvertToDataFrame(t *testing.T) {
	historyManager := NewHistoryManager("EEEE")

	histories := historyManager.NewHistory(10, "abc")

	for i := 1; i <= 5; i++ {
		sample := NewGenericRecord(float64(i))
		histories.AddRecord(sample)
	}

	allHistories := GenerateTwoWeeksData()

	fmt.Println(allHistories)

	dfOps := New(allHistories)
	filteredDf := dfOps.ToDF()
	fmt.Println("Filtered DataFrame between dates:")
	fmt.Println(filteredDf)
	startDate := time.Now().Add(-1 * time.Hour)
	endDate := time.Now()

	filteredData := dfOps.FilterDataBetween(startDate, endDate, "Timestamp")

	filteredDf = filteredData.ToDF()
	fmt.Println("Filtered DataFrame between dates:")
	fmt.Println(filteredDf)

	filteredData = dfOps.FilterDateRange("2024-04-13:15:04", "2024-04-13:16:04", "Timestamp")

	filteredDf = filteredData.ToDF()
	fmt.Println("Filtered DataFrame between dates:")
	fmt.Println(filteredDf)
}

func GenerateTwoWeeksData() []*AllHistories {
	// Initialize the list of AllHistories
	allHistories := []*AllHistories{}

	// Define the start and end times for the two-week period
	startTime := time.Now().AddDate(0, 0, -14) // Two weeks ago
	endTime := time.Now()

	// Iterate over each day in the two-week period
	for currentTime := startTime; currentTime.Before(endTime); currentTime = currentTime.Add(time.Hour * 24) {
		// Generate random data for each object UUID
		for _, uuid := range []string{"obj-1", "obj-2"} {
			// Generate random values for each time interval (5 to 30 minutes)
			for currentInterval := 5; currentInterval <= 30; currentInterval += 5 {
				// Generate a random timestamp within the current day and interval
				randomTime := currentTime.Add(time.Minute * time.Duration(currentInterval))
				// Generate a random value for the record
				randomValue := rand.Float64() * 100 // Adjust the range as needed

				// Create a new record and add it to the AllHistories list
				record := &GenericRecord[float64]{UUID: helpers.UUID(), Value: randomValue, Timestamp: randomTime}
				allHistories = append(allHistories, &AllHistories{ObjectUUID: uuid, Histories: []Record{record}})
			}
		}
	}

	return allHistories
}
