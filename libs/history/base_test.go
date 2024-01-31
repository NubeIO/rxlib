package history

import (
	"fmt"
	"testing"
	"time"
)

func TestNewGenericSample(t *testing.T) {
	history := NewGenericHistory(5)
	history2 := NewGenericHistory(5)

	// Capture the start date before adding samples
	startDate := time.Now()

	for i := 0; i < 10; i++ {
		sample := NewGenericSample(float64(i))
		time.Sleep(time.Millisecond * 500)
		history.AddSample(sample)
		sample2 := NewGenericSample(float64(i + 100))
		history2.AddSample(sample2)
	}

	// Wait for a duration longer than the total time it took to add all samples
	time.Sleep(time.Millisecond * 1000)

	duration := "2s" // Use a larger duration to cover the time span of sample addition

	fmt.Printf("Samples within the last %s:\n", duration)
	samplesByTime, err := history.GetSamplesByTime(startDate, duration)
	fmt.Printf("Samples within the error %v:\n", err)
	for i, sample := range samplesByTime {
		fmt.Printf("Record %d - UUID: %s, Value: %v, Timestamp: %v\n", i+1, sample.GetUUID(), sample.GetValue(), sample.GetTimestamp().Format(time.StampMilli))
	}

	history.GetSamples()

}
