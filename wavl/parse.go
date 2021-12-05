package wavl

import (
	"bufio"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mazzegi/wavx"

	"github.com/pkg/errors"
)

func Parse(r io.Reader) (*Project, error) {
	var commands []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		cmd := strings.Trim(scanner.Text(), " \n\r\t")
		if cmd == "" || strings.HasPrefix(cmd, "#") {
			continue
		}
		commands = append(commands, cmd)
	}
	return newParser(commands).parse()
}

func ParseFile(fileName string) (*Project, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "open file %q", fileName)
	}
	defer f.Close()
	return Parse(f)
}

type parser struct {
	commands []string
}

func newParser(commands []string) *parser {
	return &parser{
		commands: commands,
	}
}

type projectFunc func(p *Project) error

func (p *parser) parse() (*Project, error) {
	prj := NewProject()
	for _, cmd := range p.commands {
		prjFunc, err := p.parseCommand(cmd)
		if err != nil {
			return nil, errors.Wrapf(err, "parse command %q", cmd)
		}
		err = prjFunc(prj)
		if err != nil {
			return nil, errors.Wrapf(err, "exec project-func for cmd %q", cmd)
		}
	}

	return prj, nil
}

func (p *parser) parseCommand(cmd string) (projectFunc, error) {
	sl := strings.Split(cmd, " ")
	var items []string
	for _, s := range sl {
		item := strings.ToLower(strings.Trim(s, " \r\n\t"))
		if item == "" {
			continue
		}
		items = append(items, item)
	}

	prefix := firstItem(items)
	rest := spliceItems(items, 0)
	switch prefix {
	case "add":
		return p.parseAdd(rest)
	case "connect":
		return p.parseConnect(rest)
	case "output":
		return p.parseOutput(rest)
	case "sleep":
		return p.parseSleep(rest)
	case "activate":
		return p.parseActivate(rest)
	case "deactivate":
		return p.parseDeactivate(rest)
	case "set":
		return p.parseSet(rest)
	case "assign_keys":
		return p.parseAssignKeys(rest)
	case "bind_key":
		return p.parseBindKey(rest)
	default:
		return nil, errors.Errorf("invalid prefix %q", prefix)
	}
}

func (p *parser) parseAdd(items []string) (projectFunc, error) {
	comp := itemAt(items, 0)
	name := itemAt(items, 1)
	if name == "" {
		return nil, errors.Errorf("component name is empty")
	}
	rest := spliceItems(items, 1)
	switch comp {
	case "osc":
		return p.parseAddOsc(name, rest)
	case "mixer":
		return p.parseAddMixer(name, rest)
	case "filter":
		return p.parseAddFilter(name, rest)
	case "oscpool":
		return p.parseAddOscPool(name, rest)
	default:
		return nil, errors.Errorf("unknown component %q", comp)
	}
}

func (p *parser) parseConnect(items []string) (projectFunc, error) {
	var (
		from string
		to   string
	)
	err := scanItems(items, &from, &to)
	if err != nil {
		return nil, errors.Wrapf(err, "parse-connect: scan items %v", items)
	}
	toSl := strings.SplitN(to, ":", 2)
	var toName string
	var toInput = "signal"
	if len(toSl) == 1 {
		toName = toSl[0]
	} else {
		toName = toSl[0]
		toInput = toSl[1]
	}

	return func(prj *Project) error {
		return prj.Connect(from, toName, toInput)
	}, nil
}

func (p *parser) parseOutput(items []string) (projectFunc, error) {
	var (
		from string
	)
	err := scanItems(items, &from)
	if err != nil {
		return nil, errors.Wrapf(err, "parse-output: scan items %v", items)
	}

	return func(prj *Project) error {
		return prj.OutputFrom(from)
	}, nil
}

func (p *parser) parseAddOsc(name string, items []string) (projectFunc, error) {
	var (
		typ       string
		freq      float64
		ampl      float64
		overtones int
	)
	err := scanItems(items, &typ, &freq, &ampl, &overtones)
	if err != nil {
		return nil, errors.Wrapf(err, "parse-add-osc: scan items %v", items)
	}

	return func(prj *Project) error {
		return prj.AddOscillator(name, typ, freq, ampl, overtones)
	}, nil
}

func (p *parser) parseAddMixer(name string, items []string) (projectFunc, error) {
	return func(prj *Project) error {
		return prj.AddMixer(name)
	}, nil
}

func (p *parser) parseAddOscPool(name string, items []string) (projectFunc, error) {
	return func(prj *Project) error {
		return prj.AddOscillatorPool(name)
	}, nil
}

func (p *parser) parseAddFilter(name string, items []string) (projectFunc, error) {
	var (
		typ       string
		cutoff    float64
		resonance float64
	)
	err := scanItems(items, &typ, &cutoff, &resonance)
	if err != nil {
		return nil, errors.Wrapf(err, "parse-add-filter: scan items %v", items)
	}

	return func(prj *Project) error {
		return prj.AddFilter(name, typ, cutoff, resonance)
	}, nil
}

func (p *parser) parseSleep(items []string) (projectFunc, error) {
	dur, err := time.ParseDuration(firstItem(items))
	if err != nil {
		return nil, err
	}

	return func(prj *Project) error {
		prj.AddEvents(SleepEvent{Duration: dur})
		return nil
	}, nil
}

func (p *parser) parseActivate(items []string) (projectFunc, error) {
	return func(prj *Project) error {
		prj.AddEvents(
			ProjectEvent{
				projFunc: func(p *Project) { p.ActivateComponent(firstItem(items)) },
			},
		)
		return nil
	}, nil
}

func (p *parser) parseDeactivate(items []string) (projectFunc, error) {
	return func(prj *Project) error {
		prj.AddEvents(
			ProjectEvent{
				projFunc: func(p *Project) { p.DeactivateComponent(firstItem(items)) },
			},
		)
		return nil
	}, nil
}

func (p *parser) parseSet(items []string) (projectFunc, error) {
	name := firstItem(items)
	cmd := wavx.Command{}
	for _, param := range spliceItems(items, 0) {
		sl := strings.Split(param, ":")
		if len(sl) != 2 {
			return nil, errors.Errorf("invalid param %q", param)
		}
		k := sl[0]
		v := sl[1]
		cmd[k] = v
	}

	return func(prj *Project) error {
		prj.AddEvents(
			ProjectEvent{
				projFunc: func(p *Project) { p.Execute(name, cmd) },
			},
		)
		return nil
	}, nil
}

func (p *parser) parseAssignKeys(items []string) (projectFunc, error) {
	return func(prj *Project) error {
		prj.AssignKeys(firstItem(items))
		return nil
	}, nil
}

/*
bind_key + filter1 cutoff:+0.1
*/

func (p *parser) parseBindKey(items []string) (projectFunc, error) {
	var (
		keyS  string
		name  string
		param string
	)
	scanItems(items, &keyS, &name, &param)
	rs := []rune(keyS)
	var key rune
	if len(rs) == 1 {
		key = rs[0]
	}

	cmd := wavx.Command{}
	sl := strings.Split(param, ":")
	if len(sl) != 2 {
		return nil, errors.Errorf("invalid param %q", param)
	}
	k := sl[0]
	v := sl[1]
	cmd[k] = v

	return func(prj *Project) error {
		prj.BindKey(key, name, cmd)
		return nil
	}, nil
}
