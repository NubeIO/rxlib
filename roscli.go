package rxlib

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/mqttwrapper"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
)

type ROSClient interface {
	WhoIs(timeout int, targetGlobalID, requestUUID string) *runtime.CommandResponse
}

type rosClient struct {
	mqttclient mqttwrapper.MQTT
	settings   *RuntimeSettings
}

func NewRosClient(mqtt mqttwrapper.MQTT, settings *RuntimeSettings) ROSClient {
	return &rosClient{mqtt, settings}
}

func (inst *rosClient) WhoIs(timeout int, targetGlobalID, requestUUID string) *runtime.CommandResponse {
	c := &runtime.Command{
		Key:  "command",
		Args: []string{"get", "objects"},
		Body: nil,
	}

	topicPublish := fmt.Sprintf("r/req/v1/cloud/%s/plain/command/%s/%s", inst.settings.GlobalID, targetGlobalID, requestUUID)
	// r/res/v1/cloud/RX-2/plain/command/R-1/req-uuid
	topicSub := fmt.Sprintf("r/res/v1/cloud/%s/plain/command/%s/%s", targetGlobalID, inst.settings.GlobalID, requestUUID)
	resp := inst.mqttclient.RequestResponse(timeout, topicPublish, topicSub, requestUUID, c)
	if resp != nil {
		fmt.Println(resp.AsString())
		return &runtime.CommandResponse{
			Error: resp.AsString(),
		}
	}

	return &runtime.CommandResponse{
		Error: "failed to get any repose",
	}

}

func toByte(body any) ([]byte, *runtime.CommandResponse) {
	marshal, err := json.Marshal(body)
	if err != nil {
		return nil, &runtime.CommandResponse{
			Error: fmt.Sprintf("marshal err: %v", err),
		}
	}

	return marshal, nil

}
