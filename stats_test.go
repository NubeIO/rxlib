package rxlib

import (
	"fmt"
	"testing"
)

func TestNewStatsBuilder(t *testing.T) {
	builder := NewStatsBuilder()

	// Set the status and loop count using the builder.
	builder.SetStatus(StatsHalted).SetLoopCount(1)

	// Add custom stats using the builder.
	customStat1 := &CustomStatus{Name: "CustomStat1", Field: "Value1"}
	customStat2 := &CustomStatus{Name: "CustomStat2", Field: 42}
	builder.AddCustomStat("stat1", customStat1).AddCustomStat("stat2", customStat2)

	// Retrieve the object stats.
	stats := builder.GetStats()

	// Print the stats including custom stats.
	fmt.Printf("Object Status: %s\n", stats.ObjectStatus)
	fmt.Printf("Loop Count: %d\n", stats.LoopCount)

	// Print custom stats
	fmt.Println("Custom Stats:")
	for name, customStat := range stats.Custom {
		fmt.Printf("%s: %+v\n", name, customStat)
	}

	builder.IncrementLoopCount()
	builder.IncrementLoopCount()
	builder.IncrementLoopCount()
	builder.IncrementLoopCount()
	builder.IncrementLoopCount()
	builder.IncrementLoopCount()

	fmt.Printf("%+v\n", stats)

	fmt.Println(GenerateCustomUUID("node"))
}
