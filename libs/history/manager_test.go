package history

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NubeIO/rxlib/helpers/pprint"
	"github.com/sashabaranov/go-openai"
	"testing"
)

func TestNewHistoryManager(t *testing.T) {

	historyManager := NewHistoryManager("EEEE")

	// Create a new history with a limit of 10 samples
	history1 := historyManager.NewHistory(10, "abc")
	uuid1 := history1.GetUUID()

	// Add some samples to history1
	for i := 1; i <= 5; i++ {
		sample := NewGenericRecord(float64(i))
		history1.AddRecord(sample)
	}

	// Create another history with a limit of 5 samples
	history2 := historyManager.NewHistory(5, "aaa")

	//uuid2 := history2.GetUUID()

	// Add some samples to history2
	for i := 6; i <= 10; i++ {
		sample := NewGenericRecord(float64(i))
		history2.AddRecord(sample)

	}

	// Retrieve a history by UUID
	retrievedHistory := historyManager.Get(uuid1)
	fmt.Println(retrievedHistory)

	// Get a list of all histories
	allHistories := historyManager.All()

	pprint.PrintJSON(historyManager.AllHistories())

	// Get all Records across all histories
	allRecords := historyManager.AllRecords()

	// Print the retrieved history, all histories, and all Records
	fmt.Printf("Retrieved History (UUID: %s):\n", uuid1)
	fmt.Printf("Records in History1: %v\n", retrievedHistory.GetRecords())

	fmt.Println("\nAll Histories:")
	for _, history := range allHistories {
		history.GetFirst()
		fmt.Printf("UUID: %s, Record Count: %d\n", history.GetUUID(), history.RecordCount())
	}

	fmt.Println("\nAll Records:")
	for uuid, Records := range allRecords {
		fmt.Printf("UUID: %s, Records: %v\n", uuid, Records)
	}

	client := openai.NewClient("")
	histString := historyManager.AllHistories()
	marshal, err := json.Marshal(histString)
	if err != nil {
		return
	}
	message := "with this data can you give me same stats on it, like min, max avg, count, return it as json"
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("%s  %s", message, string(marshal)),
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	fmt.Println("###############")
	fmt.Println(resp.Choices[0].Message.Content)

	//
	//fmt.Println(historyManager.All())
	//// Drop a specific history by UUID
	//historyManager.Drop(uuid2)
	//
	//// Drop all histories
	//historyManager.DropAll()

}
