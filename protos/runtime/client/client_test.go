package client

import (
	"fmt"
	"github.com/NubeIO/mqttwrapper"
	"testing"
	"time"
)

var jsonString = `

{
  "objectDeploy": {
   "deleted": [
            "abc",
            "123"
        ],
        "new": [
            {
                "id": "trigger",
                "inputs": [],
                "outputs": [
                    {
                        "id": "output",
                        "name": "output",
                        "direction": "output",
                        "dataType": "float"
                    }
                ],
                "connections": [
                    {
                        "source": "triggerABC",
                        "sourceHandle": "output",
                        "target": "mathABC",
                        "targetHandle": "in-1",
                        "flowDirection": "publisher"
                    }
                ],
                "meta": {
                    "uuid": "triggerABC",
                    "name": "triggerABC",
                    "position": {
                        "positionY": -38,
                        "positionX": 155
                    }
                }
            }
        ]
    },
    "timeout": 10
    }`

func TestConvertGRPCPING(t *testing.T) {
	c, err := NewClient("grpc", 9090, 8080, nil)
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
	c, err := NewClient("http", 9090, 8080, nil)
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

func TestConvertRestMQTT(t *testing.T) {
	c, err := mqttwrapper.NewMqttClient(mqttwrapper.Config{})
	if err != nil {
		return
	}
	c.Connect()
	c.StartProcessingMessages()
	c.StartPublishRateLimiting()

	client, err := NewClient("mqtt", 9090, 8080, c)
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

func callback(string2 string, any2 *Message, err error) {
	fmt.Println("RESP")
	fmt.Println(string2, any2, err)
}
