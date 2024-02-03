package rxlib

import (
	"fmt"
	"testing"
)

func TestFormatNumberFunc(t *testing.T) {

	// Example usage:
	inputValue := 6.0
	config := &Transformations{
		//Round:   intPtr(1), // Round to 2 decimal places
		//MaxValue:      float64Ptr(10.0),
		//MinValue:      float64Ptr(1.0),
		ErrorOnMinMax: false,
		ApplyScale:    false,
		//RestrictNumber: float64Ptr(3.0), // Example: Restrict to not allow 3.0
		ValueInMin:  Float64Ptr(0), // Scaling input range
		ValueInMax:  Float64Ptr(10.0),
		ValueOutMin: Float64Ptr(0.0), // Scaling output range
		ValueOutMax: Float64Ptr(100.0),
	}

	result, err := TransformationsBuilder(&inputValue, config)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Formatted Number: %v\n", NewFloat64Ptr(result))
	}

}
