package weekly

type CombinedStats struct {
	WeeklyActive    bool              `json:"weeklyActive"`
	ExceptionActive bool              `json:"exceptionActive"`
	Weekly          []*ScheduleStatus `json:"weekly"`
	Exception       []*ScheduleStatus `json:"exception"`
}

func (ws *Schedule) CheckCombined() *CombinedStats {
	weekly := ws.CheckWeekly()
	exc := ws.CheckException()
	var weeklyActive bool
	for _, status := range weekly {
		if status.IsActive {
			weeklyActive = true
			break
		}
	}
	var exceptionActive bool
	for _, status := range exc {
		if status.IsActive {
			exceptionActive = true
			break
		}
	}

	var c = &CombinedStats{
		Weekly:          weekly,
		Exception:       exc,
		WeeklyActive:    weeklyActive,
		ExceptionActive: exceptionActive,
	}
	return c

}

func (ws *Schedule) CheckWeekly() []*ScheduleStatus {
	var out []*ScheduleStatus
	weekly := ws.checkWeekly()
	if weekly == nil {
		return nil
	}
	for _, status := range weekly {
		stats := ws.compareWeeklyStatus(status)
		out = append(out, stats)
	}
	return out

}

func (ws *Schedule) CheckException() []*ScheduleStatus {
	var out []*ScheduleStatus
	for _, timeRange := range ws.Exceptions {
		stats := ws.getExceptionStatus(timeRange)
		out = append(out, stats)
	}
	return out
}
