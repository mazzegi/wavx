package wavl

import "math"

type Note struct {
	name   string
	freq   float64
	octave int
}

func (n Note) Freq() float64 {
	tp := n.octave - 4
	return n.freq * math.Pow(2, float64(tp))
}

// key mapping
var keys2notes = map[rune]Note{
	'y': {"C4", 261.63, 4},
	's': {"C#4", 277.18, 4},
	'x': {"D4", 293.66, 4},
	'd': {"D#4", 311.13, 4},
	'c': {"E4", 329.63, 4},
	'v': {"F4", 349.23, 4},
	'g': {"F#4", 369.99, 4},
	'b': {"G4", 392.00, 4},
	'h': {"G#4", 415.30, 4},
	'n': {"A4", 440.00, 4},
	'j': {"A#4", 466.16, 4},
	'm': {"B4", 493.88, 4},
}
