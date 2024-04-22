package history

import (
	"github.com/NubeIO/rxlib/helpers"
	"time"
)

type History interface {
	AddRecord(record Record)
	AddRecords(records []Record)
	GetUUID() string
	GetObjectUUID() string
	GetRecords() []Record
	GetLast() Record
	GetFirst() Record
	GetPagination(pageNumber, pageSize int) []Record
	GetRecordsByDateRange(startDate, endDate time.Time) []Record
	GetRecordsByTime(startDate time.Time, duration string) ([]Record, error)
	DeleteRecord(sample Record)
	DeleteRecords(uuids []string)
	DeleteFirst(count int) int
	DeleteLast(count int) int
	DeleteByDateRange(startDate, endDate time.Time) int
	DeleteByTime(startDate time.Time, duration string) int
	RecordCount() int
}

type GenericHistory struct {
	UUID            string   `json:"uuid"`
	ObjectUUID      string   `json:"objectUUID"`
	Values          []Record `json:"values"`
	LimitRecordsize int      `json:"limitRecordsize"`
}

func (h *GenericHistory) RecordCount() int {
	return len(h.Values)
}

func NewGenericHistory(limitRecordsize int, objectUUID string) History {
	return &GenericHistory{UUID: helpers.UUID(), ObjectUUID: objectUUID, LimitRecordsize: limitRecordsize}
}

func (h *GenericHistory) AddRecord(sample Record) {
	h.Values = append(h.Values, sample)
	if len(h.Values) > h.LimitRecordsize {
		// Remove the oldest Records to keep the size within the limit
		removedCount := len(h.Values) - h.LimitRecordsize
		h.Values = h.Values[removedCount:]
	}
}

func (h *GenericHistory) AddRecords(records []Record) {
	for _, record := range records {
		h.AddRecord(record)
	}
}

type Record interface {
	GetUUID() string
	GetValue() interface{}
	GetTimestamp() time.Time
}

type GenericRecord[T any] struct {
	UUID      string    `json:"uuid"`
	Value     T         `json:"value"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

func NewGenericRecord[T any](value T) Record {
	s := &GenericRecord[T]{
		UUID:      helpers.UUID(),
		Value:     value,
		Timestamp: time.Now(),
	}
	return s
}

func (s *GenericRecord[T]) GetUUID() string {
	return s.UUID
}

func (s *GenericRecord[T]) GetValue() interface{} {
	return s.Value
}

func (s *GenericRecord[T]) GetTimestamp() time.Time {
	return s.Timestamp
}
