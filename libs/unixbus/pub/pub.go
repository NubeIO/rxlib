package main

import (
	"fmt"
	"github.com/NubeIO/rxlib/libs/unixbus"
)

func main() {
	eventBus := unixbus.NewUnixEventBus("user.topic.test")
	err := eventBus.Publish("Hello, World!")
	if err != nil {
		fmt.Printf("Publish error: %v\n", err)
	}
}
