package rxlib

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/mqttwrapper"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
)

type ROSClient interface {
	WhoIs(timeout int, targetGlobalID, requestUUID, start, finish, global string) *runtime.CommandResponse
	GetObjects(timeout int, targetGlobalID, requestUUID string) *runtime.CommandResponse
}

type rosClient struct {
	mqttclient mqttwrapper.MQTT
	settings   *RuntimeSettings
}

func NewRosClient(mqtt mqttwrapper.MQTT, settings *RuntimeSettings) ROSClient {
	return &rosClient{mqtt, settings}
}

type commandParams struct {
	Args   []string
	Start  string
	Finish string
	Global string
}

func (inst *rosClient) executeCommand(timeout int, targetGlobalID, requestUUID string, params *commandParams) *runtime.CommandResponse {
	c := ExtendedCommand{
		Command: &runtime.Command{
			Key:  "command",
			Args: params.Args,
			Data: map[string]string{"start": params.Start, "finish": params.Finish, "global": params.Global},
		},
	}
	topicPublish := fmt.Sprintf("r/req/v1/cloud/%s/plain/command/%s/%s", targetGlobalID, inst.settings.GlobalID, requestUUID)
	topicSub := fmt.Sprintf("r/res/v1/cloud/%s/plain/command/%s/%s", inst.settings.GlobalID, targetGlobalID, requestUUID)
	resp := inst.mqttclient.RequestResponse(timeout, topicPublish, topicSub, requestUUID, c)
	if resp != nil {
		var out *runtime.CommandResponse
		err := json.Unmarshal(resp.Body, &out)
		if err != nil {
			return &runtime.CommandResponse{
				Error: fmt.Sprintf("Error unmarshalling response: %s", err),
			}
		}
		return out
	}
	return &runtime.CommandResponse{
		Error: "failed to get any response",
	}
}

func (inst *rosClient) WhoIs(timeout int, targetGlobalID, requestUUID, start, finish, global string) *runtime.CommandResponse {
	return inst.executeCommand(timeout, targetGlobalID, requestUUID, &commandParams{
		Args:   []string{"run", "whois"},
		Start:  start,
		Finish: finish,
		Global: global,
	})
}

func (inst *rosClient) GetObjects(timeout int, targetGlobalID, requestUUID string) *runtime.CommandResponse {
	params := &commandParams{
		Args: []string{"get", "objects"},
	}
	return inst.executeCommand(timeout, targetGlobalID, requestUUID, params)
}
