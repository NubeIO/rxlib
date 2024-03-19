package main

import (
	"fmt"
	"github.com/NubeIO/rxlib/libs/unixbus"
)

func main() {
	eventBus := unixbus.NewUnixEventBus("user.topic.test")
	eventBus.Subscribe(func(data interface{}) {
		fmt.Printf("Received data: %v\n", data)
	})

	// Keep the app running to listen for messages
	select {}
}
