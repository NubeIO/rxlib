package priority

import (
	"fmt"
	"github.com/NubeIO/rxlib/libs/nils"
	"github.com/NubeIO/rxlib/unitswrapper"
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
	FallBackValueBool     *bool    `json:"fallBackValueBool"`
	PriorityToWrite       int      `json:"priorityToWrite"`
	Decimal               int      `json:"decimal"`
	OverrideValue         *float64 `json:"overrideValue"`
	OverrideValueBool     *bool    `json:"overrideValueBool"`
	OverrideValuePriority int      `json:"overrideValuePriority"`

	Transformations *Transformations
	Units           *unitswrapper.EngineeringUnits
}

func NewPrimitive(body *NewPrimitiveValue) (*DataPriorityOld, *Primitives, error) {
	if body == nil {
		return nil, nil, fmt.Errorf("body can not be empty")
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

	if body.Transformations != nil {
		p.addTransformations(body.Transformations)
		if body.Transformations.Enums == nil {
		}
	}
	res, err := p.UpdateValueAndGenerateResult(body.InitialValue, body.PriorityToWrite, body.OverrideValue, body.OverrideValuePriority)
	return res, p, err
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

type DataPriorityOld struct {
	Priority     *Priority `json:"pri,omitempty'"`
	RawValue     *float64  `json:"rawValue,omitempty"`
	Symbol       *string   `json:"symbol,omitempty"`
	RawValueBool *bool     `json:"rawValueBool,omitempty"`
}

func (p *Primitives) UpdateValueFloat(newValue float64) (*DataPriorityOld, error) {
	return p.UpdateValueAndGenerateResult(nils.ToFloat64(newValue), 2, nil, 0)
}

func (p *Primitives) UpdateValueAndGenerateResult(newValue *float64, priorityNumber int, overrideValue *float64, overridePriorityNumber int) (*DataPriorityOld, error) {
	p.inValue = newValue              // Update the initial value
	p.priorityNumber = priorityNumber // Update the initial value

	var err error
	var applyEnums bool
	if p.transformations != nil && p.transformations.Enums != nil {
		applyEnums = true
	}

	if overrideValue != nil { // Override logic
		err = p.applyUnits(false, overrideValue)
		if err != nil && err.Error() != ErrUnitsNotSupported {
			return nil, err
		}
		nv := FloatValue{Value: nils.GetFloat64(overrideValue)}
		p.priority.SetValue(nv, overridePriorityNumber)
	} else if p.inValue == nil && p.fallBackValue != nil { // Fallback logic
		nv := FloatValue{Value: nils.GetFloat64(p.fallBackValue)}
		p.priority.SetValue(nv, p.priorityNumber)
	} else {
		// Handle transformations
		err = p.applyTransformations()
		if err != nil && err.Error() != ErrTransformation {
			return nil, err
		}

		// Handle units and transformation value application
		var valueToSet = p.inValue
		if p.transformedValue != nil || p.units != nil {
			if p.units != nil {
				err = p.applyUnits(true, overrideValue)
				if err != nil && err.Error() != ErrUnitsNotSupported {
					return nil, err
				}
				if p.unitsValue != nil {
					valueToSet = p.unitsValue // units
				}
			} else {
				if p.transformedValue != nil {
					valueToSet = p.transformedValue // units
				}
			}
		} else {
			valueToSet = p.inValue
		}
		nv := FloatValue{Value: nils.GetFloat64(valueToSet)}
		p.priority.SetValue(nv, p.priorityNumber)
	}

	// Enums application
	if applyEnums {
		v := p.priority.GetHighestPriorityValue()
		if v != nil {
			s, ok := EnumValue(nils.GetFloat64(v.AsFloat()), p.transformations.Enums)
			if ok {
				p.symbol = nils.ToString(s)
			}
		}
	}

	return &DataPriorityOld{
		Priority: p.priority,
		RawValue: p.inValue,
		Symbol:   p.symbol,
	}, nil
}
