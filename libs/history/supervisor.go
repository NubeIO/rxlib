package history

type Supervisor interface {
	AddRemoteManager(serverName string, manager Manager)
	GetRemoteManager(serverName string) Manager
	All() map[string]Manager
	GetAllManagersRecordsEntries() map[string]*SupervisorRecords
}

type supervisor struct {
	supervisorMap map[string]Manager
}

func NewSupervisor() Supervisor {
	return &supervisor{
		supervisorMap: make(map[string]Manager),
	}
}

func (c *supervisor) AddRemoteManager(serverName string, manager Manager) {
	c.supervisorMap[serverName] = manager
}

func (c *supervisor) GetRemoteManager(serverName string) Manager {
	return c.supervisorMap[serverName]
}

func (c *supervisor) All() map[string]Manager {
	return c.supervisorMap
}

type SupervisorRecords struct {
	Name    string                     `json:"name"`
	Manager map[string][]*AllHistories `json:"manager"`
}

func (c *supervisor) GetAllManagersRecordsEntries() map[string]*SupervisorRecords {
	recordsEntries := make(map[string]*SupervisorRecords)

	// Iterate through all managers and collect their history records
	for serverName, manager := range c.supervisorMap {
		historyRecords := manager.AllHistories()

		// Create a new entry for each manager
		recordsEntries[serverName] = &SupervisorRecords{
			Name:    manager.GetName(),
			Manager: make(map[string][]*AllHistories),
		}
		recordsEntries[serverName].Manager[serverName] = historyRecords
	}

	return recordsEntries
}
