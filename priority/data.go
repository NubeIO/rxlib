package priority

import "time"

type PreviousValue struct {
	Value     any       `json:"value,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

type WrittenValue struct {
	FromUUID   string    `json:"fromUUID"`
	FromPortID string    `json:"fromPortID"`
	Value      any       `json:"previousValue,omitempty"`
	ValueRaw   any       `json:"previousValueRaw,omitempty"`
	Timestamp  time.Time `json:"previousValueTimestamp,omitempty"`
}
