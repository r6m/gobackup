package config

import "time"

type Report struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Message   string
	Status    string
}
