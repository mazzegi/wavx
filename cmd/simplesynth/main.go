package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mazzegi/wavx"

	"github.com/mazzegi/log"
)

func main() {
	osc1 := wavx.NewStdOscillator(wavx.StdOscillatorSine, 192, 1.0, 2)
	osc2 := wavx.NewStdOscillator(wavx.StdOscillatorSaw, 192, 1.0, 0)

	adder := wavx.NewAdder()
	adder.ConnectInput(wavx.AdderInput, osc1)
	adder.ConnectInput(wavx.AdderInput, osc2)

	filter := wavx.NewFilter(wavx.FilterModeLowPass, 0.3, 0.5)
	filter.ConnectInput(wavx.FilterInputSignal, adder)

	reverb := wavx.NewReverb(44100)
	reverb.ConnectInput(filter)

	synth := wavx.NewSynthesizer(wavx.SampelRate44100, reverb)
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

	log.Infof("synth started ... ")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("> ")
	for scanner.Scan() {
		cmd := strings.Trim(scanner.Text(), " \r\n")
		fs := strings.Fields(cmd)
		if len(fs) < 3 {
			continue
		}
		comp := fs[0]
		compProp := fs[1]
		compPropValue := fs[2]
		switch {
		case comp == "osc1" && compProp == "tone" && compPropValue == "on":
			osc1.Unmute()
		case comp == "osc1" && compProp == "tone" && compPropValue == "off":
			osc1.Mute()
		case comp == "osc1" && compProp == "type":
			osc1.ChangeType(wavx.StdOscillatorType(compPropValue))
		case comp == "osc1" && compProp == "freq":
			v, _ := strconv.ParseFloat(compPropValue, 64)
			osc1.ChangeFreq(v)
		case comp == "osc1" && compProp == "ampl":
			v, _ := strconv.ParseFloat(compPropValue, 64)
			osc1.ChangeAmpl(v)

		case comp == "osc2" && compProp == "tone" && compPropValue == "on":
			osc2.Unmute()
		case comp == "osc2" && compProp == "tone" && compPropValue == "off":
			osc2.Mute()
		case comp == "osc2" && compProp == "type":
			osc2.ChangeType(wavx.StdOscillatorType(compPropValue))
		case comp == "osc2" && compProp == "freq":
			v, _ := strconv.ParseFloat(compPropValue, 64)
			osc2.ChangeFreq(v)
		case comp == "osc2" && compProp == "ampl":
			v, _ := strconv.ParseFloat(compPropValue, 64)
			osc2.ChangeAmpl(v)

		case comp == "osc1" && compProp == "overtones":
			v, _ := strconv.ParseInt(compPropValue, 10, 64)
			osc1.ChangeOvertones(int(v))
		case comp == "osc2" && compProp == "overtones":
			v, _ := strconv.ParseInt(compPropValue, 10, 64)
			osc2.ChangeOvertones(int(v))

		case comp == "filter" && compProp == "mode":
			filter.ChangeMode(wavx.FilterMode(compPropValue))
		case comp == "filter" && compProp == "cutoff":
			v, _ := strconv.ParseFloat(compPropValue, 64)
			filter.ChangeCutoff(v)
		case comp == "filter" && compProp == "res":
			v, _ := strconv.ParseFloat(compPropValue, 64)
			filter.ChangeResonance(v)

		}
		fmt.Printf("> ")
	}

	synth.Stop()
	log.Infof("synth stopped")
}
