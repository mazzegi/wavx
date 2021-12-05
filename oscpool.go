package wavx

import (
	"math"
	"sync"

	"github.com/mazzegi/log"
)

type EnveloppedOscillator struct {
	env *FixedSustainEnvelope
	osc *StdOscillator
}

type OscillatorPool struct {
	sync.RWMutex
	Activator
	oscis         map[*EnveloppedOscillator]bool
	defaultParams StdOscillatorParams
}

func NewOscillatorPool() *OscillatorPool {
	return &OscillatorPool{
		oscis: map[*EnveloppedOscillator]bool{},
		defaultParams: StdOscillatorParams{
			Type:      StdOscillatorSaw,
			Ampl:      1.0,
			Overtones: 4,
		},
	}
}

func (p *OscillatorPool) Execute(cmd Command) {
	params := p.defaultParams
	err := ApplyCommand(cmd, &params)
	if err != nil {
		log.Warnf("apply command to std-oscillator-params")
		return
	}
	osc := NewStdOscillator(params.Type, params.Freq, params.Ampl, params.Overtones)
	env := NewFixedSustainEnvelope(
		EnvelopeParams{
			Attack:  0.8,
			Decay:   0.2,
			Sustain: 1.0,
			Release: 1.0,
		},
		0.0,
	)
	p.Lock()
	defer p.Unlock()
	eosc := &EnveloppedOscillator{
		env: env,
		osc: osc,
	}
	p.oscis[eosc] = true

	doDetune := false
	if doDetune {
		i := 0.0
		detuneStart := 2.0
		for eosc := range p.oscis {
			detune := detuneStart + i
			params := eosc.osc.Parameters()
			params.Freq += detune
			eosc.osc.ChangeParameters(params)
			i += 1.0
		}
	}

	log.Infof("added oscillator")
}

func (p *OscillatorPool) Inputs() []string {
	return []string{}
}

func (a *OscillatorPool) ConnectInput(input string, op Outputter) {
	log.Warnf("no inputs accepted")
}

func (a *OscillatorPool) Output(secs float64) float64 {
	a.RLock()
	defer a.RUnlock()
	if len(a.oscis) == 0 {
		return 0
	}
	var sum float64
	var cnt float64

	for eosc := range a.oscis {
		if !eosc.env.IsStarted() {
			eosc.env.Start(secs)
		}

		ov := eosc.osc.Output(secs)
		ev := eosc.env.Value(secs)
		sum += ov * ev
		cnt++
		if !eosc.env.IsActive() {
			delete(a.oscis, eosc)
			log.Infof("removed oscillator")
		}
	}

	var v float64
	if cnt > 0 {
		v = sum / math.Sqrt(cnt)
	}
	return v
}
