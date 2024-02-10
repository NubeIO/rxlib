package priority

import (
	"fmt"
	"github.com/NubeIO/rxlib/libs/nils"
	"github.com/NubeIO/rxlib/unitswrapper"
	"testing"
)

func isNil(v any) bool {
	if v == nil {
		return true
	}
	return false
}

func TestNewPriority(t *testing.T) {

	//nv := AnyValue{Value: "1000"}
	//pri := NewPriority(5, TypeFloat)
	//pri.SetValue(nv, 1)
	//
	//nv = AnyValue{Value: "22"}
	//pri.SetValue(nv, 2)
	//
	//v := pri.GetHighestPriorityValue().AsFloat()
	//
	//fmt.Println(nils.GetFloat64(v) + 11)

	value := nils.ToInt(24)
	ov := nils.ToInt(24)

	//cv := convert.IntPointerToFloatPointer(value)

	transformationConfig := &Transformations{
		//Enums:       enums,
		ApplyMinMax: false,
		//ApplyMinMax: true,
		MinMaxValue: &MinMaxValue{MaxOutValue: nils.ToFloat64(2)},
	}
	transformationConfig = nil
	u := &unitswrapper.EngineeringUnits{
		DecimalPlaces: 0,
		UnitCategory:  "temperature",
		Unit:          "C",
		UnitTo:        "F",
	}
	transformationConfig = nil
	ov = nil
	u = nil

	data := NewDataPriority(TypeFloat, transformationConfig, u, 0)
	apply, err := data.Apply(value, ov, TypeInt)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(apply.GetFloat())

}
