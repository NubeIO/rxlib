package unixbus

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

type EventData struct {
	Topic string
	Data  interface{}
}

type UnixEventBus struct {
	socketPath string
}

// NewUnixEventBus creates a new instance of UnixEventBus for a specific topic
func NewUnixEventBus(topic string) *UnixEventBus {
	socketPath := "/tmp/" + topic + ".sock"
	return &UnixEventBus{socketPath: socketPath}
}

// Publish sends data on the topic
func (ueb *UnixEventBus) Publish(data interface{}) error {
	event := EventData{
		Topic: ueb.socketPath, // use socketPath as topic
		Data:  data,
	}

	conn, err := net.Dial("unix", ueb.socketPath)
	if err != nil {
		return fmt.Errorf("error connecting to Unix socket: %v", err)
	}
	defer conn.Close()

	jsonData, jsonErr := json.Marshal(event)
	if jsonErr != nil {
		return fmt.Errorf("error marshalling event data: %v", jsonErr)
	}

	_, writeErr := conn.Write(append(jsonData, '\n'))
	if writeErr != nil {
		return fmt.Errorf("error writing event data to Unix socket: %v", writeErr)
	}

	return nil
}

// Subscribe listens for data on the topic
func (ueb *UnixEventBus) Subscribe(handler func(data interface{})) {
	go func() {
		conn, err := net.Listen("unix", ueb.socketPath)
		if err != nil {
			fmt.Printf("Error setting up Unix socket for subscription: %v\n", err)
			return
		}
		defer conn.Close()

		for {
			c, err := conn.Accept()
			if err != nil {
				fmt.Printf("Error accepting connection: %v\n", err)
				continue
			}

			go ueb.handleConnection(c, handler)
		}
	}()
}

func (ueb *UnixEventBus) handleConnection(conn net.Conn, handler func(data interface{})) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading from connection: %v\n", err)
			}
			break
		}

		var event EventData
		err = json.Unmarshal(line, &event)
		if err != nil {
			fmt.Printf("Error unmarshalling event data: %v\n", err)
			continue
		}

		handler(event.Data)
	}
}
