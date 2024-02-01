package history

import (
	"github.com/NubeIO/rxlib/helpers"
	"time"
)

type History interface {
	AddSample(sample Record)
	GetUUID() string
	GetObjectUUID() string
	GetRecords() []Record
	GetLast() Record
	GetFirst() Record
	GetPagination(pageNumber, pageSize int) []Record
	GetRecordsByDateRange(startDate, endDate time.Time) []Record
	GetRecordsByTime(startDate time.Time, duration string) ([]Record, error)
	DeleteSample(sample Record)
	DeleteRecords(uuids []string)
	DeleteFirst(count int) int
	DeleteLast(count int) int
	DeleteByDateRange(startDate, endDate time.Time) int
	DeleteByTime(startDate time.Time, duration string) int
	SampleCount() int
}

type Record interface {
	GetUUID() string
	GetValue() interface{}
	GetTimestamp() time.Time
}

type GenericSample[T any] struct {
	UUID      string    `json:"uuid"`
	Value     T         `json:"value"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

func NewGenericSample[T any](value T) Record {
	s := &GenericSample[T]{
		UUID:      helpers.UUID(),
		Value:     value,
		Timestamp: time.Now(),
	}
	return s
}

type GenericHistory struct {
	UUID            string   `json:"uuid"`
	ObjectUUID      string   `json:"objectUUID"`
	Values          []Record `json:"values"`
	LimitRecordsize int      `json:"limitRecordsize"`
}

func NewGenericHistory(limitRecordsize int, objectUUID string) History {
	return &GenericHistory{UUID: helpers.UUID(), ObjectUUID: objectUUID, LimitRecordsize: limitRecordsize}
}

func (h *GenericHistory) AddSample(sample Record) {
	h.Values = append(h.Values, sample)
	if len(h.Values) > h.LimitRecordsize {
		// Remove the oldest Records to keep the size within the limit
		removedCount := len(h.Values) - h.LimitRecordsize
		h.Values = h.Values[removedCount:]
	}
}

func (s *GenericSample[T]) GetUUID() string {
	return s.UUID
}

func (s *GenericSample[T]) GetValue() interface{} {
	return s.Value
}

func (s *GenericSample[T]) GetTimestamp() time.Time {
	return s.Timestamp
}

func (h *GenericHistory) SampleCount() int {
	return len(h.Values)
}
