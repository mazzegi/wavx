package main

import (
	"os"
	"os/signal"

	"github.com/mazzegi/wavx"

	"github.com/mazzegi/log"
)

func main() {
	freqLFO := wavx.NewLFO(wavx.StdOscillatorSine, 0, 10, 3)
	sineModOut := wavx.NewStdOscillator(wavx.StdOscillatorSine, 192, 1.0, 0)
	sineModOut.ConnectInput(wavx.StdOscillatorInputFreqMod, freqLFO)

	squareOsc := wavx.NewStdOscillator(wavx.StdOscillatorSquare, 128, 1.0, 0)
	squareOscAmpl := wavx.NewAmplituder(0.1)
	squareOscAmpl.ConnectInput(wavx.AmplitudeInputSignal, squareOsc)
	squareOscAmpl.ConnectInput(wavx.AmplitudeInputModulation, wavx.NewLFO(wavx.StdOscillatorSine, 0, 0.1, 0.5))

	adder := wavx.NewAdder()
	adder.ConnectInput(wavx.AdderInput, sineModOut)
	adder.ConnectInput(wavx.AdderInput, wavx.NewStdOscillator(wavx.StdOscillatorSaw, 444, 1.0, 0))
	adder.ConnectInput(wavx.AdderInput, squareOscAmpl)

	lowPass := wavx.NewFilter(wavx.FilterModeLowPass, 0.3, 0.5)
	lowPass.ConnectInput(wavx.FilterInputSignal, adder)
	lowPass.ConnectInput(wavx.FilterInputCutoffModulation, wavx.NewLFO(wavx.StdOscillatorSine, 0, 0.2, 3))
	lowPass.ConnectInput(wavx.FilterInputResonanceModulation, wavx.NewLFO(wavx.StdOscillatorSine, 0.0, 0.4, 2))

	dist := wavx.NewDistortion(0.3)
	dist.ConnectInput(wavx.DistortionInputSignal, lowPass)
	dist.ConnectInput(wavx.DistortionInputTresholdModulation, wavx.NewLFO(wavx.StdOscillatorSine, 0, 0.2, 0.5))

	synth := wavx.NewSynthesizer(wavx.SampelRate44100, dist)
	err := synth.Open()
	if err != nil {
		log.Errorf("open-synth: %v", err)
		os.Exit(1)
	}
	defer synth.Close()

	err = synth.Start()
	if err != nil {
		log.Errorf("start-synth: %v", err)
		os.Exit(1)
	}

	sigC := make(chan os.Signal)
	signal.Notify(sigC, os.Kill, os.Interrupt)
	<-sigC
	synth.Stop()
}
