package priority

import (
	"fmt"
	"github.com/NubeIO/rxlib/libs/convert"
	"github.com/NubeIO/rxlib/libs/nils"
	"github.com/NubeIO/rxlib/unitswrapper"
)

type Data struct {
	pri      *Priority
	rawValue any
	symbol   *string
}

func (d *Data) IsNil() (PriorityValue, error) {
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

func (d *Data) GetPriority() *Priority {
	return d.pri
}

func (d *Data) GetPriorityValue() (PriorityValue, error) {
	highest, err := d.IsNil()
	if err != nil {
		return nil, err
	}
	return highest, nil
}

func (d *Data) GetFloatErr() (float64, error) {
	highest, err := d.IsNil()
	if err != nil {
		return 0, err
	}
	return nils.GetFloat64(highest.AsFloat()), nil
}

func (d *Data) GetFloat() float64 {
	highest, err := d.IsNil()
	if err != nil {
		return 0
	}
	return nils.GetFloat64(highest.AsFloat())
}

func (d *Data) GetFloatPointer() *float64 {
	highest, err := d.IsNil()
	if err != nil {
		return nil
	}
	if highest == nil {
		return nil
	}
	return highest.AsFloat()
}

func (d *Data) GetSymbolErr() (string, error) {
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

func (d *Data) GetSymbol() string {
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

func (d *Data) GetSymbolPointer() *string {
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

func (d *Data) GetRawValue() any {
	return d.rawValue
}

type DataImpl struct {
	transformation   *Transformations
	transformedValue *float64
	units            *unitswrapper.EngineeringUnits
	unitsValue       *float64
	symbol           *string
	decimal          int
	priority         *Priority
	dataType         Type
	out              *Data
}

func NewDataPriority(dataType Type, transformation *Transformations, units *unitswrapper.EngineeringUnits, decimal int) *DataImpl {
	if dataType == "" {
		dataType = TypeFloat
	}
	d := &DataImpl{
		priority: NewPriority(2, dataType),
		dataType: dataType,
		decimal:  decimal,
		out:      &Data{},
	}
	d.addTransformation(transformation)
	d.addUnits(units)
	return d
}

func (d *DataImpl) Apply(value, overrideValue any, fromDataType Type) (*Data, error) {
	d.out.rawValue = value
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
		f := FloatValue{Value: nils.GetFloat64(currentValue)}
		if ov != nil {

			d.priority.SetValue(f, 1)
		} else {
			d.priority.SetValue(f, 2)
		}
		d.out.pri = d.priority

	}
	return d.out, nil
}

func (d *DataImpl) addTransformation(t *Transformations) {
	if d.transformation == nil {
		d.transformation = t
	}
}

func (d *DataImpl) applyTransformation(v *float64) error {
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

func (d *DataImpl) addUnits(u *unitswrapper.EngineeringUnits) {
	if d.units == nil {
		d.units = u
	}
}

func (d *DataImpl) applyUnits(v *float64, applyConversion bool, overrideValue *float64) error {
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
	if fromType == TypeString && isFloat {
		return true
	}
	if fromType == TypeBool && isFloat {
		return true
	}
	return false

}
