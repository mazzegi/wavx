package wavx

import (
	"io"
	"os"

	"github.com/mazzegi/log"
	"github.com/pkg/errors"
	"github.com/youpy/go-wav"
)

type WavOutputter struct {
	format  *wav.WavFormat
	samples []wav.Sample
	scale   float64
}

func NewWavOutputter(filePath string) (*WavOutputter, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "open file %q", filePath)
	}
	defer f.Close()

	wr := wav.NewReader(f)
	format, err := wr.Format()
	if err != nil {
		return nil, errors.Wrap(err, "read format")
	}

	o := &WavOutputter{
		format: format,
	}

	var max int
	err = func() error {
		for {
			smpls, err := wr.ReadSamples()
			if err != nil {
				if err == io.EOF {
					return nil
				} else {
					return errors.Wrap(err, "read-samples")
				}
			}
			for _, s := range smpls {
				m := MaxAbsInts(s.Values[:])
				if m > max {
					max = m
				}
			}
			o.samples = append(o.samples, smpls...)
		}
	}()
	if err != nil {
		return nil, err
	}

	o.scale = 1.0 / float64(max)
	log.Infof("format: channels=%d, sample-rate=%d, bits-per-sample=%d, max=%d", format.NumChannels, format.SampleRate, format.BitsPerSample, max)

	return o, nil
}

//sample-rate = 44100
func (o *WavOutputter) Output(secs float64) float64 {
	sampleIdx := RoundInt(secs * float64(o.format.SampleRate))
	sampleIdx = sampleIdx % len(o.samples)

	sample := o.samples[sampleIdx]
	return float64(sample.Values[0]) * o.scale
}
