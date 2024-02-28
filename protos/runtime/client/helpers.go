package client

import (
	"fmt"
	"strings"
)

func ExtractApiTopicPath(input string) (string, string, error) {
	// Split the input string by "/"
	parts := strings.Split(input, "/")
	if len(parts) < 4 {
		return "", "", fmt.Errorf("invalid input format")
	}

	// Get the 3rd element after splitting by "/"
	subParts := strings.Split(parts[3], "_")
	if len(subParts) != 2 {
		return "", "", fmt.Errorf("invalid input format")
	}

	// Extract the elements
	path := subParts[0]
	uuid := subParts[1]

	return path, uuid, nil
}
