package cron

import (
	"fmt"
	"testing"
	"time"
)

func Test_cronSchedule_Builder(t *testing.T) {
	var mins = 0
	var hours = 6
	var days = 1
	var month = 12
	schedule := NewBuilder()
	fmt.Println(mins, hours, days)
	c := &CronBuilder{
		Day:     "MON",
		Minutes: &mins,
		Hours:   &hours,
		//Days:  &days,
		//Months: &month,
	}
	s := schedule.Builder(c)
	fmt.Println(s, month, mins, hours, days)
}

func Test_cronSchedule_Parse(t *testing.T) {
	//schedule := NewBuilder().Hours(10).Minutes(0) // 10am
	//fmt.Println(schedule)

	schedule := NewBuilder().Days(1).MonthsRange(time.January, time.April).Hours(13, 14).MinutesInterval(10)
	fmt.Println(schedule)

	schedule = schedule.New().Day("MON").Hours(8).Minutes(0)
	fmt.Println(schedule)

	// every min
	schedule = schedule.New().Minutes(1)
	fmt.Println(schedule)

	// every day @ 6am
	schedule = schedule.New().Minutes(0).Hours(6)
	fmt.Println(schedule)

	// every wed
	schedule = schedule.New().Day("WED")
	fmt.Println(schedule)

	// At 06:00 on day-of-month 23 in June
	schedule = schedule.New().Minutes(0).Hours(6).Months(6).Days(23)
	fmt.Println(schedule)

}
