package applife

import "context"

type ControlType string

const (
	ControlTypeStarted ControlType = "running"
	ControlTypeStopped ControlType = "stopped"
)

// Logger can be used to receive control messages from the executing App processes.
type Logger interface {
	OnError(ctx context.Context, process string, err error)
	OnControl(ctx context.Context, process string, ct ControlType)
	OnInfo(ctx context.Context, msg string)
}
