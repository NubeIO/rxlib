package history

import "time"

func (h *GenericHistory) DeleteSample(sample Record) {
	for i, s := range h.Values {
		if s.GetUUID() == sample.GetUUID() {
			h.Values = append(h.Values[:i], h.Values[i+1:]...)
			return
		}
	}
}

func (h *GenericHistory) DeleteSamples(uuids []string) {
	uuidSet := make(map[string]bool)
	for _, uuid := range uuids {
		uuidSet[uuid] = true
	}

	newValues := make([]Record, 0, len(h.Values))
	for _, sample := range h.Values {
		if !uuidSet[sample.GetUUID()] {
			newValues = append(newValues, sample)
		}
	}

	h.Values = newValues
}

func (h *GenericHistory) DeleteFirst(count int) int {
	if count <= 0 {
		return 0
	}

	if count >= len(h.Values) {
		count = len(h.Values)
	}
	h.Values = h.Values[count:]
	return count
}

func (h *GenericHistory) DeleteLast(count int) int {
	if count <= 0 {
		return 0
	}

	if count >= len(h.Values) {
		count = len(h.Values)
	}
	h.Values = h.Values[:len(h.Values)-count]
	return count
}

func (h *GenericHistory) DeleteByDateRange(startDate, endDate time.Time) int {
	result := make([]Record, 0)
	for _, sample := range h.Values {
		if !sample.GetTimestamp().After(startDate) || !sample.GetTimestamp().Before(endDate) {
			result = append(result, sample)
		}
	}

	deletedCount := len(h.Values) - len(result)
	h.Values = result
	return deletedCount
}

func (h *GenericHistory) DeleteByTime(startDate time.Time, duration string) int {
	durationValue, err := time.ParseDuration(duration)
	if err != nil {
		return 0
	}

	endDate := startDate.Add(durationValue)
	result := make([]Record, 0)
	for _, sample := range h.Values {
		if !sample.GetTimestamp().After(startDate) || !sample.GetTimestamp().Before(endDate) {
			result = append(result, sample)
		}
	}

	deletedCount := len(h.Values) - len(result)
	h.Values = result
	return deletedCount
}
