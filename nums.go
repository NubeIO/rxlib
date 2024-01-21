package rxlib

import (
	"errors"
	"fmt"
	"math"
)

func FormatFloat64(value float64, decimalPlace int) string {
	format := fmt.Sprintf("%%.%df", decimalPlace)
	return fmt.Sprintf(format, value)
}

type FormatNumber struct {
	FallBackValue *float64 `json:"fallBackValue"`
	// return the result value to a decimal place if its not nil
	Round *int `json:"round"`

	// limit the result based of the min/max settings
	MaxValue      *float64 `json:"maxValue"`
	MinValue      *float64 `json:"minValue"`
	ErrorOnMinMax bool

	// throw error if we have a match
	RestrictNumber *float64 `json:"restrictNumber"` // for example don't allow number 10

	// scale the result value based on the in min/max and the out min/max
	ApplyScale  bool
	ValueInMin  *float64
	ValueInMax  *float64
	ValueOutMin *float64
	ValueOutMax *float64
}

func FormatNumberBuilder(inputValue *float64, config FormatNumber) (*float64, error) {
	if config == (FormatNumber{}) {
		return nil, errors.New("config cannot be empty")
	}
	if inputValue == nil {
		if config.FallBackValue != nil {
			return config.FallBackValue, nil
		}
		return nil, nil
	}
	input := NewFloat64Ptr(inputValue)
	if config.MinValue != nil {
		var err error
		input, err = ApplyMinConstraint(input, NewFloat64Ptr(config.MinValue), config.ErrorOnMinMax)
		if err != nil {
			return nil, err
		}
	}
	if config.MaxValue != nil {
		var err error
		input, err = ApplyMaxConstraint(input, NewFloat64Ptr(config.MaxValue), config.ErrorOnMinMax)
		if err != nil {
			return nil, err
		}
	}
	if config.ApplyScale {
		if config.ValueInMin == nil || config.ValueInMax == nil || config.ValueOutMin == nil || config.ValueOutMax == nil {
			return nil, fmt.Errorf("to apply a scale we need all the format values to vaild")
		}
		input = Scale(input, NewFloat64Ptr(config.ValueInMin), NewFloat64Ptr(config.ValueInMax), NewFloat64Ptr(config.ValueOutMin), NewFloat64Ptr(config.ValueOutMax))
	}

	if config.Round != nil {
		input = ApplyDecimalPlace(input, NewIntPtr(config.Round))
	}

	if config.RestrictNumber != nil {
		if input == NewFloat64Ptr(config.RestrictNumber) {
			return nil, fmt.Errorf(" %f is a restrict number", input)
		}
	}
	return Float64Ptr(input), nil
}

func ApplyMinConstraint(input, minNum float64, errorOnMinMax bool) (float64, error) {
	if input < minNum {
		if errorOnMinMax {
			return 0, fmt.Errorf("number: %f  is less than: %f", input, minNum)
		}
		input = minNum
	}
	return input, nil
}

func ApplyMaxConstraint(input, maxNum float64, errorOnMinMax bool) (float64, error) {
	if input > maxNum {
		if errorOnMinMax {
			return 0, fmt.Errorf("number: %f is greater than: %f", input, maxNum)
		}
		input = maxNum
	}
	return input, nil
}

// Scale returns the (float64) input value (between inputMin and inputMax) scaled to a value between outputMin and outputMax
func Scale(value float64, inMin float64, inMax float64, outMin float64, outMax float64) float64 {
	scaled := ((value-inMin)/(inMax-inMin))*(outMax-outMin) + outMin
	if scaled > math.Max(outMin, outMax) {
		return math.Max(outMin, outMax)
	} else if scaled < math.Min(outMin, outMax) {
		return math.Min(outMin, outMax)
	} else {
		return scaled
	}
}

func ApplyDecimalPlace(input float64, decimalPlace int) float64 {
	if decimalPlace != 0 {
		decimalMultiplier := math.Pow(10, float64(decimalPlace))
		input = math.Round(input*decimalMultiplier) / decimalMultiplier
	}
	return input
}

func IntPtr(i int) *int {
	return &i
}

func NewIntPtr(f *int) int {
	if f != nil {
		return *f
	}
	return 0.0 // Return a default value (you can choose a different default if needed)
}

func NewFloat64Ptr(f *float64) float64 {
	if f != nil {
		return *f
	}
	return 0.0
}

func Float64Ptr(f float64) *float64 {
	return &f
}
