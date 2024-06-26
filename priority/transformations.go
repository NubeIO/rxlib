package priority

import (
	"errors"
	"fmt"
	"github.com/NubeIO/rxlib/libs/nils"
	"github.com/NubeIO/rxlib/unitswrapper"
	"google.golang.org/protobuf/types/known/structpb"
	"math"
)

type Enums struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}

type Transformations struct {
	EnableTransformation bool     `json:"enableTransformation"`
	OverridePort         bool     `json:"overridePort"` // simply override the port value
	OverridePortValue    any      `json:"overridePortValue"`
	Enums                []*Enums `json:"enums"`
	ApplyEnum            bool     `json:"applyEnum"`
	FallBackValue        *float64 `json:"fallBackValue"`
	PermitNull           bool     `json:"permitNull"` // if true will set the default value of golang types;
	// return the result value to a decimal place if it's not nil
	Round *int `json:"round"`

	// limit the result based of the min/max settings
	ApplyMinMax   bool         `json:"applyMinMax"`
	MinMaxValue   *MinMaxValue `json:"minMaxValue"`
	ErrorOnMinMax bool         `json:"errorOnMinMax"`

	// throw error if we have a match
	RestrictNumber *float64 `json:"restrictNumber"` // for example don't allow number 10

	// scale the result value based on the in min/max and the out min/max
	ApplyScale       bool              `json:"applyScale"`
	ScaleMinMaxValue *ScaleMinMaxValue `json:"scaleMinMaxValue"`

	ApplyUnits bool                `json:"applyUnits"`
	Units      *unitswrapper.Units `json:"units"`
}

func (trans *Transformations) ApplyEngineeringUnits(v float64) (value float64, symbol string, err error) {
	if !trans.ApplyUnits {
		return 0, "", nil
	}
	u := unitswrapper.InitUnits(trans.Units)
	err = u.New(v)
	if err != nil {
		return 0, "", err
	}
	conversion, err := u.Conversion()
	if err != nil {
		return 0, "", err
	}
	return conversion, u.AsSymbol(), nil
}

type MinMaxValue struct {
	MinValue    *float64 `json:"minValue"`
	MaxValue    *float64 `json:"maxValue"`
	MinOutValue *float64 `json:"minOutValue"`
	MaxOutValue *float64 `json:"maxOutValue"`
}

type ScaleMinMaxValue struct {
	MinValue    *float64 `json:"minValue"`
	MaxValue    *float64 `json:"maxValue"`
	MinOutValue *float64 `json:"minOutValue"`
	MaxOutValue *float64 `json:"maxOutValue"`
}

func TransformationsBuilder(inputValue *float64, config *Transformations) (*float64, error) {
	if config == nil {
		return nil, errors.New("config cannot be empty")
	}
	if !config.EnableTransformation {
		return nil, nil
	}
	if inputValue == nil {
		if config.FallBackValue != nil {
			return config.FallBackValue, nil
		}
		return nil, nil
	}
	input := NewFloat64Ptr(inputValue)

	if config.RestrictNumber != nil {
		if input == NewFloat64Ptr(config.RestrictNumber) {
			return nil, fmt.Errorf(" %f is a restrict number", input)
		}
	}

	if config.ApplyScale && !config.ApplyMinMax {
		if config.ScaleMinMaxValue.MinValue == nil || config.ScaleMinMaxValue.MaxValue == nil || config.ScaleMinMaxValue.MinOutValue == nil || config.ScaleMinMaxValue.MaxOutValue == nil {
			return nil, fmt.Errorf("to apply a scale we need all the format values to vaild")
		}
		input = Scale(input, NewFloat64Ptr(config.ScaleMinMaxValue.MinValue), NewFloat64Ptr(config.ScaleMinMaxValue.MaxValue), NewFloat64Ptr(config.ScaleMinMaxValue.MinOutValue), NewFloat64Ptr(config.ScaleMinMaxValue.MaxOutValue))
	}
	if config.ApplyMinMax && !config.ApplyScale {
		// apply to input
		if config.MinMaxValue.MinValue != nil {
			var err error
			input, err = ApplyMinConstraint(input, NewFloat64Ptr(config.MinMaxValue.MinValue), config.ErrorOnMinMax)
			if err != nil {
				return nil, err
			}
		}
		if config.MinMaxValue.MaxValue != nil {
			var err error
			input, err = ApplyMaxConstraint(input, NewFloat64Ptr(config.MinMaxValue.MaxValue), config.ErrorOnMinMax)
			if err != nil {
				return nil, err
			}
		}

		// apply to out result
		if config.MinMaxValue.MinOutValue != nil {
			var err error
			input, err = ApplyMinConstraint(input, NewFloat64Ptr(config.MinMaxValue.MinOutValue), false)
			if err != nil {
				return nil, err
			}
		}
		if config.MinMaxValue.MaxOutValue != nil {
			var err error
			input, err = ApplyMaxConstraint(input, NewFloat64Ptr(config.MinMaxValue.MaxOutValue), false)
			if err != nil {
				return nil, err
			}
		}
	}

	if config.Round != nil {
		input = ApplyDecimalPlace(input, NewIntPtr(config.Round))
	}

	return Float64Ptr(input), nil
}

func EnumValue(v float64, enums []*Enums) (value string, ok bool) {
	for _, enum := range enums {
		if int(v) == enum.Key {
			if enum.Value != "" {
				return enum.Value, true
			}

		}
	}
	return "", false
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
	if decimalPlace > 0 {
		decimalMultiplier := math.Pow(10, float64(decimalPlace))
		input = math.Round(input*decimalMultiplier) / decimalMultiplier
	} else {
		input = math.Round(input)
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

// ToProtoStruct converts a Transformations instance to a google.protobuf.Struct
func ToProtoStruct(t *Transformations) (*structpb.Struct, error) {
	enumsList := make([]interface{}, len(t.Enums))
	for i, enum := range t.Enums {
		enumsList[i] = map[string]interface{}{
			"key":   enum.Key,
			"value": enum.Value,
		}
	}
	transformMap := map[string]interface{}{
		"enableTransformation": t.EnableTransformation,
		"overridePort":         t.OverridePort,
		"vverridePortValue":    t.OverridePortValue,
		"enums":                enumsList,
		"applyEnum":            t.ApplyEnum,
		"fallBackValue":        nils.GetFloat64(t.FallBackValue),
		"permitNull":           t.PermitNull,
		"round":                nils.GetInt(t.Round),
		"applyMinMax":          t.ApplyMinMax,
		"minMaxValue":          convertMinMaxValue(t.MinMaxValue),
		"errorOnMinMax":        t.ErrorOnMinMax,
		"restrictNumber":       nils.GetFloat64(t.RestrictNumber),
		"applyScale":           t.ApplyScale,
		"scaleMinMaxValue":     convertScaleMinMaxValue(t.ScaleMinMaxValue),
		"applyUnits":           t.ApplyUnits,
		"units":                convertEngineeringUnits(t.Units),
	}
	return structpb.NewStruct(transformMap)
}

// Helper function to convert MinMaxValue to map
func convertMinMaxValue(m *MinMaxValue) map[string]interface{} {
	if m == nil {
		return nil
	}
	return map[string]interface{}{
		"minValue":    nils.GetFloat64(m.MinValue),
		"maxValue":    nils.GetFloat64(m.MaxValue),
		"minOutValue": nils.GetFloat64(m.MinOutValue),
		"maxOutValue": nils.GetFloat64(m.MaxOutValue),
	}
}

// Helper function to convert ScaleMinMaxValue to map
func convertScaleMinMaxValue(s *ScaleMinMaxValue) map[string]interface{} {
	if s == nil {
		return nil
	}
	return map[string]interface{}{
		"minValue":    nils.GetFloat64(s.MinValue),
		"maxValue":    nils.GetFloat64(s.MaxValue),
		"minOutValue": nils.GetFloat64(s.MinOutValue),
		"maxOutValue": nils.GetFloat64(s.MaxOutValue),
	}
}

func convertEngineeringUnits(e *unitswrapper.Units) map[string]interface{} {
	return map[string]interface{}{
		"decimalPlaces": e.DecimalPlaces,
		"unitCategory":  e.UnitCategory,
		"unit":          e.Unit,
		"unitTo":        e.UnitTo,
	}
}
