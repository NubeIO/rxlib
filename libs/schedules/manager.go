package schedules

import (
	"encoding/json"
	"sync"
	"time"
)

// Manager is an interface for managing schedules.
type Manager interface {
	Add(sch *Schedule)
	Parse(body any) *Schedule
	ParseFromString(body string) *Schedule
	Get(name string) *Schedule
	Edit(name string, sch *Schedule)
	All() map[string]*Schedule
	Delete(name string)
}

// scheduleManager is an implementation of the Manager interface.
type scheduleManager struct {
	schedules map[string]*Schedule
	mu        sync.Mutex
}

// Parse a schedule from json
func (sm *scheduleManager) Parse(body any) *Schedule {
	var sc *Schedule
	marshal, err := json.Marshal(body)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(marshal, &sc)
	if err != nil {
		return nil
	}
	return sc
}

// ParseFromString a schedule from json
func (sm *scheduleManager) ParseFromString(body string) *Schedule {
	var sc *Schedule
	err := json.Unmarshal([]byte(body), &sc)
	if err != nil {
		return nil
	}
	return sc
}

// New creates a new instance of scheduleManager and returns it as a Manager interface.
func New() Manager {
	return &scheduleManager{
		schedules: make(map[string]*Schedule),
	}
}

// Add adds a new schedule with the given name.
func (sm *scheduleManager) Add(sch *Schedule) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sch.Exceptions == nil {
		sch.Exceptions = make(map[time.Time]Entry)
	}
	if sch.DayToTimeRanges == nil {
		sch.DayToTimeRanges = make(map[time.Weekday][]Entry)
	}
	sm.schedules[sch.Name] = sch
}

// Get retrieves a schedule by name.
func (sm *scheduleManager) Get(name string) *Schedule {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.schedules[name]
}

// Edit updates an existing schedule by name.
func (sm *scheduleManager) Edit(name string, sch *Schedule) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.schedules[name] = sch
}

// All returns all schedules as a map of name to WeeklySchedule.
func (sm *scheduleManager) All() map[string]*Schedule {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.schedules
}

// Delete removes a schedule by name.
func (sm *scheduleManager) Delete(name string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.schedules, name)
}
