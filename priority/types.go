package priority

import (
	"fmt"
	"strconv"
)

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
const TypeString = "string"

type Priority struct {
	PriorityType Type
	Values       []PriorityValue
}

func (p *Priority) SetValue(value PriorityValue, priorityNumber int) {
	if priorityNumber >= 1 && priorityNumber <= len(p.Values) {
		p.Values[priorityNumber-1] = value
	}
}

func (p *Priority) SetNull(priorityNumber int) {
	if priorityNumber >= 1 && priorityNumber <= len(p.Values) {
		p.Values[priorityNumber-1] = nil
	}
}

func (p *Priority) GetHighestPriorityValue() (PriorityValue, int) {
	for i, v := range p.Values {
		if v != nil {
			return v, i + 1
		}
	}
	return nil, 0
}

func (p *Priority) GetLowestPriorityValue() (PriorityValue, int) {
	for i := len(p.Values) - 1; i >= 0; i-- {
		if p.Values[i] != nil {
			return p.Values[i], i + 1
		}
	}
	return nil, 0
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
			jsonMap[key] = val.GetValue()
		} else {
			jsonMap[key] = nil
		}
	}

	return jsonMap
}

// PriorityValue is an interface that all priority values must satisfy.
type PriorityValue interface {
	GetValue() interface{}
	AsInt() *int
	AsFloat() *float64
	AsBool() *bool
	AsString() *string
}

type IntValue struct {
	Value int
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
	Value float64
}

func (fv FloatValue) GetValue() interface{} {
	return fv.Value
}

func (fv FloatValue) AsInt() *int {
	intValue := int(fv.Value)
	return &intValue
}

func (fv FloatValue) AsFloat() *float64 {
	return &fv.Value
}

func (fv FloatValue) AsBool() *bool {
	boolValue := fv.Value != 0.0
	return &boolValue
}

func (fv FloatValue) AsString() *string {
	strValue := strconv.FormatFloat(fv.Value, 'f', -1, 64)
	return &strValue
}

type BoolValue struct {
	Value bool
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
	Value string
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
