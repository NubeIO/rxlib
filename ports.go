package rxlib

import (
	"errors"
	"fmt"
	"github.com/NubeIO/rxlib/libs/nils"
	"github.com/NubeIO/rxlib/payload"
	"github.com/NubeIO/rxlib/priority"
	"time"
)

type Port struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	UUID     string `json:"uuid,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`

	Payload *payload.Payload `json:"payload"`

	Transformation        *priority.Transformations `json:"transformation,omitempty"`
	UsingTransformation   bool                      `json:"usingTransformation,omitempty"`
	TransformationApplied bool                      `json:"transformationApplied,omitempty"`
	OverrideApplied       bool                      `json:"overrideApplied,omitempty"`
	Direction             PortDirection             `json:"direction,omitempty"`           // input or output
	DataType              priority.Type             `json:"dataType,omitempty"`            // float, bool, string, any, json
	DisableSubscription   bool                      `json:"disableSubscription,omitempty"` // if set to true we will not set up connection as a subscriber; this would be used when a connection is used to maybe pull the data on an interval
	OnlyPublishOnCOV      bool                      `json:"onlyPublishOnCOV,omitempty"`
	PreviousValue         *priority.PreviousValue   `json:"previousValue,omitempty"`

	AllowMultipleConnections bool `json:"allowMultipleConnections,omitempty"`
	HasConnection            bool `json:"hasConnection"`
	// port position is where to show the order on the Obj and where to hide the port or not
	DefaultPosition   int  `json:"defaultPosition,omitempty"`
	Hide              bool `json:"hide,omitempty"`
	HiddenByDefault   bool `json:"hiddenByDefault,omitempty"`
	OverPositionValue int  `json:"overPositionValue,omitempty"`

	LastOk      *time.Time `json:"LastOk,omitempty"`
	OkMessage   string     `json:"okMessage,omitempty"`
	LastFail    *time.Time `json:"LastFail,omitempty"`
	FailMessage string     `json:"failMessage,omitempty"`

	OnMessage func(portID string, msg *payload.Payload) `json:"-"` // used for the evntbus

}

func (p *Port) GetID() string {
	return p.ID
}

func (p *Port) GetUUID() string {
	return p.UUID
}

func (p *Port) GetName() string {
	return p.Name
}

func (p *Port) SetHasConnection(state bool) {
	p.HasConnection = state
}

func (p *Port) GetPayload() *payload.Payload {
	if p == nil {
		return nil
	}
	return p.Payload
}

func (p *Port) GetPayloadValue() (value interface{}, isNil bool) {
	if p == nil || p.Payload == nil {
		return nil, true
	}

	switch p.GetDataType() {
	case priority.TypeFloat:
		v := p.GetValueFloatPointer()
		if v == nil {
			return nil, true
		} else {
			return nils.GetFloat64(v), false
		}
	case priority.TypeInt:
		value = p.GetValueIntPointer()
		if value == nil {
			return nil, true
		}
	case priority.TypeBool:
		value = p.GetValueBoolPointer()
		if value == nil {
			return nil, true
		}
	default:
		return nil, true
	}

	return value, false
}

func (p *Port) GetDataType() priority.Type {
	if p == nil {
		return ""
	}
	return p.DataType
}

func (p *Port) SetValueBool(v bool) {
	p.GetPayload().BoolValue = nils.ToBool(v)
}

func (p *Port) SetValueString(v string) {
	p.GetPayload().StringValue = nils.ToString(v)
}

func (p *Port) SetValueJSON(v string) {
	p.GetPayload().JsonValue = nils.ToString(v)
}

func (p *Port) SetValueInt(v int) {
	p.GetPayload().IntValue = nils.ToInt32(int32(v))
}

func (p *Port) SetDisplayValue(v string) {
	p.GetPayload().DisplayValue = nils.ToString(v)
}

func (p *Port) SetValueFloat(v float64) {
	p.GetPayload().FloatValue = nils.ToFloat64(v)
}

func (p *Port) SetValueFloatNil() {
	p.GetPayload().FloatValue = nil
}

func (p *Port) GetValueInt() int {
	return int(p.GetPayload().GetIntValue())
}

func (p *Port) GetValueIntPointer() *int {
	int32Value := p.GetPayload().IntValue
	if int32Value == nil {
		return nil
	}
	intValue := int(*int32Value)
	return &intValue
}

func (p *Port) GetValueBool() bool {
	return p.GetPayload().GetBoolValue()
}

func (p *Port) GetValueBoolPointer() *bool {
	return p.GetPayload().BoolValue
}

func (p *Port) GetValueFloat() float64 {
	return p.GetPayload().GetFloatValue()
}

func (p *Port) GetValueFloatPointer() *float64 {
	return p.GetPayload().FloatValue
}

func (p *Port) Release() {
	v, isNil := p.GetPayload().GetTransformationExistingValueFloat()
	if isNil {
		p.GetPayload().FloatValue = nil
	} else {
		p.GetPayload().FloatValue = nils.ToFloat64(v)
	}
	p.GetPayload().UnsetTransformationExistingValueFloat()
	p.OverrideApplied = false
}

func (p *Port) SetOverride(v interface{}) error {
	if p == nil {
		return errors.New("cannot override nil port")
	}
	dataType := p.GetDataType()
	if dataType == "" {
		return errors.New("data type was empty")
	}
	switch dataType {
	case priority.TypeFloat:
		if value, ok := v.(float64); ok {
			p.GetPayload().SetTransformationExistingValueFloat(nils.ToFloat64(p.GetPayload().GetFloatValue()))
			p.GetPayload().FloatValue = &value
			p.OverrideApplied = true
			return nil
		} else {
			return errors.New("failed to get a valid value")
		}
	case priority.TypeInt:
		if value, ok := v.(int); ok {
			int32Value := int32(value)
			p.GetPayload().IntValue = &int32Value
			p.OverrideApplied = true
			return nil
		} else {
			return errors.New("failed to get a valid value")
		}
	case priority.TypeBool:
		if value, ok := v.(bool); ok {
			p.GetPayload().BoolValue = &value
			p.OverrideApplied = true
			return nil
		} else {
			return errors.New("failed to get a valid value")
		}
	case priority.TypeString:
		if value, ok := v.(string); ok {
			p.GetPayload().StringValue = &value
			p.OverrideApplied = true
			return nil
		} else {
			return errors.New("failed to get a valid value")
		}
	case priority.TypeJSON:
		if value, ok := v.(string); ok {
			p.GetPayload().JsonValue = &value
			p.OverrideApplied = true
			return nil
		} else {
			return errors.New("failed to get a valid value")
		}
	}
	return errors.New(fmt.Sprintf("unknown data type: %s", dataType))
}

func (p *Port) GetHasConnection() bool {
	return p.HasConnection
}

func (p *Port) SetName(v string) string {
	p.Name = v
	return p.Name
}

func (p *Port) SetLastOk(message string) string {
	p.LastOk = nils.ToTime(time.Now())
	p.OkMessage = message
	return p.Name
}

func (p *Port) SetLastFail(message string) string {
	p.LastFail = nils.ToTime(time.Now())
	p.FailMessage = message
	return p.Name
}

func (p *Port) AddTransformation(v *priority.Transformations) error {
	p.Transformation = v
	return nil
}

func (p *Port) IsEnabled() bool {
	if p.Disabled {
		return false
	}
	return true
}

func (p *Port) IsDisabled() bool {
	if p.Disabled {
		return true
	}
	return false
}

func (p *Port) OnlyPublishCOV() bool {
	return p.OnlyPublishOnCOV
}

func (p *Port) Enable() {
	p.Disabled = false
}

func (p *Port) Disable() {
	p.Disabled = true
}

func (p *Port) SetDisableSubscription() {
	p.DisableSubscription = true
}

func (p *Port) SubscriptionDisabled() bool {
	if p.DisableSubscription {
		return true
	}
	return false
}

func (p *Port) SetPreviousValue(value any) {
	p.PreviousValue = &priority.PreviousValue{
		Value:     value,
		Timestamp: time.Now(),
	}
}

func (p *Port) GetPreviousValue() *priority.PreviousValue {
	return p.PreviousValue
}

type PortOpts struct {
	DefaultPosition          int  `json:"defaultPosition"`
	HiddenByDefault          bool `json:"hiddenByDefault,omitempty"`
	AllowMultipleConnections bool `json:"allowMultipleConnections,omitempty"`
}

func portOpts(opts ...*PortOpts) *PortOpts {
	p := &PortOpts{}
	if len(opts) > 0 {
		if opts[0] != nil {
			p.DefaultPosition = opts[0].DefaultPosition
			p.HiddenByDefault = opts[0].HiddenByDefault
			p.AllowMultipleConnections = opts[0].AllowMultipleConnections
		}
	}
	return p
}

func NewPortFloatCallBack(id string, f func(portID string, message *payload.Payload), opts ...*PortOpts) *NewPort {
	p := &NewPort{
		ID:        id,
		Name:      id,
		DataType:  priority.TypeFloat,
		OnMessage: f,
	}
	p.DefaultPosition = portOpts(opts...).DefaultPosition
	p.HiddenByDefault = portOpts(opts...).HiddenByDefault
	p.AllowMultipleConnections = portOpts(opts...).AllowMultipleConnections
	return p
}

func NewPortBoolCallBack(id string, f func(portID string, message *payload.Payload), opts ...*PortOpts) *NewPort {
	p := &NewPort{
		ID:        id,
		Name:      id,
		DataType:  priority.TypeBool,
		OnMessage: f,
	}
	p.DefaultPosition = portOpts(opts...).DefaultPosition
	p.HiddenByDefault = portOpts(opts...).HiddenByDefault
	p.AllowMultipleConnections = false
	return p
}

func NewPortAnyCallBack(id string, f func(portID string, message *payload.Payload), opts ...*PortOpts) *NewPort {
	p := &NewPort{
		ID:        id,
		Name:      id,
		DataType:  priority.TypeAny,
		OnMessage: f,
	}
	p.DefaultPosition = portOpts(opts...).DefaultPosition
	p.HiddenByDefault = portOpts(opts...).HiddenByDefault
	p.AllowMultipleConnections = portOpts(opts...).AllowMultipleConnections
	return p
}

func NewPortFloat(id string, opts ...*PortOpts) *NewPort {
	p := &NewPort{
		ID:       id,
		Name:     id,
		DataType: priority.TypeFloat,
	}
	p.DefaultPosition = portOpts(opts...).DefaultPosition
	p.HiddenByDefault = portOpts(opts...).HiddenByDefault
	p.AllowMultipleConnections = portOpts(opts...).AllowMultipleConnections
	return p
}

func NewPortBool(id string, opts ...*PortOpts) *NewPort {
	p := &NewPort{
		ID:       id,
		Name:     id,
		DataType: priority.TypeBool,
	}
	p.DefaultPosition = portOpts(opts...).DefaultPosition
	p.HiddenByDefault = portOpts(opts...).HiddenByDefault
	p.AllowMultipleConnections = portOpts(opts...).AllowMultipleConnections
	return p
}

func NewPortAny(id string, opts ...*PortOpts) *NewPort {
	p := &NewPort{
		ID:       id,
		Name:     id,
		DataType: priority.TypeAny,
	}
	p.DefaultPosition = portOpts(opts...).DefaultPosition
	p.HiddenByDefault = portOpts(opts...).HiddenByDefault
	p.AllowMultipleConnections = portOpts(opts...).AllowMultipleConnections
	return p
}

func NewPortDate(id string, opts ...*PortOpts) *NewPort {
	p := &NewPort{
		ID:       id,
		Name:     id,
		DataType: priority.TypeDate,
	}
	p.DefaultPosition = portOpts(opts...).DefaultPosition
	p.HiddenByDefault = portOpts(opts...).HiddenByDefault
	p.AllowMultipleConnections = portOpts(opts...).AllowMultipleConnections
	return p
}

func NewPortString(id string, opts ...*PortOpts) *NewPort {
	p := &NewPort{
		ID:       id,
		Name:     id,
		DataType: priority.TypeString,
	}
	p.DefaultPosition = portOpts(opts...).DefaultPosition
	p.HiddenByDefault = portOpts(opts...).HiddenByDefault
	p.AllowMultipleConnections = portOpts(opts...).AllowMultipleConnections
	return p
}

type NewPort struct {
	ID                       string
	Name                     string
	DataType                 priority.Type
	AllowMultipleConnections bool                                      `json:"allowMultipleConnections,omitempty"`
	DefaultPosition          int                                       `json:"defaultPosition"`
	HiddenByDefault          bool                                      `json:"hiddenByDefault,omitempty"`
	OnMessage                func(portID string, msg *payload.Payload) `json:"-"`
}

type Override struct {
	Value any `json:"value"`
}

type PortFormatString struct {
	ErrorOnMixMax    bool     `json:"errorOnMixMax"`
	MinLengthString  *int     `json:"minLengthString"`
	MaxLengthString  *int     `json:"maxLengthString"`
	AllowEmptyString bool     `json:"allowEmptyString,omitempty"`
	RestrictString   *float64 `json:"restrictString"` // for example don't allow # on an mqtt topic
}

// some commonly used output names
const (
	OutputName        string = "output"
	OutputStatusName  string = "status"
	OutputLastOk      string = "last ok"
	OutputLastUpdated string = "last updated"
	OutputErrorName   string = "err"
)

const (
	NameNot   string = "not"
	NameEqual string = "equal"
	NameGT    string = "gt"
	NameGTE   string = "gte"
	NameLT    string = "lt"
	NameLTE   string = "lte"
)

// some commonly used input names
const (
	InputName  string = "input"
	InName     string = "in"
	InZeroName string = "in-0"
	In1Name    string = "in-1"
	In2Name    string = "in-2"
	In3Name    string = "in-3"
	In4Name    string = "in-4"
	In5Name    string = "in-5"
	In6Name    string = "in-6"
	In7Name    string = "in-7"
	In8Name    string = "in-8"
	In9Name    string = "in-9"
	In10Name   string = "in-10"
	In11Name   string = "in-11"
	In12Name   string = "in-12"
	In13Name   string = "in-13"
	In14Name   string = "in-14"
	In15Name   string = "in-15"
	In16Name   string = "in-16"
	In17Name   string = "in-17"
	In18Name   string = "in-18"
	In19Name   string = "in-19"
	In20Name   string = "in-20"
)

type FlowDirection string

const (
	DirectionSubscriber      string = "subscriber"
	DirectionPublisher       string = "publisher"
	DirectionRequestResponse string = "request-response"
)

type PortDirection string

const (
	Input  PortDirection = "input"
	Output PortDirection = "output"
)
