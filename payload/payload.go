package payload

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"math"
	"reflect"
)

type Body struct {
	PortID                string
	DataType              string
	IsNil                 bool
	IsInOverride          bool
	TransformationApplied bool
	Data                  any
}

type Payload struct {
	FromPortID     string
	FromObjectUUID string

	// store the original value before the transformation is set
	FromOverride                      bool
	OverrideFloat                     float64
	TransformationExistingValueFloat  *float64
	TransformationExistingValueInt    *int
	TransformationExistingValueString *string
	TransformationExistingValueBool   *bool
	*runtime.PortValue
	body *Body
}

func NewPayload(body *Body) (*Payload, error) {
	if body == nil {
		return nil, fmt.Errorf("body is nil")
	}
	dataType := body.DataType
	portID := body.PortID
	p := &Payload{
		PortValue: &runtime.PortValue{
			PortID:   portID,
			DataType: dataType,
		},
		body: body,
	}
	if body.Data != nil {
		return p.ApplyData(body.Data)
	}
	return p, nil

}

func (p *Payload) ApplyData(data any) (*Payload, error) {
	body := p.body
	if body.IsNil {
		return p, nil
	}
	var err error
	var byteData []byte
	dataType := p.DataType
	if dataType == "" {
		dataType = "json"
	}
	switch dataType {
	case "bool", "boolean":
		if data.(bool) {
			byteData = []byte{1}
		} else {
			byteData = []byte{0}
		}
	case "number", "float", "float64":
		bits := math.Float64bits(data.(float64))
		byteData = make([]byte, 8)
		binary.BigEndian.PutUint64(byteData, bits)
	case "string", "str":
		byteData = []byte(data.(string))
	default:
		byteData, err = json.Marshal(&data)
		if err != nil {
			return nil, err
		}
	}
	p.Data = byteData
	return p, nil
}

func (p *Payload) ToFloat() (float64, error) {
	if p.DataType != "float" || p.IsNil {
		return 0, fmt.Errorf("data is not float64 or is nil")
	}
	if len(p.Data) < 8 {
		return 0, fmt.Errorf("not enough bytes to unmarshal float64")
	}
	bits := binary.BigEndian.Uint64(p.Data)
	return math.Float64frombits(bits), nil
}

func (p *Payload) ToFloatPointer() (*float64, error) {
	if p.IsNil {
		return nil, nil
	}
	f, err := p.ToFloat()
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (p *Payload) ToBool() (bool, error) {
	if p.DataType != "bool" || p.IsNil {
		return false, fmt.Errorf("data is not bool or is nil")
	}
	return p.Data[0] == 1, nil
}

func (p *Payload) ToBoolPointer() (*bool, error) {
	if p.IsNil {
		return nil, nil
	}
	b, err := p.ToBool()
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (p *Payload) ToString() (string, error) {
	if p.DataType != "string" || p.IsNil {
		return "", fmt.Errorf("data is not string or is nil")
	}
	return string(p.Data), nil
}

func (p *Payload) ToStringPointer() (*string, error) {
	if p.IsNil {
		return nil, nil
	}
	s, err := p.ToString()
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (p *Payload) Unmarshal(target interface{}) error {
	if p.IsNil {
		return nil
	}
	kind := reflect.TypeOf(target).Elem().Kind()
	switch p.DataType {
	case "bool":
		if kind != reflect.Bool {
			return fmt.Errorf("target type mismatch: expected bool, got %s", kind)
		}
		*target.(*bool) = p.Data[0] == 1
	case "float64":
		if kind != reflect.Float64 {
			return fmt.Errorf("target type mismatch: expected float64, got %s", kind)
		}
		bits := binary.BigEndian.Uint64(p.Data)
		*target.(*float64) = math.Float64frombits(bits)
	case "string":
		if kind != reflect.String {
			return fmt.Errorf("target type mismatch: expected string, got %s", kind)
		}
		*target.(*string) = string(p.Data)
	case "json":
		err := json.Unmarshal(p.Data, target)
		if err != nil {
			return fmt.Errorf("failed to unmarshal json: %w", err)
		}
	default:
		return fmt.Errorf("unknown data type: %s", p.DataType)
	}

	return nil
}

func (p *Payload) SetTransformationExistingValueFloat(value *float64) {
	p.TransformationExistingValueFloat = value
}

func (p *Payload) UnsetTransformationExistingValueFloat() {
	p.TransformationExistingValueFloat = nil
}

func (p *Payload) GetTransformationExistingValueFloat() (value float64, isNil bool) {
	if p.TransformationExistingValueFloat == nil {
		return 0, true
	}
	return *p.TransformationExistingValueFloat, false
}

func (p *Payload) SetTransformationExistingValueInt(value *int) {
	p.TransformationExistingValueInt = value
}

func (p *Payload) UnsetTransformationExistingValueInt() {
	p.TransformationExistingValueInt = nil
}

func (p *Payload) SetTransformationExistingValueString(value *string) {
	p.TransformationExistingValueString = value
}

func (p *Payload) UnsetTransformationExistingValueString() {
	p.TransformationExistingValueString = nil
}

func (p *Payload) SetTransformationExistingValueBool(value *bool) {
	p.TransformationExistingValueBool = value
}

func (p *Payload) UnsetTransformationExistingValueBool() {
	p.TransformationExistingValueBool = nil
}
