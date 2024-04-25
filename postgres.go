package rxlib

import "time"

type SyncOptions struct {
	AutoSync bool
	Duration time.Duration
}

func (inst *RuntimeImpl) ObjectSync(forceSync bool, opts *SyncOptions) error {
	//TODO implement me
	panic("implement me")
}

func (inst *RuntimeImpl) HistorySync(forceSync bool, opts *SyncOptions) error {
	//TODO implement me
	panic("implement me")
}
