package keys

import (
	"context"

	"github.com/mattn/go-tty"
	"github.com/mazzegi/log"
	"github.com/pkg/errors"
)

func NewListener() (*Listener, error) {
	tty, err := tty.Open()
	if err != nil {
		return nil, errors.Wrap(err, "open-ty")
	}
	return &Listener{tty: tty}, nil
}

type Listener struct {
	tty *tty.TTY
}

func (l *Listener) ListenCtx(ctx context.Context) <-chan rune {
	c := make(chan rune)
	go func() {
		log.Infof("listen for keys ...")
		defer func() {
			close(c)
			log.Infof("listener is done")
		}()

		go func() {
			for {
				r, err := l.tty.ReadRune()
				if err != nil {
					log.Errorf("read-rune: %v", err)
					return
				}
				c <- r
			}
		}()

		<-ctx.Done()
		log.Infof("ctx is done")
		l.tty.Close()
	}()
	return c
}
