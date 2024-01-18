package rxlib

import "time"

type ObjectStats struct {
	ObjectStatus ObjectStatus             `json:"objectStatus"`
	LoopCount    int                      `json:"loopCount"` // would be how many times the loop of Start() has run
	Custom       map[string]*CustomStatus `json:"custom"`
	LastUpdated  time.Time                `json:"lastUpdated"`
	TimeSince    string                   `json:"timeSince"`
}

type CustomStatus struct {
	Name  string `json:"name"`
	Field any    `json:"field"`
}

type StatsBuilder interface {
	SetStatus(status ObjectStatus) StatsBuilder
	SetLoopCount(count int) StatsBuilder
	GetStats() *ObjectStats
	AddCustomStat(name string, stat *CustomStatus) StatsBuilder
	GetCustomStat(name string) (*CustomStatus, bool)
	DeleteCustomStat(name string)
	UpdateCustomStat(name string, stat *CustomStatus)
}

type statsBuilder struct {
	ObjectStats *ObjectStats
}

func NewStatsBuilder() StatsBuilder {
	return &statsBuilder{
		ObjectStats: &ObjectStats{},
	}
}

func (builder *statsBuilder) SetStatus(status ObjectStatus) StatsBuilder {
	builder.ObjectStats.ObjectStatus = status
	builder.ObjectStats.LastUpdated = time.Now()
	return builder
}

func (builder *statsBuilder) SetLoopCount(count int) StatsBuilder {
	builder.ObjectStats.LoopCount = count
	return builder
}

func (builder *statsBuilder) GetStats() *ObjectStats {
	builder.ObjectStats.TimeSince = TimeSince(builder.ObjectStats.LastUpdated)
	return builder.ObjectStats
}

func (builder *statsBuilder) AddCustomStat(name string, stat *CustomStatus) StatsBuilder {
	if builder.ObjectStats.Custom == nil {
		builder.ObjectStats.Custom = make(map[string]*CustomStatus)
	}
	builder.ObjectStats.Custom[name] = stat
	return builder
}

func (builder *statsBuilder) GetCustomStat(name string) (*CustomStatus, bool) {
	stat, ok := builder.ObjectStats.Custom[name]
	return stat, ok
}

func (builder *statsBuilder) DeleteCustomStat(name string) {
	delete(builder.ObjectStats.Custom, name)
}

func (builder *statsBuilder) UpdateCustomStat(name string, stat *CustomStatus) {
	if _, ok := builder.ObjectStats.Custom[name]; ok {
		builder.ObjectStats.Custom[name] = stat
	}

}
