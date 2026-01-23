package cronexecutor

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

type mockTask struct {
	id      string
	execute func(ctx context.Context) error
}

func (m *mockTask) ID() string                        { return m.id }
func (m *mockTask) Execute(ctx context.Context) error { return m.execute(ctx) }

type mockExecutor struct {
	load func(ctx context.Context) ([]ITask, error)
}

func (m *mockExecutor) LoadTasks(ctx context.Context) ([]ITask, error) { return m.load(ctx) }

func TestCronExecutorEdgeCases(t *testing.T) {
	t.Run("invalid cron spec", func(t *testing.T) {
		client := NewClient()
		executor := &mockExecutor{
			load: func(ctx context.Context) ([]ITask, error) {
				return []ITask{}, nil
			},
		}

		err := client.Register("invalid spec", executor)
		if err == nil {
			t.Error("expected error for invalid cron spec")
		}
	})

	t.Run("nil task in list", func(t *testing.T) {
		var executed int32
		executor := &mockExecutor{
			load: func(ctx context.Context) ([]ITask, error) {
				return []ITask{
					nil,
					&mockTask{id: "valid_task", execute: func(ctx context.Context) error {
						atomic.StoreInt32(&executed, 1)
						return nil
					}},
				}, nil
			},
		}

		client := NewClient()
		err := client.Register("* * * * * *", executor)
		if err != nil {
			t.Fatalf("failed to register: %v", err)
		}
		client.Start()

		time.Sleep(1100 * time.Millisecond)
		client.Stop()

		if atomic.LoadInt32(&executed) != 1 {
			t.Error("valid task should have executed despite nil task in list")
		}
	})

	t.Run("empty task ID", func(t *testing.T) {
		var validExecuted int32
		executor := &mockExecutor{
			load: func(ctx context.Context) ([]ITask, error) {
				return []ITask{
					&mockTask{id: "", execute: func(ctx context.Context) error {
						t.Error("task with empty ID should not execute")
						return nil
					}},
					&mockTask{id: "valid_id", execute: func(ctx context.Context) error {
						atomic.StoreInt32(&validExecuted, 1)
						return nil
					}},
				}, nil
			},
		}

		client := NewClient()
		err := client.Register("* * * * * *", executor)
		if err != nil {
			t.Fatalf("failed to register: %v", err)
		}
		client.Start()

		time.Sleep(1100 * time.Millisecond)
		client.Stop()

		if atomic.LoadInt32(&validExecuted) != 1 {
			t.Error("valid task should have executed")
		}
	})

	t.Run("context cancelled before Run", func(t *testing.T) {
		var executed int32
		executor := &mockExecutor{
			load: func(ctx context.Context) ([]ITask, error) {
				atomic.AddInt32(&executed, 1)
				return []ITask{}, nil
			},
		}

		client := NewClient()
		err := client.Register("* * * * * *", executor)
		if err != nil {
			t.Fatalf("failed to register: %v", err)
		}
		client.Start()

		// Stop immediately to cancel context
		client.Stop()

		// LoadTasks should not have been called (or called very few times)
		if atomic.LoadInt32(&executed) > 2 {
			t.Errorf("LoadTasks should not be called after Stop, got %d calls", atomic.LoadInt32(&executed))
		}
	})

	t.Run("LoadTasks returns error", func(t *testing.T) {
		executor := &mockExecutor{
			load: func(ctx context.Context) ([]ITask, error) {
				return nil, errors.New("database error")
			},
		}

		client := NewClient()
		err := client.Register("* * * * * *", executor)
		if err != nil {
			t.Fatalf("failed to register: %v", err)
		}
		client.Start()

		time.Sleep(1100 * time.Millisecond)
		// Should not panic or crash
		client.Stop()
	})
}

func TestCronExecutorHardening(t *testing.T) {
	t.Run("panic recovery", func(t *testing.T) {
		recovered := make(chan struct{}, 1)
		executor := &mockExecutor{
			load: func(ctx context.Context) ([]ITask, error) {
				return []ITask{
					&mockTask{id: "panic_task", execute: func(ctx context.Context) error {
						panic("intentional panic")
					}},
				}, nil
			},
		}

		client := NewClient()
		err := client.Register("* * * * * *", executor)
		if err != nil {
			t.Fatalf("failed to register: %v", err)
		}
		client.Start()

		time.Sleep(1100 * time.Millisecond)
		client.Stop()

		// Test that task can run again after panic
		executor.load = func(ctx context.Context) ([]ITask, error) {
			return []ITask{
				&mockTask{id: "panic_task", execute: func(ctx context.Context) error {
					recovered <- struct{}{}
					return nil
				}},
			}, nil
		}

		client2 := NewClient()
		err = client2.Register("* * * * * *", executor)
		if err != nil {
			t.Fatalf("failed to register: %v", err)
		}
		client2.Start()

		select {
		case <-recovered:
			// Success
		case <-time.After(2 * time.Second):
			t.Error("task did not run again after panic")
		}
		client2.Stop()
	})

	t.Run("context cancellation", func(t *testing.T) {
		cancelled := make(chan struct{}, 1)
		executor := &mockExecutor{
			load: func(ctx context.Context) ([]ITask, error) {
				return []ITask{
					&mockTask{id: "context_task", execute: func(ctx context.Context) error {
						<-ctx.Done()
						cancelled <- struct{}{}
						return nil
					}},
				}, nil
			},
		}

		client := NewClient()
		err := client.Register("* * * * * *", executor)
		if err != nil {
			t.Fatalf("failed to register: %v", err)
		}
		client.Start()

		time.Sleep(1100 * time.Millisecond)
		client.Stop()

		select {
		case <-cancelled:
			// Success
		case <-time.After(1 * time.Second):
			t.Error("task did not receive context cancellation")
		}
	})

	t.Run("concurrency limit enforced", func(t *testing.T) {
		var activeTasks int32
		var maxParallel int32

		executor := &mockExecutor{
			load: func(ctx context.Context) ([]ITask, error) {
				var tasks []ITask
				for i := 0; i < 10; i++ {
					id := time.Now().Format(time.RFC3339Nano)
					tasks = append(tasks, &mockTask{id: id, execute: func(ctx context.Context) error {
						current := atomic.AddInt32(&activeTasks, 1)
						if current > atomic.LoadInt32(&maxParallel) {
							atomic.StoreInt32(&maxParallel, current)
						}
						time.Sleep(200 * time.Millisecond)
						atomic.AddInt32(&activeTasks, -1)
						return nil
					}})
					time.Sleep(1 * time.Millisecond)
				}
				return tasks, nil
			},
		}

		client := NewClient()
		err := client.Register("* * * * * *", executor, RegisterOptions{Concurrency: 3})
		if err != nil {
			t.Fatalf("failed to register: %v", err)
		}
		client.Start()

		time.Sleep(1500 * time.Millisecond)
		client.Stop()

		finalMax := atomic.LoadInt32(&maxParallel)
		if finalMax > 3 {
			t.Errorf("expected max 3 concurrent tasks, got %d", finalMax)
		}
		if finalMax == 0 {
			t.Error("expected some tasks to execute")
		}
	})

	t.Run("parallel execution", func(t *testing.T) {
		var activeTasks int32
		var maxParallel int32
		var totalExecutions int32

		executor := &mockExecutor{
			load: func(ctx context.Context) ([]ITask, error) {
				return []ITask{
					&mockTask{id: "task1", execute: func(ctx context.Context) error {
						atomic.AddInt32(&totalExecutions, 1)
						current := atomic.AddInt32(&activeTasks, 1)
						if current > atomic.LoadInt32(&maxParallel) {
							atomic.StoreInt32(&maxParallel, current)
						}
						time.Sleep(200 * time.Millisecond)
						atomic.AddInt32(&activeTasks, -1)
						return nil
					}},
					&mockTask{id: "task2", execute: func(ctx context.Context) error {
						atomic.AddInt32(&totalExecutions, 1)
						current := atomic.AddInt32(&activeTasks, 1)
						if current > atomic.LoadInt32(&maxParallel) {
							atomic.StoreInt32(&maxParallel, current)
						}
						time.Sleep(200 * time.Millisecond)
						atomic.AddInt32(&activeTasks, -1)
						return nil
					}},
				}, nil
			},
		}

		client := NewClient()
		err := client.Register("* * * * * *", executor)
		if err != nil {
			t.Fatalf("failed to register: %v", err)
		}
		client.Start()

		time.Sleep(1500 * time.Millisecond)
		client.Stop()

		if atomic.LoadInt32(&maxParallel) < 2 {
			t.Errorf("expected parallel execution, max parallel: %d", atomic.LoadInt32(&maxParallel))
		}
		if atomic.LoadInt32(&totalExecutions) < 2 {
			t.Errorf("expected multiple executions, got: %d", atomic.LoadInt32(&totalExecutions))
		}
	})
}
