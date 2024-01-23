package rxlib

import (
	"errors"
	"fmt"
	units "github.com/NubeIO/engineering-units"
	"sync"
)

type EngineeringUnits struct {
	UnitsLib      units.Units `json:"-"`
	DecimalPlaces int         `json:"decimalPlaces"` // 2.345 to 2.3
	UnitCategory  string      `json:"unitCategory"`
	Unit          string      `json:"unit"`   // from temp °C, if the user just set this we can apply the unit
	UnitTo        string      `json:"unitTo"` // to temp °F
	value         float64
	unitsInit     units.Unit
	mux           sync.Mutex
}

// InitUnits for usage you then need to use the New()
func InitUnits(eu *EngineeringUnits) *EngineeringUnits {
	eu.UnitsLib = units.New()
	return eu
}

func (eu *EngineeringUnits) New(input float64) error {
	eu.mux.Lock()
	defer eu.mux.Unlock()
	eu.value = input
	if eu.Unit == "" {
		return errors.New("unit can not be empty")
	}
	c, err := eu.UnitsLib.Conversion(eu.UnitCategory, eu.Unit, input)
	if eu.UnitTo != "" {
		err := c.CheckUnit(eu.UnitTo)
		if err != nil {
			return fmt.Errorf("to unit err: %v", err)
		}
	}
	eu.unitsInit = c
	return err
}

// Conversion will do the conversion; eg temp=f to temp-c
func (eu *EngineeringUnits) Conversion() (float64, error) {
	eu.mux.Lock()
	defer eu.mux.Unlock()

	if eu.Unit == "" {
		return 0, errors.New("unit can not be empty")
	}
	if eu.unitsInit == nil {
		return 0, errors.New("unitsInit can not be empty")
	}
	return eu.unitsInit.ChangeUnit(eu.UnitTo), nil
}

// AsSymbol will do no conversion; but return as a string with its symbol
func (eu *EngineeringUnits) AsSymbol() string {
	eu.mux.Lock()
	defer eu.mux.Unlock()

	if eu.unitsInit == nil {
		return "error"
	}
	return eu.unitsInit.AsSymbol()
}

// AsSymbolWithDecimal will do no conversion; but return as a string with its symbol
func (eu *EngineeringUnits) AsSymbolWithDecimal() string {
	eu.mux.Lock()
	defer eu.mux.Unlock()

	if eu.unitsInit == nil {
		return "error"
	}
	return eu.unitsInit.AsSymbolWithDecimal(eu.DecimalPlaces)
}

func (eu *EngineeringUnits) ChangeUnitAsSymbol() string {
	eu.mux.Lock()
	defer eu.mux.Unlock()

	if eu.unitsInit == nil {
		return "error"
	}
	return eu.unitsInit.ChangeUnitAsSymbol(eu.UnitTo, eu.DecimalPlaces)
}
