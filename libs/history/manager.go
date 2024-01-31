package history

import (
	"sync"
	"time"
)

type Manager interface {
	// NewHistory creates a new history with the specified limitSampleSize.
	NewHistory(limitSampleSize int) History

	// Get retrieves a history by its UUID.
	Get(uuid string) History

	// All returns a slice of all available histories.
	All() []History

	// AllSamples returns a map of all samples across all histories.
	AllSamples() map[string][]Record

	// Drop deletes a history by its UUID.
	Drop(uuid string)

	// DropAll will delete all the histories
	DropAll()

	// GetPagination retrieves paginated samples for all histories managed by the Manager.
	GetPagination(pageNumber, pageSize int) map[string][]Record

	// GetSamplesByDateRange retrieves samples within a specified date range for all histories managed by the Manager.
	GetSamplesByDateRange(startDate, endDate time.Time) map[string][]Record

	// GetSamplesByTime retrieves samples within a specified time duration for all histories managed by the Manager.
	// It takes a startDate and duration string (e.g., "10s", "1h") as input and returns a map where keys are UUIDs
	GetSamplesByTime(startDate time.Time, duration string) (map[string][]Record, error)

	// DeleteSamples deletes samples from specified histories based on UUIDs for all histories managed by the Manager.
	DeleteSamples(uuids map[string]string)
}

type historyManager struct {
	histories map[string]History
	mu        sync.RWMutex
}

func NewHistoryManager() Manager {
	return &historyManager{
		histories: make(map[string]History),
	}
}

func (hm *historyManager) NewHistory(limitSampleSize int) History {
	history := NewGenericHistory(limitSampleSize)
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hm.histories[history.GetUUID()] = history
	return history
}

func (hm *historyManager) Get(uuid string) History {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	return hm.histories[uuid]
}

func (hm *historyManager) All() []History {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	histories := make([]History, 0, len(hm.histories))
	for _, history := range hm.histories {
		histories = append(histories, history)
	}
	return histories
}

func (hm *historyManager) AllSamples() map[string][]Record {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	samples := make(map[string][]Record)
	for uuid, history := range hm.histories {
		samples[uuid] = history.GetSamples()
	}
	return samples
}

func (hm *historyManager) Drop(uuid string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	delete(hm.histories, uuid)
}

func (hm *historyManager) DropAll() {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hm.histories = make(map[string]History)
}

func (hm *historyManager) GetPagination(pageNumber, pageSize int) map[string][]Record {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	pagination := make(map[string][]Record)

	for uuid, history := range hm.histories {
		samples := history.GetPagination(pageNumber, pageSize)
		pagination[uuid] = samples
	}

	return pagination
}

func (hm *historyManager) GetSamplesByDateRange(startDate, endDate time.Time) map[string][]Record {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	result := make(map[string][]Record)

	for uuid, history := range hm.histories {
		samples := history.GetSamplesByDateRange(startDate, endDate)
		result[uuid] = samples
	}

	return result
}

func (hm *historyManager) GetSamplesByTime(startDate time.Time, duration string) (map[string][]Record, error) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	result := make(map[string][]Record)
	_, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}

	for uuid, history := range hm.histories {
		samples, err := history.GetSamplesByTime(startDate, duration)
		if err != nil {
			return nil, err
		}
		result[uuid] = samples
	}

	return result, nil
}

func (hm *historyManager) DeleteSamples(uuids map[string]string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	for uuid := range uuids {
		if history, exists := hm.histories[uuid]; exists {
			uuidSlice := []string{uuid}
			history.DeleteSamples(uuidSlice)
		}
	}
}
