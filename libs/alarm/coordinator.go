package alarm

type Coordinator interface {
	AddRemoteManager(serverName string, manager AlarmManager)
	GetRemoteManager(serverName string) AlarmManager
	All() map[string]AlarmManager
	GetAllManagersTransactionsEntries() map[string]*CoordinatorTransactions
}

type alarmCoordinator struct {
	alarmMap map[string]AlarmManager
}

func NewCoordinator() Coordinator {
	return &alarmCoordinator{
		alarmMap: make(map[string]AlarmManager),
	}
}

func (a *alarmCoordinator) AddRemoteManager(serverName string, manager AlarmManager) {
	a.alarmMap[serverName] = manager
}

func (a *alarmCoordinator) GetRemoteManager(serverName string) AlarmManager {
	return a.alarmMap[serverName]
}

func (a *alarmCoordinator) All() map[string]AlarmManager {
	return a.alarmMap
}

type CoordinatorTransactions struct {
	Title   string                         `json:"title"`
	Manager map[string][]*TransactionEntry `json:"manager"`
}

func (a *alarmCoordinator) GetAllManagersTransactionsEntries() map[string]*CoordinatorTransactions {
	transactionsEntries := make(map[string]*CoordinatorTransactions)

	// Iterate through all managed AlarmManagers and collect their transaction entries
	for serverName, manager := range a.alarmMap {
		alarmTransactionsEntries := manager.GetAllTransactionsEntries()

		// Merge alarmTransactionsEntries into transactionsEntries
		for key, value := range alarmTransactionsEntries {
			if _, exists := transactionsEntries[key]; !exists {
				transactionsEntries[key] = &CoordinatorTransactions{
					Title:   manager.GetTitle(),
					Manager: make(map[string][]*TransactionEntry),
				}
			}
			transactionsEntries[key].Manager[serverName] = value
		}
	}

	return transactionsEntries
}
