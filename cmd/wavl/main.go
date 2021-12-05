package main

import (
	"bytes"
	"context"
	"os"
	"os/signal"

	"github.com/mazzegi/log"
	"github.com/mazzegi/wavx/wavl"
)

func handleErr(err error) {
	if err == nil {
		return
	}
	log.Errorf("%v", err)
	os.Exit(1)
}

var commands = `
add osc osc1 sine 110 1 1
add osc osc2 saw  220 1 1
add mixer mix
connect osc1 mix:signal
connect osc2 mix:signal

output mix
sleep 2s
deactivate osc2
sleep 2s
activate osc2

set osc1 freq:110
sleep 500ms
set osc1 freq:+20
sleep 500ms
set osc1 freq:+20
sleep 500ms
set osc1 freq:+20
sleep 500ms
set osc1 freq:-10
sleep 500ms
set osc1 freq:-10
`

var commandsOvertones = `
add osc osc1 saw 110 1 3

output osc1

sleep 2s
set osc1 overtones:1
sleep 500ms
set osc1 overtones:2
sleep 500ms
set osc1 overtones:3
sleep 500ms
set osc1 overtones:4
sleep 500ms
set osc1 overtones:5
sleep 500ms
set osc1 overtones:6

sleep 1s
set osc1 freq:+40
`

var commandsKeys = `
add osc osc1 saw 110 1 1
#set osc1 overtones:2

add filter filter1 highpass 0.3 0.5
connect osc1 filter1:signal

assign_keys osc1
bind_key + filter1 cutoff:+0.01
bind_key - filter1 cutoff:-0.01

bind_key ö filter1 resonance:+0.01
bind_key ä filter1 resonance:-0.01

output filter1
sleep 2s
`

var commandsKeysPool = `
add oscpool pool

#add filter filter1 highpass 0.3 0.5
#connect pool filter1:signal

assign_keys pool
#bind_key + filter1 cutoff:+0.01
#bind_key - filter1 cutoff:-0.01
#bind_key ö filter1 resonance:+0.01
#bind_key ä filter1 resonance:-0.01

#output filter1
output pool
sleep 2s
`

func main() {
	//buf := bytes.NewBufferString(commandsKeys)
	buf := bytes.NewBufferString(commandsKeysPool)

	prj, err := wavl.Parse(buf)
	handleErr(err)

	prj.Start()
	defer prj.Stop()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	prj.Loop(ctx)

	//<-ctx.Done()
	log.Infof("done")
}
