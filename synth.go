package wavx

import (
	"github.com/gordonklaus/portaudio"
	"github.com/mazzegi/log"
	"github.com/pkg/errors"
)

type Outputter interface {
	Output(secs float64) float64
}

type InputOutputter interface {
	Outputter
	Activate()
	Deactivate()
	Execute(cmd Command)
	Inputs() []string
	ConnectInput(input string, op Outputter)
}

const SampelRate44100 = 44100

type Synthesizer struct {
	stream     *portaudio.Stream
	steps      uint64
	sampleRate int
	outputter  Outputter
}

func NewSynthesizer(sampleRate int, outputter Outputter) *Synthesizer {
	s := &Synthesizer{
		sampleRate: sampleRate,
		outputter:  outputter,
	}

	return s
}

func (s *Synthesizer) Open() error {
	err := portaudio.Initialize()
	if err != nil {
		return errors.Wrap(err, "portaudio: initialize")
	}

	s.stream, err = portaudio.OpenDefaultStream(0, 1, float64(s.sampleRate), 0, func(out [][]float32) {
		for i := range out[0] {
			out[0][i] = s.Next()
		}
	})
	if err != nil {
		s.Close()
		return errors.Wrap(err, "portaudio: open-default-stream")
	}
	return nil
}

func (s *Synthesizer) Close() error {
	err := s.Stop()
	if err != nil {
		log.Errorf("portaudio: stop-stream: %v", err)
	}
	if s.stream != nil {
		err := s.stream.Close()
		if err != nil {
			log.Errorf("portaudio: close-stream: %v", err)
		}
		s.stream = nil
	}

	err = portaudio.Terminate()
	if err != nil {
		return errors.Wrap(err, "portaudio: terminate")
	}
	return nil
}

func (s *Synthesizer) Start() error {
	if s.stream == nil {
		return errors.Errorf("stream is not open")
	}
	return s.stream.Start()
}

func (s *Synthesizer) Stop() error {
	if s.stream == nil {
		return errors.Errorf("stream is not open")
	}
	return s.stream.Stop()
}

func (s *Synthesizer) Next() float32 {
	secs := float64(s.steps) / float64(s.sampleRate)
	v := s.outputter.Output(secs)
	s.steps++
	return float32(v)
}
