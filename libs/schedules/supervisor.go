package schedules

type Supervisor interface {
	AddRemoteManager(serverName string, manager ScheduleManager)
	GetRemoteManager(serverName string) ScheduleManager
	All() map[string]ScheduleManager
}

type supervisor struct {
	supervisorMap map[string]ScheduleManager
}

func NewSupervisor() Supervisor {
	return &supervisor{
		supervisorMap: make(map[string]ScheduleManager),
	}
}

func (a *supervisor) AddRemoteManager(serverName string, manager ScheduleManager) {
	a.supervisorMap[serverName] = manager
}

func (a *supervisor) GetRemoteManager(serverName string) ScheduleManager {
	return a.supervisorMap[serverName]
}

func (a *supervisor) All() map[string]ScheduleManager {
	return a.supervisorMap
}

type SupervisorSchedules struct {
	Title   string                        `json:"title"`
	Manager map[string][]*ScheduleManager `json:"manager"`
}
