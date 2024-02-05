package priority

import (
	"strconv"
)

// PriorityValue is an interface that all priority values must satisfy.
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
