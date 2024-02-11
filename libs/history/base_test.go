package history

import (
	"fmt"
	"testing"
	"time"
)

func TestNewGenericSample(t *testing.T) {
	history := NewGenericHistory(5, "abc")
	history2 := NewGenericHistory(5, "abc")

	// Capture the start date before adding Records
	startDate := time.Now()

	for i := 0; i < 10; i++ {
		sample := NewGenericRecord(float64(i))
		time.Sleep(time.Millisecond * 500)
		history.AddRecord(sample)
		sample2 := NewGenericRecord(float64(i + 100))
		history2.AddRecord(sample2)
	}

	// Wait for a duration longer than the total time it took to add all Records
	time.Sleep(time.Millisecond * 1000)

	duration := "2s" // Use a larger duration to cover the time span of sample addition

	fmt.Printf("Records within the last %s:\n", duration)
	RecordsByTime, err := history.GetRecordsByTime(startDate, duration)
	fmt.Printf("Records within the error %v:\n", err)
	for i, sample := range RecordsByTime {
		fmt.Printf("Record %d - UUID: %s, Values: %v, Timestamp: %v\n", i+1, sample.GetUUID(), sample.GetValue(), sample.GetTimestamp().Format(time.StampMilli))
	}

	history.GetRecords()

}
