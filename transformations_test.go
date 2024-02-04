package rxlib

import (
	"fmt"
	"testing"
)

func TestTransformationsBuilder(t *testing.T) {

	transformationConfig := &Transformations{
		ApplyScale:     true,
		ValueInMin:     Float64Ptr(0),
		ValueInMax:     Float64Ptr(10),
		ValueOutMin:    Float64Ptr(0),
		ValueOutMax:    Float64Ptr(100),
		RestrictNumber: Float64Ptr(1.01),
	}

	inputValue := 1100.0

	// Apply the transformation
	result, err := TransformationsBuilder(&inputValue, transformationConfig)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Scaled Result: %f\n", *result)
	}

}
