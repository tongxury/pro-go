package biz

const (
	WorkflowStatusCreated   = "created"
	WorkflowStatusRunning   = "running"
	WorkflowStatusFailed    = "failed"
	WorkflowStatusCompleted = "completed"
	WorkflowStatusCanceled  = "canceled"
	WorkflowStatusPaused    = "paused"
)

const (
	JobStatusWaiting    = "waiting"
	JobStatusRunning    = "running"
	JobStatusConfirming = "confirming"
	JobStatusFailed     = "failed"
	JobStatusCompleted  = "completed"
	JobStatusCanceled   = "canceled"
)

const (
	ExecuteStatusWaiting                  = "waiting"
	ExecuteStatusCompleted                = "completed"
	ExecuteStatusCompletedAndSkipNextJobs = "completed_and_skip_next"
	ExecuteStatusFailed                   = "failed"
	ExecuteStatusRunning                  = "running"
)
