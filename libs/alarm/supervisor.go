package alarm

type Supervisor interface {
	AddRemoteManager(serverName string, manager AlarmManager)
	GetRemoteManager(serverName string) AlarmManager
	All() map[string]AlarmManager
	GetAllManagersTransactionsEntries() map[string]*SupervisorTransactions
}

type supervisor struct {
	supervisorMap map[string]AlarmManager
}

func NewSupervisor() Supervisor {
	return &supervisor{
		supervisorMap: make(map[string]AlarmManager),
	}
}

func (a *supervisor) AddRemoteManager(serverName string, manager AlarmManager) {
	a.supervisorMap[serverName] = manager
}

func (a *supervisor) GetRemoteManager(serverName string) AlarmManager {
	return a.supervisorMap[serverName]
}

func (a *supervisor) All() map[string]AlarmManager {
	return a.supervisorMap
}

type SupervisorTransactions struct {
	Title   string                         `json:"title"`
	Manager map[string][]*TransactionEntry `json:"manager"`
}

func (a *supervisor) GetAllManagersTransactionsEntries() map[string]*SupervisorTransactions {
	transactionsEntries := make(map[string]*SupervisorTransactions)

	// Iterate through all managed AlarmManagers and collect their transaction entries
	for serverName, manager := range a.supervisorMap {
		alarmTransactionsEntries := manager.GetAllTransactionsEntries()

		// Merge alarmTransactionsEntries into transactionsEntries
		for key, value := range alarmTransactionsEntries {
			if _, exists := transactionsEntries[key]; !exists {
				transactionsEntries[key] = &SupervisorTransactions{
					Title:   manager.GetTitle(),
					Manager: make(map[string][]*TransactionEntry),
				}
			}
			transactionsEntries[key].Manager[serverName] = value
		}
	}

	return transactionsEntries
}
