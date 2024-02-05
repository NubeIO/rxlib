package unitswrapper

import (
	"fmt"
	"testing"
)

func TestInitUnits(t *testing.T) {
	u := InitUnits(&EngineeringUnits{
		DecimalPlaces: 1,
		UnitCategory:  "temperature",
		Unit:          "C",
		UnitTo:        "F",
	})
	v := 22.5
	err := u.New(v)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	fmt.Printf("Converted %s \n", u.ChangeUnitAsSymbol())
	c, err := u.Conversion()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("Converted %f \n", c)
	//fmt.Printf("Converted %s\n", unit.AsSymbolWithDecimal(1))
}
