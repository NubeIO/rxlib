package history

import (
	"fmt"
	"testing"
	"time"
)

func TestConvertToDataFrame(t *testing.T) {
	historyManager := NewHistoryManager("EEEE")

	histories := historyManager.NewHistory(10, "abc")

	for i := 1; i <= 5; i++ {
		sample := NewGenericRecord(float64(i))
		histories.AddRecord(sample)
	}
	now := time.Now()
	allHistories := []*AllHistories{
		{
			ObjectUUID: "obj-1",
			Histories: []Record{
				&GenericRecord[float64]{UUID: "rec-1", Value: 15.5, Timestamp: now.Add(-1 * time.Hour)},
				&GenericRecord[float64]{UUID: "rec-2", Value: 20.5, Timestamp: now.Add(-30 * time.Minute)},
			},
		},
		{
			ObjectUUID: "obj-2",
			Histories: []Record{
				&GenericRecord[float64]{UUID: "rec-3", Value: 5.5, Timestamp: now.Add(-2 * time.Hour)},
				&GenericRecord[float64]{UUID: "rec-4", Value: 10.5, Timestamp: now.Add(-3 * time.Hour)},
			},
		},
	}

	dfOps := New(allHistories)
	filteredDf := dfOps.GetDF()
	fmt.Println("Filtered DataFrame between dates:")
	fmt.Println(filteredDf)
	startDate := time.Now().Add(-1 * time.Hour)
	endDate := time.Now()

	filteredData := dfOps.FilterDataBetween(startDate, endDate, "Timestamp")

	filteredDf = filteredData.GetDF()
	fmt.Println("Filtered DataFrame between dates:")
	fmt.Println(filteredDf)

	filteredData = dfOps.FilterDateRange("2024-04-13:15:04", "2024-04-13:16:04", "Timestamp")

	filteredDf = filteredData.GetDF()
	fmt.Println("Filtered DataFrame between dates:")
	fmt.Println(filteredDf)
}
