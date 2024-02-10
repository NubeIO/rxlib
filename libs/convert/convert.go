package convert

import (
	"strconv"
)

// AnyToBool converts any type to bool.
// It returns false if the conversion is not possible.
func AnyToBool(value interface{}) bool {
	switch v := value.(type) {
	case int:
		return v != 0
	case int8:
		return v != 0
	case int16:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	case uint:
		return v != 0
	case uint8:
		return v != 0
	case uint16:
		return v != 0
	case uint32:
		return v != 0
	case uint64:
		return v != 0
	case float32:
		return v != 0.0
	case float64:
		return v != 0.0
	case bool:
		return v
	case string:
		switch v {
		case "true", "1", "t", "T", "TRUE":
			return true
		default:
			return false
		}
	default:
		return false
	}
}

// AnyToBoolPointer converts any type to *bool.
// It returns nil if the conversion is not possible.
func AnyToBoolPointer(value interface{}) *bool {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case *bool:
		return v
	case *int:
		if v != nil {
			b := *v != 0
			return &b
		}
	case *float64:
		if v != nil {
			b := *v != 0.0
			return &b
		}
	case *string:
		if v != nil {
			if boolValue, err := strconv.ParseBool(*v); err == nil {
				return &boolValue
			}
		}
	}

	b := AnyToBool(value)
	return &b
}

// AnyToInt converts any type to int.
// It returns 0 if the conversion is not possible.
func AnyToInt(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return i
	default:
		return 0
	}
}

// AnyToIntPointer converts any type to *int.
// It returns nil if the conversion is not possible.
func AnyToIntPointer(value interface{}) *int {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case *int:
		return v
	case *float64:
		if v != nil {
			i := int(*v)
			return &i
		}
	case *bool:
		if v != nil {
			i := 0
			if *v {
				i = 1
			}
			return &i
		}
	case *string:
		if v != nil {
			if intValue, err := strconv.Atoi(*v); err == nil {
				return &intValue
			}
		}
	}

	i := AnyToInt(value)
	return &i
}

//-------------FLOAT----------------------

// AnyToFloat converts any type to float64.
// It returns 0.0 if the conversion is not possible.
func AnyToFloat(value interface{}) float64 {
	switch v := value.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case bool:
		if v {
			return 1.0
		}
		return 0.0
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0.0
		}
		return f
	default:
		return 0.0
	}
}

// AnyToFloatPointer converts any type to *float64.
// It returns nil if the conversion is not possible.
func AnyToFloatPointer(value interface{}) *float64 {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case *int:
		if v != nil {
			f := float64(*v)
			return &f
		}
	case *float64:
		return v
	case *bool:
		if v != nil && *v {
			f := 1.0
			return &f
		}
	case *string:
		if v != nil {
			if floatValue, err := strconv.ParseFloat(*v, 64); err == nil {
				return &floatValue
			}
		}
	}

	f := AnyToFloat(value)
	return &f
}

// StringToFloat converts a string to float64.
func StringToFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// BoolToFloat converts a bool to float64.
func BoolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

// IntToFloat converts an int to float64.
func IntToFloat(i int) float64 {
	return float64(i)
}

// StringPointerToFloatPointer converts a pointer to a string to *float64.
func StringPointerToFloatPointer(s *string) (*float64, error) {
	if s == nil {
		return nil, nil
	}
	f, err := strconv.ParseFloat(*s, 64)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// BoolPointerToFloatPointer converts a pointer to a bool to *float64.
func BoolPointerToFloatPointer(b *bool) *float64 {
	if b == nil {
		return nil
	}
	if *b {
		f := 1.0
		return &f
	}
	f := 0.0
	return &f
}

// IntPointerToFloatPointer converts a pointer to an int to *float64.
func IntPointerToFloatPointer(i *int) *float64 {
	if i == nil {
		return nil
	}
	f := float64(*i)
	return &f
}

// StringToInt converts a string to int.
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// BoolToInt converts a bool to int.
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// FloatToInt converts a float64 to int.
func FloatToInt(f float64) int {
	return int(f)
}

// StringPointerToInt converts a pointer to a string to *int.
func StringPointerToInt(s *string) (*int, error) {
	if s == nil {
		return nil, nil
	}
	i, err := strconv.Atoi(*s)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

// BoolPointerToIntPointer converts a pointer to a bool to *int.
func BoolPointerToIntPointer(b *bool) *int {
	if b == nil {
		return nil
	}
	if *b {
		i := 1
		return &i
	}
	i := 0
	return &i
}

// FloatPointerToIntPointer converts a pointer to a float64 to *int.
func FloatPointerToIntPointer(f *float64) *int {
	if f == nil {
		return nil
	}
	i := int(*f)
	return &i
}

// StringToBool converts a string to bool.
// It returns true for "true", "1", "t", "T", "TRUE", and false otherwise.
func StringToBool(s string) bool {
	switch s {
	case "true", "1", "t", "T", "TRUE":
		return true
	default:
		return false
	}
}

// IntToBool converts an int to bool.
// It returns true for non-zero values and false otherwise.
func IntToBool(i int) bool {
	return i != 0
}

// FloatToBool converts a float64 to bool.
// It returns true for non-zero values and false otherwise.
func FloatToBool(f float64) bool {
	return f != 0.0
}

// StringPointerToBoolPointer converts a pointer to a string to *bool.
// It returns true for "true", "1", "t", "T", "TRUE", and false otherwise.
func StringPointerToBoolPointer(s *string) *bool {
	if s == nil {
		return nil
	}
	b := StringToBool(*s)
	return &b
}

// IntPointerToBoolPointer converts a pointer to an int to *bool.
// It returns true for non-zero values and false otherwise.
func IntPointerToBoolPointer(i *int) *bool {
	if i == nil {
		return nil
	}
	b := IntToBool(*i)
	return &b
}

// FloatPointerToBoolPointer converts a pointer to a float64 to *bool.
// It returns true for non-zero values and false otherwise.
func FloatPointerToBoolPointer(f *float64) *bool {
	if f == nil {
		return nil
	}
	b := FloatToBool(*f)
	return &b
}

// IntToString converts an int to string.
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// BoolToString converts a bool to string.
func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

// FloatToString converts a float64 to string.
func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// IntPointerToStringPointer converts a pointer to an int to *string.
func IntPointerToStringPointer(i *int) *string {
	if i == nil {
		return nil
	}
	s := IntToString(*i)
	return &s
}

// BoolPointerToStringPointer converts a pointer to a bool to *string.
func BoolPointerToStringPointer(b *bool) *string {
	if b == nil {
		return nil
	}
	s := BoolToString(*b)
	return &s
}

// FloatPointerToStringPointer converts a pointer to a float64 to *string.
func FloatPointerToStringPointer(f *float64) *string {
	if f == nil {
		return nil
	}
	s := FloatToString(*f)
	return &s
}
