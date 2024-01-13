package applife

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// App collection of processes which will be executed and terminated with the app host process.
type App struct {
	processes []Process

	ctx      context.Context
	fnCancel context.CancelFunc

	logger Logger
}

// NewApp creates new App collection. Use AddProcess to add processes to the app which will be executed once Run
// is called.
func NewApp(ctx context.Context, logger Logger) *App {
	ctx, fnCancel := context.WithCancel(ctx)
	app := App{
		processes: make([]Process, 0),
		ctx:       ctx,
		fnCancel:  fnCancel,
		logger:    logger,
	}

	return &app
}

// AddProcess adds process to app that will be executed once Run is called.
func (a *App) AddProcess(name string, fnProcess FnProcess) {
	chStopped := make(chan struct{})
	process := Process{
		Name:      name,
		fnProcess: fnProcess,
		chStopped: chStopped,
	}

	a.processes = append(a.processes, process)
}

// Run starts all added processes and blocks until host process sends terminate/interrupt signal.
func (a *App) Run() {
	for _, p := range a.processes {
		p := p
		go func() {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						a.logError(a.ctx, p.Name, err)
					} else {
						err = errors.New("panic occurred")
						a.logError(a.ctx, p.Name, err)
					}
				}
				a.logControl(a.ctx, p.Name, ControlTypeStopped)
				p.Stopped()
			}()
			a.logControl(a.ctx, p.Name, ControlTypeStarted)
			p.fnProcess(a.ctx)
		}()
	}

	a.waitForExit()
	a.waitForProcesses()
}

// waitForExit blocks until host process signalled interrupt or terminate. Cancels context of each running process.
func (a *App) waitForExit() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	a.logInfo(a.ctx, fmt.Sprintf("receive os signal %s", sig))
	a.fnCancel()
}

// waitForProcesses waits until each process signalled on chStopped
func (a *App) waitForProcesses() {
	for _, p := range a.processes {
		<-p.chStopped
	}
}

func (a *App) logError(ctx context.Context, process string, err error) {
	if a.logger != nil {
		a.logger.OnError(ctx, process, err)
	}
}

func (a *App) logControl(ctx context.Context, process string, ct ControlType) {
	if a.logger != nil {
		a.logger.OnControl(ctx, process, ct)
	}
}

func (a *App) logInfo(ctx context.Context, msg string) {
	if a.logger != nil {
		a.logger.OnInfo(ctx, msg)
	}
}
