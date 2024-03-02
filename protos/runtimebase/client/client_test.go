package client

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/mqttwrapper"
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/helpers/pprint"
	"testing"
	"time"
)

func TestConvertGRPCPING(t *testing.T) {
	c, err := NewClient("", "grpc", 9090, 8080, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := c.Ping(nil, callback)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
	time.Sleep(time.Second * 2)
}

func TestConvertRestPING(t *testing.T) {
	c, err := NewClient("", "http", 9090, 8080, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	command := rxlib.CommandPing()
	_, err = c.Command(&Opts{}, command, callbackCommand)
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(time.Second * 2)
}

func TestConvertRestMQTT(t *testing.T) {
	c, err := mqttwrapper.NewMqttClient(mqttwrapper.Config{})
	if err != nil {
		return
	}
	c.Connect()
	c.StartProcessingMessages()
	c.StartPublishRateLimiting()

	client, err := NewClient("", "mqtt", 9090, 8080, c)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Ping(nil, callback)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
	time.Sleep(time.Second * 2)
}

func callbackCommand(string2 string, any2 *rxlib.CommandResponse, err error) {
	msg := &Message{}
	err = json.Unmarshal(any2.Any, &msg)
	if err != nil {
		return
	}
	pprint.PrintJSON(msg)
}

func callback(string2 string, any2 *Message, err error) {
	fmt.Println("RESP")
	fmt.Println(string2, any2, err)
}
