package wavl

import (
	"context"
	"time"

	"github.com/mazzegi/log"
)

type Looper struct {
	events []Event
}

func NewLooper(events []Event) *Looper {
	return &Looper{
		events: events,
	}
}

func (l *Looper) Run(ctx context.Context, prj *Project) {
	if len(l.events) == 0 {
		log.Infof("looper: no events")
		return
	}
	curr := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
			l.handle(ctx, prj, l.events[curr])
			curr++
			if curr >= len(l.events) {
				curr = 0
			}
		}

	}
}

func (l *Looper) handle(ctx context.Context, prj *Project, e Event) {
	switch e := e.(type) {
	case SleepEvent:
		l.handleSleep(ctx, e)
	case ProjectEvent:
		e.projFunc(prj)
	}
}

func (l *Looper) handleSleep(ctx context.Context, e SleepEvent) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(e.Duration):
		return
	}
}
