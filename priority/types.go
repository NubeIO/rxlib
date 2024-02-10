package priority

import (
	"fmt"
	"github.com/NubeIO/rxlib/libs/convert"
	"strconv"
)

/*
	nv := FloatValue{Value: 11}
	pri := NewPriority(5, TypeFloat)
	pri.SetValue(nv, 1)

	nv = FloatValue{Value: 22}
	pri.SetValue(nv, 2)

	fmt.Println(pri.ToMap())
*/

func NewPriority(count int, valueType Type) *Priority {
	return &Priority{
		PriorityType: valueType,
		Values:       make([]PriorityValue, count),
	}
}

type Type string

const TypeBool = "bool"
const TypeInt = "int"
const TypeFloat = "float"
const TypeAny = "any"
const TypeString = "string"

type Priority struct {
	PriorityType Type            `json:"priorityType"`
	Values       []PriorityValue `json:"values"`
}

func (p *Priority) Count() int {
	return len(p.Values)
}

func (p *Priority) SetValue(value PriorityValue, priorityNumber int) {
	if priorityNumber >= 1 && priorityNumber <= len(p.Values) {
		p.Values[priorityNumber-1] = value
	}
}

func (p *Priority) SetValueFloat(value float64, priorityNumber int) {
	nv := FloatValue{Value: value}
	p.SetValue(nv, priorityNumber)
}

func (p *Priority) SetNull(priorityNumber int) {
	if priorityNumber >= 1 && priorityNumber <= len(p.Values) {
		p.Values[priorityNumber-1] = nil
	}
}

func (p *Priority) GetHighestPriorityValueFallback(fallback PriorityValue) PriorityValue {
	for _, v := range p.Values {
		if v != nil {
			return v
		}
	}
	return fallback
}

func (p *Priority) GetHighestPriority() (PriorityValue, int) {
	for i, v := range p.Values {
		if v != nil {
			return v, i + 1
		}
	}
	return nil, 0
}

func (p *Priority) GetHighestPriorityValue() PriorityValue {
	for _, v := range p.Values {
		if v != nil {
			return v
		}
	}
	return nil
}

func (p *Priority) GetLowestPriority() (PriorityValue, int) {
	for i := len(p.Values) - 1; i >= 0; i-- {
		if p.Values[i] != nil {
			return p.Values[i], i + 1
		}
	}
	return nil, 0
}

func (p *Priority) GetLowestPriorityValue() PriorityValue {
	for i := len(p.Values) - 1; i >= 0; i-- {
		if p.Values[i] != nil {
			return p.Values[i]
		}
	}
	return nil
}

func (p *Priority) GetByPriorityNumber(priorityNumber int) PriorityValue {
	if priorityNumber >= 1 && priorityNumber <= len(p.Values) {
		return p.Values[priorityNumber-1]
	}
	return nil
}

func (p *Priority) ToMap() map[string]interface{} {
	jsonMap := make(map[string]interface{})
	for i, val := range p.Values {
		key := fmt.Sprintf("p%d", i+1) // Keys like _1, _2, ..., _16
		if val != nil {
			switch v := val.(type) {
			case AnyValue:
				jsonMap[key] = v.GetValue() // Use the GetValue method of AnyValue
			default:
				jsonMap[key] = val.GetValue()
			}
		} else {
			jsonMap[key] = nil
		}
	}

	return jsonMap
}

// PriorityValue is an interface that all pri values must satisfy.
type PriorityValue interface {
	GetValue() interface{}
	AsInt() *int
	AsFloat() *float64
	AsBool() *bool
	AsString() *string
}

type IntValue struct {
	Value int `json:"value"`
}

func (iv IntValue) GetValue() interface{} {
	return iv.Value
}

func (iv IntValue) AsInt() *int {
	return &iv.Value
}

func (iv IntValue) AsFloat() *float64 {
	floatValue := float64(iv.Value)
	return &floatValue
}

func (iv IntValue) AsBool() *bool {
	boolValue := iv.Value != 0
	return &boolValue
}

func (iv IntValue) AsString() *string {
	strValue := strconv.Itoa(iv.Value)
	return &strValue
}

type FloatValue struct {
	Value float64 `json:"value"`
}

func (iv FloatValue) GetValue() interface{} {
	return iv.Value
}

func (iv FloatValue) AsInt() *int {
	intValue := int(iv.Value)
	return &intValue
}

func (iv FloatValue) AsFloat() *float64 {
	return &iv.Value
}

func (iv FloatValue) AsBool() *bool {
	boolValue := iv.Value != 0.0
	return &boolValue
}

func (iv FloatValue) AsString() *string {
	strValue := strconv.FormatFloat(iv.Value, 'f', -1, 64)
	return &strValue
}

type BoolValue struct {
	Value bool `json:"value"`
}

func (bv BoolValue) GetValue() interface{} {
	return bv.Value
}

func (bv BoolValue) AsInt() *int {
	intValue := 0
	if bv.Value {
		intValue = 1
	}
	return &intValue
}

func (bv BoolValue) AsFloat() *float64 {
	floatValue := 0.0
	if bv.Value {
		floatValue = 1.0
	}
	return &floatValue
}

func (bv BoolValue) AsBool() *bool {
	return &bv.Value
}

func (bv BoolValue) AsString() *string {
	strValue := strconv.FormatBool(bv.Value)
	return &strValue
}

type StringValue struct {
	Value string `json:"value"`
}

func (sv StringValue) GetValue() interface{} {
	return sv.Value
}

func (sv StringValue) AsInt() *int {
	intValue, err := strconv.Atoi(sv.Value)
	if err != nil {
		return nil
	}
	return &intValue
}

func (sv StringValue) AsFloat() *float64 {
	floatValue, err := strconv.ParseFloat(sv.Value, 64)
	if err != nil {
		return nil
	}
	return &floatValue
}

func (sv StringValue) AsBool() *bool {
	boolValue, err := strconv.ParseBool(sv.Value)
	if err != nil {
		return nil
	}
	return &boolValue
}

func (sv StringValue) AsString() *string {
	return &sv.Value
}

type AnyValue struct {
	Value interface{} `json:"value"`
}

func (av AnyValue) GetValue() interface{} {
	return av.Value
}

func (av AnyValue) AsInt() *int {
	return convert.AnyToIntPointer(av.Value)
}

func (av AnyValue) AsFloat() *float64 {
	return convert.AnyToFloatPointer(av.Value)
}

func (av AnyValue) AsBool() *bool {
	return convert.AnyToBoolPointer(av.Value)

}

func (av AnyValue) AsString() *string {
	if av.Value == nil {
		return nil
	}
	s := fmt.Sprint(av.Value)
	return &s
}
