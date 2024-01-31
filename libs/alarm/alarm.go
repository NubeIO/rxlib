package alarm

import (
	"github.com/NubeIO/rxlib/helpers"
	"sync"
	"time"
)

type AlarmManager interface {
	NewAlarm(limitSize int) Alarm
	Get(uuid string) Alarm
	All() map[string]Alarm
	Drop(uuid string)
	DropAll()
	GetAllTransactionsEntries() map[string][]*TransactionEntry
	AllTransactions() map[string][]Transaction
	GetTransactionPagination(pageNumber, pageSize int) map[string][]Transaction
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

func (m *manager) GetTransactionPagination(pageNumber, pageSize int) map[string][]Transaction {
	return nil
}

func (m *manager) NewAlarm(limitSize int) Alarm {
	alarm := NewAlarm(limitSize)
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

type Alarm interface {
	AddTransaction(body *AddTransaction, t Transaction)
	SetStatus(status AlarmStatus)
	GetUUID() string
	GetTransactions() []Transaction
	GetLast() Transaction
	GetFirst() Transaction
	GetAllTransactionsEntries() map[string][]*TransactionEntry
	GetPagination(pageNumber, pageSize int) []Transaction
	GetByDateRange(startDate, endDate time.Time) []Transaction
	GetByTime(startDate time.Time, duration string) ([]Transaction, error)
	DeleteTransaction(uuid string)
	DeleteTransactions(uuids []string)
	DeleteFirst(count int) int
	DeleteLast(count int) int
	DeleteByDateRange(startDate, endDate time.Time) int
	DeleteByTime(startDate time.Time, duration string) int
	TransactionCount() int
}

func NewAlarm(limitSampleSize int) Alarm {
	return &AlarmEntry{UUID: helpers.UUID(), LimitTransactionCount: limitSampleSize}
}

type AlarmEntry struct {
	UUID                  string              `json:"uuid"`
	ObjectID              string              `json:"objectID"`             // Device
	ObjectUUID            string              `json:"objectUUID,omitempty"` // dev_abc123
	Type                  string              `json:"type"`                 // Ping
	Status                AlarmStatus         `json:"status"`               // Active
	Notified              bool                `json:"notified,omitempty"`
	NotifiedAt            time.Time           `json:"notified_at,omitempty"`
	CreatedAt             time.Time           `json:"created_at,omitempty"`
	LastUpdated           time.Time           `json:"last_updated,omitempty"`
	Transactions          []*TransactionEntry `json:"transactions,omitempty"`
	LimitTransactionCount int                 `json:"limitTransactionCount"`
}

type AddTransaction struct {
	Status   AlarmStatus   `json:"status"`   // Active
	Severity AlarmSeverity `json:"severity"` // Crucial
	Target   string        `json:"target,omitempty"`
	Title    string        `json:"title,omitempty"`
	Body     string        `json:"body,omitempty"`
}

func NewTransactionBody(status AlarmStatus, severity AlarmSeverity, title, body string) *AddTransaction {
	return &AddTransaction{
		Status:   status,
		Severity: severity,
		Target:   "",
		Title:    title,
		Body:     body,
	}
}

func (a *AlarmEntry) AddTransaction(body *AddTransaction, t Transaction) {
	if te, ok := t.(*TransactionEntry); ok {
		te.AlarmUUID = a.GetUUID() // Set the AlarmUUID to the UUID of the AlarmEntry
		te.SetStatus(body.Status)
		te.SetBody(body.Body)
		te.SetTitle(body.Title)
		te.SetSeverity(body.Severity)
		te.lastUpdated()
		te.createdAt()
		a.Transactions = append(a.Transactions, te)
		// Update the alarm status and last updated timestamp here
		a.Status = a.calculateAlarmStatus()
		a.LastUpdated = time.Now()
	}
}

func (a *AlarmEntry) calculateAlarmStatus() AlarmStatus {

	for _, t := range a.Transactions {
		if t.Status == AlarmStatusClosed {
			return AlarmStatusClosed
		}
	}

	return AlarmStatusActive
}

func transactionToTransactionEntry(alarmUUID string, t Transaction) *TransactionEntry {
	return &TransactionEntry{
		UUID:        t.GetUUID(),
		AlarmUUID:   alarmUUID,
		Status:      t.GetStatus(),
		Severity:    t.GetSeverity(),
		Target:      t.GetTarget(),
		Title:       t.GetTitle(),
		Body:        t.GetBody(),
		CreatedAt:   t.GetCreatedAt(),
		LastUpdated: t.GetLastUpdated(),
	}
}

func (a *AlarmEntry) GetAllTransactionsEntries() map[string][]*TransactionEntry {
	transactions := make(map[string][]*TransactionEntry)
	for _, t := range a.Transactions {
		transactionEntry := transactionToTransactionEntry(a.GetUUID(), t)
		transactions[a.UUID] = append(transactions[a.UUID], transactionEntry)
	}

	return transactions
}

func (a *AlarmEntry) lastUpdated() {
	a.LastUpdated = time.Now()
}

func (a *AlarmEntry) SetStatus(status AlarmStatus) {
	a.lastUpdated()
	a.Status = status
}

func (a *AlarmEntry) GetUUID() string {
	return a.UUID
}

func (a *AlarmEntry) GetTransactions() []Transaction {
	transactions := make([]Transaction, len(a.Transactions))
	for i, t := range a.Transactions {
		transactions[i] = t
	}
	return transactions
}

func (a *AlarmEntry) GetLast() Transaction {
	if len(a.Transactions) > 0 {
		return a.Transactions[len(a.Transactions)-1]
	}
	return nil
}

func (a *AlarmEntry) GetFirst() Transaction {
	if len(a.Transactions) > 0 {
		return a.Transactions[0]
	}
	return nil
}

func (a *AlarmEntry) GetPagination(pageNumber, pageSize int) []Transaction {
	startIndex := (pageNumber - 1) * pageSize
	endIndex := startIndex + pageSize
	if startIndex < 0 {
		startIndex = 0
	}
	if endIndex > len(a.Transactions) {
		endIndex = len(a.Transactions)
	}

	paginatedTransactions := make([]Transaction, endIndex-startIndex)

	for i, t := range a.Transactions[startIndex:endIndex] {
		paginatedTransactions[i] = t
	}

	return paginatedTransactions
}

func (a *AlarmEntry) GetByDateRange(startDate, endDate time.Time) []Transaction {
	var result []Transaction
	for _, t := range a.Transactions {
		if t.GetCreatedAt().After(startDate) && t.GetCreatedAt().Before(endDate) {
			result = append(result, t)
		}
	}
	return result
}

func (a *AlarmEntry) GetByTime(startDate time.Time, duration string) ([]Transaction, error) {
	durationValue, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}

	endDate := startDate.Add(durationValue)
	result := make([]Transaction, 0)
	for _, t := range a.Transactions {
		if t.GetCreatedAt().After(startDate) && t.GetCreatedAt().Before(endDate) {
			result = append(result, t)
		}
	}
	return result, nil
}

func (a *AlarmEntry) DeleteTransaction(uuid string) {
	var newTransactions []*TransactionEntry
	for _, t := range a.Transactions {
		if t.GetUUID() != uuid {
			newTransactions = append(newTransactions, t)
		}
	}
	a.Transactions = newTransactions
}

func (a *AlarmEntry) DeleteTransactions(uuids []string) {
	for _, uuid := range uuids {
		a.DeleteTransaction(uuid)
	}
}

func (a *AlarmEntry) DeleteFirst(count int) int {
	if count <= 0 {
		return 0
	}
	if count >= len(a.Transactions) {
		count = len(a.Transactions)
	}
	a.Transactions = a.Transactions[count:]
	return count
}

func (a *AlarmEntry) DeleteLast(count int) int {
	if count <= 0 {
		return 0
	}
	if count >= len(a.Transactions) {
		count = len(a.Transactions)
	}
	a.Transactions = a.Transactions[:len(a.Transactions)-count]
	return count
}

func (a *AlarmEntry) DeleteByDateRange(startDate, endDate time.Time) int {
	var deletedCount int
	var newTransactions []*TransactionEntry
	for _, t := range a.Transactions {
		if t.GetCreatedAt().Before(startDate) || t.GetCreatedAt().After(endDate) {
			newTransactions = append(newTransactions, t)
		} else {
			deletedCount++
		}
	}
	a.Transactions = newTransactions
	return deletedCount
}

func (a *AlarmEntry) DeleteByTime(startDate time.Time, duration string) int {
	durationValue, err := time.ParseDuration(duration)
	if err != nil {
		return 0
	}

	endDate := startDate.Add(durationValue)
	result := make([]*TransactionEntry, 0)
	deletedCount := 0

	for _, t := range a.Transactions {
		if t.GetCreatedAt().Before(startDate) || t.GetCreatedAt().After(endDate) {
			result = append(result, t)
		} else {
			deletedCount++
		}
	}

	a.Transactions = result
	return deletedCount
}

func (a *AlarmEntry) TransactionCount() int {
	return len(a.Transactions)
}

type Transaction interface {
	GetAlarmUUID() string
	GetUUID() string
	GetStatus() AlarmStatus
	GetSeverity() AlarmSeverity
	GetTarget() string
	GetTitle() string
	GetBody() string
	GetCreatedAt() time.Time
	GetLastUpdated() time.Time

	SetTitle(title string)
	SetStatus(state AlarmStatus)
	SetSeverity(s AlarmSeverity)
	SetBody(body string)
}

type TransactionEntry struct {
	Title       string        `json:"title,omitempty"`
	Status      AlarmStatus   `json:"status"`   // Active
	Severity    AlarmSeverity `json:"severity"` // Crucial
	Target      string        `json:"target,omitempty"`
	Body        string        `json:"body,omitempty"`
	UUID        string        `json:"uuid"`
	AlarmUUID   string        `json:"alarmUUID"`
	CreatedAt   time.Time     `json:"createdAt,omitempty"`
	LastUpdated time.Time     `json:"lastUpdated,omitempty"`
}

func (t *TransactionEntry) SetTitle(title string) {
	t.Title = title
}

func (t *TransactionEntry) SetSeverity(s AlarmSeverity) {
	t.Severity = s
}

func (t *TransactionEntry) SetBody(body string) {
	t.Body = body
}

func NewTransaction() Transaction {
	s := &TransactionEntry{
		UUID: helpers.UUID(),
	}
	s.createdAt()
	return s
}

func (t *TransactionEntry) lastUpdated() {
	t.LastUpdated = time.Now()
}
func (t *TransactionEntry) createdAt() {
	t.CreatedAt = time.Now()
}

func (t *TransactionEntry) SetStatus(status AlarmStatus) {
	t.lastUpdated()
	t.Status = status
}

func (t *TransactionEntry) GetStatus() AlarmStatus {
	return t.Status
}

func (t *TransactionEntry) GetSeverity() AlarmSeverity {
	return t.Severity
}

func (t *TransactionEntry) GetTarget() string {
	return t.Target
}

func (t *TransactionEntry) GetTitle() string {
	return t.Title
}

func (t *TransactionEntry) GetBody() string {
	return t.Body
}

func (t *TransactionEntry) GetAlarmUUID() string {
	return t.AlarmUUID
}

func (t *TransactionEntry) GetUUID() string {
	return t.UUID

}

func (t *TransactionEntry) GetCreatedAt() time.Time {
	return t.CreatedAt
}
func (t *TransactionEntry) GetLastUpdated() time.Time {
	return t.LastUpdated
}
