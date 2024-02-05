package nils

import (
	"fmt"
	"testing"
)

func TestToGet(t *testing.T) {

	floatVal, floatPtr, ok := ToAny[string]("1")
	fmt.Println(floatVal, floatPtr, ok)
}

func TestToAnyIntToFloat64(t *testing.T) {
	var intValue int = 42
	expectedFloat := float64(42)

	floatVal, floatPtr, ok := ToAny[float64](intValue)
	if !ok {
		t.Fatalf("Expected successful conversion, got failure")
	}
	if floatVal != expectedFloat {
		t.Errorf("Expected %v, got %v", expectedFloat, floatVal)
	}
	if floatPtr == nil || *floatPtr != expectedFloat {
		t.Errorf("Expected pointer to %v, got %v", expectedFloat, floatPtr)
	}
}

func TestToAnyBoolToString(t *testing.T) {
	var boolValue bool = true
	expectedString := "true"

	strVal, strPtr, ok := ToAny[string](boolValue)
	if !ok {
		t.Fatalf("Expected successful conversion, got failure")
	}
	if strVal != expectedString {
		t.Errorf("Expected %v, got %v", expectedString, strVal)
	}
	if strPtr == nil || *strPtr != expectedString {
		t.Errorf("Expected pointer to %v, got %v", expectedString, strPtr)
	}
}

func TestToAnyBoolPtrToInt(t *testing.T) {
	boolValue := true
	boolPtr := &boolValue
	expectedInt := 1 // Assuming true converts to 1, and false to 0

	intVal, intPtr, ok := ToAny[int](boolPtr)
	if !ok {
		t.Fatalf("Expected successful conversion, got failure")
	}
	if intVal != expectedInt {
		t.Errorf("Expected %v, got %v", expectedInt, intVal)
	}
	if intPtr == nil || *intPtr != expectedInt {
		t.Errorf("Expected pointer to %v, got %v", expectedInt, intPtr)
	}

	// Test with a nil pointer
	var nilBoolPtr *bool
	intVal, intPtr, ok = ToAny[int](nilBoolPtr)
	if ok {
		t.Errorf("Expected conversion to fail for nil pointer, but it succeeded")
	}
}

func TestToAnyIntPtrToFloat64(t *testing.T) {
	intValue := 42
	intPtr := &intValue
	expectedFloat := float64(42)

	floatVal, floatPtr, ok := ToAny[float64](intPtr)
	if !ok {
		t.Fatalf("Expected successful conversion, got failure")
	}
	if floatVal != expectedFloat {
		t.Errorf("Expected %v, got %v", expectedFloat, floatVal)
	}
	if floatPtr == nil || *floatPtr != expectedFloat {
		t.Errorf("Expected pointer to %v, got %v", expectedFloat, floatPtr)
	}

	var nilIntPtr *int
	floatVal, floatPtr, ok = ToAny[float64](nilIntPtr)
	if ok {
		t.Errorf("Expected conversion to fail for nil pointer, but it succeeded")
	}
}
