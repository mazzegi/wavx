package wavx

type LFO struct {
	osc       *StdOscillator
	offset    float64
	amplitude float64
}

func NewLFO(typ StdOscillatorType, offset, amplitude, freq float64) *LFO {
	osc := NewStdOscillator(typ, freq, 1.0, 0)
	return &LFO{
		osc:       osc,
		offset:    offset,
		amplitude: amplitude,
	}
}

func (lfo *LFO) Output(secs float64) float64 {
	return lfo.offset + lfo.amplitude*lfo.osc.Output(secs)
}
