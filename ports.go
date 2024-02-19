package rxlib

import (
	"fmt"
	"github.com/NubeIO/rxlib/libs/nils"
	"github.com/NubeIO/rxlib/priority"
	"github.com/NubeIO/rxlib/unitswrapper"
	"time"
)

type Port struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	Disabled bool   `json:"disabled"`

	// Input/Output port values
	Values              *priority.Value        `json:"-"`    // value should be used for anything
	DataDisplay         *priority.PriorityData `json:"data"` // only used for when its called over rest
	dataPriority        *priority.DataPriority
	Transformation      *priority.Transformations      `json:"transformation"`
	Units               *unitswrapper.EngineeringUnits `json:"units"`
	Direction           PortDirection                  `json:"direction"`           // input or output
	DataType            priority.Type                  `json:"dataType"`            // float, bool, string, any, json
	DisableSubscription bool                           `json:"disableSubscription"` // if set to true we will not set up connection as a subscriber; this would be used when a connection is used to maybe pull the data on an interval
	OnlyPublishOnCOV    bool                           `json:"onlyPublishOnCOV"`
	PreviousValue       *priority.PreviousValue        `json:"previousValue,omitempty"`

	AllowMultipleConnections bool `json:"allowMultipleConnections,omitempty"`

	// port position is where to show the order on the Obj and where to hide the port or not
	DefaultPosition   int  `json:"defaultPosition"`
	Hide              bool `json:"hide,omitempty"`
	HiddenByDefault   bool `json:"hiddenByDefault,omitempty"`
	OverPositionValue int  `json:"overPositionValue,omitempty"`

	LastOk      *time.Time `json:"LastOk,omitempty"`
	OkMessage   string     `json:"okMessage,omitempty"`
	LastFail    *time.Time `json:"LastFail,omitempty"`
	FailMessage string     `json:"failMessage,omitempty"`

	OnMessage func(msg *Payload) `json:"-"` // used for the evntbus

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

func (p *Port) initPriority(dataType priority.Type, decimal int) {
	if p.dataPriority == nil {
		p.dataPriority = priority.NewValuePriority(dataType, nil, nil, decimal)
	}
}

func (p *Port) InitPriority(dataType priority.Type, decimal int) {
	p.initPriority(dataType, decimal)
}

func (p *Port) AddTransformation(v *priority.Transformations) error {
	if p.dataPriority == nil {
		return fmt.Errorf("data priority is empty")
	}
	p.dataPriority.AddTransformation(v)
	return nil

}

func (p *Port) AddUnits(v *unitswrapper.EngineeringUnits) error {
	if p.dataPriority == nil {
		return fmt.Errorf("data priority is empty")
	}
	p.dataPriority.AddUnits(v)
	return nil
}

func (p *Port) AddEnums(v []*priority.Enums) error {
	if p.dataPriority == nil {
		return fmt.Errorf("data priority is empty")
	}
	if p.Transformation == nil {
		return fmt.Errorf("transformation is empty")
	}
	p.Transformation.Enums = v
	return nil
}

func (p *Port) GetValue() *priority.Value {
	return p.Values
}

func (p *Port) GetHighestPriority() any {
	if p == nil {
		return nil
	}
	if p.Values == nil {
		return nil
	}
	return p.Values.GetHighestPriority()
}

func (p *Port) GetValueDisplay() *priority.PriorityData {
	if p == nil {
		return nil
	}
	if p.Values == nil {
		return nil
	}
	return p.Values.PriorityData()
}

func (p *Port) Write(value any) error {
	if p == nil {
		return fmt.Errorf("port is nil")
	}
	d, err := p.dataPriority.Apply(value, nil, p.GetDataType())
	p.Values = d
	return err
}

func (p *Port) WritePriority(value any, fromDataType priority.Type) error {
	d, err := p.dataPriority.Apply(value, nil, fromDataType)
	p.Values = d
	return err
}

func (p *Port) OverrideValue(value any) (*priority.Value, error) {
	d, err := p.dataPriority.Apply(nil, value, p.DataType)
	p.Values = d
	return p.Values, err
}

func (p *Port) ReleaseOverride() error {
	d, err := p.dataPriority.Apply(nil, nil, p.DataType)
	if err != nil {
		return err
	}
	p.Values = d
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

func (p *Port) GetDataType() priority.Type {
	if p == nil {
		return ""
	}
	return p.DataType
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

func NewPortFloatCallBack(id string, f func(message *Payload), opts ...*PortOpts) *NewPort {
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

func NewPortAnyCallBack(id string, f func(message *Payload), opts ...*PortOpts) *NewPort {
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
	AllowMultipleConnections bool               `json:"allowMultipleConnections,omitempty"`
	DefaultPosition          int                `json:"defaultPosition"`
	HiddenByDefault          bool               `json:"hiddenByDefault,omitempty"`
	OnMessage                func(msg *Payload) `json:"-"`
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

// some commonly used input names
const (
	InputName string = "input"
	In1Name   string = "in-1"
	In2Name   string = "in-2"
)

type FlowDirection string

const (
	DirectionSubscriber      FlowDirection = "subscriber"
	DirectionPublisher       FlowDirection = "publisher"
	DirectionRequestResponse FlowDirection = "request-response"
)

type PortDirection string

const (
	Input  PortDirection = "input"
	Output PortDirection = "output"
)
