package wavx

type Reverb struct {
	inputSignal Outputter
	comb1       *CombFilter
	comb2       *CombFilter
	comb3       *CombFilter
	comb4       *CombFilter
	allPass1    *AllPassFilter
	allPass2    *AllPassFilter
}

func NewReverb(sampleRate float64) *Reverb {
	delay := 78.9
	decay := 0.45
	return &Reverb{
		comb1:    NewCombFilter(decay, delay, sampleRate),
		comb2:    NewCombFilter(decay-0.1313, delay-11.73, sampleRate),
		comb3:    NewCombFilter(decay-0.2743, delay+19.31, sampleRate),
		comb4:    NewCombFilter(decay-0.31, delay-7.79, sampleRate),
		allPass1: NewAllPassFilter(0.13, 89.27, sampleRate),
		allPass2: NewAllPassFilter(0.13, 89.27, sampleRate),
	}
}

func (d *Reverb) ConnectInput(op Outputter) {
	d.inputSignal = op

}

func (d *Reverb) Output(secs float64) float64 {
	if d.inputSignal == nil {
		return 0
	}
	newVal := d.inputSignal.Output(secs)
	c1 := d.comb1.Next(newVal)
	c2 := d.comb2.Next(newVal)
	c3 := d.comb3.Next(newVal)
	c4 := d.comb4.Next(newVal)
	c := (c1 + c2 + c3 + c4) / 4

	c = d.allPass1.Next(c)
	c = d.allPass2.Next(c)

	return c
}
