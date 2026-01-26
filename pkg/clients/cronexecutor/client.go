package cronexecutor

import (
	"context"
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
)

// Client is the entry point for the cron task executor.
type Client struct {
	c      *cron.Cron
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// NewClient creates a new Client instance with second-precision cron support.
func NewClient() *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		c:      cron.New(cron.WithSeconds()),
		ctx:    ctx,
		cancel: cancel,
	}
}

// JobWrapper is a shim that implements the cron.Job interface.
type JobWrapper struct {
	state     sync.Map        // Tracks currently running task IDs to prevent duplicates.
	executor  IExecutor       // The executor that loads tasks.
	wg        *sync.WaitGroup // Reference to the Client's WaitGroup for tracking in-flight tasks.
	semaphore chan struct{}   // Limits concurrent tasks for this specific job.
	ctx       context.Context // Shared context from Client
}

// Run is called by the cron scheduler.
func (t *JobWrapper) Run() {
	// 0. Early return if context is already cancelled (e.g., during shutdown).
	select {
	case <-t.ctx.Done():
		return
	default:
	}

	// 1. Fetch available tasks using the shared context.
	tasks, err := t.executor.LoadTasks(t.ctx)
	if err != nil {
		fmt.Printf("failed to load tasks: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		return
	}

	for _, task := range tasks {
		// Validate task is not nil
		if task == nil {
			fmt.Printf("warning: skipping nil task\n")
			continue
		}

		taskID := task.ID()

		// Validate task ID is not empty
		if taskID == "" {
			fmt.Printf("warning: skipping task with empty ID\n")
			continue
		}

		// 2. Check if this specific task is already running (Single-Instance Enforcement).
		if _, loaded := t.state.LoadOrStore(taskID, struct{}{}); loaded {
			//fmt.Printf("warning: skipping task with already loaded ID: %s\n", taskID)
			continue
		}

		// 3. Try to acquire semaphore (Concurrency Limit).
		select {
		case t.semaphore <- struct{}{}:
			// Successfully acquired
		default:
			// Concurrency limit reached, clean up state for this ID and skip.
			t.state.Delete(taskID)
			continue
		}

		// 4. Start execution in a new goroutine.
		t.wg.Add(1)
		go func(tk ITask) {
			defer t.wg.Done()
			defer func() {
				// Panic recovery to ensure resources are always released.
				if r := recover(); r != nil {
					fmt.Printf("task %s panicked: %v\n", tk.ID(), r)
				}
				<-t.semaphore           // Release semaphore
				t.state.Delete(tk.ID()) // Release single-instance lock
			}()

			// Pass the shared context to the task, allowing it to respond to shutdown signals.
			if err := tk.Execute(t.ctx); err != nil {
				fmt.Printf("task %s execution failed: %v\n", tk.ID(), err)
			}
		}(task)
	}
}

// RegisterOptions allows configuring the behavior of a registered executor.
type RegisterOptions struct {
	Concurrency int // Maximum number of concurrent tasks. Defaults to 20 if <= 0.
}

// Register adds a new executor to the cron schedule with optional configuration.
// Returns an error if the cron spec is invalid.
func (t *Client) Register(spec string, executor IExecutor, opts ...RegisterOptions) error {
	concurrency := 100
	if len(opts) > 0 && opts[0].Concurrency > 0 {
		concurrency = opts[0].Concurrency
	}

	jw := &JobWrapper{
		executor:  executor,
		wg:        &t.wg,
		semaphore: make(chan struct{}, concurrency),
		ctx:       t.ctx,
	}

	_, err := t.c.AddJob(spec, jw)

	return err
}

// Start begins the cron scheduler.
func (t *Client) Start() {
	t.c.Start()
}

// Stop gracefully shuts down the scheduler and waits for all in-flight tasks to complete.
func (t *Client) Stop() {
	// 1. Signal all running tasks to abort via context.
	t.cancel()

	// 2. Tell the cron scheduler to stop starting new jobs.
	stopCtx := t.c.Stop()
	<-stopCtx.Done()

	// 3. Wait for all background goroutines to finish.
	t.wg.Wait()
}
