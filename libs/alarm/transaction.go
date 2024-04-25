package alarm

import (
	"github.com/NubeIO/rxlib/helpers"
	"time"
)

type Transaction interface {
	GetAlarmUUID() string
	GetUUID() string
	GetStatus() Status
	GetSeverity() Severity
	GetTarget() string
	GetTitle() string
	GetBody() string
	GetCreatedAt() time.Time
	GetLastUpdated() time.Time

	SetTitle(title string)
	SetStatus(state Status)
	SetSeverity(s Severity)
	SetBody(body string)
}

type TransactionEntry struct {
	Title       string    `json:"title,omitempty"`
	Status      Status    `json:"status"`   // Active
	Severity    Severity  `json:"severity"` // Crucial
	Target      string    `json:"target,omitempty"`
	Body        string    `json:"body,omitempty"`
	UUID        string    `json:"uuid"`
	AlarmUUID   string    `json:"alarmUUID"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
}

func (t *TransactionEntry) SetTitle(title string) {
	t.Title = title
}

func (t *TransactionEntry) SetSeverity(s Severity) {
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

func (t *TransactionEntry) SetStatus(status Status) {
	t.lastUpdated()
	t.Status = status
}

func (t *TransactionEntry) GetStatus() Status {
	return t.Status
}

func (t *TransactionEntry) GetSeverity() Severity {
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
