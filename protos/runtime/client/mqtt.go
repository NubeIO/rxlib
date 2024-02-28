package client

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/helpers"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type MQTTClient struct {
	mqtt.Client
	requests map[string]chan *MQTTPayload
}

func (m *MQTTClient) Close() error {
	//TODO implement me
	panic("implement me")
}

type MQTTPayload struct {
	RequestUUID string      `json:"requestUUID"`
	Payload     interface{} `json:"payload"`
}

func newMQTTClient() (*MQTTClient, error) {
	// Create MQTT client options
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	// Create MQTT client
	client := mqtt.NewClient(opts)
	// Connect to the broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("error connecting to MQTT broker: %v", token.Error())
	}

	return &MQTTClient{
		Client:   client,
		requests: make(map[string]chan *MQTTPayload),
	}, nil
}

func (m *MQTTClient) Ping(opts *Opts, callback func(string, *Message, error)) (string, error) {
	requestTopic := "ros/api/RX-1/ping"
	payloadData := &rxlib.Command{
		SenderGlobalID: "RX-1",
		Key:            "ping",
	}
	return m.RequestResponse(requestTopic, payloadData, func(uuid string, payload *Payload, err error) {
		var message *Message
		if err == nil && payload != nil {
			err = json.Unmarshal(payload.Payload, &message)
		}
		callback(uuid, message, err)
	})
}

type Payload struct {
	Payload []byte
	Topic   string
	UUID    string
}

func (m *MQTTClient) RequestResponse(requestTopic string, payloadData interface{}, callback func(string, *Payload, error)) (string, error) {
	newUUID := helpers.UUID()
	requestTopicWithUUID := fmt.Sprintf("%s_%s", requestTopic, newUUID)
	respTopicWithUUID := fmt.Sprintf("%s/response", requestTopicWithUUID)
	// Channel to signal the receipt of the message
	done := make(chan struct{})

	// Subscribe to the response topic
	token := m.Subscribe(respTopicWithUUID, 0, func(client mqtt.Client, msg mqtt.Message) {
		response := &Payload{
			Payload: msg.Payload(),
			Topic:   msg.Topic(),
		}
		_, requestUUID, err := ExtractApiTopicPath(msg.Topic())
		if err != nil {
			return
		}
		if requestUUID == newUUID {
			// Handle the response
			callback(requestUUID, response, nil)
			close(done) // Signal that the message has been received
			return
		}
	})
	token.Wait()
	if token.Error() != nil {
		return "", token.Error()
	}
	defer m.Unsubscribe(respTopicWithUUID)

	// Marshal the payload to JSON
	marshaledPayload, err := json.Marshal(payloadData)
	if err != nil {
		return "", err
	}

	// Publish the request
	m.Publish(requestTopicWithUUID, 0, false, marshaledPayload)

	// Wait for response or timeout
	select {
	case <-done:
		// Message received
	case <-time.After(2 * time.Second):
		// Timeout occurred
		callback("", nil, fmt.Errorf("timeout occurred"))
	}

	return newUUID, nil
}

func (m *MQTTClient) ObjectsDeploy(object *rxlib.Deploy, opts *Opts, callback func(*Callback, error)) (string, error) {
	//uuid := uuid.New().String()
	//
	//requestTopic := fmt.Sprintf("objects/deploy/%s/request", uuid)
	//responseTopic := fmt.Sprintf("objects/deploy/%s/response", uuid)
	//
	//// Subscribe to the response topic
	//token := m.Subscribe(responseTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
	//	var response MQTTPayload
	//	if err := json.Unmarshal(msg.Payload(), &response); err != nil {
	//		// Handle JSON unmarshal error
	//		return
	//	}
	//	if response.RequestUUID == uuid {
	//		// Handle the response
	//		callback(response.RequestUUID, response.Payload, nil)
	//	}
	//})
	//token.Wait()
	//if token.Error() != nil {
	//	return "", token.Error()
	//}
	//defer m.Unsubscribe(responseTopic)
	//
	//// Marshal the payload to JSON
	//marshaledPayload, err := json.Marshal(MQTTPayload{
	//	RequestUUID: uuid,
	//	Payload:     object,
	//})
	//if err != nil {
	//	return "", err
	//}
	//
	//// Publish the request
	//m.Publish(requestTopic, 0, false, marshaledPayload)
	//
	//// Wait for response or timeout
	//select {
	//case <-time.After(2 * time.Second):
	//	// Timeout occurred
	//	callback("", nil, fmt.Errorf("timeout occurred"))
	//}

	return "uuid", nil
}
