package schedules

import "sync"

// ScheduleManager is an interface for managing schedules schedules.
type ScheduleManager interface {
	Add(name string, sch *Schedule)
	Get(name string) *Schedule
	Edit(name string, sch *Schedule)
	All() map[string]*Schedule
	Delete(name string)
}

// scheduleManager is an implementation of the ScheduleManager interface.
type scheduleManager struct {
	schedules map[string]*Schedule
	mu        sync.Mutex
}

// New creates a new instance of scheduleManager and returns it as a ScheduleManager interface.
func New() ScheduleManager {
	return &scheduleManager{
		schedules: make(map[string]*Schedule),
	}
}

// Add adds a new schedule with the given name.
func (sm *scheduleManager) Add(name string, sch *Schedule) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.schedules[name] = sch
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
