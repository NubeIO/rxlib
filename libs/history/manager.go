package history

import (
	"sync"
	"time"
)

type Manager interface {
	// NewHistory creates a new history with the specified limitSampleSize.
	NewHistory(limitSampleSize int, objectUUID string) History

	GetName() string

	// Get retrieves a history by its UUID.
	Get(uuid string) History

	// AllHistories returns a slice of all available histories. with some stats
	AllHistories() []*AllHistories

	// All returns a slice of all available histories.
	All() []History

	// AllRecords returns a map of all samples across all histories.
	AllRecords() map[string][]Record

	// Drop deletes a history by its UUID.
	Drop(uuid string)

	// DropAll will delete all the histories
	DropAll()

	// GetPagination retrieves paginated samples for all histories managed by the Manager.
	GetPagination(pageNumber, pageSize int) map[string][]Record

	// GetRecordsByDateRange retrieves samples within a specified date range for all histories managed by the Manager.
	GetRecordsByDateRange(startDate, endDate time.Time) map[string][]Record

	// GetRecordsByTime retrieves samples within a specified time duration for all histories managed by the Manager.
	// It takes a startDate and duration string (e.g., "10s", "1h") as input and returns a map where keys are UUIDs
	GetRecordsByTime(startDate time.Time, duration string) (map[string][]Record, error)

	// DeleteRecords deletes samples from specified histories based on UUIDs for all histories managed by the Manager.
	DeleteRecords(uuids map[string]string)
}

type historyManager struct {
	name      string
	histories map[string]History
	mu        sync.RWMutex
}

func (hm *historyManager) GetName() string {
	return hm.name
}

func NewHistoryManager(name string) Manager {
	if name == "" {
		panic("history manager name can not be empty")
	}
	return &historyManager{
		name:      name,
		histories: make(map[string]History),
	}
}

func (hm *historyManager) NewHistory(limitSampleSize int, objectUUID string) History {
	history := NewGenericHistory(limitSampleSize, objectUUID)
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

type AllHistories struct {
	ObjectUUID  string   `json:"objectUUID"`
	HistoryUUID string   `json:"historyUUID"`
	Count       int      `json:"count"`
	Histories   []Record `json:"histories"`
}

func (hm *historyManager) AllHistories() []*AllHistories {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	histories := make([]*AllHistories, 0, len(hm.histories))
	for _, history := range hm.histories {
		h := &AllHistories{
			HistoryUUID: history.GetUUID(),
			ObjectUUID:  history.GetObjectUUID(),
			Count:       len(history.GetRecords()),
			Histories:   history.GetRecords(),
		}
		histories = append(histories, h)
	}
	return histories
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

func (hm *historyManager) AllRecords() map[string][]Record {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	samples := make(map[string][]Record)
	for uuid, history := range hm.histories {
		samples[uuid] = history.GetRecords()
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

func (hm *historyManager) GetRecordsByDateRange(startDate, endDate time.Time) map[string][]Record {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	result := make(map[string][]Record)

	for uuid, history := range hm.histories {
		samples := history.GetRecordsByDateRange(startDate, endDate)
		result[uuid] = samples
	}

	return result
}

func (hm *historyManager) GetRecordsByTime(startDate time.Time, duration string) (map[string][]Record, error) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	result := make(map[string][]Record)
	_, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}

	for uuid, history := range hm.histories {
		samples, err := history.GetRecordsByTime(startDate, duration)
		if err != nil {
			return nil, err
		}
		result[uuid] = samples
	}

	return result, nil
}

func (hm *historyManager) DeleteRecords(uuids map[string]string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	for uuid := range uuids {
		if history, exists := hm.histories[uuid]; exists {
			uuidSlice := []string{uuid}
			history.DeleteRecords(uuidSlice)
		}
	}
}
