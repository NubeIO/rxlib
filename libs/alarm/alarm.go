package alarm

import (
	"github.com/NubeIO/rxlib/helpers"
	"time"
)

type Alarm interface {
	AddTransaction(body *AddTransaction, t Transaction)
	SetStatus(status AlarmStatus)
	GetTitle() string
	GetUUID() string
	GetObjectUUID() string
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

type AddAlarm struct {
	Title      string `json:"title"`
	ObjectType string `json:"objectType"`           // Device
	ObjectUUID string `json:"objectUUID,omitempty"` // dev_abc123
}

func NewAlarm(limitSampleSize int, alarmBody *AddAlarm) Alarm {
	checksAddAlarm(alarmBody)
	return &AlarmEntry{UUID: helpers.UUID(), Title: alarmBody.Title, ObjectType: alarmBody.ObjectType, ObjectUUID: alarmBody.ObjectUUID, LimitTransactionCount: limitSampleSize}
}

func checksAddAlarm(alarmBody *AddAlarm) {
	if alarmBody == nil {
		panic("add alarm alarmBody is empty")
	}
	if alarmBody.Title == "" {
		panic("add alarm alarmBody.Title is empty")
	}
	if alarmBody.ObjectType == "" {
		panic("add alarm alarmBody.ObjectType is empty")
	}
	if alarmBody.ObjectUUID == "" {
		panic("add alarm alarmBody.ObjectUUID is empty")
	}
}

type AlarmEntry struct {
	Title                 string              `json:"title"`
	UUID                  string              `json:"uuid"`
	ObjectType            string              `json:"objectType"`           // Device
	ObjectUUID            string              `json:"objectUUID,omitempty"` // dev_abc123
	AlarmType             AlarmType           `json:"alarmType"`            // Ping
	Status                AlarmStatus         `json:"status"`               // Active
	Notified              bool                `json:"notified,omitempty"`
	NotifiedAt            time.Time           `json:"notified_at,omitempty"`
	CreatedAt             time.Time           `json:"created_at,omitempty"`
	LastUpdated           time.Time           `json:"last_updated,omitempty"`
	Transactions          []*TransactionEntry `json:"transactions,omitempty"`
	LimitTransactionCount int                 `json:"limitTransactionCount"`
}

func (a *AlarmEntry) GetTitle() string {
	return a.Title
}

func (a *AlarmEntry) GetObjectUUID() string {
	return a.ObjectUUID
}

type AddTransaction struct {
	Status   AlarmStatus   `json:"status"`   // Active
	Severity AlarmSeverity `json:"severity"` // Crucial
	Title    string        `json:"title,omitempty"`
	Body     string        `json:"body,omitempty"`
}

func NewTransactionBody(status AlarmStatus, severity AlarmSeverity, title, body string) *AddTransaction {
	return &AddTransaction{
		Status:   status,
		Severity: severity,
		Title:    title,
		Body:     body,
	}
}

func (a *AlarmEntry) AddTransaction(body *AddTransaction, t Transaction) {
	checksAddTransaction(body)
	if te, ok := t.(*TransactionEntry); ok {
		te.AlarmUUID = a.GetUUID() // Set the AlarmUUID to the UUID of the AlarmEntry
		te.SetStatus(body.Status)
		te.SetBody(body.Body)
		te.SetTitle(body.Title)
		te.SetSeverity(body.Severity)
		te.lastUpdated()
		te.createdAt()
		a.Transactions = append(a.Transactions, te)
		a.Status = a.calculateAlarmStatus()
		a.LastUpdated = time.Now()
	}
}

func checksAddTransaction(body *AddTransaction) {
	if body == nil {
		panic("add alarm AddTransaction is empty")
	}
	if body.Status == "" {
		panic("add alarm AddTransaction.Status is empty")
	}
	if body.Severity == "" {
		panic("add alarm AddTransaction.Severity is empty")
	}
	if body.Title == "" {
		panic("add alarm AddTransaction.Title is empty")
	}
	if body.Body == "" {
		panic("add alarm AddTransaction.Body is empty")
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
