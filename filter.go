package wavx

import (
	"sync"

	"github.com/mazzegi/log"
)

type FilterMode string

const (
	FilterModeLowPass  FilterMode = "lowpass"
	FilterModeHighPass FilterMode = "highpass"
	FilterModeBandPass FilterMode = "bandpass"
)

type FilterParams struct {
	Mode      FilterMode
	Cutoff    float64
	Resonance float64
}

type Filter struct {
	mx                sync.RWMutex
	Params            FilterParams
	inputSignal       Outputter
	inputCutoffMod    Outputter
	inputResonanceMod Outputter
	buf0, buf1        float64
	Activator
}

func NewFilter(mode FilterMode, baseCutoff float64, baseResonance float64) *Filter {
	return &Filter{
		Params: FilterParams{
			Mode:      mode,
			Cutoff:    baseCutoff,
			Resonance: baseResonance,
		},
	}
}

const (
	FilterInputSignal              = "signal"
	FilterInputCutoffModulation    = "cutoff-modulation"
	FilterInputResonanceModulation = "resonance-modulation"
)

func (f *Filter) Inputs() []string {
	return []string{
		FilterInputSignal,
		FilterInputCutoffModulation,
		FilterInputResonanceModulation,
	}
}

func (a *Filter) Execute(cmd Command) {
	params := a.Parameters()
	err := ApplyCommand(cmd, &params)
	if err != nil {
		log.Warnf("apply-command: %v", err)
		return
	}
	a.ChangeParameters(params)
	log.Infof("filter: cmd %s => %v", cmd, params)
}

func (f *Filter) ConnectInput(input string, op Outputter) {
	switch input {
	case FilterInputSignal:
		f.inputSignal = op
	case FilterInputCutoffModulation:
		f.inputCutoffMod = op
	case FilterInputResonanceModulation:
		f.inputResonanceMod = op
	default:
		log.Warnf("no such input %q", input)
	}
}

func (f *Filter) Parameters() FilterParams {
	f.mx.RLock()
	defer f.mx.RUnlock()
	return f.Params
}

func (f *Filter) ChangeParameters(params FilterParams) {
	f.mx.Lock()
	defer f.mx.Unlock()
	f.Params = params
}

func (f *Filter) ChangeMode(mode FilterMode) {
	f.mx.Lock()
	defer f.mx.Unlock()
	f.Params.Mode = mode
}

func (f *Filter) ChangeCutoff(cutoff float64) {
	f.mx.Lock()
	defer f.mx.Unlock()
	f.Params.Cutoff = cutoff
}

func (f *Filter) ChangeResonance(res float64) {
	f.mx.Lock()
	defer f.mx.Unlock()
	f.Params.Resonance = res
}

func (f *Filter) Output(secs float64) float64 {
	if f.inputSignal == nil {
		return 0
	}
	newVal := f.inputSignal.Output(secs)
	if !f.IsActive() {
		return newVal
	}

	params := f.Parameters()
	mode, cutoff, resonance := params.Mode, params.Cutoff, params.Resonance
	if f.inputCutoffMod != nil {
		cutoff += f.inputCutoffMod.Output(secs)
	}
	if f.inputResonanceMod != nil {
		resonance += f.inputResonanceMod.Output(secs)
	}
	if cutoff < 0 {
		cutoff = 0
	} else if cutoff >= 1.0 {
		cutoff = 0.99
	}
	feedback := resonance + resonance/(1.0-cutoff)

	f.buf0 += cutoff * (newVal - f.buf0 + feedback*(f.buf0-f.buf1))
	f.buf1 += cutoff * (f.buf0 - f.buf1)
	switch mode {
	case FilterModeLowPass:
		return f.buf1
	case FilterModeHighPass:
		return newVal - f.buf0
	case FilterModeBandPass:
		return f.buf0 - f.buf1
	default:
		return newVal
	}
}
