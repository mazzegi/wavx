package wavx

type Chorus struct {
	inputSignal Outputter
	inputLFO    Outputter
	delay       float64
	sampleRate  float64
	comb        *CombFilter
}

func NewChorus(sampleRate float64, delayMSec float64) *Chorus {
	return &Chorus{
		comb:       NewCombFilter(0.5, delayMSec, sampleRate),
		sampleRate: sampleRate,
		delay:      delayMSec,
	}
}

func (d *Chorus) ConnectInput(op Outputter) {
	d.inputSignal = op
}

func (d *Chorus) ConnectInputLFO(op Outputter) {
	d.inputLFO = op
}

func (d *Chorus) Output(secs float64) float64 {
	if d.inputSignal == nil {
		return 0
	}
	if d.inputLFO != nil {
		lfo := d.inputLFO.Output(secs)
		d.comb.Change(lfo, d.sampleRate)
	}
	newVal := d.inputSignal.Output(secs)
	c := d.comb.Next(newVal)
	return c
}
