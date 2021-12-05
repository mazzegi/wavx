package main

import (
	"os"
	"os/signal"

	"github.com/mazzegi/wavx"

	"github.com/mazzegi/log"
)

func handleErr(err error) {
	if err == nil {
		return
	}
	log.Errorf("%v", err)
	os.Exit(1)
}

func main() {
	wo, err := wavx.NewWavOutputter("../../samples/ihaveno.wav")
	handleErr(err)

	// lowPass := wavx.NewFilter(wavx.FilterModeLowPass, 0.1, 0.8)
	// lowPass.ConnectInput(wavx.FilterInputSignal, wo)

	// dist := wavx.NewDistortion(0.1)
	// dist.ConnectInput(wavx.DistortionInputSignal, lowPass)

	// reverb := wavx.NewReverb(44100)
	// reverb.ConnectInput(wo)
	lfo := wavx.NewLFO(wavx.StdOscillatorSine, 2, 0.5, 50)

	mfin := wavx.NewChorus(44100, 200)
	mfin.ConnectInput(wo)
	mfin.ConnectInputLFO(lfo)

	// reverb := wavx.NewReverb(44100)
	// reverb.ConnectInput(mfin)

	synth := wavx.NewSynthesizer(wavx.SampelRate44100, mfin)
	err = synth.Open()
	handleErr(err)

	defer synth.Close()

	err = synth.Start()
	handleErr(err)

	sigC := make(chan os.Signal)
	signal.Notify(sigC, os.Kill, os.Interrupt)
	<-sigC
	synth.Stop()
}
