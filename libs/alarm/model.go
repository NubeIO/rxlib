package alarm

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
	AlarmSeverityCritical AlarmSeverity = "critical"
	AlarmSeverityMinor    AlarmSeverity = "minor"
	AlarmSeverityInfo     AlarmSeverity = "info"
	AlarmSeverityWarning  AlarmSeverity = "warning"
	AlarmSeverityError    AlarmSeverity = "error"
)

const (
	AlarmTargetMobile AlarmTarget = "mobile"
	AlarmTargetNone   AlarmTarget = "none"
)
