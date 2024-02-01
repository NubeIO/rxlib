package history

import "time"

func (h *GenericHistory) GetUUID() string {
	return h.UUID
}

func (h *GenericHistory) GetObjectUUID() string {
	return h.ObjectUUID
}

func (h *GenericHistory) GetRecords() []Record {
	return h.Values
}

func (h *GenericHistory) GetLast() Record {
	if len(h.Values) == 0 {
		return nil
	}
	return h.Values[len(h.Values)-1]
}

func (h *GenericHistory) GetFirst() Record {
	if len(h.Values) == 0 {
		return nil
	}
	return h.Values[0]
}

func (h *GenericHistory) GetPagination(pageNumber, pageSize int) []Record {
	if pageNumber <= 0 || pageSize <= 0 {
		return nil
	}

	startIndex := (pageNumber - 1) * pageSize
	if startIndex >= len(h.Values) {
		return nil
	}

	endIndex := startIndex + pageSize
	if endIndex > len(h.Values) {
		endIndex = len(h.Values)
	}

	return h.Values[startIndex:endIndex]
}

func (h *GenericHistory) GetRecordsByDateRange(startDate, endDate time.Time) []Record {
	result := make([]Record, 0)
	for _, sample := range h.Values {
		if sample.GetTimestamp().After(startDate) && sample.GetTimestamp().Before(endDate) {
			result = append(result, sample)
		}
	}
	return result
}

// GetRecordsByTime
// such as "300s", "-1.5h" or "2h45m".
// Valid time units are  "ms", "s", "m", "h".
// eg; GetRecordsByTime(time.Now(), "10m")
func (h *GenericHistory) GetRecordsByTime(startDate time.Time, duration string) ([]Record, error) {
	durationValue, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}

	endDate := startDate.Add(durationValue)
	result := make([]Record, 0)
	for _, sample := range h.Values {
		if sample.GetTimestamp().After(startDate) && sample.GetTimestamp().Before(endDate) {
			result = append(result, sample)
		}
	}
	return result, nil
}
