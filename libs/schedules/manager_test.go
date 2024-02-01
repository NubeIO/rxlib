package schedules

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	sm := New()

	location, err := time.LoadLocation("Australia/Sydney")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(location)

	weeklySchedule := NewSchedule(true, location)
	// Add time ranges for Sunday
	err = weeklySchedule.AddTimeRange(time.Tuesday, "04:34", "11:00")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = weeklySchedule.AddTimeRange(time.Wednesday, "19:00", "19:35")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	name := "new"
	sm.Add(name, weeklySchedule)
	//now := time.Now().UTC()
	sch := sm.Get(name).CheckCombined()

	fmt.Println(sch)

}

func TestExp(t *testing.T) {
	sm := New()

	location, err := time.LoadLocation("Australia/Sydney")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(location)

	weeklySchedule := NewSchedule(true, location)

	err = weeklySchedule.AddException("2024-01-30 09:00:00", "2024-02-12 05:00:00")
	if err != nil {
		fmt.Println("Error adding exception:", err)
		return
	}
	name := "new"
	sm.Add(name, weeklySchedule)
	weeklySchedule.CheckCombined()
	sch := sm.Get(name).CheckCombined()

	fmt.Println(sch)

}
