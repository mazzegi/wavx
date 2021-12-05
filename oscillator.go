package wavx

import (
	"math"
	"sync"

	"github.com/mazzegi/log"
)

//
type StdOscillatorType string

const (
	StdOscillatorSine     StdOscillatorType = "sine"
	StdOscillatorSquare   StdOscillatorType = "square"
	StdOscillatorSaw      StdOscillatorType = "saw"
	StdOscillatorTriangle StdOscillatorType = "triangle"
)

const (
	StdOscillatorInputFreqMod = "frequency-modulation"
)

type StdOscillatorParams struct {
	Type      StdOscillatorType
	Freq      float64
	Ampl      float64
	Overtones int
}

type StdOscillator struct {
	Params       StdOscillatorParams
	FreqModInput Outputter
	mx           sync.RWMutex
	Muted        bool
}

func NewStdOscillator(typ StdOscillatorType, baseFreq float64, baseAmpl float64, overtones int) *StdOscillator {
	return &StdOscillator{
		Params: StdOscillatorParams{
			Type:      typ,
			Freq:      baseFreq,
			Ampl:      baseAmpl,
			Overtones: overtones,
		},
	}
}

func (o *StdOscillator) Activate() {
	o.Unmute()
}

func (o *StdOscillator) Deactivate() {
	o.Mute()
}

func (o *StdOscillator) Execute(cmd Command) {
	params := o.Parameters()
	err := ApplyCommand(cmd, &params)
	if err != nil {
		log.Warnf("apply-command: %v", err)
		return
	}
	o.ChangeParameters(params)
}

func (o *StdOscillator) Inputs() []string {
	return []string{
		StdOscillatorInputFreqMod,
	}
}

func (o *StdOscillator) ConnectInput(input string, op Outputter) {
	switch input {
	case StdOscillatorInputFreqMod:
		o.FreqModInput = op
	default:
		log.Warnf("no such input %q", input)
	}
}

func (o *StdOscillator) IsMuted() bool {
	o.mx.RLock()
	defer o.mx.RUnlock()
	return o.Muted
}

func (o *StdOscillator) Mute() {
	o.mx.Lock()
	defer o.mx.Unlock()
	o.Muted = true
}

func (o *StdOscillator) Unmute() {
	o.mx.Lock()
	defer o.mx.Unlock()
	o.Muted = false
}

//
func (o *StdOscillator) Parameters() StdOscillatorParams {
	o.mx.RLock()
	defer o.mx.RUnlock()
	return o.Params
}

func (o *StdOscillator) ChangeParameters(params StdOscillatorParams) {
	o.mx.Lock()
	defer o.mx.Unlock()
	o.Params = params
}

func (o *StdOscillator) ChangeType(typ StdOscillatorType) {
	o.mx.Lock()
	defer o.mx.Unlock()
	o.Params.Type = typ
}

func (o *StdOscillator) ChangeFreq(freq float64) {
	o.mx.Lock()
	defer o.mx.Unlock()
	o.Params.Freq = freq
}

func (o *StdOscillator) ChangeAmpl(ampl float64) {
	o.mx.Lock()
	defer o.mx.Unlock()
	o.Params.Ampl = ampl
}

func (o *StdOscillator) ChangeOvertones(n int) {
	o.mx.Lock()
	defer o.mx.Unlock()
	o.Params.Overtones = n
}

//

func (o *StdOscillator) normalizedInPeriod(secs float64, freq, freqMod float64) float64 {
	x := secs*freq + freqMod
	x = x - math.Floor(x)
	return x
}

func (o *StdOscillator) Output(secs float64) float64 {
	if o.IsMuted() {
		return 0
	}

	var freqMod float64
	if o.FreqModInput != nil {
		freqMod = o.FreqModInput.Output(secs)
	}
	params := o.Parameters()

	v := o.calc(secs, params.Type, params.Freq, freqMod, params.Ampl)
	for i := 0; i < params.Overtones; i++ {
		vo := o.calc(secs, params.Type, params.Freq*float64(i+2), 0, params.Ampl)
		v += vo
	}
	if params.Overtones > 0 {
		v /= float64(params.Overtones)
	}

	return v
}

func (o *StdOscillator) calc(secs float64, typ StdOscillatorType, freq, freqMod, ampl float64) float64 {
	var v float64
	switch typ {
	case StdOscillatorSine:
		v = math.Sin(2.0*math.Pi*secs*freq + freqMod)
	case StdOscillatorSquare:
		x := o.normalizedInPeriod(secs, freq, freqMod)
		if x > 0.5 {
			v = 1
		} else {
			v = -1
		}
	case StdOscillatorSaw:
		x := o.normalizedInPeriod(secs, freq, freqMod)
		v = -1.0 + 2*x
	case StdOscillatorTriangle:
		x := o.normalizedInPeriod(secs, freq, freqMod)
		if x < 0.5 {
			v = 1 - x*4
		} else {
			v = -1 + (x-0.5)*4
		}
	default:
		v = 0
	}
	return ampl * v
}
