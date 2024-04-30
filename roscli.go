package rxlib

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/mqttwrapper"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"sync"
)

type ROSClient interface {
	WhoIs(timeout int, targetGlobalID, requestUUID, start, finish, global string) *runtime.CommandResponse
	GetObjects(timeout int, targetGlobalID, requestUUID string) *runtime.CommandResponse
	RQL(timeout int, targetGlobalID, requestUUID, script string) *runtime.CommandResponse
	BulkRQL(timeout int, requestUUID, script string, targetGlobalID ...string) *runtime.CommandResponse
	GlobalRQL(bufferDuration int, requestUUID, script string) *runtime.CommandResponse
}

type rosClient struct {
	mqttclient mqttwrapper.MQTT
	settings   *RuntimeSettings
}

func NewRosClient(mqtt mqttwrapper.MQTT, settings *RuntimeSettings) ROSClient {
	return &rosClient{mqtt, settings}
}

func (inst *rosClient) BulkRQL(timeout int, requestUUID, script string, targetGlobalID ...string) *runtime.CommandResponse {
	c := ExtendedCommand{
		Command: &runtime.Command{
			Key:  "rql",
			Data: map[string]string{"script": script},
		},
	}
	return inst.executeCommandBulk(timeout, targetGlobalID, requestUUID, c)
}

// executeCommandBulk sends a command to multiple targets and collects responses.
func (inst *rosClient) executeCommandBulk(timeout int, targetGlobalIDs []string, requestUUID string, c ExtendedCommand) *runtime.CommandResponse {
	var wg sync.WaitGroup
	responses := make([]*runtime.CommandResponse, len(targetGlobalIDs))
	wg.Add(len(targetGlobalIDs))

	for i, targetID := range targetGlobalIDs {
		go func(index int, targetGlobalID string) {
			defer wg.Done()
			responses[index] = inst.executeCommand(timeout, targetGlobalID, requestUUID, c)
		}(i, targetID)
	}

	wg.Wait()
	return &runtime.CommandResponse{
		Response: responses,
	}
}

func (inst *rosClient) executeCommand(timeout int, targetGlobalID, requestUUID string, c ExtendedCommand) *runtime.CommandResponse {
	topicPublish := fmt.Sprintf("r/req/v1/cloud/%s/plain/command/%s/%s", targetGlobalID, inst.settings.GlobalID, requestUUID)
	topicSub := fmt.Sprintf("r/res/v1/cloud/%s/plain/command/%s/%s", inst.settings.GlobalID, targetGlobalID, requestUUID)
	resp := inst.mqttclient.RequestResponse(timeout, topicPublish, topicSub, requestUUID, c)
	if resp != nil {
		var out *runtime.CommandResponse
		err := json.Unmarshal(resp.Body, &out)
		if err != nil {
			return &runtime.CommandResponse{
				TypeError: fmt.Sprintf("Error unmarshalling response: %s", err),
			}
		}
		return out
	}
	return &runtime.CommandResponse{
		TypeError: "failed to get any response",
	}
}

func (inst *rosClient) WhoIs(timeout int, targetGlobalID, requestUUID, start, finish, global string) *runtime.CommandResponse {
	c := ExtendedCommand{
		Command: &runtime.Command{
			Key:  "command",
			Args: []string{"run", "whois"},
			Data: map[string]string{"start": start, "finish": finish, "global": global},
		},
	}
	return inst.executeCommand(timeout, targetGlobalID, requestUUID, c)
}

func (inst *rosClient) GetObjects(timeout int, targetGlobalID, requestUUID string) *runtime.CommandResponse {
	c := ExtendedCommand{
		Command: &runtime.Command{
			Key:  "command",
			Args: []string{"get", "objects"},
			Data: map[string]string{"as": "json"},
		},
	}
	return inst.executeCommand(timeout, targetGlobalID, requestUUID, c)
}

func (inst *rosClient) RQL(timeout int, targetGlobalID, requestUUID, script string) *runtime.CommandResponse {
	c := ExtendedCommand{
		Command: &runtime.Command{
			Key:  "rql",
			Data: map[string]string{"script": script},
		},
	}
	return inst.executeCommand(timeout, targetGlobalID, requestUUID, c)
}

func (inst *rosClient) GlobalRQL(bufferDuration int, requestUUID, script string) *runtime.CommandResponse {
	c := ExtendedCommand{
		Command: &runtime.Command{
			Key:  "rql",
			Data: map[string]string{"script": script},
		},
	}
	topicPublish := fmt.Sprintf("r/req/v1/cloud/%s/plain/command/%s/%s", "global", inst.settings.GlobalID, requestUUID)
	topicSub := fmt.Sprintf("r/res/v1/cloud/%s/plain/command/+/%s", inst.settings.GlobalID, requestUUID)
	resp := inst.mqttclient.RequestResponseStream(bufferDuration, topicPublish, topicSub, requestUUID, c)
	if resp != nil {
		var out = &runtime.CommandResponse{Response: []*runtime.CommandResponse{}}
		for _, response := range resp {
			var parsed *runtime.CommandResponse
			err := json.Unmarshal(response.Body, &parsed)
			if err == nil {
				out.Response = append(out.Response, parsed)
			}
		}
		return out
	}
	return nil

}
