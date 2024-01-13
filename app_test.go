package applife

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	tl := new(testLog)
	ctx := context.Background()

	ap := NewApp(ctx, tl)

	processes := []string{"process 1", "process 2", "process 3"}

	for _, p := range processes {
		ap.AddProcess(p, func(ctx context.Context) {
			<-ctx.Done()
		})
	}

	chDone := make(chan struct{})

	go func() {
		ap.Run()
		chDone <- struct{}{}
	}()

	time.Sleep(time.Millisecond * 200)

	sendSignal(t, syscall.SIGTERM)

	select {
	case <-chDone:
	case <-time.After(time.Millisecond * 200):
		t.Error("timeout waiting for termination")
	}

	assert.Equal(t, len(processes), len(tl.Started))
	assert.Equal(t, len(processes), len(tl.Stopped))
	for _, p := range processes {
		assert.Contains(t, tl.Started, p)
		assert.Contains(t, tl.Stopped, p)
	}
}

func sendSignal(t *testing.T, sig os.Signal) {
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	err = process.Signal(sig)
	if err != nil {
		t.Fatal(err)
	}
}

type testLog struct {
	T       *testing.T
	Started []string
	Stopped []string
}

func (t *testLog) OnError(_ context.Context, process string, err error) {
	fmt.Println("error", process, err)
}

func (t *testLog) OnControl(_ context.Context, process string, ct ControlType) {
	fmt.Println(process, ct)
	if ct == ControlTypeStarted {
		t.Started = append(t.Started, process)
	} else if ct == ControlTypeStopped {
		t.Stopped = append(t.Stopped, process)
	}
}

func (t *testLog) OnInfo(_ context.Context, msg string) {
	fmt.Println(msg)
}
