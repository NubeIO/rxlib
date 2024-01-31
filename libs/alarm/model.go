package alarm

import "time"

type Alarm2 struct {
	UUID         string          `json:"uuid"`
	ObjectID     string          `json:"objectID"`             // Device
	ObjectUUID   string          `json:"objectUUID,omitempty"` // dev_abc123
	Type         string          `json:"type"`                 // Ping
	Status       string          `json:"status" gorm:"index"`  // Active
	Notified     bool            `json:"notified,omitempty"`
	NotifiedAt   time.Time       `json:"notified_at,omitempty"`
	CreatedAt    time.Time       `json:"created_at,omitempty"`
	LastUpdated  time.Time       `json:"last_updated,omitempty"`
	Transactions []*Transaction2 `json:"transactions,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

type Transaction2 struct {
	UUID      string     `json:"uuid"`
	AlarmUUID string     `json:"alertUUID"`
	Status    string     `json:"status"`   // Active
	Severity  string     `json:"severity"` // Crucial
	Target    string     `json:"target,omitempty"`
	Title     string     `json:"title,omitempty"`
	Body      string     `json:"body,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

type AlarmStatus string
type AlarmEntityType string
type AlarmType string
type AlarmSeverity string
type AlarmTarget string

const (
	AlarmStatusActive       AlarmStatus = "active"
	AlarmStatusAcknowledged AlarmStatus = "acknowledged"
	AlarmStatusClosed       AlarmStatus = "closed"
)

const (
	AlarmEntityTypeGateway AlarmEntityType = "gateway"
	AlarmEntityTypeNetwork AlarmEntityType = "network"
	AlarmEntityTypeDevice  AlarmEntityType = "device"
	AlarmEntityTypePoint   AlarmEntityType = "point"
	AlarmEntityTypeService AlarmEntityType = "service"
)

const (
	AlarmTypePing      AlarmType = "ping"
	AlarmTypeFault     AlarmType = "fault"
	AlarmTypeThreshold AlarmType = "threshold"
	AlarmTypeFlatLine  AlarmType = "flat-line"
)

const (
	AlarmSeverityCrucial AlarmSeverity = "crucial"
	AlarmSeverityMinor   AlarmSeverity = "minor"
	AlarmSeverityInfo    AlarmSeverity = "info"
	AlarmSeverityWarning AlarmSeverity = "warning"
)

const (
	AlarmTargetMobile AlarmTarget = "mobile"
	AlarmTargetNone   AlarmTarget = "none"
)
