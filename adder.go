package wavx

import (
	"sync"

	"github.com/mazzegi/log"
)

const (
	AdderInput = "signal"
)

type Adder struct {
	inputs []Outputter
	mx     sync.RWMutex
	Activator
}

func NewAdder() *Adder {
	return &Adder{}
}

func (a *Adder) Execute(cmd Command) {

}

func (a *Adder) Inputs() []string {
	return []string{
		AdderInput,
	}
}

func (a *Adder) ConnectInput(input string, op Outputter) {
	switch input {
	case AdderInput:
		a.inputs = append(a.inputs, op)
	default:
		log.Warnf("no such input %q", input)
	}
}

func (a *Adder) Output(secs float64) float64 {
	if !a.IsActive() || len(a.inputs) == 0 {
		return 0
	}
	var sum float64
	for _, in := range a.inputs {
		sum += in.Output(secs)
	}
	return sum / float64(len(a.inputs))
}
