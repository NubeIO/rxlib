package priority

import "time"

type DataValue struct {
	Value         any  `json:"value"`
	Disable       bool `json:"disable"`
	OverrideValue any  `json:"overrideValue"`
}

type PreviousValue struct {
	PreviousValue          any       `json:"previousValue,omitempty"`
	PreviousValueRaw       any       `json:"previousValueRaw,omitempty"`
	PreviousValueTimestamp time.Time `json:"previousValueTimestamp,omitempty"`
}

type WrittenValue struct {
	FromUUID   string    `json:"fromUUID"`
	FromPortID string    `json:"fromPortID"`
	Value      any       `json:"previousValue,omitempty"`
	ValueRaw   any       `json:"previousValueRaw,omitempty"`
	Timestamp  time.Time `json:"previousValueTimestamp,omitempty"`
}
