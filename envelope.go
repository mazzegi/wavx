package wavx

import "math"

const envEps = float64(0.0001)

// EnvelopeParams holds all envelope params; attack, decay and release are in seconds
type EnvelopeParams struct {
	Attack  float64
	Decay   float64
	Sustain float64
	Release float64
}

type Envelope struct {
	params     EnvelopeParams
	isActive   bool
	startedAt  float64
	isStarted  bool
	releasedAt float64
	isReleased bool
}

func NewEnvelope(ps EnvelopeParams) *Envelope {
	return &Envelope{
		params: ps,
	}
}

func (e *Envelope) IsActive() bool {
	return e.isActive
}

func (e *Envelope) IsStarted() bool {
	return e.isStarted
}

func (e *Envelope) Start(secs float64) {
	e.startedAt = secs
	e.isStarted = true
	e.isActive = true
}

func (e *Envelope) Release(secs float64) {
	if !e.isActive {
		return
	}
	e.releasedAt = secs
	e.isReleased = true
}

func (e *Envelope) Value(secs float64) float64 {
	if e.isReleased {
		dur := secs - e.releasedAt
		v := -e.params.Sustain*dur/e.params.Release + e.params.Sustain
		if math.Abs(v) < envEps {
			e.isActive = false
		}
		return v
	} else if e.isStarted {
		dur := secs - e.startedAt
		switch {
		case dur < e.params.Attack:
			return dur / e.params.Attack
		case dur < e.params.Attack+e.params.Decay:
			return 1.0 - (1.0-e.params.Sustain)*(dur-e.params.Attack)/e.params.Decay
		default:
			return e.params.Sustain
		}
	} else {
		return 0
	}
}

//
type FixedSustainEnvelope struct {
	*Envelope
	sustainDuration float64
}

func NewFixedSustainEnvelope(ps EnvelopeParams, sustainDuration float64) *FixedSustainEnvelope {
	return &FixedSustainEnvelope{
		Envelope:        NewEnvelope(ps),
		sustainDuration: sustainDuration,
	}
}

func (e *FixedSustainEnvelope) Value(secs float64) float64 {
	if !e.isActive {
		return 0
	}
	if !e.isReleased {
		//auto release after sustain duration
		dur := secs - e.startedAt
		if dur >= e.params.Attack+e.params.Decay+e.sustainDuration {
			e.Release(secs)
		}
	}
	return e.Envelope.Value(secs)
}
