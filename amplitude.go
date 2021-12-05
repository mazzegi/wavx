package wavx

import "github.com/mazzegi/log"

type Amplituder struct {
	baseAmplitude   float64
	inputSignal     Outputter
	inputModulation Outputter
}

func NewAmplituder(baseAmpl float64) *Amplituder {
	return &Amplituder{
		baseAmplitude: baseAmpl,
	}
}

const (
	AmplitudeInputSignal     = "signal"
	AmplitudeInputModulation = "modulation"
)

func (a *Amplituder) Inputs() []string {
	return []string{
		AmplitudeInputSignal,
		AmplitudeInputModulation,
	}
}

func (a *Amplituder) ConnectInput(input string, op Outputter) {
	switch input {
	case AmplitudeInputSignal:
		a.inputSignal = op
	case AmplitudeInputModulation:
		a.inputModulation = op
	default:
		log.Warnf("no such input %q", input)
	}
}

func (a *Amplituder) Output(secs float64) float64 {
	if a.inputSignal == nil {
		return 0
	}
	ampl := a.baseAmplitude
	if a.inputModulation != nil {
		ampl += a.inputModulation.Output(secs)
	}
	return ampl * a.inputSignal.Output(secs)
}
