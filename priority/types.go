package priority

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
	return nil
}

func (iv IntValue) AsBool() *bool {
	return nil
}

func (iv IntValue) AsString() *string {
	return nil
}

type FloatValue struct {
	Value float64
}

func (fv FloatValue) GetValue() interface{} {
	return fv.Value
}

func (fv FloatValue) AsInt() *int {
	return nil
}

func (fv FloatValue) AsFloat() *float64 {
	return &fv.Value
}

func (fv FloatValue) AsBool() *bool {
	return nil
}

func (fv FloatValue) AsString() *string {
	return nil
}

type BoolValue struct {
	Value bool
}

func (bv BoolValue) AsInt() *int {
	return nil
}

func (bv BoolValue) AsFloat() *float64 {
	return nil
}

func (bv BoolValue) AsBool() *bool {
	return &bv.Value
}

func (bv BoolValue) AsString() *string {
	return nil
}

func (bv BoolValue) GetValue() interface{} {
	return bv.Value
}

type StringValue struct {
	Value string
}

func (sv StringValue) GetValue() interface{} {
	return sv.Value
}

func (sv StringValue) AsInt() *int {
	return nil
}

func (sv StringValue) AsFloat() *float64 {
	return nil
}

func (sv StringValue) AsBool() *bool {
	return nil
}

func (sv StringValue) AsString() *string {
	return &sv.Value
}
