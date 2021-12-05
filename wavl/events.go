package wavl

import "time"

type Event interface{}

type SleepEvent struct {
	time.Duration
}

type ProjectEvent struct {
	projFunc func(p *Project)
}
