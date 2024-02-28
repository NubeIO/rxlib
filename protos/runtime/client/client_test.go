package client

import (
	"fmt"
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

func TestConvertRestPING(t *testing.T) {
	c, err := NewClient("http", 9090, 8080)
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
	c, err := NewClient("mqtt", 9090, 8080)
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

func callback(string2 string, any2 *Message, err error) {
	fmt.Println("RESP")
	fmt.Println(string2, any2, err)
}
