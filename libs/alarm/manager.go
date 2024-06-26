package alarm

import (
	"sync"
	"time"
)

type Manager interface {
	GetTitle() string
	NewAlarm(limitTransactionSize int, alarmBody *AddAlarm) Alarm
	Get(uuid string) Alarm
	All() map[string]Alarm
	Drop(uuid string)
	DropAll()
	Entries() map[string][]*TransactionEntry
	AllTransactions() map[string][]Transaction
	GetTransactionByDateRange(startDate, endDate time.Time) map[string][]Transaction
	GetTransactionByTime(startDate time.Time, duration string) (map[string][]Transaction, error)
	DeleteTransactions(uuids map[string]string)
}

func NewAlarmManager(title string) Manager {
	return &manager{
		title:    title,
		alarmMap: make(map[string]Alarm),
	}
}

type manager struct {
	title    string
	alarmMap map[string]Alarm
	mu       sync.RWMutex
}

func (m *manager) GetTitle() string {
	return m.title
}

func (m *manager) NewAlarm(limitTransactionSize int, alarmBody *AddAlarm) Alarm {
	alarm := NewAlarm(limitTransactionSize, alarmBody)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alarmMap[alarm.GetUUID()] = alarm
	if limitTransactionSize <= 0 {
		limitTransactionSize = 100
	}
	if limitTransactionSize >= 100000 {
		limitTransactionSize = 100000
	}
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

func (m *manager) Entries() map[string][]*TransactionEntry {
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
