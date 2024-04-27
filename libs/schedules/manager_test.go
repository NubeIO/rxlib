package schedules

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers/pprint"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	sm := New()

	sm.Add(&Schedule{
		Name:            "test",
		UseUTC:          false,
		SetTimezone:     "",
		setTimezone:     nil,
		DayToTimeRanges: nil,
		Exceptions:      nil,
	})

	weeklySchedule := NewSchedule("test", false, "Australia/Sydney")
	// Add time ranges for Sunday
	err := weeklySchedule.AddTimeRange(time.Sunday, "07:34", "11:00")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = weeklySchedule.AddTimeRange(time.Wednesday, "19:00", "19:35")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = weeklySchedule.AddException("2024-01-30 09:00:00", "2024-02-12 05:00:00")

	sm.Add(weeklySchedule)
	//now := time.Now().UTC()
	sch := sm.Get("test")

	pprint.PrintJSON(sch)

}

func TestExp(t *testing.T) {
	sm := New()

	weeklySchedule := NewSchedule("test", true, "Australia/Sydney")

	err := weeklySchedule.AddException("2024-01-30 09:00:00", "2024-02-12 05:00:00")
	if err != nil {
		fmt.Println("Error adding exception:", err)
		return
	}
	sm.Add(weeklySchedule)
	weeklySchedule.CheckCombined()
	sch := sm.Get("test").CheckCombined()

	fmt.Println(sch)

}

func TestParse(t *testing.T) {
	sm := New()

	s := sm.ParseFromString(schString)
	fmt.Println(s)
	sm.Add(s)

	run := sm.Get("test").CheckCombined()

	pprint.PrintJSON(run)

}

var schString = `
{
    "name": "test",
    "useUTC": false,
    "setTimezone": "Australia/Sydney",
    "dayToTimeRanges": {
        "0": [
            {
                "Start": "2024-04-28T07:34:00+10:00",
                "Stop": "2024-04-28T11:00:00+10:00"
            }
        ],
        "3": [
            {
                "Start": "2024-04-28T19:00:00+10:00",
                "Stop": "2024-04-28T19:35:00+10:00"
            }
        ]
    },
    "exceptions": {
        "2024-01-30T09:00:00Z": {
            "Start": "2024-01-30T09:00:00Z",
            "Stop": "2024-02-12T05:00:00Z"
        }
    }
}
`
