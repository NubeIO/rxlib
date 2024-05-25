package client

import (
	"fmt"
	"github.com/NubeIO/mqttwrapper"
	"testing"
	"time"
)

func TestConvertGRPCPING(t *testing.T) {
	c, err := NewClient("", "grpc", 9090, 1770, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := c.Ping(&Opts{Timeout: defaultTimeout}, callback)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
	time.Sleep(time.Second * 2)
}

func TestConvertRestMQTT(t *testing.T) {
	c, err := mqttwrapper.NewMqttClient(mqttwrapper.Config{})
	if err != nil {
		return
	}
	c.Connect()
	c.StartProcessingMessages()

	client, err := NewClient("", "mqtt", 9090, 1770, c)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Ping(&Opts{Timeout: defaultTimeout}, callback)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
	time.Sleep(time.Second * 2)
}

func callback(string2 string, any2 *Message, err error) {
	fmt.Println("RESP")
	fmt.Println(string2, any2, err)
}
