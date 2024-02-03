package rubix

type Supervisor interface {
	AddRemoteManager(serverName string, manager Manager)
	GetRemoteManager(serverName string) Manager
	All() map[string]Manager
}

type SupervisorImpl struct {
	supervisorMap map[string]Manager
}

func NewSupervisor() Supervisor {
	return &SupervisorImpl{
		supervisorMap: make(map[string]Manager),
	}
}

func (a *SupervisorImpl) AddRemoteManager(serverName string, manager Manager) {
	a.supervisorMap[serverName] = manager
}

func (a *SupervisorImpl) GetRemoteManager(serverName string) Manager {
	return a.supervisorMap[serverName]
}

func (a *SupervisorImpl) All() map[string]Manager {
	return a.supervisorMap
}
