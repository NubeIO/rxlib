package client

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/mqttwrapper"
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/rxlib/helpers"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"time"
)

type MQTTClient struct {
	mqttClient mqttwrapper.MQTT
	requests   map[string]chan *MQTTPayload
}

func (m *MQTTClient) Command(opts *Opts, command *rxlib.ExtendedCommand, callback func(string, *rxlib.CommandResponse, error)) (string, error) {
	if opts == nil {
		return "", fmt.Errorf("opts body can not be empty")
	}
	requestTopic := fmt.Sprintf("ros/api/%s/command", opts.TargetGlobalID)
	newUUID := helpers.UUID()
	go func() {
		m.RequestResponse(newUUID, requestTopic, command, func(uuid string, payload *Payload, err error) {
			var message *rxlib.CommandResponse
			if err == nil && payload != nil {
				err = json.Unmarshal(payload.Payload, &message)
			}
			callback(newUUID, message, err)
		})
	}()
	// Return the newUUID immediately
	return newUUID, nil
}

func (m *MQTTClient) Close() error {
	//TODO implement me
	panic("implement me")
}

type MQTTPayload struct {
	RequestUUID string      `json:"requestUUID"`
	Payload     interface{} `json:"payload"`
}

func newMQTTClient(mqtt mqttwrapper.MQTT) (*MQTTClient, error) {
	// Create MQTT client options
	return &MQTTClient{
		mqttClient: mqtt,
		requests:   make(map[string]chan *MQTTPayload),
	}, nil
}

func (m *MQTTClient) Ping(opts *Opts, callback func(string, *Message, error)) (string, error) {
	if opts == nil {
		return "", fmt.Errorf("opts body can not be empty")
	}
	requestTopic := fmt.Sprintf("ros/api/%s/ping", opts.TargetGlobalID)
	payloadData := &rxlib.ExtendedCommand{
		Command: &runtime.Command{
			SenderGlobalID: opts.SenderGlobalID,
			Key:            "ping",
		},
	}
	newUUID := helpers.UUID()

	// Start a goroutine to handle the request and callback asynchronously
	go func() {
		m.RequestResponse(newUUID, requestTopic, payloadData, func(uuid string, payload *Payload, err error) {
			var message *Message
			if err == nil && payload != nil {
				err = json.Unmarshal(payload.Payload, &message)
			}
			callback(newUUID, message, err)
		})
	}()

	// Return the newUUID immediately
	return newUUID, nil
}

type Payload struct {
	Payload []byte
	Topic   string
	UUID    string
}

func (m *MQTTClient) RequestResponse(newUUID, requestTopic string, payloadData interface{}, callback func(string, *Payload, error)) (string, error) {

	requestTopicWithUUID := fmt.Sprintf("%s_%s", requestTopic, newUUID)
	respTopicWithUUID := fmt.Sprintf("%s/response", requestTopicWithUUID)
	// Channel to signal the receipt of the message
	done := make(chan struct{})

	err := m.mqttClient.Subscribe(respTopicWithUUID, func(topic string, payload []byte) {
		response := &Payload{
			Payload: payload,
			Topic:   topic,
		}
		_, requestUUID, err := ExtractApiTopicPath(topic)
		if err != nil {
			return
		}
		if requestUUID == newUUID {
			callback(requestUUID, response, nil)
			close(done) // Signal that the message has been received
			return
		}
	})
	if err != nil {
		return "", err
	}

	defer m.mqttClient.Unsubscribe(respTopicWithUUID)

	marshaledPayload, err := json.Marshal(payloadData)
	if err != nil {
		return "", err
	}

	m.mqttClient.Publish(requestTopicWithUUID, marshaledPayload)

	select {
	case <-done:
		// Message received
	case <-time.After(2 * time.Second):
		// Timeout occurred
		callback("", nil, fmt.Errorf("timeout occurred"))
	}

	return newUUID, nil
}

//func (m *MQTTClient) RequestResponse(requestTopic string, payloadData interface{}, callback func(string, *Payload, error)) (string, error) {
//	newUUID := helpers.UUID()
//	requestTopicWithUUID := fmt.Sprintf("%s_%s", requestTopic, newUUID)
//	respTopicWithUUID := fmt.Sprintf("%s/response", requestTopicWithUUID)
//	// Channel to signal the receipt of the message
//	done := make(chan struct{})
//
//	// Subscribe to the response topic
//	token := m.Subscribe(respTopicWithUUID, 0, func(client mqtt.Client, msg mqtt.Message) {
//		response := &Payload{
//			Payload: msg.Payload(),
//			Topic:   msg.Topic(),
//		}
//		_, requestUUID, err := ExtractApiTopicPath(msg.Topic())
//		if err != nil {
//			return
//		}
//		if requestUUID == newUUID {
//			// Handle the response
//			callback(requestUUID, response, nil)
//			close(done) // Signal that the message has been received
//			return
//		}
//	})
//	token.Wait()
//	if token.Error() != nil {
//		return "", token.Error()
//	}
//	defer m.Unsubscribe(respTopicWithUUID)
//
//	// Marshal the payload to JSON
//	marshaledPayload, err := json.Marshal(payloadData)
//	if err != nil {
//		return "", err
//	}
//
//	// Publish the request
//	m.Publish(requestTopicWithUUID, 0, false, marshaledPayload)
//
//	// Wait for response or timeout
//	select {
//	case <-done:
//		// Message received
//	case <-time.After(2 * time.Second):
//		// Timeout occurred
//		callback("", nil, fmt.Errorf("timeout occurred"))
//	}
//
//	return newUUID, nil
//}

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
