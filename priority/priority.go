package priority

import (
	"fmt"
	"github.com/NubeIO/rxlib/libs/convert"
	"github.com/NubeIO/rxlib/libs/nils"
	"github.com/NubeIO/rxlib/unitswrapper"
)

type PriorityTable struct {
	P1 any `json:"1"`
	P2 any `json:"2"`
}

type PriorityData struct {
	Priority        *PriorityTable `json:"priority,omitempty"`
	HighestPriority any            `json:"highestPriority,omitempty"`
	Symbol          *string        `json:"symbol,omitempty"`
	DataType        Type           `json:"dataType,omitempty"`
	RawValue        any            `json:"rawValue,omitempty"`
}

type Value struct {
	pri      *Priority
	rawValue any
	symbol   *string
	dataType Type
}

func (d *Value) IsNil() (PriorityValue, error) {
	if d == nil {
		return nil, fmt.Errorf("no vaild data to return")
	}
	if d.pri == nil {
		return nil, fmt.Errorf("no vaild data to return")
	}
	h := d.pri.GetHighestPriorityValue()
	if h == nil {
		return nil, fmt.Errorf("all values are empty")
	}
	return h, nil
}

func (d *Value) IsNull() bool {
	if d == nil {
		return true
	}
	if d.pri == nil {
		return true
	}
	h := d.pri.GetHighestPriorityValue()
	if h == nil {
		return true
	}
	return false
}

func (d *Value) PriorityData() *PriorityData {
	if d.pri == nil {
		return nil
	}
	var p1Value any
	var p2Value any
	p1 := d.pri.GetByPriorityNumber(1)
	if p1 != nil {
		p1Value = p1.GetValue()
	}
	p2 := d.pri.GetByPriorityNumber(2)
	if p2 != nil {
		p2Value = p2.GetValue()
	}

	return &PriorityData{
		Priority: &PriorityTable{
			P1: p1Value,
			P2: p2Value,
		},
		HighestPriority: d.GetHighestPriority(),
		Symbol:          d.GetSymbolPointer(),
		DataType:        d.GetType(),
		RawValue:        d.GetRawValue(),
	}
}

func (d *Value) GetPriority() *Priority {
	return d.pri
}

func (d *Value) GetPriorityValue() (PriorityValue, error) {
	highest, err := d.IsNil()
	if err != nil {
		return nil, err
	}
	return highest, nil
}

func (d *Value) GetHighestPriority() any {
	if d.pri == nil {
		return nil
	}
	v, _ := d.pri.GetHighestPriority()
	if v == nil {
		return nil
	}
	return v.GetValue()
}

func (d *Value) GetFloatErr() (float64, error) {
	highest, err := d.IsNil()
	if err != nil {
		return 0, err
	}
	return nils.GetFloat64(highest.AsFloat()), nil
}

func (d *Value) GetFloat() float64 {
	highest, err := d.IsNil()
	if err != nil {
		return 0
	}
	return nils.GetFloat64(highest.AsFloat())
}

func (d *Value) GetType() Type {
	return d.dataType
}

func (d *Value) IsTypeFloat() bool {
	if d.dataType == TypeFloat {
		return true
	}
	return false
}

func (d *Value) IsTypeNumber() bool {
	if d.dataType == TypeFloat || d.dataType == TypeInt {
		return true
	}
	return false
}

func (d *Value) GetFloatPointer() *float64 {
	highest, err := d.IsNil()
	if err != nil {
		return nil
	}
	if highest == nil {
		return nil
	}
	return highest.AsFloat()
}

func (d *Value) GetSymbolErr() (string, error) {
	_, err := d.IsNil()
	if err != nil {
		return "", err
	}
	s := d.symbol
	if s == nil {
		return "", fmt.Errorf("no vaild symbol")
	}
	return nils.GetString(s), nil
}

func (d *Value) GetSymbol() string {
	_, err := d.IsNil()
	if err != nil {
		return ""
	}
	s := d.symbol
	if s == nil {
		return ""
	}
	return nils.GetString(s)
}

func (d *Value) GetSymbolPointer() *string {
	_, err := d.IsNil()
	if err != nil {
		return nil
	}
	s := d.symbol
	if s == nil {
		return nil
	}
	return s
}

func (d *Value) GetRawValue() any {
	return d.rawValue
}

type DataPriority struct {
	transformation   *Transformations
	transformedValue *float64
	units            *unitswrapper.EngineeringUnits
	unitsValue       *float64
	symbol           *string
	decimal          int
	priority         *Priority
	dataType         Type
	out              *Value
}

func NewValuePriority(dataType Type, transformation *Transformations, units *unitswrapper.EngineeringUnits, decimal int) *DataPriority {
	if dataType == "" {
		dataType = TypeFloat
	}
	d := &DataPriority{
		priority: NewPriority(2, dataType),
		dataType: dataType,
		decimal:  decimal,
		out:      &Value{},
	}
	d.AddTransformation(transformation)
	d.AddUnits(units)
	return d
}

func (d *DataPriority) Apply(value, overrideValue any, fromDataType Type) (*Value, error) {
	if value == nil && overrideValue == nil { // release an override
		_, pri := d.priority.GetHighestPriority()
		if pri == 1 {
			d.priority.SetNull(1)
		}
		return d.out, nil
	}

	d.out.rawValue = value
	d.out.dataType = d.dataType
	if tryFloat(fromDataType, d.dataType) {
		cv := convert.AnyToFloatPointer(value)
		ov := convert.AnyToFloatPointer(overrideValue)
		err := d.applyTransformation(cv)
		if err != nil {
			return nil, err
		}
		err = d.applyUnits(cv, true, ov)
		if err != nil {
			return nil, err
		}
		var currentValue *float64
		if d.transformedValue != nil {
			currentValue = d.transformedValue
		}
		if d.unitsValue != nil {
			currentValue = d.unitsValue
		}
		if d.transformedValue == nil && d.unitsValue == nil {
			currentValue = cv
		}

		if ov != nil {
			f := FloatValue{Value: nils.GetFloat64(ov)}
			d.priority.SetValue(f, 1)
		}
		if currentValue != nil {
			f := FloatValue{Value: nils.GetFloat64(currentValue)}
			d.priority.SetValue(f, 2)
		}
		d.out.pri = d.priority
	} else if d.dataType == TypeString {
		f := StringValue{Value: fmt.Sprint(value)}
		d.priority.SetValue(f, 2)
		d.out.pri = d.priority
	} else if d.dataType == TypeDate {
		f := StringValue{Value: fmt.Sprint(value)}
		d.priority.SetValue(f, 2)
		d.out.pri = d.priority
	} else {
		f := AnyValue{Value: value}
		d.priority.SetValue(f, 2)
		d.out.pri = d.priority
	}
	return d.out, nil
}

func (d *DataPriority) AddTransformation(t *Transformations) {
	if d.transformation == nil {
		d.transformation = t
	}
}

func (d *DataPriority) applyTransformation(v *float64) error {
	if v == nil {
		return nil
	}
	if d.transformation == nil {
		return nil
	}
	out, err := TransformationsBuilder(v, d.transformation)
	if err != nil {
		return err
	}
	d.transformedValue = out
	return nil
}

func (d *DataPriority) AddUnits(u *unitswrapper.EngineeringUnits) {
	if d.units == nil {
		d.units = u
	}
}

func (d *DataPriority) applyUnits(v *float64, applyConversion bool, overrideValue *float64) error {
	if d.units == nil {
		return nil
	}
	if d.units.Unit == "" {
		return fmt.Errorf(ErrUnitsNotSupported)
	}

	if d.transformation != nil {
		if d.transformedValue != nil {
			v = d.transformedValue // if transformations where applied then use them
		}
	}

	if overrideValue != nil {
		v = overrideValue // apply override
	}
	value := nils.GetFloat64(v)
	err := d.units.New(value)
	if err != nil {
		return err
	}
	if d.units.UnitTo != "" { // assume we do a conversion
		if applyConversion {
			converted, err := d.units.Conversion()
			if err != nil {
				return err
			}
			d.unitsValue = nils.ToFloat64(converted)
			d.symbol = nils.ToString(fmt.Sprintf("%s", overriden(converted, d.decimal, d.units.UnitTo)))
		} else {
			if overrideValue != nil {
				d.symbol = nils.ToString(fmt.Sprintf("%s (overridden)", overriden(value, d.decimal, d.units.UnitTo)))
			}
		}

	} else { // no conversion but apply symbol
		if overrideValue != nil {
			d.symbol = nils.ToString(fmt.Sprintf("%s (overridden)", overriden(value, d.decimal, d.units.Unit)))
		} else {
			d.symbol = nils.ToString(overriden(value, d.decimal, d.units.Unit))
		}

	}
	d.out.symbol = d.symbol
	return nil
}

func tryFloat(fromType, toType Type) bool {
	var isFloat bool
	if toType == TypeFloat {
		isFloat = true
	}
	if fromType == TypeFloat && isFloat {
		return true
	}
	if fromType == TypeInt && isFloat {
		return true
	}
	if fromType == TypeBool && isFloat {
		return true
	}
	return false

}
