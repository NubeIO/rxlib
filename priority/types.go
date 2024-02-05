package priority

import (
	"fmt"
	"github.com/NubeIO/rxlib/libs/nils"
	"github.com/NubeIO/rxlib/unitswrapper"
	"strconv"
)

// Primitives for an output this would be the value we send out, this is also the input value (for the output/input apply transformations for output apply units)
type Primitives struct {
	priority          *Priority
	transformations   *Transformations
	units             *unitswrapper.EngineeringUnits
	transformedValue  *float64
	unitsValue        *float64
	symbol            *string
	priorityNumber    int
	inValue           *float64
	fallBackValue     *float64
	decimal           int
	inValueBool       *bool
	fallBackValueBool *bool
}

type NewPrimitiveValue struct {
	PriorityCount         int      `json:"priorityCount"`
	ValueType             Type     `json:"valueType"`
	InitialValue          *float64 `json:"initialValue"`
	FallBackValue         *float64 `json:"fallBackValue"`
	PriorityToWrite       int      `json:"priorityToWrite"`
	Decimal               int      `json:"decimal"`
	OverrideValue         *float64 `json:"overrideValue"`
	OverrideValuePriority int      `json:"overrideValuePriority"`

	Transformations *Transformations
	Units           *unitswrapper.EngineeringUnits
}

func NewPrimitive(body *NewPrimitiveValue) (*PrimitivesResult, error) {
	if body == nil {
		return nil, fmt.Errorf("body can not be empty")
	}
	if body.PriorityCount < 1 {
		body.PriorityCount = 2
	}
	if body.ValueType == "" {
		body.ValueType = TypeFloat
	}
	if body.PriorityToWrite < 1 {
		body.PriorityToWrite = 2
	}
	if body.PriorityToWrite < 1 {
		body.PriorityToWrite = 2
	}
	priorityInstance := NewPriority(body.PriorityCount, body.ValueType)
	p := &Primitives{
		priority:       priorityInstance,
		inValue:        body.InitialValue,
		priorityNumber: body.PriorityToWrite,
		fallBackValue:  body.FallBackValue,
		decimal:        body.Decimal,
	}
	var byPassUnits bool
	if body.Transformations != nil {
		p.addTransformations(body.Transformations)
		if body.Transformations.Enums == nil {
			byPassUnits = true
		}
	}
	if body.Units != nil {
		if body.Units.Unit != "" {
			if byPassUnits {
				u := unitswrapper.InitUnits(body.Units)
				p.addUnits(u)
			}
		}
	}
	return p.Result(body.OverrideValue, body.OverrideValuePriority)
}

func (p *Primitives) addTransformations(t *Transformations) {
	if t == nil {
		return
	}
	p.transformations = t
}

const (
	ErrTransformation    = "transformations not provided"
	ErrUnitsNotSupported = "units not provided"
	ErrUnitsEmptyValue   = "value is empty"
)

func (p *Primitives) applyTransformations() error {
	if p.transformations == nil {
		return fmt.Errorf(ErrTransformation)
	}
	transformationFormed, err := TransformationsBuilder(p.inValue, p.transformations)
	if err != nil {
		return err
	}
	p.transformedValue = transformationFormed
	return nil
}

func (p *Primitives) addUnits(u *unitswrapper.EngineeringUnits) {
	p.units = u
}

func (p *Primitives) applyUnits(applyConversion bool, overrideValue *float64) error {
	if p.units == nil {
		return fmt.Errorf(ErrUnitsNotSupported)
	}
	if p.units.Unit == "" {
		return fmt.Errorf(ErrUnitsNotSupported)
	}
	var v = p.inValue
	if p.transformations != nil {
		if p.transformedValue != nil {
			v = p.transformedValue // if transformations where applied then use them
		}
	}
	if overrideValue != nil {
		v = overrideValue
	}
	value := nils.GetFloat64(v)
	err := p.units.New(value)
	if err != nil {
		return err
	}

	if p.units.UnitTo != "" { // assume we do a conversion
		if applyConversion {
			fmt.Println(1111)
			converted, err := p.units.Conversion()
			if err != nil {
				return err
			}
			p.unitsValue = nils.ToFloat64(converted)
			p.symbol = nils.ToString(fmt.Sprintf("%s", overriden(converted, p.decimal, p.units.UnitTo)))
		} else {
			if overrideValue != nil {
				p.symbol = nils.ToString(fmt.Sprintf("%s (overridden)", overriden(value, p.decimal, p.units.UnitTo)))
			}
		}

	} else { // no conversion but apply symbol
		if overrideValue != nil {
			p.symbol = nils.ToString(fmt.Sprintf("%s (overridden)", overriden(value, p.decimal, p.units.Unit)))
		} else {
			p.symbol = nils.ToString(overriden(value, p.decimal, p.units.Unit))
		}

	}
	return nil
}

func overriden(value float64, decimalPlace int, unit string) string {
	format := fmt.Sprintf("%%.%df", decimalPlace)
	return fmt.Sprintf(format+" %s", value, unit)
}

type PrimitivesResult struct {
	Priority     *Priority `json:"priority,omitempty'"`
	RawValue     *float64  `json:"rawValue,omitempty"`
	Symbol       *string   `json:"symbol,omitempty"`
	RawValueBool *bool     `json:"rawValueBool,omitempty"`
}

func (p *Primitives) ResultBool() (*PrimitivesResult, error) {
	if p.inValueBool == nil { // pass on fallback
		if p.fallBackValueBool != nil {
			newValue := BoolValue{Value: nils.GetBool(p.fallBackValueBool)}
			p.priority.SetValue(newValue, p.priorityNumber)
			return &PrimitivesResult{
				Priority: p.priority,
			}, nil
		}
	}
	return nil, nil
}

func (p *Primitives) Result(overrideValue *float64, priorityNumber int) (*PrimitivesResult, error) {
	var err error
	var applyEnums bool
	if p.transformations != nil {
		if p.transformations.Enums != nil {
			applyEnums = true
		}
	}
	if overrideValue != nil { // override
		err = p.applyUnits(false, overrideValue)
		if err != nil {
			if err.Error() == ErrUnitsNotSupported {
			} else {
				return nil, err
			}
		}
		newValue := FloatValue{Value: nils.GetFloat64(overrideValue)}
		p.priority.SetValue(newValue, priorityNumber)
	}

	if p.inValue == nil { // pass on fallback
		if p.fallBackValue != nil {
			newValue := FloatValue{Value: nils.GetFloat64(p.fallBackValue)}
			p.priority.SetValue(newValue, p.priorityNumber)
			return &PrimitivesResult{
				Priority: p.priority,
			}, nil
		}
	}

	err = p.applyTransformations()
	if err != nil {
		if err.Error() == ErrTransformation {
		} else {
			return nil, err
		}
	}

	if p.transformedValue != nil || p.units != nil {
		if p.units != nil {

			if overrideValue == nil {
				err = p.applyUnits(true, overrideValue)
				if err != nil {
					if err.Error() == ErrUnitsNotSupported {
					} else {
						return nil, err
					}
				}
				newValue := FloatValue{Value: nils.GetFloat64(p.unitsValue)}
				p.priority.SetValue(newValue, p.priorityNumber)
			}
		} else {
			newValue := FloatValue{Value: nils.GetFloat64(p.transformedValue)}
			p.priority.SetValue(newValue, p.priorityNumber)
		}
	} else { // no units or transformations

		newValue := FloatValue{Value: nils.GetFloat64(p.inValue)}
		p.priority.SetValue(newValue, p.priorityNumber)
	}

	if applyEnums { // apply enums
		v, _ := p.priority.GetHighestPriorityValue()
		if v != nil {
			s, ok := EnumValue(nils.GetFloat64(v.AsFloat()), p.transformations.Enums)
			if ok {
				p.symbol = nils.ToString(s)
			}
		}
	}

	return &PrimitivesResult{
		Priority: p.priority,
		RawValue: p.inValue,
		Symbol:   p.symbol,
	}, nil
}

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
	PriorityType Type            `json:"priorityType"`
	Values       []PriorityValue `json:"values"`
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

func (p *Priority) GetHighestPriorityValueFallback(fallback PriorityValue) PriorityValue {
	for _, v := range p.Values {
		if v != nil {
			return v
		}
	}
	return fallback
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
