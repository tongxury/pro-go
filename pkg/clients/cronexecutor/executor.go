package cronexecutor

import "context"

// IExecutor defines the interface for fetching multiple tasks to be executed.
type IExecutor interface {
	LoadTasks(ctx context.Context) ([]ITask, error)
}

type Executor struct {
	function func(ctx context.Context) ([]ITask, error)
}

func (t Executor) LoadTasks(ctx context.Context) ([]ITask, error) {
	return t.function(ctx)
}

func NewExecutor(f func(ctx context.Context) ([]ITask, error)) IExecutor {
	return Executor{
		function: f,
	}
}
