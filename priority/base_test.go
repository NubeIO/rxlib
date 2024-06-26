package priority

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers/pprint"
	"github.com/NubeIO/rxlib/libs/nils"
	"github.com/NubeIO/rxlib/unitswrapper"
	"testing"
)

func TestNewAsOutput(t *testing.T) {

	rawValue := nils.ToFloat64(24)
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
		ApplyMinMax: false,
		MinMaxValue: &MinMaxValue{MaxOutValue: nils.ToFloat64(2)},
	}
	u := &unitswrapper.EngineeringUnits{
		DecimalPlaces: 0,
		UnitCategory:  "temperature",
		Unit:          "C",
		UnitTo:        "F",
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
	pprint.PrintJSON(transformationConfig)
	dataPri, prim, err := NewPrimitive(body)
	if err != nil {
		fmt.Println(err)
	}

	pprint.PrintJSON(dataPri)

	rawValue = nils.ToFloat64(20)
	result, err := prim.UpdateValueAndGenerateResult(rawValue, 2, nil, 0)
	if err != nil {
		return
	}

	rawValue = nils.ToFloat64(25)
	result, err = prim.UpdateValueFloat(25)
	if err != nil {
		fmt.Println(err)
		return
	}

	pprint.PrintJSON(result)

}

func TestNewAsOutputNull(t *testing.T) {

	body := &NewPrimitiveValue{}

	dataPri, _, err := NewPrimitive(body)
	if err != nil {
		fmt.Println(err)
	}
	a := dataPri.Priority.GetHighestPriorityValue().AsFloat()
	if a == nil {
		panic(111)
	}
	fmt.Println(nils.GetFloat64(a))
}
