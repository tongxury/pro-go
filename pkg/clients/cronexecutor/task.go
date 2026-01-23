package cronexecutor

import (
	"context"
)

// ITask defines the interface for an individual task.
type ITask interface {
	ID() string
	Execute(ctx context.Context) error
}

type Task struct {
	key      string
	function func(ctx context.Context) error
}

func (t Task) ID() string {
	return t.key
}

func (t Task) Execute(ctx context.Context) error {
	return t.function(ctx)
}

func NewTask(key string, function func(ctx context.Context) error) ITask {
	return Task{
		key:      key,
		function: function,
	}
}
