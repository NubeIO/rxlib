package rxlib

import "time"

type EventbusOpts struct {
	TTL time.Duration `json:"ttl"` // time to live timeout in seconds
}
