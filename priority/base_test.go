package priority

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers/pprint"
	"testing"
)

func TestNewPriority(t *testing.T) {
	p := NewPriority()

	intVal := 10
	p.SetValue(IntValue{Value: intVal}, 1)
	p.SetValue(IntValue{Value: intVal + 10}, 2)
	// Get highest priority value
	highestValue, highestPriority := p.GetHighestPriorityValue()
	if highestValue != nil {
		fmt.Printf("Highest priority value: %v at priority %d\n", highestValue.GetValue(), highestPriority)
	} else {
		fmt.Println("No value set.")
	}
	p.SetNull(1)

	// Get highest priority value
	highestValue, highestPriority = p.GetHighestPriorityValue()
	if highestValue != nil {
		fmt.Printf("Highest priority value: %v at priority %d\n", highestValue.GetValue(), highestPriority)
	} else {
		fmt.Println("No value set.")
	}

	pprint.PrintJSON(p.ToMap())

}
