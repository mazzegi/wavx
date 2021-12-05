package wavx

import (
	"fmt"
	"testing"
)

func TestEnvelope(t *testing.T) {
	env := NewFixedSustainEnvelope(
		EnvelopeParams{
			Attack:  0.2,
			Decay:   0.3,
			Sustain: 0.7,
			Release: 4,
		},
		2.0,
	)
	s := float64(0)
	env.Start(s)
	for ; s < 10.0; s += 0.1 {
		fmt.Printf("%.2f sec => %.4f\n", s, env.Value(s))
		if !env.IsActive() {
			fmt.Printf("envelope is done\n")
			break
		}
	}
}
