package priority

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers/pprint"
	"testing"
)

func TestNewPriority(t *testing.T) {
	p := NewPriority(16, TypeFloat)

	p.SetValue(FloatValue{Value: 120}, 1)
	//p.SetValue(StringValue{Value: "2222"}, 2)
	p.SetValue(StringValue{Value: "sorry cant convert me"}, 3)
	highestValue, highestPriority := p.GetHighestPriorityValue()

	//intValue := ApplyMin(*highestValue.AsInt(), 20)
	//fmt.Printf("Min Int Value: %d\n", intValue)

	fmt.Println(*highestValue.AsInt(), "VVVVVVVVVVVV")

	//fmt.Println(n)

	// Get highest priority value
	if highestValue != nil {
		fmt.Printf("Highest priority value: %v at priority %d\n", highestValue.GetValue(), highestPriority)
	} else {
		fmt.Println("No value set.")
	}

	//p.SetNull(1)

	// Get highest priority value
	highestValue, highestPriority = p.GetHighestPriorityValue()
	if highestValue != nil {
		fmt.Printf("Highest priority value: %v at priority %d\n", highestValue.GetValue(), highestPriority)
	} else {
		fmt.Println("No value set.")
	}

	pprint.PrintJSON(p.ToMap())

}
