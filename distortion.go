package wavx

import (
	"github.com/mazzegi/log"
)

type Distortion struct {
	baseThreshold     float64
	inputSignal       Outputter
	inputThresholdMod Outputter
}

func NewDistortion(baseThreshold float64) *Distortion {
	return &Distortion{
		baseThreshold: baseThreshold,
	}
}

const (
	DistortionInputSignal             = "signal"
	DistortionInputTresholdModulation = "threshold-modulation"
)

func (d *Distortion) Inputs() []string {
	return []string{
		DistortionInputSignal,
		DistortionInputTresholdModulation,
	}
}

func (d *Distortion) ConnectInput(input string, op Outputter) {
	switch input {
	case DistortionInputSignal:
		d.inputSignal = op
	case DistortionInputTresholdModulation:
		d.inputThresholdMod = op
	default:
		log.Warnf("no such input %q", input)
	}
}

func (d *Distortion) Output(secs float64) float64 {
	if d.inputSignal == nil {
		return 0
	}
	newVal := d.inputSignal.Output(secs)
	threshold := d.baseThreshold
	if d.inputThresholdMod != nil {
		threshold += d.inputThresholdMod.Output(secs)
	}
	if threshold < 0 {
		threshold = 0
	} else if threshold > 1 {
		threshold = 1
	}

	if newVal > threshold {
		return threshold
	} else if newVal < -threshold {
		return -threshold
	}
	return newVal
}
