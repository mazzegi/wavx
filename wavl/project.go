package wavl

import (
	"context"
	"fmt"

	"github.com/mazzegi/wavx"
	"github.com/mazzegi/wavx/wavl/keys"

	"github.com/mazzegi/log"
	"github.com/pkg/errors"
)

type KeyBinding struct {
	Key       rune
	Component wavx.InputOutputter
	Command   wavx.Command
}

type Project struct {
	sampleRate      int
	components      map[string]wavx.InputOutputter
	outputFrom      wavx.InputOutputter
	synth           *wavx.Synthesizer
	events          []Event
	assignedKeyComp wavx.InputOutputter
	keyBindings     map[rune]KeyBinding
}

func NewProject() *Project {
	p := &Project{
		sampleRate:  wavx.SampelRate44100,
		components:  map[string]wavx.InputOutputter{},
		keyBindings: map[rune]KeyBinding{},
	}

	return p
}

func (p *Project) addComponent(name string, comp wavx.InputOutputter) error {
	if _, ok := p.components[name]; ok {
		return errors.Errorf("components with name %q already exists", name)
	}
	p.components[name] = comp
	return nil
}

func (p *Project) AddOscillatorPool(name string) error {
	return p.addComponent(name, wavx.NewOscillatorPool())
}

func (p *Project) AddOscillator(name string, typ string, freq float64, ampl float64, overtones int) error {
	return p.addComponent(name, wavx.NewStdOscillator(wavx.StdOscillatorType(typ), freq, ampl, overtones))
}

func (p *Project) AddMixer(name string) error {
	return p.addComponent(name, wavx.NewAdder())
}

func (p *Project) AddFilter(name string, typ string, cutoff float64, resonance float64) error {
	return p.addComponent(name, wavx.NewFilter(wavx.FilterMode(typ), cutoff, resonance))
}

func (p *Project) Connect(fromName string, toName string, input string) error {
	from, ok := p.components[fromName]
	if !ok {
		return errors.Errorf("no such component %q", fromName)
	}
	to, ok := p.components[toName]
	if !ok {
		return errors.Errorf("no such component %q", toName)
	}
	to.ConnectInput(input, from)
	return nil
}

func (p *Project) OutputFrom(name string) error {
	in, ok := p.components[name]
	if !ok {
		return errors.Errorf("no such component %q", name)
	}
	p.outputFrom = in
	return nil
}

func (p *Project) AddEvents(es ...Event) {
	p.events = append(p.events, es...)
}

func (p *Project) Start() error {
	if p.synth != nil {
		return errors.Errorf("synth is already running")
	}
	if p.outputFrom == nil {
		return errors.Errorf("no output from set")
	}
	p.synth = wavx.NewSynthesizer(p.sampleRate, p.outputFrom)
	err := p.synth.Open()
	if err != nil {
		return errors.Wrap(err, "open synth")
	}
	err = p.synth.Start()
	if err != nil {
		return errors.Wrap(err, "start synth")
	}
	return nil
}

func (p *Project) Loop(ctx context.Context) {
	l := NewLooper(p.events)
	go l.Run(ctx, p)
	if p.assignedKeyComp != nil {
		active := false
		p.assignedKeyComp.Deactivate()

		kl, err := keys.NewListener()
		if err != nil {
			panic(fmt.Sprintf("new key-listener: %v", err))
		}
		log.Infof("listen to events")
		kc := kl.ListenCtx(ctx)
		var currKey rune
		oct := 4
		for e := range kc {
			log.Infof("evt: %q", string(e))
			if e == currKey {
				p.assignedKeyComp.Deactivate()
				active = false
				currKey = rune(0)
				continue
			}

			octSwitch := true
			switch e {
			case '1':
				oct = 1
			case '2':
				oct = 2
			case '3':
				oct = 3
			case '4':
				oct = 4
			case '5':
				oct = 5
			case '6':
				oct = 6
			case '7':
				oct = 7
			case '8':
				oct = 8
			case '9':
				oct = 9
			default:
				octSwitch = false
			}

			if note, ok := keys2notes[e]; ok {
				if !active {
					p.assignedKeyComp.Activate()
					active = true
				}
				note.octave = oct
				p.assignedKeyComp.Execute(wavx.Command{
					"freq": fmt.Sprintf("%f", note.Freq()),
				})
				currKey = e
			} else if octSwitch {
				if note, ok := keys2notes[currKey]; ok {
					note.octave = oct
					p.assignedKeyComp.Execute(wavx.Command{
						"freq": fmt.Sprintf("%f", note.Freq()),
					})
				}
			} else if e == ' ' {
				p.assignedKeyComp.Deactivate()
				active = false
			} else if kb, ok := p.keyBindings[e]; ok {
				log.Infof("send command to comp")
				kb.Component.Execute(kb.Command)
			}
		}
	} else {
		<-ctx.Done()
	}
}

func (p *Project) Stop() error {
	if p.synth == nil {
		return errors.Errorf("synth is not running")
	}
	defer func() { p.synth = nil }()
	return p.synth.Close()
}

func (p *Project) ActivateComponent(compName string) error {
	comp, ok := p.components[compName]
	if !ok {
		return errors.Errorf("no such component %q", compName)
	}
	comp.Activate()
	return nil
}

func (p *Project) DeactivateComponent(compName string) error {
	comp, ok := p.components[compName]
	if !ok {
		return errors.Errorf("no such component %q", compName)
	}
	comp.Deactivate()
	return nil
}

func (p *Project) Execute(compName string, cmd wavx.Command) error {
	comp, ok := p.components[compName]
	if !ok {
		return errors.Errorf("no such component %q", compName)
	}
	comp.Execute(cmd)
	return nil
}

func (p *Project) AssignKeys(compName string) error {
	comp, ok := p.components[compName]
	if !ok {
		return errors.Errorf("no such component %q", compName)
	}
	p.assignedKeyComp = comp
	return nil
}

func (p *Project) BindKey(key rune, compName string, cmd wavx.Command) error {
	comp, ok := p.components[compName]
	if !ok {
		return errors.Errorf("no such component %q", compName)
	}
	p.keyBindings[key] = KeyBinding{
		Key:       key,
		Component: comp,
		Command:   cmd,
	}
	return nil
}
