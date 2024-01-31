package history

import (
	"fmt"
	"testing"
)

func TestNewHistoryManager(t *testing.T) {

	historyManager := NewHistoryManager()

	// Create a new history with a limit of 10 samples
	history1 := historyManager.NewHistory(10)
	uuid1 := history1.GetUUID()

	// Add some samples to history1
	for i := 1; i <= 5; i++ {
		sample := NewGenericSample(float64(i))
		history1.AddSample(sample)
	}

	// Create another history with a limit of 5 samples
	history2 := historyManager.NewHistory(5)

	uuid2 := history2.GetUUID()

	// Add some samples to history2
	for i := 6; i <= 10; i++ {
		sample := NewGenericSample(float64(i))
		history2.AddSample(sample)

	}

	// Retrieve a history by UUID
	retrievedHistory := historyManager.Get(uuid1)
	fmt.Println(retrievedHistory)

	// Get a list of all histories
	allHistories := historyManager.All()

	// Get all samples across all histories
	allSamples := historyManager.AllSamples()

	// Print the retrieved history, all histories, and all samples
	fmt.Printf("Retrieved History (UUID: %s):\n", uuid1)
	fmt.Printf("Samples in History1: %v\n", retrievedHistory.GetSamples())

	fmt.Println("\nAll Histories:")
	for _, history := range allHistories {
		history.GetFirst()
		fmt.Printf("UUID: %s, Record Count: %d\n", history.GetUUID(), history.SampleCount())
	}

	fmt.Println("\nAll Samples:")
	for uuid, samples := range allSamples {
		fmt.Printf("UUID: %s, Samples: %v\n", uuid, samples)
	}

	fmt.Println(historyManager.All())
	// Drop a specific history by UUID
	historyManager.Drop(uuid2)

	// Drop all histories
	historyManager.DropAll()

}
