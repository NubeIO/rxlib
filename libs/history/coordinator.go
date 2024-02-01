package history

type Coordinator interface {
	AddRemoteManager(serverName string, manager Manager)
	GetRemoteManager(serverName string) Manager
	All() map[string]Manager
	GetAllManagersRecordsEntries() map[string]*CoordinatorRecords
}

type coordinator struct {
	coordinatorMap map[string]Manager
}

func NewCoordinator() Coordinator {
	return &coordinator{
		coordinatorMap: make(map[string]Manager),
	}
}

func (c *coordinator) AddRemoteManager(serverName string, manager Manager) {
	c.coordinatorMap[serverName] = manager
}

func (c *coordinator) GetRemoteManager(serverName string) Manager {
	return c.coordinatorMap[serverName]
}

func (c *coordinator) All() map[string]Manager {
	return c.coordinatorMap
}

type CoordinatorRecords struct {
	Name    string                     `json:"name"`
	Manager map[string][]*AllHistories `json:"manager"`
}

func (c *coordinator) GetAllManagersRecordsEntries() map[string]*CoordinatorRecords {
	recordsEntries := make(map[string]*CoordinatorRecords)

	// Iterate through all managers and collect their history records
	for serverName, manager := range c.coordinatorMap {
		historyRecords := manager.AllHistories()

		// Create a new entry for each manager
		recordsEntries[serverName] = &CoordinatorRecords{
			Name:    manager.GetName(),
			Manager: make(map[string][]*AllHistories),
		}
		recordsEntries[serverName].Manager[serverName] = historyRecords
	}

	return recordsEntries
}
