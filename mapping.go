package rxlib

import (
	"encoding/json"
	"github.com/NubeIO/rxlib/libs/convert"
	"github.com/NubeIO/rxlib/priority"
	"time"
)

/*
	+----------------------+
	|                      |
	|     	  ROS-1        | (either can be publisher/receiver)
	|                      |
	+----------------------+
			   |
			   |
			   |
			   |
			   |
	+------------------------+
	|     	  ROS-2          |
	+------------------------+


publisher;
If set as the publisher, this means that we are watching object and will post each value when an COV event occurs
eg; ROS-1 will publish each COV value to ROS-2

receiver;
Is on the receiving end of the COV event

dual;
Is when we want 2-way read/write; as in keeping each object in sync both ways

pull/push-interval;
COV will be disabled, and we will simply request the data on a CRON
eg; cloud to pull on 5 min interval from edge (cloud: pull-requester, edge: pull-responder)
eg; edge to push to cloud every 1min (cloud: push-receiver, edge: push-initiator)

Push; would be the mapping that is connected to the object we are monitoring to send the value on each CRON
eg; ROS-1 will push the value to ROS-2 every 15min

Pull;
eg; ROS-1 will request/pull the value to ROS-2 every 15min
*/

type MappingType string

const (
	MappingTypePublisher         MappingType = "publisher"
	MappingTypeReceiver          MappingType = "receiver"
	MappingTypeDual              MappingType = "dual"
	MappingIntervalPullRequester MappingType = "pull-requester" // needs a cron added
	MappingIntervalPullResponder MappingType = "pull-responder" // disable the subscription on the eventbus
	MappingIntervalPushInitiator MappingType = "push-initiator" // disable the subscription on the eventbus, needs a cron added
	MappingIntervalPushReceiver  MappingType = "push-receiver"
)

type Mapping struct {
	GlobalID          string           `json:"globalID"`
	LeaderNetworkUUID string           `json:"leaderNetworkUUID,omitempty"`
	LeaderMappingUUID string           `json:"leaderMappingUUID,omitempty"`
	LeaderDataType    priority.Type    `json:"leaderDataType,omitempty"` // make it easy for an Obj to decode in incoming data; eg string, map[], user
	TargetMappingUUID string           `json:"targetMappingUUID"`
	TargetNetworkUUID string           `json:"targetNetworkUUID"`
	TargetDataType    priority.Type    `json:"targetDataType,omitempty"`
	ResponseUUID      string           `json:"responseUUID"`
	Type              MappingType      `json:"type"`
	MappingInterval   *MappingInterval `json:"mappingInterval,omitempty"`
}
type MappingInterval struct {
	Disable         bool          `json:"disable"`
	DurationSeconds time.Duration `json:"durationSeconds"`
}

// MappingPayload is what is sent over the network (rest/mqtt)
type MappingPayload struct {
	LeaderMappingUUID string                 `json:"leaderUUID"`
	LeaderNetworkUUID string                 `json:"leaderNetworkUUID"`
	TargetMappingUUID string                 `json:"targetMappingUUID"`
	TargetNetworkUUID string                 `json:"targetNetworkUUID"`
	PriorityData      *priority.PriorityData `json:"priorityData,omitempty"`
	Response          MappingPayloadState    `json:"response,omitempty"`
}

type MappingPayloadState string

const (
	MappingPayloadStatusSent           MappingPayloadState = "sent"
	MappingPayloadStatusOffline        MappingPayloadState = "broker is offline"
	MappingPayloadStatusFailedToSend   MappingPayloadState = "failed to send"
	MappingPayloadStatusSending        MappingPayloadState = "sending"
	MappingResponseOk                  MappingPayloadState = "ok"
	MappingResponseInvalidPayload      MappingPayloadState = "invalid payload"
	MappingResponseMappingNotFound     MappingPayloadState = "mapping not found"
	MappingResponseObjectNoConnections MappingPayloadState = "mapping has no connections"
)

func NewMappingPayload(leaderMappingUUID, leaderNetworkUUID, targetMappingUUID, targetNetworkUUID string) *MappingPayload {
	return &MappingPayload{
		LeaderMappingUUID: leaderMappingUUID,
		LeaderNetworkUUID: leaderNetworkUUID,
		TargetMappingUUID: targetMappingUUID,
		TargetNetworkUUID: targetNetworkUUID,
	}
}

func (m *MappingPayload) GetMappingPayloadState() MappingPayloadState {
	return m.Response
}

func (m *MappingPayload) MappingPayloadStatusSending() {
	m.Response = MappingPayloadStatusSending
}

func (m *MappingPayload) MappingPayloadStatusSent() {
	m.Response = MappingPayloadStatusSent
}

func (m *MappingPayload) MappingResponseMappingNotFound() {
	m.Response = MappingResponseMappingNotFound
}

func (m *MappingPayload) MappingResponseInvalidPayload() {
	m.Response = MappingResponseInvalidPayload
}

func (m *MappingPayload) MappingResponseOk() {
	m.Response = MappingResponseOk
}

func (m *MappingPayload) MappingResponseObjectNoConnections() {
	m.Response = MappingResponseObjectNoConnections
}

func IsMappingStateIsOk(state MappingPayloadState) bool {
	if state == MappingResponseOk {
		return true
	}
	return false

}

func (m *MappingPayload) MappingStateIsOk() bool {
	if IsMappingStateIsOk(m.Response) {
		return true
	}
	return false

}

func (m *MappingPayload) GetLeaderMappingUUID() string {
	return m.LeaderMappingUUID
}

func (m *MappingPayload) GetTargetMappingUUID() string {
	return m.TargetMappingUUID
}

func (m *MappingPayload) AddData(data *priority.PriorityData) {
	m.PriorityData = data
}

func (m *MappingPayload) GetData() *priority.PriorityData {
	return m.PriorityData
}

func (m *MappingPayload) DataType() priority.Type {
	if m.PriorityData == nil {
		return ""
	}
	return m.PriorityData.DataType
}

func (m *MappingPayload) GetPriority() *priority.PriorityTable {
	if m.PriorityData == nil {
		return nil
	}
	return m.PriorityData.Priority
}

func (m *MappingPayload) HighestPriority() any {
	if m.PriorityData == nil {
		return nil
	}
	return m.PriorityData.HighestPriority
}

func (m *MappingPayload) HighestAsFloatPointer() *float64 {
	if m.PriorityData == nil {
		return nil
	}
	v := m.HighestPriority()
	if v == nil {
		return nil
	}
	if m.IsTypeNumber() {
		return convert.ConvertToFloatPtr(v)
	}
	return nil

}

func (m *MappingPayload) IsTypeNumber() bool {
	if m.DataType() == priority.TypeFloat || m.DataType() == priority.TypeInt {
		return true
	}
	return false
}

func UnmarshalMappingByte(resp []byte) (*Mapping, error) {
	payload := &Mapping{}
	err := json.Unmarshal(resp, &payload)
	return payload, err
}

func UnmarshalMappingArrayByte(resp []byte) ([]*Mapping, error) {
	var payload []*Mapping
	err := json.Unmarshal(resp, &payload)
	return payload, err
}

func UnmarshalMappingPayloadArrayByte(resp []byte) ([]*MappingPayload, error) {
	var payload []*MappingPayload
	err := json.Unmarshal(resp, &payload)
	return payload, err
}

func GetByTargetNetworkUUID(mappings []*MappingPayload) map[string][]*MappingPayload {
	result := make(map[string][]*MappingPayload)
	for _, mapping := range mappings {
		result[mapping.TargetNetworkUUID] = append(result[mapping.TargetNetworkUUID], mapping)
	}
	return result
}

func GetByLeaderNetworkUUID(mappings []*MappingPayload) map[string][]*MappingPayload {
	result := make(map[string][]*MappingPayload)
	for _, mapping := range mappings {
		result[mapping.LeaderNetworkUUID] = append(result[mapping.LeaderNetworkUUID], mapping)
	}
	return result
}

func MappingsToPayloads(mappings []*MappingPayload) []*Payload {
	var payloads []*Payload
	for _, mapping := range mappings {
		payloads = append(payloads, &Payload{MappingPayload: mapping})
	}
	return payloads
}

func GetMappingsFromPayloads(payloads []*Payload) []*MappingPayload {
	var out []*MappingPayload
	for _, payload := range payloads {
		out = append(out, payload.GetMappingPayload())
	}
	return out
}
