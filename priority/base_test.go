package priority

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers/pprint"
	"github.com/NubeIO/rxlib/libs/nils"
	"github.com/NubeIO/rxlib/unitswrapper"
	"testing"
)

func TestNewAsOutput(t *testing.T) {

	rawValue := nils.ToFloat64(100)
	fallback := nils.ToFloat64(0)
	enums := []*Enums{
		&Enums{
			Key:   0,
			Value: "off",
		},
		&Enums{
			Key:   1,
			Value: "on",
		},
		&Enums{
			Key:   2,
			Value: "manual",
		},
	}
	fmt.Println(enums[0].Key)
	transformationConfig := &Transformations{
		//Enums:       enums,
		ApplyMinMax: true,
		MinMaxValue: &MinMaxValue{MaxOutValue: nils.ToFloat64(2)},
	}
	u := &unitswrapper.EngineeringUnits{
		DecimalPlaces: 0,
		UnitCategory:  "temperature",
		Unit:          "",
		UnitTo:        "",
	}
	body := &NewPrimitiveValue{
		PriorityCount:   0,
		ValueType:       "",
		InitialValue:    rawValue,
		FallBackValue:   fallback,
		PriorityToWrite: 0,
		Decimal:         1,
		//OverrideValue:         Float64Ptr(300),
		OverrideValuePriority: 1,
		Transformations:       transformationConfig,
		Units:                 u,
	}
	resp, err := NewPrimitive(body)
	if err != nil {
		fmt.Println(err)
	}

	pprint.PrintJSON(resp)

}
