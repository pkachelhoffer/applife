package applife

import "context"

type FnProcess func(ctx context.Context)

type Process struct {
	Name      string
	fnProcess FnProcess

	chStopped chan struct{}
}

func (p Process) Stopped() {
	p.chStopped <- struct{}{}
}
