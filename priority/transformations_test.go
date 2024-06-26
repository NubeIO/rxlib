package priority

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers/pprint"
	"github.com/NubeIO/rxlib/libs/nils"
	"testing"
)

func TestTransformationsBuilderApplyMinMax(t *testing.T) {
	transformationConfig := &Transformations{
		ApplyMinMax: true,
		MinMaxValue: &MinMaxValue{MaxOutValue: nils.ToFloat64(2)},
	}
	builder, err := TransformationsBuilder(nils.ToFloat64(1.1), transformationConfig)
	if err != nil {
		return
	}
	fmt.Println(nils.GetFloat64(builder))
}

func TestTransformationsBuilderEnums(t *testing.T) {
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
	transformationConfig := &Transformations{
		Enums: enums,
	}
	pprint.PrintJSON(transformationConfig)
	builder, err := TransformationsBuilder(nils.ToFloat64(1.1), transformationConfig)
	if err != nil {
		return
	}
	fmt.Println(nils.GetFloat64(builder))
}
