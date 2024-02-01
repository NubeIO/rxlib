package schedules

type Coordinator interface {
	AddRemoteManager(serverName string, manager ScheduleManager)
	GetRemoteManager(serverName string) ScheduleManager
	All() map[string]ScheduleManager
}

type coordinator struct {
	coordinatorMap map[string]ScheduleManager
}

func NewCoordinator() Coordinator {
	return &coordinator{
		coordinatorMap: make(map[string]ScheduleManager),
	}
}

func (a *coordinator) AddRemoteManager(serverName string, manager ScheduleManager) {
	a.coordinatorMap[serverName] = manager
}

func (a *coordinator) GetRemoteManager(serverName string) ScheduleManager {
	return a.coordinatorMap[serverName]
}

func (a *coordinator) All() map[string]ScheduleManager {
	return a.coordinatorMap
}

type CoordinatorSchedules struct {
	Title   string                        `json:"title"`
	Manager map[string][]*ScheduleManager `json:"manager"`
}
