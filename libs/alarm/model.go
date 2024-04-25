package alarm

type Status string
type EntityType string
type Type string
type Severity string
type Target string

const (
	StatusActive       Status = "active"
	StatusAcknowledged Status = "acknowledged"
	StatusClosed       Status = "closed"
)

const (
	EntityTypeGateway EntityType = "gateway"
	EntityTypeNetwork EntityType = "network"
	EntityTypeDevice  EntityType = "device"
	EntityTypePoint   EntityType = "point"
	EntityTypeService EntityType = "service"
)

const (
	TypePing      Type = "ping"
	TypeFault     Type = "fault"
	ThresholdType Type = "threshold"
	FlatLineType  Type = "flat-line"
)

const (
	SeverityCritical Severity = "critical"
	SeverityMinor    Severity = "minor"
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityError    Severity = "error"
)

const (
	TargetMobile Target = "mobile"
	TargetNone   Target = "none"
)
