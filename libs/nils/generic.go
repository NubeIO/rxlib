package nils

import (
	"reflect"
	"strconv"
)

// To returns a pointer to the passed value.
func To[T any](t T) *T {
	return &t
}

// ToOrNil returns a pointer to the passed value, or nil, if the passed value is a zero value.
// If the passed value has `IsZero() bool` method (for example, time.Time instance),
// it is used to determine if the value is zero.
func ToOrNil[T comparable](t T) *T {
	if z, ok := any(t).(interface{ IsZero() bool }); ok {
		if z.IsZero() {
			return nil
		}
		return &t
	}

	var zero T
	if t == zero {
		return nil
	}
	return &t
}

// Get returns the value from the passed pointer or the zero value if the pointer is nil.
func Get[T any](t *T) T {
	if t == nil {
		var zero T
		return zero
	}
	return *t
}

func ToAny[O any](value interface{}) (okValue O, pointer *O, ok bool) {
	var output O

	if value == nil {
		return output, nil, false
	}

	inputValue := reflect.ValueOf(value)

	// Dereference if the value is a pointer
	if inputValue.Kind() == reflect.Ptr {
		if inputValue.IsNil() {
			return output, nil, false
		}
		inputValue = inputValue.Elem()
	}

	// Specific conversion logic for bool to string
	if inputValue.Kind() == reflect.Bool && reflect.TypeOf(output).Kind() == reflect.String {
		boolValue := inputValue.Bool()
		stringValue := strconv.FormatBool(boolValue)
		convertedValue := any(stringValue).(O)
		return convertedValue, &convertedValue, true
	}

	// Generic conversion for other types
	outputValue := reflect.ValueOf(&output).Elem()
	if inputValue.Type().ConvertibleTo(outputValue.Type()) {
		convertedValue := inputValue.Convert(outputValue.Type())
		output = convertedValue.Interface().(O)
		return output, &output, true
	}

	// If it's not convertible, return the zero value and false
	return output, nil, false
}
