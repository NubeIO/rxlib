package rxlib

import "time"

type ObjectStats struct {
	ObjectStatus ObjectStatus             `json:"status"`
	Loaded       ObjectStatus             `json:"loaded"`
	LoopCount    uint                     `json:"loopCount,omitempty"` // would be how many times the loop of Start() has run
	Custom       map[string]*CustomStatus `json:"custom,omitempty"`
	LastUpdated  *time.Time               `json:"lastUpdated,omitempty"`
	TimeSince    string                   `json:"timeSince,omitempty"`
}

type CustomStatus struct {
	Name  string `json:"name"`
	Field any    `json:"field"`
}
