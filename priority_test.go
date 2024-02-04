package rxlib

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers/pprint"
	"github.com/NubeIO/rxlib/priority"
	"testing"
)

// Example test function
func TestPriority_GetHighestPriorityValue(t *testing.T) {
	port := Port{
		PriorityValue: priority.Priority{},
	}

	p := port.PriorityValue
	p.SetValue(priority.IntValue{2}, 1)
	p.SetValue(priority.FloatValue{20.5}, 2)
	p.SetValue(priority.BoolValue{true}, 3)
	p.SetValue(priority.StringValue{"test"}, 4)

	value, _ := p.GetHighestPriorityValue()
	if value == nil {
		t.Errorf("Expected a value, got nil")
	} else {
		fmt.Println("Highest priority value:", value.GetValue())
	}
	fmt.Println(*value.AsInt() + 10)

	pprint.PrintJSON(p.ToMap())
}
