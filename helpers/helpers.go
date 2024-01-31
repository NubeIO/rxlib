package helpers

import (
	"encoding/json"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"strings"
	"time"
)

func PrintJSON(x interface{}) {
	ioWriter := os.Stdout
	w := json.NewEncoder(ioWriter)
	w.SetIndent("", "    ")
	w.Encode(x)
}

func UUID(prefix ...string) string {
	u, err := uuid.NewUUID()
	if err != nil {
		return time.Now().Format(time.StampNano)
	}
	// Convert the UUID to a string and remove hyphens
	uuidString := strings.ReplaceAll(u.String(), "-", "")
	uuidString = shuffleCharacters(uuidString)
	shortUUID := uuidString[:16]
	return shortUUID
}

func shuffleCharacters(word string) string {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Convert the word to a slice of characters
	characters := []rune(word)

	// Randomly shuffle the characters using Fisher-Yates algorithm
	for i := len(characters) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		characters[i], characters[j] = characters[j], characters[i]
	}

	// Convert the shuffled slice back to a string
	shuffledWord := string(characters)

	return shuffledWord
}
