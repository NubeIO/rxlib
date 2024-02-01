package alarm

import (
	"sync"
	"time"
)

type AlarmManager interface {
	NewAlarm(limitSize int, alarmBody *AddAlarm) Alarm
	Get(uuid string) Alarm
	All() map[string]Alarm
	Drop(uuid string)
	DropAll()
	GetAllTransactionsEntries() map[string][]*TransactionEntry
	AllTransactions() map[string][]Transaction
	GetTransactionByDateRange(startDate, endDate time.Time) map[string][]Transaction
	GetTransactionByTime(startDate time.Time, duration string) (map[string][]Transaction, error)
	DeleteTransactions(uuids map[string]string)
}

func NewAlarmManager() AlarmManager {
	return &manager{
		alarmMap: make(map[string]Alarm),
	}
}

type manager struct {
	alarmMap map[string]Alarm
	mu       sync.RWMutex
}

func (m *manager) NewAlarm(limitSize int, alarmBody *AddAlarm) Alarm {
	alarm := NewAlarm(limitSize, alarmBody)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alarmMap[alarm.GetUUID()] = alarm
	return alarm
}

func (m *manager) Get(uuid string) Alarm {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.alarmMap[uuid]
}

func (m *manager) All() map[string]Alarm {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.alarmMap
}

func (m *manager) Drop(uuid string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.alarmMap, uuid)
}

func (m *manager) DropAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alarmMap = make(map[string]Alarm)
}

func (m *manager) GetAllTransactionsEntries() map[string][]*TransactionEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	transactionsEntries := make(map[string][]*TransactionEntry)
	for _, alarm := range m.alarmMap {
		transactions := alarm.GetTransactions()
		transactionEntries := make([]*TransactionEntry, len(transactions))
		for i, t := range transactions {
			transactionEntries[i] = transactionToTransactionEntry(alarm.GetUUID(), t)
		}
		transactionsEntries[alarm.GetUUID()] = transactionEntries
	}
	return transactionsEntries
}

func (m *manager) AllTransactions() map[string][]Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()
	transactions := make(map[string][]Transaction)
	for _, alarm := range m.alarmMap {
		transactions[alarm.GetUUID()] = alarm.GetTransactions()
	}
	return transactions
}

func (m *manager) GetPagination(pageNumber, pageSize int) map[string][]Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()
	transactions := make(map[string][]Transaction)
	for _, alarm := range m.alarmMap {
		transactions[alarm.GetUUID()] = alarm.GetPagination(pageNumber, pageSize)
	}
	return transactions
}

func (m *manager) GetTransactionByDateRange(startDate, endDate time.Time) map[string][]Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()
	transactions := make(map[string][]Transaction)
	for _, alarm := range m.alarmMap {
		transactions[alarm.GetUUID()] = alarm.GetByDateRange(startDate, endDate)
	}
	return transactions
}

func (m *manager) GetTransactionByTime(startDate time.Time, duration string) (map[string][]Transaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	transactions := make(map[string][]Transaction)
	for _, alarm := range m.alarmMap {
		t, err := alarm.GetByTime(startDate, duration)
		if err != nil {
			return nil, err
		}
		transactions[alarm.GetUUID()] = t
	}
	return transactions, nil
}

func (m *manager) DeleteTransactions(uuids map[string]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for uuid, _ := range uuids {
		alarm, ok := m.alarmMap[uuid]
		if ok {
			alarm.DeleteTransaction(uuid)
		}
	}
}
