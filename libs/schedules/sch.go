package schedules

import (
	"time"
)

// TimeRange represents a start and stop time.
type TimeRange struct {
	Start time.Time
	Stop  time.Time
}

type Schedule struct {
	Name            string `json:"name"`
	UseUTC          bool   `json:"useUTC"`
	SetTimezone     string `json:"setTimezone"`
	setTimezone     *time.Location
	DayToTimeRanges map[time.Weekday][]TimeRange `json:"dayToTimeRanges"`
	Exceptions      map[time.Time]TimeRange      `json:"exceptions"`
}

// NewSchedule creates a new WeeklySchedule.
func NewSchedule(name string, useUTC bool, timezone string) *Schedule {
	var loc *time.Location
	if timezone != "" {
		loc, _ = time.LoadLocation(timezone)
	}
	return &Schedule{
		Name:            name,
		UseUTC:          useUTC,
		SetTimezone:     timezone,
		setTimezone:     loc,
		DayToTimeRanges: make(map[time.Weekday][]TimeRange),
		Exceptions:      make(map[time.Time]TimeRange),
	}
}

func (ws *Schedule) AddTimeRange(day time.Weekday, start, stop string) error {
	layout := "15:04"
	currentTime := time.Now()
	if ws.UseUTC {
		currentTime = currentTime.UTC()
	} else if ws.setTimezone != nil {
		currentTime = currentTime.In(ws.setTimezone)
	}

	startTime, err := time.Parse(layout, start)
	if err != nil {
		return err
	}
	stopTime, err := time.Parse(layout, stop)
	if err != nil {
		return err
	}

	// Set the date components for start and stop times
	startTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), startTime.Hour(), startTime.Minute(), 0, 0, ws.setTimezone)
	stopTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), stopTime.Hour(), stopTime.Minute(), 0, 0, ws.setTimezone)

	ws.DayToTimeRanges[day] = append(ws.DayToTimeRanges[day], TimeRange{Start: startTime, Stop: stopTime})
	return nil
}

func (ws *Schedule) compareWeeklyStatus(status *ScheduleStatus) *ScheduleStatus {
	currentTime := time.Now()
	if status.IsUTC {
		currentTime = currentTime.UTC()
	}

	// Adjust the start and end times to the current or next week
	for status.StartDate.Before(currentTime) && status.EndDate.Before(currentTime) {
		status.StartDate = status.StartDate.Add(7 * 24 * time.Hour)
		status.EndDate = status.EndDate.Add(7 * 24 * time.Hour)
	}

	// Determine if currently active and set next start and stop times
	if currentTime.Before(status.StartDate) {
		status.IsActive = false
		status.NextStart = status.StartDate
		status.NextStop = status.EndDate
	} else if currentTime.After(status.StartDate) && currentTime.Before(status.EndDate) {
		status.IsActive = true
		status.NextStart = status.StartDate.Add(7 * 24 * time.Hour) // Start time of the next week's range
		status.NextStop = status.EndDate                            // Stop time of the current range
	} else {
		status.IsActive = false
		status.NextStart = status.StartDate.Add(7 * 24 * time.Hour) // Start time of the next week's range
		status.NextStop = status.EndDate.Add(7 * 24 * time.Hour)    // Stop time of the next week's range
	}

	return status
}

type ScheduleStatus struct {
	IsUTC     bool
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool
	NextStart time.Time // show the next start compared to time.now and StartDate or if its IsUTC == true then time.now.utc()
	NextStop  time.Time // show the next stop compared to the EndDate and in UTC or not UTC
}

// checkWeekly returns the status of each time range in the schedules schedule.
func (ws *Schedule) checkWeekly() []*ScheduleStatus {

	currentTime := time.Now()

	if ws.UseUTC {
		currentTime = currentTime.UTC()
	} else if ws.setTimezone != nil {
		currentTime = currentTime.In(ws.setTimezone)
	}

	var weeklyStatuses []*ScheduleStatus
	for day, ranges := range ws.DayToTimeRanges {
		for _, tr := range ranges {
			startDate := ws.getAdjustedTime(tr.Start, day, currentTime)
			endDate := ws.getAdjustedTime(tr.Stop, day, currentTime)

			isActive := currentTime.After(startDate) && currentTime.Before(endDate)
			weeklyStatuses = append(weeklyStatuses, &ScheduleStatus{
				IsUTC:     ws.UseUTC,
				StartDate: startDate,
				EndDate:   endDate,
				IsActive:  isActive,
			})
		}
	}
	return weeklyStatuses
}

// getAdjustedTime adjusts the given time to the correct weekday and timezone.
func (ws *Schedule) getAdjustedTime(t time.Time, day time.Weekday, reference time.Time) time.Time {
	adjusted := time.Date(reference.Year(), reference.Month(), reference.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), reference.Location())
	// Calculate the difference in days and adjust the date
	dayDiff := int(day - reference.Weekday())
	if dayDiff < 0 {
		dayDiff += 7 // Ensure positive difference
	}
	return adjusted.AddDate(0, 0, dayDiff)
}

// getExceptionStatus returns the status of a specific exception.
func (ws *Schedule) getExceptionStatus(ex TimeRange) *ScheduleStatus {
	now := time.Now()

	if ws.UseUTC {
		now = now.UTC()
	} else if ws.setTimezone != nil {
		now = now.In(ws.setTimezone)
	}

	isActive := now.After(ex.Start) && now.Before(ex.Stop)
	nextStart := ex.Start
	nextStop := ex.Stop

	// If the exception is in the future, next start and stop are the start and stop of the exception.
	// If the exception is currently active, next start and stop are the end of the exception.
	// If the exception is past, set next start and stop to zero time (as there's no next active period for a past exception).
	if now.After(ex.Stop) {
		nextStart = time.Time{}
		nextStop = time.Time{}
	}

	return &ScheduleStatus{
		StartDate: ex.Start,
		EndDate:   ex.Stop,
		IsActive:  isActive,
		NextStart: nextStart,
		NextStop:  nextStop,
	}
}

// AddException adds an exception time range to the schedule and returns the status of all exceptions.
func (ws *Schedule) AddException(start, stop string) error {
	startTime, err := time.Parse("2006-01-02 15:04:05", start)
	if err != nil {
		return err
	}
	stopTime, err := time.Parse("2006-01-02 15:04:05", stop)
	if err != nil {
		return err
	}
	ws.Exceptions[startTime] = TimeRange{Start: startTime, Stop: stopTime}
	return nil
}
