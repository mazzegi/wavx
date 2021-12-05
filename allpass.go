package wavx

type AllPassFilter struct {
	Gain         float64
	DelayMsec    float64
	SamplingRate float64
	bufferX      *LinearBuffer
	bufferY      *LinearBuffer
}

func NewAllPassFilter(gain, delayMsec, samplingRate float64) *AllPassFilter {
	delSamples := int(samplingRate / 1000.0 * delayMsec)
	return &AllPassFilter{
		Gain:         gain,
		DelayMsec:    delayMsec,
		SamplingRate: samplingRate,
		bufferX:      NewLinearBuffer(delSamples),
		bufferY:      NewLinearBuffer(delSamples),
	}
}

func (f *AllPassFilter) Next(x float64) float64 {
	y := -f.Gain*x + f.bufferX.FirstValue() + f.Gain*f.bufferY.FirstValue()
	f.bufferX.Push(x)
	f.bufferY.Push(y)
	return y
}
