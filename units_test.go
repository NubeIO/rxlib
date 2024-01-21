package rxlib

import (
	"fmt"
	"testing"
)

func TestInitUnits(t *testing.T) {
	u := InitUnits(&EngineeringUnits{
		DecimalPlaces: 1,
		UnitCategory:  "temperature",
		Unit:          "F",
		UnitTo:        "C",
	})
	v := 75.1
	err := u.New(v)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	fmt.Printf("AsSymbol %s \n", u.AsSymbolWithDecimal())
	c, err := u.Conversion()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("AsSymbol %f \n", c)
	//fmt.Printf("AsSymbol %s\n", unit.AsSymbolWithDecimal(1))
}
