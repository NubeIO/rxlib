package alarm

import (
	"github.com/NubeIO/rxlib/helpers"
	"time"
)

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
