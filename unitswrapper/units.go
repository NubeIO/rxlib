package unitswrapper

import (
	"errors"
	"fmt"
	units "github.com/NubeIO/engineering-units"
	"sync"
)

type EngineeringUnits struct {
	DecimalPlaces int    `json:"decimalPlaces"` // 2.345 to 2.3
	UnitCategory  string `json:"unitCategory"`
	Unit          string `json:"unit"`   // from temp °C, if the user just set this we can apply the unit
	UnitTo        string `json:"unitTo"` // to temp °F
	value         float64
	unitsLib      units.Units
	unitLib       units.Unit
	mux           sync.Mutex
}

// InitUnits for usage you then need to use the New()
func InitUnits(eu *EngineeringUnits) *EngineeringUnits {
	eu.unitsLib = units.New()
	return eu
}

func (eu *EngineeringUnits) New(input float64) error {
	eu.mux.Lock()
	defer eu.mux.Unlock()
	eu.value = input
	if eu.Unit == "" {
		return errors.New("unit can not be empty")
	}
	c, err := eu.unitsLib.Conversion(eu.UnitCategory, eu.Unit, input)
	if err != nil {
		return err
	}
	eu.unitLib = c
	err = eu.unitLib.CheckUnit(eu.Unit)
	if err != nil {
		return fmt.Errorf("unit err: %v", err)
	}
	if eu.UnitTo != "" {
		err := eu.unitLib.CheckUnit(eu.UnitTo)
		if err != nil {
			return fmt.Errorf("to unit err: %v", err)
		}
	}
	return err
}

// Conversion will do the conversion; eg temp=f to temp-c
func (eu *EngineeringUnits) Conversion() (float64, error) {
	eu.mux.Lock()
	defer eu.mux.Unlock()

	if eu.Unit == "" {
		return 0, errors.New("unit can not be empty")
	}
	if eu.unitLib == nil {
		return 0, errors.New("unitLib can not be empty")
	}
	return eu.unitLib.ChangeUnit(eu.UnitTo), nil
}

// AsSymbol will do no conversion; but return as a string with its symbol
func (eu *EngineeringUnits) AsSymbol() string {
	eu.mux.Lock()
	defer eu.mux.Unlock()

	if eu.unitLib == nil {
		return "error"
	}
	return eu.unitLib.AsSymbol()
}

// AsSymbolWithDecimal will do no conversion; but return as a string with its symbol
func (eu *EngineeringUnits) AsSymbolWithDecimal() string {
	eu.mux.Lock()
	defer eu.mux.Unlock()

	if eu.unitLib == nil {
		return "error"
	}
	return eu.unitLib.AsSymbolWithDecimal(eu.DecimalPlaces)
}

func (eu *EngineeringUnits) ChangeUnitAsSymbol() string {
	eu.mux.Lock()
	defer eu.mux.Unlock()

	if eu.unitLib == nil {
		return "error"
	}
	return eu.unitLib.ChangeUnitAsSymbol(eu.UnitTo, eu.DecimalPlaces)
}
