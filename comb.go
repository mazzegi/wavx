package wavx

type CombFilter struct {
	Gain         float64
	DelayMsec    float64
	SamplingRate float64
	buffer       *LinearBuffer
}

func NewCombFilter(gain, delayMsec, samplingRate float64) *CombFilter {
	delSamples := int(samplingRate / 1000.0 * delayMsec)
	return &CombFilter{
		Gain:         gain,
		DelayMsec:    delayMsec,
		SamplingRate: samplingRate,
		buffer:       NewLinearBuffer(delSamples),
	}
}

func (f *CombFilter) Change(delayMsec, samplingRate float64) {
	f.buffer.ChangeMaxSize(int(samplingRate / 1000.0 * delayMsec))
}

func (f *CombFilter) Next(x float64) float64 {
	y := x + f.Gain*f.buffer.FirstValue()
	f.buffer.Push(y)
	return y
}
